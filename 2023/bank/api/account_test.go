package api

import (
	"testing"
)

func TestSomething(t *testing.T) {
	testEnv := NewTestEnv(config)

	testEnv.db.Query("SELECT 1")

	testEnv2 := NewTestEnv(config)
	testEnv.db.Query("SELECT 1")

	testEnv.purge()
	testEnv2.purge()
}
