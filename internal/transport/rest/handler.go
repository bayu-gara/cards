package rest

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	//usecase
	deckuc "github.com/bayu-gara/cards/internal/usecase/deck"

	//domain
	"github.com/bayu-gara/cards/internal/domain/model"

	//external
	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"
)

type RESTServer struct {
	Port int
}

type CreateDeckResponse struct {
	ID        string `json:"deck_id"`
	Shuffled  bool   `json:"shuffled"`
	Remaining int    `json:"remaining"`
}

type OpenDeckResponse struct {
	ID        string       `json:"deck_id"`
	Shuffled  bool         `json:"shuffled"`
	Remaining int          `json:"remaining"`
	Card      []model.Card `json:"cards"`
}

type DrawCardResponse struct {
	Card []model.Card `json:"cards"`
}

func (rs RESTServer) Serve() error {
	err := http.ListenAndServe(fmt.Sprintf(":%d", rs.Port), getHandler())
	if errors.Is(err, http.ErrServerClosed) {
		log.Println("server closed")
	}

	return err
}

func getHandler() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/v1/deck/create", createDeck).Methods("POST")
	router.HandleFunc("/v1/deck/open", opendDeck).Methods("GET")
	router.HandleFunc("/v1/deck/draw", drawCard).Methods("GET")

	return router
}

func createDeck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		err error
	)

	queryParams := r.URL.Query()
	shuffleStr := queryParams.Get("shuffle")
	var shuffle bool
	if shuffleStr == "" {
		shuffle = false
	} else {
		shuffle, err = strconv.ParseBool(shuffleStr)
		if err != nil {
			http.Error(w, "Failed read shuffle param", http.StatusBadRequest)
		}
	}

	cardsStr := queryParams.Get("cards")
	var cards []string
	if cardsStr == "" {
		cards = nil
	} else {
		cards = strings.Split(cardsStr, ",")
	}

	result, err := deckuc.UsecaseObj.CreateDeck(ctx, shuffle, cards)
	if err != nil {
		log.Printf("There is an error : %v", err)
		http.Error(w, "Failed to create deck", http.StatusInternalServerError)
		return
	}

	response := CreateDeckResponse{
		ID:        result.ID,
		Shuffled:  result.IsShuffled,
		Remaining: result.RemainingCard,
	}

	writeSuccess(w, response)
}

func opendDeck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		err error
	)

	queryParams := r.URL.Query()
	deckID := queryParams.Get("deck_id")
	if deckID == "" {
		http.Error(w, "Deck id should not be empty", http.StatusBadRequest)
	}

	result, err := deckuc.UsecaseObj.OpenDeck(ctx, deckID)
	if err != nil {
		log.Printf("There is an error : %v", err)
		http.Error(w, "Failed to open deck", http.StatusInternalServerError)
		return
	}

	response := OpenDeckResponse{
		ID:        result.ID,
		Shuffled:  result.IsShuffled,
		Remaining: result.RemainingCard,
		Card:      result.Cards,
	}

	writeSuccess(w, response)
}

func drawCard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		err error
	)

	queryParams := r.URL.Query()
	deckID := queryParams.Get("deck_id")
	if deckID == "" {
		http.Error(w, "Deck id should not be empty", http.StatusBadRequest)
	}

	countStr := queryParams.Get("count")
	var count int
	if countStr == "" {
		http.Error(w, "Count should be greater than zero", http.StatusBadRequest)
		return
	} else {
		count, err = strconv.Atoi(countStr)
		if err != nil {
			http.Error(w, "Count should be a number", http.StatusBadRequest)
			return
		}

		if count < 1 {
			http.Error(w, "Count should be greater than zero", http.StatusBadRequest)
			return
		}
	}

	result, err := deckuc.UsecaseObj.DrawCard(ctx, deckID, count)
	if err != nil {
		log.Printf("There is an error : %v", err)
		http.Error(w, "Failed to draw card", http.StatusInternalServerError)
		return
	}

	response := DrawCardResponse{
		Card: result.Cards,
	}

	writeSuccess(w, response)
}

func writeSuccess(w http.ResponseWriter, data interface{}) {
	jsonRes, err := jsoniter.Marshal(data)
	if err != nil {
		http.Error(w, "Failed to encode data to json", http.StatusInternalServerError)
		return
	}

	// Set content type header
	w.Header().Set("Content-Type", "application/json")

	// Write JSON response
	w.WriteHeader(http.StatusOK)
	w.Write(jsonRes)
}
