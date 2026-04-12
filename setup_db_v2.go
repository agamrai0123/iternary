package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"

	_ "github.com/lib/pq"
)

func tryConnect(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func main() {
	fmt.Println("╔════════════════════════════════════════╗")
	fmt.Println("║  Phase B Week 1 - Database Setup      ║")
	fmt.Println("║  Day 2: PostgreSQL Initialization    ║")
	fmt.Println("╚════════════════════════════════════════╝\n")

	connectionStrings := []struct {
		name    string
		connStr string
	}{
		{"Default (host connection)", "postgresql://localhost/postgres?sslmode=disable"},
		{"Explicit postgres user", "postgresql://postgres@localhost/postgres?sslmode=disable"},
		{"With password 'postgres'", "postgresql://postgres:postgres@localhost/postgres?sslmode=disable"},
		{"With password ''", "postgresql://postgres:@localhost/postgres?sslmode=disable"},
		{"TCP with port", "host=localhost port=5432 user=postgres sslmode=disable"},
		{"TCP with password", "host=localhost port=5432 user=postgres password=postgres sslmode=disable"},
	}

	var db *sql.DB
	var err error

	fmt.Println("[CONNECTING] Trying various PostgreSQL connection methods...")
	for _, cs := range connectionStrings {
		fmt.Printf("\n  Trying: %s\n", cs.name)
		db, err = tryConnect(cs.connStr)
		if err == nil {
			fmt.Printf("  ✓ SUCCESS: Connected!\n")
			goto connected
		}
		fmt.Printf("  ✗ Failed: %v\n", err)
	}

	fmt.Println("\n❌ Failed to connect to PostgreSQL with all methods")
	os.Exit(1)

connected:
	defer db.Close()

	fmt.Println("\n[STEP 1/3] Creating database 'itinerary_production'...")
	result, err := db.Exec("CREATE DATABASE itinerary_production;")
	if err != nil {
		if err.Error()[0:10] != "pq: error" {
			fmt.Printf("  Note: %v (may already exist)\n", err)
		}
	} else {
		fmt.Printf("  ✓ Database created: %v\n", result)
	}

	fmt.Println("\n[STEP 2/3] Connecting to itinerary_production database...")
	var db2 *sql.DB

	// Try connecting to the new database
	for _, cs := range []string{
		"postgresql://localhost/itinerary_production?sslmode=disable",
		"postgresql://postgres@localhost/itinerary_production?sslmode=disable",
		"postgresql://postgres:postgres@localhost/itinerary_production?sslmode=disable",
		"host=localhost port=5432 user=postgres password=postgres dbname=itinerary_production sslmode=disable",
	} {
		db2, err = tryConnect(cs)
		if err == nil {
			fmt.Println("  ✓ Connected to itinerary_production")
			break
		}
	}

	if db2 == nil {
		fmt.Printf("  ✗ Failed to connect to itinerary_production: %v\n", err)
		os.Exit(1)
	}
	defer db2.Close()

	fmt.Println("\n[STEP 3/3] Applying PostgreSQL schema...")

	// Find schema file
	schemaPath := "migrations/001_initial_schema.sql"
	if _, err := os.Stat(schemaPath); err != nil {
		schemaPath = "../migrations/001_initial_schema.sql"
	}

	schema, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		fmt.Printf("  ✗ Failed to read schema: %v\n", err)
		os.Exit(1)
	}

	_, err = db2.Exec(string(schema))
	if err != nil {
		fmt.Printf("  ✗ Failed to apply schema: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("  ✓ Schema applied successfully")

	// Verify
	fmt.Println("\n[VERIFY] Checking created tables...")
	rows, err := db2.Query(`
		SELECT tablename FROM pg_tables 
		WHERE schemaname = 'public' 
		ORDER BY tablename
	`)
	if err != nil {
		fmt.Printf("  ✗ Query failed: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			fmt.Printf("  ✗ Scan failed: %v\n", err)
			os.Exit(1)
		}
		tables = append(tables, tableName)
	}

	fmt.Printf("  ✓ Created %d tables:\n", len(tables))
	for _, t := range tables {
		fmt.Printf("    - %s\n", t)
	}

	// Count indexes
	row := db2.QueryRow(`
		SELECT COUNT(*) FROM pg_indexes WHERE schemaname = 'public'
	`)
	var indexCount int
	if err := row.Scan(&indexCount); err != nil {
		fmt.Printf("  ✗ Index count query failed: %v\n", err)
	} else {
		fmt.Printf("  ✓ Created %d indexes\n", indexCount)
	}

	fmt.Println("\n╔════════════════════════════════════════╗")
	fmt.Println("║  ✅ DATABASE INITIALIZATION COMPLETE  ║")
	fmt.Println("╚════════════════════════════════════════╝\n")

	fmt.Println("Database Details:")
	fmt.Println("  Name: itinerary_production")
	fmt.Println("  Tables: " + fmt.Sprintf("%d", len(tables)))
	fmt.Println("  Indexes: " + fmt.Sprintf("%d", indexCount))
	fmt.Println("\nReady for Day 2 data migration!")
}
