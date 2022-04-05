package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
)

var migrateCount *int

func dsn() string {
	dbhost := "127.0.0.1"
	dbport := os.Getenv("PG_PORT")
	dbusername := os.Getenv("PG_USER")
	dbpassword := os.Getenv("PG_PASS")
	dbname := os.Getenv("PG_DB")
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbusername, dbpassword, dbhost, dbport, dbname)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate the db schema",
	Long:  `Migrate the db schema required by auth`,
	Run: func(cmd *cobra.Command, args []string) {
		m, err := migrate.New(
			"file://db/migrations",
			dsn())
		if err != nil {
			log.Fatal("Failed to create migration:" + err.Error())
		}

		if *migrateCount == 0 {
			// All

			err := m.Up()
			if err == migrate.ErrNoChange {
				log.Printf("No Change")
				return
			}

			if err != nil {
				log.Fatal("Failed migration:" + err.Error())
			}
			log.Printf("Successfully Migrated All\n")

		} else {
			err := m.Steps(*migrateCount)

			if err == migrate.ErrNoChange {
				log.Printf("No Change")
				return
			}

			if err != nil {
				log.Fatal("failed migration:" + err.Error())
			}
			log.Printf("successfully migrated %d\n", *migrateCount)
		}
	},
}

func setupMigrate() {
	migrateCount = migrateCmd.Flags().IntP("count", "n", 0, "migrate --count [+-]N")
}
