package seeders

import (
	"database/sql"
	"fmt"
	"log"
)

func SeedUsers(db *sql.DB) {
	_, err := db.Exec(`
        INSERT INTO users (name, email, age, created_at)
VALUES
    ('John Doe', 'john.doe@example.com', 30, now()),
    ('Jane Smith', 'jane.smith@example.com', 25, now()),
    ('Emily Johnson', 'emily.johnson@example.com', 22, now()),
    ('Michael Brown', 'michael.brown@example.com', 35, now()),
    ('Linda Williams', 'linda.williams@example.com', 28, now()),
    ('David Miller', 'david.miller@example.com', 40, now());
    `)
	if err != nil {
		log.Fatalf("Error seeding products table: %v", err)
	}
	fmt.Println("user data seeded successfully")
}
