package defaults

import (
	"database/sql"
	"fmt"
	"regexp"
	"testing"

	//external
	"github.com/DATA-DOG/go-sqlmock"
)

func TestInitSQLiteTables(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	type args struct {
		connection *sql.DB
	}
	tests := []struct {
		name string
		args args
		mock func()
	}{
		{
			name: "Success",
			args: args{
				connection: db,
			},
			mock: func() {
				sqlMock.ExpectExec(regexp.QuoteMeta(`
					CREATE TABLE IF NOT EXISTS card (
						code TEXT PRIMARY KEY,
						suit TEXT,
						'value' TEXT
					)
				`)).WillReturnResult(sqlmock.NewResult(1, 1))

				sqlMock.ExpectExec(regexp.QuoteMeta(`
					CREATE TABLE IF NOT EXISTS deck (
						id TEXT PRIMARY KEY,
						card_codes TEXT,
						is_shuffled INTEGER
					)
				`)).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			InitSQLiteTables(tt.args.connection)

			if err := sqlMock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestInsertDefaultCards(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	type args struct {
		connection *sql.DB
	}
	tests := []struct {
		name string
		args args
		mock func()
	}{
		{
			name: "Success",
			args: args{
				connection: db,
			},
			mock: func() {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
				sqlMock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM card")).WillReturnRows(rows)

				sqlMock.ExpectBegin()
				stmtMock := sqlMock.ExpectPrepare(regexp.QuoteMeta("INSERT INTO card (code, suit, 'value') VALUES (?,?,?)"))

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
					stmtMock.ExpectExec().WithArgs("A"+code, suit, "ACE").WillReturnResult(sqlmock.NewResult(1, 1))

					for j := 2; j <= 10; j++ {
						stmtMock.ExpectExec().WithArgs(fmt.Sprintf("%d%s", j, code), suit, fmt.Sprintf("%d", j)).WillReturnResult(sqlmock.NewResult(1, 1))
					}

					//insert JACK
					stmtMock.ExpectExec().WithArgs("J"+code, suit, "JACK").WillReturnResult(sqlmock.NewResult(1, 1))

					//insert QUEEN
					stmtMock.ExpectExec().WithArgs("Q"+code, suit, "QUEEN").WillReturnResult(sqlmock.NewResult(1, 1))

					//insert KING
					stmtMock.ExpectExec().WithArgs("K"+code, suit, "KING").WillReturnResult(sqlmock.NewResult(1, 1))
				}

				sqlMock.ExpectCommit()
				stmtMock.WillBeClosed()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			InsertDefaultCards(tt.args.connection)

			if err := sqlMock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}
