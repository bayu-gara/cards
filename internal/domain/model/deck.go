package model

type Deck struct {
	ID         string   `json:"id"`
	CardCodes  []string `json:"card_codes"`
	IsShuffled bool     `json:"is_shuffled"`
}
