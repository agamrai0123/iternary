package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	// Try connecting with different passwords
	passwords := []string{"postgres", "password", "", "postgres_secure_2026"}
	var db *sql.DB
	var err error

	fmt.Println("Attempting to connect to PostgreSQL...")
	for _, pwd := range passwords {
		psqlInfo := fmt.Sprintf("host=localhost port=5432 user=postgres password=%s sslmode=disable", pwd)
		db, err = sql.Open("postgres", psqlInfo)
		if err == nil && db.Ping() == nil {
			fmt.Printf("✓ Connected with password: [%s]\n", pwd)
			break
		}
		if db != nil {
			db.Close()
		}
	}

	if db == nil || db.Ping() != nil {
		log.Fatalf("Failed to connect to PostgreSQL with any password")
	}

	defer db.Close()

	// Create database
	fmt.Println("\n[1/3] Creating database 'itinerary_production'...")
	_, err = db.Exec("CREATE DATABASE itinerary_production;")
	if err != nil {
		fmt.Printf("Note: %v (may already exist)\n", err)
	} else {
		fmt.Println("✓ Database created")
	}

	// Connect to new database
	fmt.Println("\n[2/3] Connecting to itinerary_production database...")
	psqlInfo := fmt.Sprintf("host=localhost port=5432 user=postgres password=postgres dbname=itinerary_production sslmode=disable")
	db2, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Printf("Connection error: %v\n", err)
		// Try with other passwords
		for _, pwd := range passwords {
			psqlInfo = fmt.Sprintf("host=localhost port=5432 user=postgres password=%s dbname=itinerary_production sslmode=disable", pwd)
			db2, err = sql.Open("postgres", psqlInfo)
			if err == nil && db2.Ping() == nil {
				fmt.Printf("✓ Connected with password: [%s]\n", pwd)
				break
			}
			if db2 != nil {
				db2.Close()
			}
		}
	}

	if db2 == nil || db2.Ping() != nil {
		log.Fatalf("Failed to connect to itinerary_production: %v", err)
	}
	defer db2.Close()
	fmt.Println("✓ Connected to itinerary_production")

	// Apply schema
	fmt.Println("\n[3/3] Applying PostgreSQL schema...")
	schema, err := ioutil.ReadFile("migrations/001_initial_schema.sql")
	if err != nil {
		log.Fatalf("Failed to read schema file: %v", err)
	}

	_, err = db2.Exec(string(schema))
	if err != nil {
		log.Fatalf("Failed to apply schema: %v", err)
	}
	fmt.Println("✓ Schema applied successfully")

	// Verify tables
	fmt.Println("\n[VERIFY] Checking created tables...")
	rows, err := db2.Query(`
		SELECT tablename FROM pg_tables 
		WHERE schemaname = 'public' 
		ORDER BY tablename
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			log.Fatal(err)
		}
		tables = append(tables, tableName)
	}

	fmt.Printf("\n✓ Created %d tables:\n", len(tables))
	for _, t := range tables {
		fmt.Printf("  - %s\n", t)
	}

	// Verify indexes
	rows2, err := db2.Query(`
		SELECT COUNT(*) as index_count 
		FROM pg_indexes 
		WHERE schemaname = 'public'
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows2.Close()

	var indexCount int
	if rows2.Next() {
		rows2.Scan(&indexCount)
	}
	fmt.Printf("\n✓ Created %d indexes\n", indexCount)

	fmt.Println("\n✅ Database initialization complete!")
}
