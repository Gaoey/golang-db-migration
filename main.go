package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
	"tutorial/dbvcs/services/payment"
	"tutorial/dbvcs/services/transfer"

	"github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	migrateSQL "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/viper"
)

func initialize() {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("cannot read in viper config:%s", err)
	}

	viper.AutomaticEnv()
}

func initDB() (*sql.DB, error) {
	conf := mysql.Config{
		DBName:               viper.GetString("mariadb.database"),
		User:                 viper.GetString("mariadb.username"),
		Passwd:               viper.GetString("mariadb.password"),
		Net:                  "tcp",
		Addr:                 viper.GetString("mariadb.host") + ":" + viper.GetString("mariadb.port"),
		AllowNativePasswords: true,
		Timeout:              viper.GetDuration("mariadb.timeout"),
		ReadTimeout:          viper.GetDuration("mariadb.readtimeout"),
		WriteTimeout:         viper.GetDuration("mariadb.writetimeout"),
		ParseTime:            viper.GetBool("mariadb.parsetime"),
		MultiStatements:      viper.GetBool("mariadb.multistatement"),
		Loc:                  time.Local,
	}

	db, err := sql.Open("mysql", conf.FormatDSN())
	if err != nil {
		return nil, err
	}

	return db, nil
}

var tableList = []string{"payment", "transfer"}

func main() {
	initialize()

	forceFlag := flag.Int("force", -1, "force version")
	stepFlag := flag.String("step", "", "step up or down")
	tableFlag := flag.String("table", "", "table selected")
	onlyStateFlag := flag.Bool("only-state", false, "migrate only database state")

	flag.Parse()

	if *stepFlag == "" {
		log.Fatalf("\nmandatory command -step \nexample -step=\"up\" or  -step=\"down\"")
	}

	if *tableFlag == "" && !contains(tableList, *tableFlag) {
		log.Fatalf("\nplease selected table 'payment' or 'transfer' to migration \nexample -table=\"payment\"")
	}

	db, err := initDB()
	if err != nil {
		log.Fatalf("could not connect to the MySQL database... %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("could not ping DB... %v", err)
	}

	driver, err := migrateSQL.WithInstance(db, &migrateSQL.Config{})
	if err != nil {
		log.Fatalf("could not start sql migration... %v", err)
	}

	migrationFile := fmt.Sprintf("file://migration/%s", *tableFlag)
	m, err := migrate.NewWithDatabaseInstance(migrationFile, "mysql", driver)
	if err != nil {
		log.Fatalf("migration failed... %v", err)
	}

	if *forceFlag >= 0 {
		m.Force(*forceFlag)
		log.Println("Forced to version %f", *forceFlag)
	}

	// up schema migration
	if *stepFlag == "up" && !*onlyStateFlag {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("An error occurred while syncing the database.. %v", err)
		}
		log.Println("Database up migrated")
	}

	// down schema migration
	if *stepFlag == "down" && !*onlyStateFlag {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("An error occurred while syncing the database.. %v", err)
		}
		log.Println("Database down migrated")
	}

	dbStateMigration(*tableFlag, *stepFlag)

	os.Exit(0)
}

func dbStateMigration(table string, step string) {
	if table == "payment" && step == "up" {
		payment.Up()
	}

	if table == "payment" && step == "down" {
		payment.Down()
	}

	if table == "transfer" && step == "up" {
		transfer.Up()
	}

	if table == "transfer" && step == "down" {
		transfer.Down()
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
