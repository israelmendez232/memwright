package model

import (
	"encoding/json"
	"time"

	"memwright/api/internal/srs"
)

type Deck struct {
	ID          int64      `json:"id" db:"id"`
	UserID      int64      `json:"user_id" db:"user_id"`
	ParentID    *int64     `json:"parent_id,omitempty" db:"parent_id"`
	Name        string     `json:"name" db:"name"`
	Description string     `json:"description" db:"description"`
	Algorithm   string     `json:"algorithm" db:"algorithm"`
	SRSConfig   *SRSConfig `json:"srs_config" db:"srs_config"`
	Position    int        `json:"position" db:"position"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

type SRSConfig struct {
	SM2 *srs.SM2Config `json:"sm2,omitempty"`
}

func (c *SRSConfig) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, c)
}

func (c SRSConfig) Value() (interface{}, error) {
	if c.SM2 == nil {
		return nil, nil
	}
	return json.Marshal(c)
}

func (d *Deck) GetSM2Config() srs.SM2Config {
	if d.SRSConfig != nil && d.SRSConfig.SM2 != nil {
		return *d.SRSConfig.SM2
	}
	return srs.SM2Config{}
}

const (
	AlgorithmSM2  = "sm2"
	AlgorithmFSRS = "fsrs"
)
