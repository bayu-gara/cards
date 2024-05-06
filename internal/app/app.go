package app

import (
	"fmt"

	//domain
	"github.com/bayu-gara/cards/internal/domain/defaults"
	repo "github.com/bayu-gara/cards/internal/domain/repository"

	//usecase
	deckuc "github.com/bayu-gara/cards/internal/usecase/deck"

	//transport
	"github.com/bayu-gara/cards/internal/transport"

	//lib
	"github.com/bayu-gara/cards/pkg/config"
	db "github.com/bayu-gara/cards/pkg/database"
)

func Run(mode string) error {
	cfg := config.GetConfig()

	//init domain
	connector := db.GetConnector(cfg.Database.Engine)
	if connector == nil {
		return fmt.Errorf("invalid database engine: %s", cfg.Database.Engine)
	}

	dbConnection, err := connector.Connect(cfg.Database)
	if err != nil {
		return err
	}
	defer dbConnection.Close()

	defaults.InitSQLiteTables(dbConnection)
	defaults.InsertDefaultCards(dbConnection)

	cardRepo := repo.NewSQLCardRepository(dbConnection)
	deckRepo := repo.NewSQLDeckRepository(dbConnection)

	//init usecase
	deckuc.InitUsecase(cardRepo, deckRepo)

	server, err := transport.GetHandler(mode)
	if err != nil {
		return err
	}

	return server.Serve()
}
