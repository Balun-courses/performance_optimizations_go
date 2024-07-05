package entity

// no clean arch boundaries, just example service with single models

type (
	CardNumber string
	CardToken  string
)

type Card struct {
	CardNumber CardNumber `json:"card_number"`
	CardToken  CardToken  `json:"card_token"`
}
