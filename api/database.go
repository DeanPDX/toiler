package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

var dbPool *pgxpool.Pool

func initializeDB(connectionString string) {
	fmt.Println("\nInitializing DB")
	pool, err := pgxpool.Connect(context.Background(), connectionString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	dbPool = pool

	applyMigrations()
}

func closeDB() {
	dbPool.Close()
}

func applyMigrations() {
	fmt.Println("\nRunning migrations")
	// Create migrations table if not exists
	mustExec(`CREATE TABLE IF NOT EXISTS db_migrations (
		file_name varchar (255),
		applied_on timestamp not null
	)`)

	// TODO: move this to configuration.
	migrationsFolder := "./database/migrations"

	files, err := os.ReadDir(migrationsFolder)
	if err != nil {
		log.Fatal(err)
	}

	applied := 0
	skipped := 0

	for _, f := range files {
		fileName := f.Name()
		// If we've already run this migration, nothing to do.
		if migrationRun(fileName) {
			skipped++
			continue
		}
		applied++
		// Open file and read contents
		data, err := os.ReadFile(migrationsFolder + "/" + fileName)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("\nApplying ", fileName)
		// Execute the SQL in the file and blow up if anything goes wrong.
		mustExec(string(data))
		// Add to migrations table so we know we ran it
		logMigrationRun(fileName)
	}

	fmt.Println("\nApplied", applied, "migrations and skipped", skipped)
}

// Execute SQL. If something fails, log.Fatal. Use with caution when something really must run without fail for the app to continue execution.
func mustExec(sql string) {
	if _, err := dbPool.Exec(context.Background(), sql); err != nil {
		log.Fatal(err)
	}
}

func migrationRun(fileName string) bool {
	var count int
	err := dbPool.QueryRow(context.Background(), `select count(*) from db_migrations where file_name = $1;`, fileName).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	return count > 0
}

func logMigrationRun(fileName string) {
	_, err := dbPool.Exec(context.Background(), `insert into db_migrations(file_name, applied_on) values ($1, $2);`, fileName, time.Now())
	if err != nil {
		log.Fatal(err)
	}
}
