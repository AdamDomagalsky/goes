package api

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	db "github.com/AdamDomagalsky/goes/2023/bank/db/sqlc"
	"github.com/AdamDomagalsky/goes/2023/bank/util"
	"github.com/gin-gonic/gin"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq" // blank import: side-effect init pg driver

	"github.com/ory/dockertest/v3"
)

var (
	config  util.Config
	testEnv *TestEnv
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	testEnv = setup()
	exitCode := m.Run()
	testEnv.teardown()
	os.Exit(exitCode)
}

type TestEnv struct {
	config   util.Config
	db       *sql.DB
	server   *Server
	pool     *dockertest.Pool
	resource *dockertest.Resource
}

func (enviroment *TestEnv) teardown() {
	// You can't defer this because os.Exit doesn't care for defer
	if err := enviroment.pool.Purge(enviroment.resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}

func setup() *TestEnv {
	if testEnv != nil {
		testEnv.teardown()
	}
	var err error
	config, err = util.LoadConfig("..")
	if err != nil {
		log.Fatal("Cannot load env config:", err)
	}

	pool, conn, resource := setupTestDB(config)
	err = migrateUp(conn, config.DATABASE_NAME)
	if err != nil {
		fmt.Printf("failed migrateUp: %+v\n", err)
		err = pool.Purge(resource)
		if err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
		log.Fatalf("Could not migrate up: %s", err)

	}
	store := db.NewStore(conn)
	server := NewServer(store)

	return &TestEnv{
		config:   config,
		db:       conn,
		server:   server,
		pool:     pool,
		resource: resource,
	}
}

func setupTestDB(config util.Config) (pool *dockertest.Pool, db *sql.DB, resource *dockertest.Resource) {
	containerEnvs := []string{
		fmt.Sprintf("POSTGRES_PASSWORD=%s", config.DATABASE_PASSWORD),
		fmt.Sprintf("POSTGRES_USER=%s", config.DATABASE_USERNAME),
		fmt.Sprintf("POSTGRES_DB=%s", config.DATABASE_NAME),
		fmt.Sprintf("POSTGRES_HOST=%s", config.DATABASE_HOST),
		fmt.Sprintf("POSTGRES_PORT=%s", config.DATABASE_PORT),
	}

	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err = pool.Run("postgres", "15", containerEnvs)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		var err error
		datasourceName := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable",
			config.DATABASE_DRVIER,
			config.DATABASE_USERNAME,
			config.DATABASE_PASSWORD,
			config.DATABASE_HOST,
			resource.GetPort(fmt.Sprintf("%s/tcp", config.DATABASE_PORT)),
			config.DATABASE_NAME)
		db, err = sql.Open(config.DATABASE_DRVIER, datasourceName)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	return
}

func migrateUp(db *sql.DB, databaseName string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	mm, err := migrate.NewWithDatabaseInstance(
		"file://../db/migrations",
		databaseName,
		driver)
	if err != nil {
		return err
	}
	return mm.Up()
}
