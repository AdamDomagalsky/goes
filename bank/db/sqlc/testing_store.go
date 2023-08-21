package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/AdamDomagalsky/goes/2023/bank/util"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/lib/pq" // blank import: side-effect init pg driver
	"github.com/ory/dockertest"

	// "github.com/ory/dockertest"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq" // blank import: side-effect init pg driver
)

type TestStore struct {
	Config   util.Config
	db       *sql.DB
	Store    Store
	pool     *dockertest.Pool
	resource *dockertest.Resource
}

func (enviroment *TestStore) Teardown() {
	// You can't defer this because os.Exit doesn't care for defer
	if err := enviroment.pool.Purge(enviroment.resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}

// Similar to db.NewStore() but additionally:
// - creates a new test database via docker
// - runs all migrations up
// - provide a teardown function to clean up the database
func NewStoreTestStore(path string, testStore *TestStore) *TestStore {
	var config util.Config
	var err error

	if testStore != nil {
		testStore.Teardown()
		config = testStore.Config
	} else {
		config, err = util.LoadConfig(path)
		if err != nil {
			log.Fatal("Cannot load env config:", err)
		}
	}

	pool, conn, resource := setupTestDB(config)
	err = MigrateUp(conn, config.DATABASE_NAME)
	if err != nil {
		fmt.Printf("failed migrateUp: %+v\n", err)
		err = pool.Purge(resource)
		if err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
		log.Fatalf("Could not migrate up: %s", err)

	}
	store := NewStore(conn)

	return &TestStore{
		Config:   config,
		db:       conn,
		Store:    store,
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

func MigrateUp(db *sql.DB, databaseName string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	rootPath, err := getProjectRootPath()
	if err != nil {
		log.Fatal(err)
	}
	targetDir := "migrations"
	foundPath, err := findDirectory(rootPath, targetDir)
	if err != nil {
		log.Fatal(err)
	}
	mm, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", foundPath),
		databaseName,
		driver)
	if err != nil {
		return err
	}
	return mm.Up()
}

// findDirectory searches for a directory with the given name within the given rootPath.
func findDirectory(rootPath string, targetDir string) (string, error) {
	var foundPath string

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && info.Name() == targetDir {
			foundPath = path
			return filepath.SkipDir // Skip searching further within this directory
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	if foundPath != "" {
		return foundPath, nil
	} else {
		return "", fmt.Errorf("directory %s not found in %s", targetDir, rootPath)
	}
}

// returns root of the project based based on .env file
func getProjectRootPath() (string, error) {
	currentPath, _ := os.Getwd()
	for {
		configPath := filepath.Join(currentPath, ".env")
		if _, err := os.Stat(configPath); err == nil {
			return currentPath, nil
		}

		parentPath := filepath.Dir(currentPath)
		if parentPath == currentPath {
			break // Reached the filesystem root without finding the file
		}
		currentPath = parentPath
	}

	return "", fmt.Errorf("project root not found")
}
