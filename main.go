package main

import (
	"database/sql"
	"flag"
	"log"
	"os"
	"strconv"
	"time"

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

func main() {
	initialize()

	forceFlag := flag.String("force", "-1", "force version")
	stepFlag := flag.String("step", "", "up version")

	flag.Parse()

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

	m, err := migrate.NewWithDatabaseInstance("file://migration/", "mysql", driver)
	if err != nil {
		log.Fatalf("migration failed... %v", err)
	}

	// force version
	forceVersion, err := strconv.Atoi(*forceFlag)
	if err != nil {
		log.Fatalf("force version not correctly", err)
	}

	if forceVersion >= 0 {
		m.Force(forceVersion)
		log.Println("Forced to version %f", forceVersion)
	}

	// up migration

	if *stepFlag == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("An error occurred while syncing the database.. %v", err)
		}
		log.Println("Database up migrated")
	}

	// up migration
	if *stepFlag == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("An error occurred while syncing the database.. %v", err)
		}
		log.Println("Database down migrated")
	}

	os.Exit(0)
}
