package repository

import (
	"context"
	"database/sql"
	"strings"

	"github.com/bayu-gara/cards/internal/domain/model"
)

const (
	cardCodeSeparator = ","
)

type DeckRepository interface {
	Insert(ctx context.Context, deck model.Deck) error
	UpdateByID(ctx context.Context, deck model.Deck) error
	FindByID(ctx context.Context, id string) (model.Deck, error)
}

func NewSQLDeckRepository(db *sql.DB) DeckRepository {
	return &SQLDeckRepo{
		DB: db,
	}
}

type SQLDeckRepo struct {
	DB *sql.DB
}

func (repo *SQLDeckRepo) Insert(ctx context.Context, deck model.Deck) (err error) {
	query := "INSERT INTO deck(id, card_codes, is_shuffled) VALUES(?,?,?)"
	_, err = repo.DB.ExecContext(ctx, query, deck.ID, strings.Join(deck.CardCodes, cardCodeSeparator), deck.IsShuffled)
	if err != nil {
		return err
	}

	return nil
}

func (repo *SQLDeckRepo) UpdateByID(ctx context.Context, deck model.Deck) (err error) {
	query := "UPDATE deck SET card_codes = ? WHERE id = ?"
	_, err = repo.DB.ExecContext(ctx, query, strings.Join(deck.CardCodes, cardCodeSeparator), deck.ID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *SQLDeckRepo) FindByID(ctx context.Context, id string) (deck model.Deck, err error) {
	query := "SELECT id, card_codes, is_shuffled FROM deck WHERE id=?"
	rows, err := repo.DB.QueryContext(ctx, query, id)
	if err != nil {
		return deck, err
	}
	defer rows.Close()

	rawDeck := struct {
		id           string
		cardCodesStr string
		isShuffled   bool
	}{
		id:           "",
		cardCodesStr: "",
		isShuffled:   false,
	}

	for rows.Next() {
		err := rows.Scan(&rawDeck.id, &rawDeck.cardCodesStr, &rawDeck.isShuffled)
		if err != nil {
			return deck, err
		}
	}

	deck.ID = rawDeck.id
	deck.CardCodes = strings.Split(rawDeck.cardCodesStr, cardCodeSeparator)
	deck.IsShuffled = rawDeck.isShuffled

	return deck, nil
}
