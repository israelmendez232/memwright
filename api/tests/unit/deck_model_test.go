package unit

import (
	"encoding/json"
	"testing"

	"memwright/api/internal/model"
	"memwright/api/internal/srs"
)

func TestSRSConfig_Scan(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
		check   func(*model.SRSConfig) bool
	}{
		{
			name:    "nil value",
			input:   nil,
			wantErr: false,
			check:   func(c *model.SRSConfig) bool { return c.SM2 == nil },
		},
		{
			name:    "valid JSON",
			input:   []byte(`{"sm2":{"initial_ease_factor":2.5,"min_ease_factor":1.3}}`),
			wantErr: false,
			check: func(c *model.SRSConfig) bool {
				return c.SM2 != nil && c.SM2.InitialEaseFactor == 2.5 && c.SM2.MinEaseFactor == 1.3
			},
		},
		{
			name:    "empty JSON object",
			input:   []byte(`{}`),
			wantErr: false,
			check:   func(c *model.SRSConfig) bool { return c.SM2 == nil },
		},
		{
			name:    "non-byte value",
			input:   "not bytes",
			wantErr: false,
			check:   func(c *model.SRSConfig) bool { return c.SM2 == nil },
		},
		{
			name:    "invalid JSON",
			input:   []byte(`{invalid}`),
			wantErr: true,
			check:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &model.SRSConfig{}
			err := c.Scan(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.check != nil && !tt.check(c) {
				t.Errorf("Scan() result check failed")
			}
		})
	}
}

func TestSRSConfig_Value(t *testing.T) {
	tests := []struct {
		name      string
		config    model.SRSConfig
		wantNil   bool
		wantValid bool
	}{
		{
			name:      "nil SM2 config",
			config:    model.SRSConfig{SM2: nil},
			wantNil:   true,
			wantValid: false,
		},
		{
			name: "valid SM2 config",
			config: model.SRSConfig{
				SM2: &srs.SM2Config{
					InitialEaseFactor: 2.5,
					MinEaseFactor:     1.3,
				},
			},
			wantNil:   false,
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := tt.config.Value()
			if err != nil {
				t.Errorf("Value() error = %v", err)
				return
			}

			if tt.wantNil {
				if value != nil {
					t.Errorf("Value() = %v, want nil", value)
				}
				return
			}

			if !tt.wantValid {
				return
			}

			bytes, ok := value.([]byte)
			if !ok {
				t.Errorf("Value() returned non-bytes type")
				return
			}

			var parsed model.SRSConfig
			if err := json.Unmarshal(bytes, &parsed); err != nil {
				t.Errorf("Value() returned invalid JSON: %v", err)
				return
			}

			if parsed.SM2 == nil || parsed.SM2.InitialEaseFactor != tt.config.SM2.InitialEaseFactor {
				t.Errorf("Value() round-trip failed")
			}
		})
	}
}

func TestDeck_GetSM2Config(t *testing.T) {
	tests := []struct {
		name   string
		deck   model.Deck
		expect srs.SM2Config
	}{
		{
			name:   "nil SRSConfig",
			deck:   model.Deck{SRSConfig: nil},
			expect: srs.SM2Config{},
		},
		{
			name:   "nil SM2 in SRSConfig",
			deck:   model.Deck{SRSConfig: &model.SRSConfig{SM2: nil}},
			expect: srs.SM2Config{},
		},
		{
			name: "valid SM2 config",
			deck: model.Deck{
				SRSConfig: &model.SRSConfig{
					SM2: &srs.SM2Config{
						InitialEaseFactor:  2.5,
						MinEaseFactor:      1.3,
						MaxEaseFactor:      3.0,
						EaseDecrement:      0.2,
						EaseIncrement:      0.15,
						EasyBonusMultipler: 1.3,
						GraduatingInterval: 1,
						MasteredThreshold:  21,
					},
				},
			},
			expect: srs.SM2Config{
				InitialEaseFactor:  2.5,
				MinEaseFactor:      1.3,
				MaxEaseFactor:      3.0,
				EaseDecrement:      0.2,
				EaseIncrement:      0.15,
				EasyBonusMultipler: 1.3,
				GraduatingInterval: 1,
				MasteredThreshold:  21,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.deck.GetSM2Config()
			if result != tt.expect {
				t.Errorf("GetSM2Config() = %+v, want %+v", result, tt.expect)
			}
		})
	}
}

func TestAlgorithmConstants(t *testing.T) {
	if model.AlgorithmSM2 != "sm2" {
		t.Errorf("AlgorithmSM2 = %q, want %q", model.AlgorithmSM2, "sm2")
	}
	if model.AlgorithmFSRS != "fsrs" {
		t.Errorf("AlgorithmFSRS = %q, want %q", model.AlgorithmFSRS, "fsrs")
	}
}
