package deck

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"

	//domain
	model "github.com/bayu-gara/cards/internal/domain/model"
	repo "github.com/bayu-gara/cards/internal/domain/repository"

	//lib
	"github.com/bayu-gara/cards/pkg/util"

	//external
	"github.com/golang/mock/gomock"
)

func unmock() {
	utilShuffleSliceString = util.ShuffleSliceString
}

func TestInitUsecase(t *testing.T) {
	type args struct {
		cardRepo repo.CardRepository
		deckRepo repo.DeckRepository
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Success",
			args: args{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitUsecase(tt.args.cardRepo, tt.args.deckRepo)
		})
	}
}

func TestDeckUsecase_CreateDeck(t *testing.T) {
	ctl := gomock.NewController(t)
	cardRepoMock := repo.NewMockCardRepository(ctl)
	deckRepoMock := repo.NewMockDeckRepository(ctl)
	defer ctl.Finish()

	type fields struct {
		cardRepo repo.CardRepository
		deckRepo repo.DeckRepository
		drawLock sync.Mutex
	}
	type args struct {
		ctx         context.Context
		shuffle     bool
		wantedCards []string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		mock       func()
		wantResult Deck
		wantErr    bool
	}{
		{
			name: "Error Find All",
			fields: fields{
				cardRepo: cardRepoMock,
				deckRepo: deckRepoMock,
				drawLock: sync.Mutex{},
			},
			args: args{
				ctx:         context.Background(),
				shuffle:     false,
				wantedCards: nil,
			},
			mock: func() {
				cardRepoMock.EXPECT().FindAll(gomock.Any()).Return(nil, errors.New("connection issue"))
			},
			wantResult: Deck{},
			wantErr:    true,
		},
		{
			name: "Success zero wanted cards",
			fields: fields{
				cardRepo: cardRepoMock,
				deckRepo: deckRepoMock,
				drawLock: sync.Mutex{},
			},
			args: args{
				ctx:         context.Background(),
				shuffle:     true,
				wantedCards: nil,
			},
			mock: func() {
				cardRepoMock.EXPECT().FindAll(gomock.Any()).Return([]model.Card{
					{
						Code:  "AS",
						Suit:  "SPADES",
						Value: "ACE",
					},
					{
						Code:  "2S",
						Suit:  "SPADES",
						Value: "2",
					},
					{
						Code:  "3S",
						Suit:  "SPADES",
						Value: "3",
					},
					{
						Code:  "4S",
						Suit:  "SPADES",
						Value: "4",
					},
					{
						Code:  "5S",
						Suit:  "SPADES",
						Value: "5",
					},
					{
						Code:  "6S",
						Suit:  "SPADES",
						Value: "6",
					},
					{
						Code:  "7S",
						Suit:  "SPADES",
						Value: "7",
					},
					{
						Code:  "8S",
						Suit:  "SPADES",
						Value: "8",
					},
					{
						Code:  "9S",
						Suit:  "SPADES",
						Value: "9",
					},
					{
						Code:  "10S",
						Suit:  "SPADES",
						Value: "10",
					},
					{
						Code:  "JS",
						Suit:  "SPADES",
						Value: "JACK",
					},
					{
						Code:  "QS",
						Suit:  "SPADES",
						Value: "QUEEN",
					},
					{
						Code:  "KS",
						Suit:  "SPADES",
						Value: "KING",
					},
					{
						Code:  "AH",
						Suit:  "HEARTS",
						Value: "ACE",
					},
					{
						Code:  "2H",
						Suit:  "HEARTS",
						Value: "2",
					},
					{
						Code:  "3H",
						Suit:  "HEARTS",
						Value: "3",
					},
					{
						Code:  "4H",
						Suit:  "HEARTS",
						Value: "4",
					},
					{
						Code:  "5H",
						Suit:  "HEARTS",
						Value: "5",
					},
					{
						Code:  "6H",
						Suit:  "HEARTS",
						Value: "6",
					},
					{
						Code:  "7H",
						Suit:  "HEARTS",
						Value: "7",
					},
					{
						Code:  "8H",
						Suit:  "HEARTS",
						Value: "8",
					},
					{
						Code:  "9H",
						Suit:  "HEARTS",
						Value: "9",
					},
					{
						Code:  "10H",
						Suit:  "HEARTS",
						Value: "10",
					},
					{
						Code:  "JH",
						Suit:  "HEARTS",
						Value: "JACK",
					},
					{
						Code:  "QH",
						Suit:  "HEARTS",
						Value: "QUEEN",
					},
					{
						Code:  "KH",
						Suit:  "HEARTS",
						Value: "KING",
					},
					{
						Code:  "AC",
						Suit:  "CLUBS",
						Value: "ACE",
					},
					{
						Code:  "2C",
						Suit:  "CLUBS",
						Value: "2",
					},
					{
						Code:  "3C",
						Suit:  "CLUBS",
						Value: "3",
					},
					{
						Code:  "4C",
						Suit:  "CLUBS",
						Value: "4",
					},
					{
						Code:  "5C",
						Suit:  "CLUBS",
						Value: "5",
					},
					{
						Code:  "6C",
						Suit:  "CLUBS",
						Value: "6",
					},
					{
						Code:  "7C",
						Suit:  "CLUBS",
						Value: "7",
					},
					{
						Code:  "8C",
						Suit:  "CLUBS",
						Value: "8",
					},
					{
						Code:  "9C",
						Suit:  "CLUBS",
						Value: "9",
					},
					{
						Code:  "10C",
						Suit:  "CLUBS",
						Value: "10",
					},
					{
						Code:  "JC",
						Suit:  "CLUBS",
						Value: "JACK",
					},
					{
						Code:  "QC",
						Suit:  "CLUBS",
						Value: "QUEEN",
					},
					{
						Code:  "KC",
						Suit:  "CLUBS",
						Value: "KING",
					},
					{
						Code:  "AD",
						Suit:  "DIAMONDS",
						Value: "ACE",
					},
					{
						Code:  "2D",
						Suit:  "DIAMONDS",
						Value: "2",
					},
					{
						Code:  "3D",
						Suit:  "DIAMONDS",
						Value: "3",
					},
					{
						Code:  "4D",
						Suit:  "DIAMONDS",
						Value: "4",
					},
					{
						Code:  "5D",
						Suit:  "DIAMONDS",
						Value: "5",
					},
					{
						Code:  "6D",
						Suit:  "DIAMONDS",
						Value: "6",
					},
					{
						Code:  "7D",
						Suit:  "DIAMONDS",
						Value: "7",
					},
					{
						Code:  "8D",
						Suit:  "DIAMONDS",
						Value: "8",
					},
					{
						Code:  "9D",
						Suit:  "DIAMONDS",
						Value: "9",
					},
					{
						Code:  "10D",
						Suit:  "DIAMONDS",
						Value: "10",
					},
					{
						Code:  "JD",
						Suit:  "DIAMONDS",
						Value: "JACK",
					},
					{
						Code:  "QD",
						Suit:  "DIAMONDS",
						Value: "QUEEN",
					},
					{
						Code:  "KD",
						Suit:  "DIAMONDS",
						Value: "KING",
					},
				}, nil)

				utilShuffleSliceString = func(slice []string) {}

				deckRepoMock.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantResult: Deck{
				ID:         "123",
				IsShuffled: true,
				Cards: []model.Card{
					{
						Code:  "AS",
						Suit:  "SPADES",
						Value: "ACE",
					},
					{
						Code:  "2S",
						Suit:  "SPADES",
						Value: "2",
					},
					{
						Code:  "3S",
						Suit:  "SPADES",
						Value: "3",
					},
					{
						Code:  "4S",
						Suit:  "SPADES",
						Value: "4",
					},
					{
						Code:  "5S",
						Suit:  "SPADES",
						Value: "5",
					},
					{
						Code:  "6S",
						Suit:  "SPADES",
						Value: "6",
					},
					{
						Code:  "7S",
						Suit:  "SPADES",
						Value: "7",
					},
					{
						Code:  "8S",
						Suit:  "SPADES",
						Value: "8",
					},
					{
						Code:  "9S",
						Suit:  "SPADES",
						Value: "9",
					},
					{
						Code:  "10S",
						Suit:  "SPADES",
						Value: "10",
					},
					{
						Code:  "JS",
						Suit:  "SPADES",
						Value: "JACK",
					},
					{
						Code:  "QS",
						Suit:  "SPADES",
						Value: "QUEEN",
					},
					{
						Code:  "KS",
						Suit:  "SPADES",
						Value: "KING",
					},
					{
						Code:  "AH",
						Suit:  "HEARTS",
						Value: "ACE",
					},
					{
						Code:  "2H",
						Suit:  "HEARTS",
						Value: "2",
					},
					{
						Code:  "3H",
						Suit:  "HEARTS",
						Value: "3",
					},
					{
						Code:  "4H",
						Suit:  "HEARTS",
						Value: "4",
					},
					{
						Code:  "5H",
						Suit:  "HEARTS",
						Value: "5",
					},
					{
						Code:  "6H",
						Suit:  "HEARTS",
						Value: "6",
					},
					{
						Code:  "7H",
						Suit:  "HEARTS",
						Value: "7",
					},
					{
						Code:  "8H",
						Suit:  "HEARTS",
						Value: "8",
					},
					{
						Code:  "9H",
						Suit:  "HEARTS",
						Value: "9",
					},
					{
						Code:  "10H",
						Suit:  "HEARTS",
						Value: "10",
					},
					{
						Code:  "JH",
						Suit:  "HEARTS",
						Value: "JACK",
					},
					{
						Code:  "QH",
						Suit:  "HEARTS",
						Value: "QUEEN",
					},
					{
						Code:  "KH",
						Suit:  "HEARTS",
						Value: "KING",
					},
					{
						Code:  "AC",
						Suit:  "CLUBS",
						Value: "ACE",
					},
					{
						Code:  "2C",
						Suit:  "CLUBS",
						Value: "2",
					},
					{
						Code:  "3C",
						Suit:  "CLUBS",
						Value: "3",
					},
					{
						Code:  "4C",
						Suit:  "CLUBS",
						Value: "4",
					},
					{
						Code:  "5C",
						Suit:  "CLUBS",
						Value: "5",
					},
					{
						Code:  "6C",
						Suit:  "CLUBS",
						Value: "6",
					},
					{
						Code:  "7C",
						Suit:  "CLUBS",
						Value: "7",
					},
					{
						Code:  "8C",
						Suit:  "CLUBS",
						Value: "8",
					},
					{
						Code:  "9C",
						Suit:  "CLUBS",
						Value: "9",
					},
					{
						Code:  "10C",
						Suit:  "CLUBS",
						Value: "10",
					},
					{
						Code:  "JC",
						Suit:  "CLUBS",
						Value: "JACK",
					},
					{
						Code:  "QC",
						Suit:  "CLUBS",
						Value: "QUEEN",
					},
					{
						Code:  "KC",
						Suit:  "CLUBS",
						Value: "KING",
					},
					{
						Code:  "AD",
						Suit:  "DIAMONDS",
						Value: "ACE",
					},
					{
						Code:  "2D",
						Suit:  "DIAMONDS",
						Value: "2",
					},
					{
						Code:  "3D",
						Suit:  "DIAMONDS",
						Value: "3",
					},
					{
						Code:  "4D",
						Suit:  "DIAMONDS",
						Value: "4",
					},
					{
						Code:  "5D",
						Suit:  "DIAMONDS",
						Value: "5",
					},
					{
						Code:  "6D",
						Suit:  "DIAMONDS",
						Value: "6",
					},
					{
						Code:  "7D",
						Suit:  "DIAMONDS",
						Value: "7",
					},
					{
						Code:  "8D",
						Suit:  "DIAMONDS",
						Value: "8",
					},
					{
						Code:  "9D",
						Suit:  "DIAMONDS",
						Value: "9",
					},
					{
						Code:  "10D",
						Suit:  "DIAMONDS",
						Value: "10",
					},
					{
						Code:  "JD",
						Suit:  "DIAMONDS",
						Value: "JACK",
					},
					{
						Code:  "QD",
						Suit:  "DIAMONDS",
						Value: "QUEEN",
					},
					{
						Code:  "KD",
						Suit:  "DIAMONDS",
						Value: "KING",
					},
				},
				RemainingCard: 52,
			},
			wantErr: false,
		},
		{
			name: "Success non empty wanted cards",
			fields: fields{
				cardRepo: cardRepoMock,
				deckRepo: deckRepoMock,
				drawLock: sync.Mutex{},
			},
			args: args{
				ctx:         context.Background(),
				shuffle:     false,
				wantedCards: []string{"AS", "AH", "AC", "AD"},
			},
			mock: func() {
				cardRepoMock.EXPECT().FindByCodes(gomock.Any(), []string{"AS", "AH", "AC", "AD"}).Return([]model.Card{
					{
						Code:  "AS",
						Suit:  "SPADES",
						Value: "ACE",
					},
					{
						Code:  "AH",
						Suit:  "HEARTS",
						Value: "ACE",
					},
					{
						Code:  "AC",
						Suit:  "CLUBS",
						Value: "ACE",
					},
					{
						Code:  "AD",
						Suit:  "DIAMONDS",
						Value: "ACE",
					},
				}, nil)

				deckRepoMock.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantResult: Deck{
				ID:         "123",
				IsShuffled: false,
				Cards: []model.Card{
					{
						Code:  "AS",
						Suit:  "SPADES",
						Value: "ACE",
					},
					{
						Code:  "AH",
						Suit:  "HEARTS",
						Value: "ACE",
					},
					{
						Code:  "AC",
						Suit:  "CLUBS",
						Value: "ACE",
					},
					{
						Code:  "AD",
						Suit:  "DIAMONDS",
						Value: "ACE",
					},
				},
				RemainingCard: 4,
			},
			wantErr: false,
		},
		/*
			{
				name: "Zero wanted cards",
				fields: fields{
					cardRepo: cardRepoMock,
					deckRepo: deckRepoMock,
					drawLock: sync.Mutex{},
				},
				args: args{
					ctx:     context.Background(),
					shuffle: false,
					wantedCards: []string{
						"AS", "2S", "3S", "4S", "5S", "6S", "7S", "8S", "9S", "10S", "JS", "QS", "KS",
						"AH", "2H", "3H", "4H", "5H", "6H", "7H", "8H", "9H", "10H", "JH", "QH", "KH",
						"AC", "2C", "3C", "4C", "5C", "6C", "7C", "8C", "9C", "10C", "JC", "QC", "KC",
						"AD", "2D", "3D", "4D", "5D", "6D", "7D", "8D", "9D", "10D", "JD", "QD", "KD",
					},
				},
			},
		*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			defer unmock()

			du := &DeckUsecase{
				cardRepo: tt.fields.cardRepo,
				deckRepo: tt.fields.deckRepo,
				drawLock: tt.fields.drawLock,
			}
			gotResult, err := du.CreateDeck(tt.args.ctx, tt.args.shuffle, tt.args.wantedCards)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeckUsecase.CreateDeck() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotResult.ID = "123"
			tt.wantResult.ID = "123"
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("DeckUsecase.CreateDeck() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestDeckUsecase_OpenDeck(t *testing.T) {
	ctl := gomock.NewController(t)
	cardRepoMock := repo.NewMockCardRepository(ctl)
	deckRepoMock := repo.NewMockDeckRepository(ctl)
	defer ctl.Finish()

	type fields struct {
		cardRepo repo.CardRepository
		deckRepo repo.DeckRepository
		drawLock sync.Mutex
	}
	type args struct {
		ctx    context.Context
		deckID string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		mock       func()
		wantResult Deck
		wantErr    bool
	}{
		{
			name: "Empty deck id",
			fields: fields{
				cardRepo: cardRepoMock,
				deckRepo: deckRepoMock,
				drawLock: sync.Mutex{},
			},
			args: args{
				ctx:    context.Background(),
				deckID: "",
			},
			mock:       func() {},
			wantResult: Deck{},
			wantErr:    true,
		},
		{
			name: "Error find by id",
			fields: fields{
				cardRepo: cardRepoMock,
				deckRepo: deckRepoMock,
				drawLock: sync.Mutex{},
			},
			args: args{
				ctx:    context.Background(),
				deckID: "-1",
			},
			mock: func() {
				deckRepoMock.EXPECT().FindByID(gomock.Any(), "-1").Return(model.Deck{}, errors.New("connection issue"))
			},
			wantResult: Deck{},
			wantErr:    true,
		},
		{
			name: "Error find by codes",
			fields: fields{
				cardRepo: cardRepoMock,
				deckRepo: deckRepoMock,
				drawLock: sync.Mutex{},
			},
			args: args{
				ctx:    context.Background(),
				deckID: "1",
			},
			mock: func() {
				deckRepoMock.EXPECT().FindByID(gomock.Any(), "1").Return(model.Deck{
					ID:         "1",
					CardCodes:  []string{"AS", "AH", "2H"},
					IsShuffled: false,
				}, nil)

				cardRepoMock.EXPECT().FindByCodes(gomock.Any(), []string{"AS", "AH", "2H"}).Return(nil, errors.New("connection issue"))
			},
			wantResult: Deck{},
			wantErr:    true,
		},
		{
			name: "Success",
			fields: fields{
				cardRepo: cardRepoMock,
				deckRepo: deckRepoMock,
				drawLock: sync.Mutex{},
			},
			args: args{
				ctx:    context.Background(),
				deckID: "1",
			},
			mock: func() {
				deckRepoMock.EXPECT().FindByID(gomock.Any(), "1").Return(model.Deck{
					ID:         "1",
					CardCodes:  []string{"AS", "AH", "2H"},
					IsShuffled: false,
				}, nil)

				cardRepoMock.EXPECT().FindByCodes(gomock.Any(), []string{"AS", "AH", "2H"}).Return([]model.Card{
					{
						Code:  "AH",
						Suit:  "HEARTS",
						Value: "ACE",
					},
					{
						Code:  "2H",
						Suit:  "HEARTS",
						Value: "2",
					},
					{
						Code:  "AS",
						Suit:  "SPADES",
						Value: "ACE",
					},
				}, nil)
			},
			wantResult: Deck{
				ID:            "1",
				IsShuffled:    false,
				RemainingCard: 3,
				Cards: []model.Card{
					{
						Code:  "AS",
						Suit:  "SPADES",
						Value: "ACE",
					},
					{
						Code:  "AH",
						Suit:  "HEARTS",
						Value: "ACE",
					},
					{
						Code:  "2H",
						Suit:  "HEARTS",
						Value: "2",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			defer unmock()

			du := &DeckUsecase{
				cardRepo: tt.fields.cardRepo,
				deckRepo: tt.fields.deckRepo,
				drawLock: tt.fields.drawLock,
			}
			gotResult, err := du.OpenDeck(tt.args.ctx, tt.args.deckID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeckUsecase.OpenDeck() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("DeckUsecase.OpenDeck() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestDeckUsecase_DrawCard(t *testing.T) {
	ctl := gomock.NewController(t)
	cardRepoMock := repo.NewMockCardRepository(ctl)
	deckRepoMock := repo.NewMockDeckRepository(ctl)
	defer ctl.Finish()

	type fields struct {
		cardRepo repo.CardRepository
		deckRepo repo.DeckRepository
		drawLock sync.Mutex
	}
	type args struct {
		ctx    context.Context
		deckID string
		count  int
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		mock       func()
		wantResult Deck
		wantErr    bool
	}{
		{
			name: "Empty deck id",
			fields: fields{
				cardRepo: cardRepoMock,
				deckRepo: deckRepoMock,
				drawLock: sync.Mutex{},
			},
			args: args{
				ctx:    context.Background(),
				deckID: "",
				count:  0,
			},
			mock:       func() {},
			wantResult: Deck{},
			wantErr:    true,
		},
		{
			name: "Count < 1",
			fields: fields{
				cardRepo: cardRepoMock,
				deckRepo: deckRepoMock,
				drawLock: sync.Mutex{},
			},
			args: args{
				ctx:    context.Background(),
				deckID: "1",
				count:  0,
			},
			mock:       func() {},
			wantResult: Deck{},
			wantErr:    true,
		},
		{
			name: "Error find by ID",
			fields: fields{
				cardRepo: cardRepoMock,
				deckRepo: deckRepoMock,
				drawLock: sync.Mutex{},
			},
			args: args{
				ctx:    context.Background(),
				deckID: "1",
				count:  1,
			},
			mock: func() {
				deckRepoMock.EXPECT().FindByID(gomock.Any(), "1").Return(model.Deck{}, errors.New("connection issue"))
			},
			wantResult: Deck{},
			wantErr:    true,
		},
		{
			name: "Error find by codes",
			fields: fields{
				cardRepo: cardRepoMock,
				deckRepo: deckRepoMock,
				drawLock: sync.Mutex{},
			},
			args: args{
				ctx:    context.Background(),
				deckID: "1",
				count:  1,
			},
			mock: func() {
				deckRepoMock.EXPECT().FindByID(gomock.Any(), "1").Return(model.Deck{
					ID:         "1",
					CardCodes:  []string{"AS", "AH", "2H"},
					IsShuffled: false,
				}, nil)

				cardRepoMock.EXPECT().FindByCodes(gomock.Any(), []string{"AS", "AH", "2H"}).Return(nil, errors.New("connection issue"))
			},
			wantResult: Deck{},
			wantErr:    true,
		},
		{
			name: "No cards left",
			fields: fields{
				cardRepo: cardRepoMock,
				deckRepo: deckRepoMock,
				drawLock: sync.Mutex{},
			},
			args: args{
				ctx:    context.Background(),
				deckID: "1",
				count:  1,
			},
			mock: func() {
				deckRepoMock.EXPECT().FindByID(gomock.Any(), "1").Return(model.Deck{
					ID:         "1",
					CardCodes:  []string{},
					IsShuffled: false,
				}, nil)

				cardRepoMock.EXPECT().FindByCodes(gomock.Any(), []string{}).Return([]model.Card{}, nil)
			},
			wantResult: Deck{
				ID:            "1",
				IsShuffled:    false,
				RemainingCard: 0,
				Cards:         nil,
			},
			wantErr: false,
		},
		{
			name: "Success",
			fields: fields{
				cardRepo: cardRepoMock,
				deckRepo: deckRepoMock,
				drawLock: sync.Mutex{},
			},
			args: args{
				ctx:    context.Background(),
				deckID: "1",
				count:  1,
			},
			mock: func() {
				deckRepoMock.EXPECT().FindByID(gomock.Any(), "1").Return(model.Deck{
					ID:         "1",
					CardCodes:  []string{"AS", "AH", "2H"},
					IsShuffled: false,
				}, nil)

				cardRepoMock.EXPECT().FindByCodes(gomock.Any(), []string{"AS", "AH", "2H"}).Return([]model.Card{
					{
						Code:  "AH",
						Suit:  "HEARTS",
						Value: "ACE",
					},
					{
						Code:  "2H",
						Suit:  "HEARTS",
						Value: "2",
					},
					{
						Code:  "AS",
						Suit:  "SPADES",
						Value: "ACE",
					},
				}, nil)

				deckRepoMock.EXPECT().UpdateByID(gomock.Any(), model.Deck{
					ID:         "1",
					CardCodes:  []string{"AH", "2H"},
					IsShuffled: false,
				})
			},
			wantResult: Deck{
				ID:            "1",
				IsShuffled:    false,
				RemainingCard: 2,
				Cards: []model.Card{
					{
						Code:  "AS",
						Suit:  "SPADES",
						Value: "ACE",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			defer unmock()

			du := &DeckUsecase{
				cardRepo: tt.fields.cardRepo,
				deckRepo: tt.fields.deckRepo,
				drawLock: tt.fields.drawLock,
			}
			gotResult, err := du.DrawCard(tt.args.ctx, tt.args.deckID, tt.args.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeckUsecase.DrawCard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("DeckUsecase.DrawCard() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_isValidCodes(t *testing.T) {
	type args struct {
		codes []string
		cards []model.Card
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Empty cards",
			args: args{
				codes: []string{"2S", "7H", "9H"},
				cards: nil,
			},
			want: false,
		},
		{
			name: "Contain invalid codes",
			args: args{
				codes: []string{"2S", "7H", "9H"},
				cards: []model.Card{
					{
						Code:  "2S",
						Suit:  "SPADES",
						Value: "2",
					},
					{
						Code:  "3S",
						Suit:  "SPADES",
						Value: "3",
					},
					{
						Code:  "9H",
						Suit:  "HEARTS",
						Value: "9",
					},
				},
			},
			want: false,
		},
		{
			name: "Valid",
			args: args{
				codes: []string{"2S", "7H", "9H"},
				cards: []model.Card{
					{
						Code:  "2S",
						Suit:  "SPADES",
						Value: "2",
					},
					{
						Code:  "7H",
						Suit:  "HEARTS",
						Value: "7",
					},
					{
						Code:  "9H",
						Suit:  "HEARTS",
						Value: "9",
					},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidCodes(tt.args.codes, tt.args.cards); got != tt.want {
				t.Errorf("isValidCodes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filterDuplicateCodes(t *testing.T) {
	type args struct {
		codes []string
	}
	tests := []struct {
		name       string
		args       args
		wantResult []string
	}{
		{
			name: "Empry codes",
			args: args{
				codes: nil,
			},
			wantResult: nil,
		},
		{
			name: "Duplicate code filtered",
			args: args{
				codes: []string{
					"AS", "2H", "3H", "AS", "KC",
				},
			},
			wantResult: []string{
				"AS", "2H", "3H", "KC",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := filterDuplicateCodes(tt.args.codes); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("filterDuplicateCodes() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_sortCardsByDefaultCodes(t *testing.T) {
	type args struct {
		cards []model.Card
	}
	tests := []struct {
		name       string
		args       args
		wantResult []model.Card
	}{
		{
			name: "Sorted",
			args: args{
				cards: []model.Card{
					{
						Code:  "KC",
						Suit:  "CLUBS",
						Value: "KING",
					},
					{
						Code:  "JS",
						Suit:  "SPADES",
						Value: "JACK",
					},
					{
						Code:  "AH",
						Suit:  "HEARTS",
						Value: "ACE",
					},
				},
			},
			wantResult: []model.Card{
				{
					Code:  "JS",
					Suit:  "SPADES",
					Value: "JACK",
				},
				{
					Code:  "AH",
					Suit:  "HEARTS",
					Value: "ACE",
				},
				{
					Code:  "KC",
					Suit:  "CLUBS",
					Value: "KING",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := sortCardsByDefaultCodes(tt.args.cards); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("sortCardsByDefaultCodes() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_sortCardsByCodes(t *testing.T) {
	type args struct {
		codes []string
		cards []model.Card
	}
	tests := []struct {
		name       string
		args       args
		wantResult []model.Card
	}{
		{
			name: "Empty cards",
			args: args{
				codes: []string{"AH", "KC", "JS"},
				cards: nil,
			},
			wantResult: nil,
		},
		{
			name: "Sorted",
			args: args{
				codes: []string{"AH", "KC", "JS"},
				cards: []model.Card{
					{
						Code:  "KC",
						Suit:  "CLUBS",
						Value: "KING",
					},
					{
						Code:  "JS",
						Suit:  "SPADES",
						Value: "JACK",
					},
					{
						Code:  "AH",
						Suit:  "HEARTS",
						Value: "ACE",
					},
				},
			},
			wantResult: []model.Card{
				{
					Code:  "AH",
					Suit:  "HEARTS",
					Value: "ACE",
				},
				{
					Code:  "KC",
					Suit:  "CLUBS",
					Value: "KING",
				},
				{
					Code:  "JS",
					Suit:  "SPADES",
					Value: "JACK",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := sortCardsByCodes(tt.args.codes, tt.args.cards); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("sortCardsByCodes() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
