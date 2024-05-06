package repository

import (
	"context"
	"database/sql"
	"strings"

	"github.com/bayu-gara/cards/internal/domain/model"
)

type CardRepository interface {
	Insert(ctx context.Context, card model.Card) error
	FindByCode(ctx context.Context, code string) (model.Card, error)
	FindByCodes(ctx context.Context, codes []string) ([]model.Card, error)
	FindAll(ctx context.Context) ([]model.Card, error)
}

func NewSQLCardRepository(db *sql.DB) CardRepository {
	return &SQLCardRepo{
		DB: db,
	}
}

type SQLCardRepo struct {
	DB *sql.DB
}

func (repo *SQLCardRepo) Insert(ctx context.Context, card model.Card) (err error) {
	query := "INSERT INTO card(code, suit, `value`) VALUES(?,?,?)"
	_, err = repo.DB.ExecContext(ctx, query, card.Code, card.Suit, card.Value)
	if err != nil {
		return err
	}

	return nil
}

func (repo *SQLCardRepo) FindByCode(ctx context.Context, code string) (card model.Card, err error) {
	query := "SELECT code, suit, `value` FROM card WHERE code=?"
	rows, err := repo.DB.QueryContext(ctx, query, code)
	if err != nil {
		return card, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&card.Code, &card.Suit, &card.Value)
		if err != nil {
			return card, err
		}
	}
	return card, nil
}

func (repo *SQLCardRepo) FindByCodes(ctx context.Context, codes []string) (cards []model.Card, err error) {
	query := "SELECT * FROM card WHERE code IN (" + strings.Join(make([]string, len(codes)), "?,") + "?)"

	args := make([]interface{}, len(codes))
	for i, code := range codes {
		args[i] = code
	}

	rows, err := repo.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var card model.Card
		err := rows.Scan(&card.Code, &card.Suit, &card.Value)
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}
	return cards, nil
}

func (repo *SQLCardRepo) FindAll(ctx context.Context) (cards []model.Card, err error) {
	query := "SELECT code, suit, `value` FROM card"
	rows, err := repo.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var card model.Card
		err := rows.Scan(&card.Code, &card.Suit, &card.Value)
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}

	return cards, nil
}
