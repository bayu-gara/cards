package app

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	//internal
	"github.com/bayu-gara/cards/internal/transport"

	//lib
	"github.com/bayu-gara/cards/pkg/config"
	database "github.com/bayu-gara/cards/pkg/database"

	//external
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
)

func TestRun(t *testing.T) {
	ctl := gomock.NewController(t)
	sqlConnectorMock := database.NewMockSQLConnector(ctl)
	serverMock := transport.NewMockServer(ctl)
	defer ctl.Finish()

	db, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	type args struct {
		mode string
	}
	tests := []struct {
		name    string
		args    args
		mock    func()
		wantErr bool
	}{
		{
			name: "Invalid DB Engine",
			args: args{
				mode: "rest",
			},
			mock: func() {
				config.GetConfig = func() config.Config {
					return config.Config{
						Database: config.DatabaseConfig{
							Engine: "mysql",
						},
					}
				}
			},
			wantErr: true,
		},
		{
			name: "Failed connect to DB",
			args: args{
				mode: "rest",
			},
			mock: func() {
				config.GetConfig = func() config.Config {
					return config.Config{
						Database: config.DatabaseConfig{
							Engine: "sqlite3",
						},
					}
				}

				database.GetConnector = func(engine string) database.SQLConnector {
					return sqlConnectorMock
				}

				sqlConnectorMock.EXPECT().Connect(config.DatabaseConfig{
					Engine: "sqlite3",
				}).Return(nil, errors.New("connection issue"))
			},
			wantErr: true,
		},
		{
			name: "Invalid transport mode",
			args: args{
				mode: "rest",
			},
			mock: func() {
				config.GetConfig = func() config.Config {
					return config.Config{
						Database: config.DatabaseConfig{
							Engine: "sqlite3",
						},
					}
				}

				database.GetConnector = func(engine string) database.SQLConnector {
					return sqlConnectorMock
				}

				sqlConnectorMock.EXPECT().Connect(config.DatabaseConfig{
					Engine: "sqlite3",
				}).Return(db, nil)

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

				transport.GetHandler = func(mode string) (server transport.Server, err error) {
					return nil, errors.New("Invalid mode")
				}

				sqlMock.ExpectClose()
			},
			wantErr: true,
		},
		{
			name: "Success",
			args: args{
				mode: "rest",
			},
			mock: func() {
				config.GetConfig = func() config.Config {
					return config.Config{
						Database: config.DatabaseConfig{
							Engine: "sqlite3",
						},
					}
				}

				db, sqlMock, err = sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				database.GetConnector = func(engine string) database.SQLConnector {
					return sqlConnectorMock
				}

				sqlConnectorMock.EXPECT().Connect(config.DatabaseConfig{
					Engine: "sqlite3",
				}).Return(db, nil)

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

				transport.GetHandler = func(mode string) (server transport.Server, err error) {
					return serverMock, nil
				}

				serverMock.EXPECT().Serve().Return(nil)
				sqlMock.ExpectClose()
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			if err := Run(tt.args.mode); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := sqlMock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}
