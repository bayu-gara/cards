package repository

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"regexp"
	"testing"

	//domain
	"github.com/bayu-gara/cards/internal/domain/model"

	//external
	"github.com/DATA-DOG/go-sqlmock"
)

func TestNewSQLCardRepository(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	type args struct {
		db *sql.DB
	}
	tests := []struct {
		name string
		args args
		want CardRepository
	}{
		{
			name: "Success",
			args: args{db: db},
			want: &SQLCardRepo{DB: db},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSQLCardRepository(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSQLCardRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLCardRepo_Insert(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		ctx  context.Context
		card model.Card
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func()
		wantErr bool
	}{
		{
			name: "Failed to insert",
			fields: fields{
				DB: db,
			},
			args: args{
				ctx: context.Background(),
				card: model.Card{
					Code:  "AS",
					Suit:  "SPADES",
					Value: "ACE",
				},
			},
			mock: func() {
				sqlMock.ExpectExec(
					"INSERT INTO card",
				).WithArgs("AS", "SPADES", "ACE").WillReturnError(errors.New("connection issue"))
			},
			wantErr: true,
		},
		{
			name: "Success insert",
			fields: fields{
				DB: db,
			},
			args: args{
				ctx: context.Background(),
				card: model.Card{
					Code:  "AS",
					Suit:  "SPADES",
					Value: "ACE",
				},
			},
			mock: func() {
				sqlMock.ExpectExec(
					"INSERT INTO card",
				).WithArgs("AS", "SPADES", "ACE").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			repo := &SQLCardRepo{
				DB: tt.fields.DB,
			}
			if err := repo.Insert(tt.args.ctx, tt.args.card); (err != nil) != tt.wantErr {
				t.Errorf("SQLCardRepo.Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := sqlMock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestSQLCardRepo_FindByCode(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		ctx  context.Context
		code string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		mock     func()
		wantCard model.Card
		wantErr  bool
	}{
		{
			name: "Error fetching data",
			fields: fields{
				DB: db,
			},
			args: args{
				ctx:  context.Background(),
				code: "AS",
			},
			mock: func() {
				sqlMock.ExpectQuery(
					"SELECT code, suit, `value` FROM card WHERE code=?",
				).WithArgs("AS").WillReturnError(errors.New("connection issue"))
			},
			wantCard: model.Card{
				Code:  "",
				Suit:  "",
				Value: "",
			},
			wantErr: true,
		},
		{
			name: "Success fetching data",
			fields: fields{
				DB: db,
			},
			args: args{
				ctx:  context.Background(),
				code: "AS",
			},
			mock: func() {
				rows := sqlmock.NewRows([]string{"code", "suit", "value"}).AddRow("AS", "SPADES", "ACE")
				sqlMock.ExpectQuery(
					"SELECT code, suit, `value` FROM card WHERE code=?",
				).WithArgs("AS").WillReturnRows(rows)
			},
			wantCard: model.Card{
				Code:  "AS",
				Suit:  "SPADES",
				Value: "ACE",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			repo := &SQLCardRepo{
				DB: tt.fields.DB,
			}
			gotCard, err := repo.FindByCode(tt.args.ctx, tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLCardRepo.FindByCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCard, tt.wantCard) {
				t.Errorf("SQLCardRepo.FindByCode() = %v, want %v", gotCard, tt.wantCard)
			}
			if err := sqlMock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestSQLCardRepo_FindByCodes(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		ctx   context.Context
		codes []string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		mock      func()
		wantCards []model.Card
		wantErr   bool
	}{
		{
			name: "Error fetching data",
			fields: fields{
				DB: db,
			},
			args: args{
				ctx:   context.Background(),
				codes: []string{"AS", "AC"},
			},
			mock: func() {
				sqlMock.ExpectQuery(
					regexp.QuoteMeta("SELECT * FROM card WHERE code IN (?,?)"),
				).WithArgs("AS", "AC").WillReturnError(errors.New("connection issue"))
			},
			wantCards: nil,
			wantErr:   true,
		},
		{
			name: "Success fetching data",
			fields: fields{
				DB: db,
			},
			args: args{
				ctx:   context.Background(),
				codes: []string{"AS", "AC"},
			},
			mock: func() {
				rows := sqlmock.NewRows([]string{"code", "suit", "value"})
				rows.AddRow("AS", "SPADES", "ACE")
				rows.AddRow("AC", "CLUBS", "ACE")
				sqlMock.ExpectQuery(
					regexp.QuoteMeta("SELECT * FROM card WHERE code IN (?,?)"),
				).WithArgs("AS", "AC").WillReturnRows(rows)
			},
			wantCards: []model.Card{
				{
					Code:  "AS",
					Suit:  "SPADES",
					Value: "ACE",
				},
				{
					Code:  "AC",
					Suit:  "CLUBS",
					Value: "ACE",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			repo := &SQLCardRepo{
				DB: tt.fields.DB,
			}
			gotCards, err := repo.FindByCodes(tt.args.ctx, tt.args.codes)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLCardRepo.FindByCodes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCards, tt.wantCards) {
				t.Errorf("SQLCardRepo.FindByCodes() = %v, want %v", gotCards, tt.wantCards)
			}
			if err := sqlMock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestSQLCardRepo_FindAll(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		mock      func()
		wantCards []model.Card
		wantErr   bool
	}{
		{
			name: "Error fetching data",
			fields: fields{
				DB: db,
			},
			args: args{
				ctx: context.Background(),
			},
			mock: func() {
				sqlMock.ExpectQuery(
					regexp.QuoteMeta("SELECT code, suit, `value` FROM card"),
				).WillReturnError(errors.New("connection issue"))
			},
			wantCards: nil,
			wantErr:   true,
		},
		{
			name: "Success fetching data",
			fields: fields{
				DB: db,
			},
			args: args{
				ctx: context.Background(),
			},
			mock: func() {
				rows := sqlmock.NewRows([]string{"code", "suit", "value"})
				rows.AddRow("AS", "SPADES", "ACE")
				rows.AddRow("AC", "CLUBS", "ACE")
				sqlMock.ExpectQuery(
					regexp.QuoteMeta("SELECT code, suit, `value` FROM card"),
				).WillReturnRows(rows)
			},
			wantCards: []model.Card{
				{
					Code:  "AS",
					Suit:  "SPADES",
					Value: "ACE",
				},
				{
					Code:  "AC",
					Suit:  "CLUBS",
					Value: "ACE",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			repo := &SQLCardRepo{
				DB: tt.fields.DB,
			}
			gotCards, err := repo.FindAll(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLCardRepo.FindAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCards, tt.wantCards) {
				t.Errorf("SQLCardRepo.FindAll() = %v, want %v", gotCards, tt.wantCards)
			}
			if err := sqlMock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}
