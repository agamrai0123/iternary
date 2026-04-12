package main

import (
	"database/sql"
	"fmt"
	"log"
)

func verifyDataIntegrity(sqliteDB, pgDB *sql.DB) {
	// Tables to verify
	tables := []string{
		"users",
		"destinations",
		"itineraries",
		"itinerary_items",
		"comments",
		"likes",
		"user_plans",
		"user_trips",
		"trip_segments",
		"trip_photos",
		"trip_reviews",
		"user_trip_posts",
	}

	fmt.Println("  Data Integrity Check:")
	allMatch := true

	for _, table := range tables {
		var sqliteCount, pgCount int64

		// Count SQLite
		err := sqliteDB.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", table)).Scan(&sqliteCount)
		if err != nil {
			log.Printf("    ✗ %s: Error querying SQLite: %v", table, err)
			allMatch = false
			continue
		}

		// Count PostgreSQL
		err = pgDB.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", table)).Scan(&pgCount)
		if err != nil {
			log.Printf("    ✗ %s: Error querying PostgreSQL: %v", table, err)
			allMatch = false
			continue
		}

		// Compare
		if sqliteCount == pgCount {
			fmt.Printf("    ✓ %s: %d rows verified (SQLite: %d → PostgreSQL: %d)\n", table, pgCount, sqliteCount, pgCount)
		} else {
			fmt.Printf("    ✗ %s: MISMATCH! SQLite: %d vs PostgreSQL: %d\n", table, sqliteCount, pgCount)
			allMatch = false
		}
	}

	if allMatch {
		fmt.Println("  ✓ All tables verified successfully!")
	} else {
		fmt.Println("  ⚠ Some verification issues detected (check warnings above)")
	}
}

func validateRelationships(pgDB *sql.DB) {
	fmt.Println("  Relationship Validation:")

	// Check referential integrity
	checks := []struct {
		name  string
		query string
	}{
		{
			name: "Itineraries have valid users",
			query: `SELECT COUNT(*) FROM itineraries i 
					WHERE NOT EXISTS (SELECT 1 FROM users u WHERE u.id = i.user_id)`,
		},
		{
			name: "Itineraries have valid destinations",
			query: `SELECT COUNT(*) FROM itineraries i 
					WHERE NOT EXISTS (SELECT 1 FROM destinations d WHERE d.id = i.destination_id)`,
		},
		{
			name: "Itinerary items have valid itineraries",
			query: `SELECT COUNT(*) FROM itinerary_items ii 
					WHERE NOT EXISTS (SELECT 1 FROM itineraries i WHERE i.id = ii.itinerary_id)`,
		},
		{
			name: "Comments have valid users",
			query: `SELECT COUNT(*) FROM comments c 
					WHERE NOT EXISTS (SELECT 1 FROM users u WHERE u.id = c.user_id)`,
		},
		{
			name: "Comments have valid itineraries",
			query: `SELECT COUNT(*) FROM comments c 
					WHERE NOT EXISTS (SELECT 1 FROM itineraries i WHERE i.id = c.itinerary_id)`,
		},
		{
			name: "Likes have valid users",
			query: `SELECT COUNT(*) FROM likes l 
					WHERE NOT EXISTS (SELECT 1 FROM users u WHERE u.id = l.user_id)`,
		},
		{
			name: "Likes have valid itineraries",
			query: `SELECT COUNT(*) FROM likes l 
					WHERE NOT EXISTS (SELECT 1 FROM itineraries i WHERE i.id = l.itinerary_id)`,
		},
		{
			name: "User plans have valid users",
			query: `SELECT COUNT(*) FROM user_plans up 
					WHERE NOT EXISTS (SELECT 1 FROM users u WHERE u.id = up.user_id)`,
		},
		{
			name: "User trips have valid users",
			query: `SELECT COUNT(*) FROM user_trips ut 
					WHERE NOT EXISTS (SELECT 1 FROM users u WHERE u.id = ut.user_id)`,
		},
		{
			name: "Trip segments have valid user trips",
			query: `SELECT COUNT(*) FROM trip_segments ts 
					WHERE NOT EXISTS (SELECT 1 FROM user_trips ut WHERE ut.id = ts.user_trip_id)`,
		},
		{
			name: "Trip photos have valid user trips",
			query: `SELECT COUNT(*) FROM trip_photos tp 
					WHERE NOT EXISTS (SELECT 1 FROM user_trips ut WHERE ut.id = tp.user_trip_id)`,
		},
		{
			name: "Trip reviews have valid user trips",
			query: `SELECT COUNT(*) FROM trip_reviews tr 
					WHERE NOT EXISTS (SELECT 1 FROM user_trips ut WHERE ut.id = tr.user_trip_id)`,
		},
		{
			name: "User trip posts have valid user trips",
			query: `SELECT COUNT(*) FROM user_trip_posts utp 
					WHERE NOT EXISTS (SELECT 1 FROM user_trips ut WHERE ut.id = utp.user_trip_id)`,
		},
		{
			name: "User trip posts have valid users",
			query: `SELECT COUNT(*) FROM user_trip_posts utp 
					WHERE NOT EXISTS (SELECT 1 FROM users u WHERE u.id = utp.user_id)`,
		},
	}

	allValid := true
	for _, check := range checks {
		var orphans int64
		err := pgDB.QueryRow(check.query).Scan(&orphans)
		if err != nil {
			log.Printf("    ✗ %s: Error running check: %v", check.name, err)
			allValid = false
			continue
		}

		if orphans == 0 {
			fmt.Printf("    ✓ %s\n", check.name)
		} else {
			fmt.Printf("    ✗ %s: Found %d orphaned records!\n", check.name, orphans)
			allValid = false
		}
	}

	if allValid {
		fmt.Println("  ✓ All relationships validated successfully!")
	} else {
		fmt.Println("  ⚠ Some relationship issues detected (check warnings above)")
	}
}
