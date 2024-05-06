package deck

import (
	"context"
	"errors"
	"sync"

	//domain
	model "github.com/bayu-gara/cards/internal/domain/model"
	repo "github.com/bayu-gara/cards/internal/domain/repository"

	//lib
	"github.com/bayu-gara/cards/pkg/util"

	//external
	"github.com/google/uuid"
)

var (
	UsecaseObj Usecase

	utilShuffleSliceString = util.ShuffleSliceString
)

type Deck struct {
	ID            string
	IsShuffled    bool
	RemainingCard int
	Cards         []model.Card
}

type Usecase interface {
	CreateDeck(ctx context.Context, shuffle bool, cards []string) (result Deck, err error)
	OpenDeck(ctx context.Context, deckID string) (result Deck, err error)
	DrawCard(ctx context.Context, deckID string, count int) (result Deck, err error)
}

type DeckUsecase struct {
	cardRepo repo.CardRepository
	deckRepo repo.DeckRepository

	drawLock sync.Mutex
}

func InitUsecase(cardRepo repo.CardRepository, deckRepo repo.DeckRepository) {
	UsecaseObj = &DeckUsecase{
		cardRepo: cardRepo,
		deckRepo: deckRepo,
	}
}

func (du *DeckUsecase) CreateDeck(ctx context.Context, shuffle bool, wantedCards []string) (result Deck, err error) {
	var cards []model.Card
	if len(wantedCards) == 0 {
		cards, err = du.cardRepo.FindAll(ctx)
		if err != nil {
			return result, err
		}

		cards = sortCardsByDefaultCodes(cards)

		wantedCards = make([]string, len(cards))
		for i, card := range cards {
			wantedCards[i] = card.Code
		}
	} else {
		wantedCards = filterDuplicateCodes(wantedCards)

		cards, err = du.cardRepo.FindByCodes(ctx, wantedCards)
		if err != nil {
			return result, err
		}

		if !isValidCodes(wantedCards, cards) {
			return result, errors.New("invalid cards")
		}

		cards = sortCardsByCodes(wantedCards, cards)
	}

	if shuffle {
		utilShuffleSliceString(wantedCards)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return result, err
	}

	result.ID = id.String()
	result.Cards = cards
	result.IsShuffled = shuffle
	result.RemainingCard = len(cards)

	err = du.deckRepo.Insert(ctx, model.Deck{
		ID:         result.ID,
		CardCodes:  wantedCards,
		IsShuffled: shuffle,
	})
	if err != nil {
		return Deck{}, err
	}

	return result, nil
}

func (du *DeckUsecase) OpenDeck(ctx context.Context, deckID string) (result Deck, err error) {
	if deckID == "" {
		return result, errors.New("empty deck ID")
	}

	deckModel, err := du.deckRepo.FindByID(ctx, deckID)
	if err != nil {
		return result, err
	}

	cards, err := du.cardRepo.FindByCodes(ctx, deckModel.CardCodes)
	if err != nil {
		return result, err
	}

	sortedCards := sortCardsByCodes(deckModel.CardCodes, cards)

	result.ID = deckModel.ID
	result.Cards = sortedCards
	result.IsShuffled = deckModel.IsShuffled
	result.RemainingCard = len(sortedCards)

	return result, nil
}

func (du *DeckUsecase) DrawCard(ctx context.Context, deckID string, count int) (result Deck, err error) {
	if deckID == "" {
		return result, errors.New("empty deck ID")
	}

	if count < 1 {
		return result, errors.New("number of cards to draw should greater than zero")
	}

	du.drawLock.Lock()
	defer du.drawLock.Unlock()

	deckModel, err := du.deckRepo.FindByID(ctx, deckID)
	if err != nil {
		return result, err
	}

	cards, err := du.cardRepo.FindByCodes(ctx, deckModel.CardCodes)
	if err != nil {
		return result, err
	}

	cards = sortCardsByCodes(deckModel.CardCodes, cards)

	numberOfCards := len(cards)
	if numberOfCards == 0 {
		result.ID = deckModel.ID
		result.IsShuffled = deckModel.IsShuffled
		return result, nil
	}

	var drawnCards []model.Card
	var remainingCodes []string
	if count > numberOfCards {
		drawnCards = cards
		remainingCodes = nil
	} else {
		drawnCards = cards[:count]
		remainingCodes = deckModel.CardCodes[count:]
	}

	err = du.deckRepo.UpdateByID(ctx, model.Deck{
		ID:         deckModel.ID,
		CardCodes:  remainingCodes,
		IsShuffled: deckModel.IsShuffled,
	})
	if err != nil {
		return result, err
	}

	result.ID = deckModel.ID
	result.Cards = drawnCards
	result.IsShuffled = deckModel.IsShuffled
	result.RemainingCard = len(remainingCodes)

	return result, nil
}

func isValidCodes(codes []string, cards []model.Card) bool {
	if len(cards) == 0 {
		return false
	}

	codesMap := make(map[string]bool)
	for _, card := range cards {
		codesMap[card.Code] = true
	}

	for _, code := range codes {
		if _, ok := codesMap[code]; !ok {
			return false
		}
	}

	return true
}

func filterDuplicateCodes(codes []string) (result []string) {
	if len(codes) == 0 {
		return result
	}

	codesMap := make(map[string]bool)
	for _, code := range codes {
		if _, ok := codesMap[code]; !ok {
			codesMap[code] = true
			result = append(result, code)
		}
	}

	return result
}

func sortCardsByDefaultCodes(cards []model.Card) (result []model.Card) {
	codes := []string{
		"AS", "2S", "3S", "4S", "5S", "6S", "7S", "8S", "9S", "10S", "JS", "QS", "KS",
		"AH", "2H", "3H", "4H", "5H", "6H", "7H", "8H", "9H", "10H", "JH", "QH", "KH",
		"AC", "2C", "3C", "4C", "5C", "6C", "7C", "8C", "9C", "10C", "JC", "QC", "KC",
		"AD", "2D", "3D", "4D", "5D", "6D", "7D", "8D", "9D", "10D", "JD", "QD", "KD",
	}

	return sortCardsByCodes(codes, cards)
}

func sortCardsByCodes(codes []string, cards []model.Card) (result []model.Card) {
	if len(cards) == 0 {
		return result
	}

	codesMap := make(map[string]model.Card)
	for _, card := range cards {
		codesMap[card.Code] = card
	}

	for _, code := range codes {
		if card, ok := codesMap[code]; ok {
			result = append(result, card)
		}
	}

	return result
}
