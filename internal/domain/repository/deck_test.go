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

func TestNewSQLDeckRepository(t *testing.T) {
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
		want DeckRepository
	}{
		{
			name: "Success",
			args: args{db: db},
			want: &SQLDeckRepo{DB: db},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSQLDeckRepository(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSQLDeckRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLDeckRepo_Insert(t *testing.T) {
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
		deck model.Deck
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
				deck: model.Deck{
					ID:         "123",
					CardCodes:  []string{"AS", "AC"},
					IsShuffled: false,
				},
			},
			mock: func() {
				sqlMock.ExpectExec(
					"INSERT INTO deck",
				).WithArgs("123", "AS,AC", false).WillReturnError(errors.New("connection issue"))
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
				deck: model.Deck{
					ID:         "123",
					CardCodes:  []string{"AS", "AC"},
					IsShuffled: false,
				},
			},
			mock: func() {
				sqlMock.ExpectExec(
					"INSERT INTO deck",
				).WithArgs("123", "AS,AC", false).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			repo := &SQLDeckRepo{
				DB: tt.fields.DB,
			}
			if err := repo.Insert(tt.args.ctx, tt.args.deck); (err != nil) != tt.wantErr {
				t.Errorf("SQLDeckRepo.Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := sqlMock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestSQLDeckRepo_UpdateByID(t *testing.T) {
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
		deck model.Deck
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func()
		wantErr bool
	}{
		{
			name: "Failed to update",
			fields: fields{
				DB: db,
			},
			args: args{
				ctx: context.Background(),
				deck: model.Deck{
					ID:         "123",
					CardCodes:  []string{"AS", "AC"},
					IsShuffled: false,
				},
			},
			mock: func() {
				sqlMock.ExpectExec(
					"UPDATE deck",
				).WithArgs("AS,AC", "123").WillReturnError(errors.New("connection issue"))
			},
			wantErr: true,
		},
		{
			name: "Success update",
			fields: fields{
				DB: db,
			},
			args: args{
				ctx: context.Background(),
				deck: model.Deck{
					ID:         "123",
					CardCodes:  []string{"AS", "AC"},
					IsShuffled: false,
				},
			},
			mock: func() {
				sqlMock.ExpectExec(
					"UPDATE deck",
				).WithArgs("AS,AC", "123").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			repo := &SQLDeckRepo{
				DB: tt.fields.DB,
			}
			if err := repo.UpdateByID(tt.args.ctx, tt.args.deck); (err != nil) != tt.wantErr {
				t.Errorf("SQLDeckRepo.UpdateByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := sqlMock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestSQLDeckRepo_FindByID(t *testing.T) {
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
		id  string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		mock     func()
		wantDeck model.Deck
		wantErr  bool
	}{
		{
			name: "Failed to fetch data",
			fields: fields{
				DB: db,
			},
			args: args{
				ctx: context.Background(),
				id:  "123",
			},
			mock: func() {
				sqlMock.ExpectQuery(
					regexp.QuoteMeta("SELECT id, card_codes, is_shuffled FROM deck WHERE id=?"),
				).WithArgs("123").WillReturnError(errors.New("connection issue"))
			},
			wantErr: true,
		},
		{
			name: "Success fetch data",
			fields: fields{
				DB: db,
			},
			args: args{
				ctx: context.Background(),
				id:  "123",
			},
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "card_codes", "is_shuffled"}).AddRow("123", "AS,AC", false)
				sqlMock.ExpectQuery(
					regexp.QuoteMeta("SELECT id, card_codes, is_shuffled FROM deck WHERE id=?"),
				).WithArgs("123").WillReturnRows(rows)
			},
			wantDeck: model.Deck{
				ID:         "123",
				CardCodes:  []string{"AS", "AC"},
				IsShuffled: false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			repo := &SQLDeckRepo{
				DB: tt.fields.DB,
			}
			gotDeck, err := repo.FindByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLDeckRepo.FindByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDeck, tt.wantDeck) {
				t.Errorf("SQLDeckRepo.FindByID() = %v, want %v", gotDeck, tt.wantDeck)
			}
			if err := sqlMock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}
