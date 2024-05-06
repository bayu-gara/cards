package defaults

import (
	"database/sql"
	"fmt"
	"log"
)

func InitSQLiteTables(connection *sql.DB) {
	_, err := connection.Exec(`
		CREATE TABLE IF NOT EXISTS card (
			code TEXT PRIMARY KEY,
			suit TEXT,
			'value' TEXT
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = connection.Exec(`
		CREATE TABLE IF NOT EXISTS deck (
			id TEXT PRIMARY KEY,
			card_codes TEXT,
			is_shuffled INTEGER
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func InsertDefaultCards(connection *sql.DB) {
	var count int
	err := connection.QueryRow("SELECT COUNT(*) FROM card").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	if count > 0 {
		return
	}

	// Begin a transaction
	tx, err := connection.Begin()
	if err != nil {
		log.Fatal(err)
	}

	// Prepare the INSERT statement within the transaction
	stmt, err := tx.Prepare("INSERT INTO card (code, suit, 'value') VALUES (?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for i := 0; i < 4; i++ {
		code := ""
		suit := ""
		switch i {
		case 0:
			code = "S"
			suit = "SPADES"
		case 1:
			code = "H"
			suit = "HEARTS"
		case 2:
			code = "C"
			suit = "CLUBS"
		case 3:
			code = "D"
			suit = "DIAMONDS"
		}

		//insert ACE
		_, err = stmt.Exec("A"+code, suit, "ACE")
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}

		for j := 2; j <= 10; j++ {
			_, err = stmt.Exec(fmt.Sprintf("%d%s", j, code), suit, fmt.Sprintf("%d", j))
			if err != nil {
				tx.Rollback()
				log.Fatal(err)
			}
		}

		//insert JACK
		_, err = stmt.Exec("J"+code, suit, "JACK")
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}

		//insert QUEEN
		_, err = stmt.Exec("Q"+code, suit, "QUEEN")
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}

		//insert KING
		_, err = stmt.Exec("K"+code, suit, "KING")
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
	}

	// Commit the transaction if everything executed successfully
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
}
