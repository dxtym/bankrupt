package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/dxtym/bankrupt/utils"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

// test db conenction
func TestMain(m *testing.M) {
	config, err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	testDB, err = sql.Open(config.Driver, config.Source)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}
