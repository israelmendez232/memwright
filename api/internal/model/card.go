package model

import "time"

type CardType string

const (
	CardTypeBasic          CardType = "basic"
	CardTypeCloze          CardType = "cloze"
	CardTypeMCQ            CardType = "mcq"
	CardTypeImageOcclusion CardType = "image_occlusion"
	CardTypeAudio          CardType = "audio"
	CardTypeReverse        CardType = "reverse"
)

type Card struct {
	ID        int64     `json:"id" db:"id"`
	DeckID    int64     `json:"deck_id" db:"deck_id"`
	Type      CardType  `json:"type" db:"type"`
	Front     string    `json:"front" db:"front"`
	Back      string    `json:"back" db:"back"`
	Extra     string    `json:"extra,omitempty" db:"extra"`
	Tags      []string  `json:"tags,omitempty" db:"tags"`
	Position  int       `json:"position" db:"position"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
