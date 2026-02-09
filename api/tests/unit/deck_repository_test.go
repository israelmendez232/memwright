package unit

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"memwright/api/internal/model"
	"memwright/api/internal/repository"
	"memwright/api/internal/srs"
)

type mockDB struct {
	queryRowFunc func(ctx context.Context, query string, args ...interface{}) *sql.Row
	queryFunc    func(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	execFunc     func(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

func (m *mockDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	if m.queryRowFunc != nil {
		return m.queryRowFunc(ctx, query, args...)
	}
	return nil
}

func (m *mockDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if m.queryFunc != nil {
		return m.queryFunc(ctx, query, args...)
	}
	return nil, nil
}

func (m *mockDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if m.execFunc != nil {
		return m.execFunc(ctx, query, args...)
	}
	return nil, nil
}

type mockResult struct {
	lastInsertId int64
	rowsAffected int64
	err          error
}

func (m *mockResult) LastInsertId() (int64, error) {
	return m.lastInsertId, m.err
}

func (m *mockResult) RowsAffected() (int64, error) {
	return m.rowsAffected, m.err
}

func TestNewDeckRepository(t *testing.T) {
	db := &mockDB{}
	repo := repository.NewDeckRepository(db)
	if repo == nil {
		t.Error("NewDeckRepository() returned nil")
	}
}

func TestDeckRepository_Delete_NotFound(t *testing.T) {
	db := &mockDB{
		execFunc: func(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
			return &mockResult{rowsAffected: 0}, nil
		},
	}

	repo := repository.NewDeckRepository(db)
	err := repo.Delete(context.Background(), 999)

	if !errors.Is(err, model.ErrNotFound) {
		t.Errorf("Delete() error = %v, want ErrNotFound", err)
	}
}

func TestDeckRepository_Delete_Success(t *testing.T) {
	db := &mockDB{
		execFunc: func(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
			return &mockResult{rowsAffected: 1}, nil
		},
	}

	repo := repository.NewDeckRepository(db)
	err := repo.Delete(context.Background(), 1)

	if err != nil {
		t.Errorf("Delete() error = %v, want nil", err)
	}
}

func TestDeckRepository_Delete_DBError(t *testing.T) {
	expectedErr := errors.New("database error")
	db := &mockDB{
		execFunc: func(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
			return nil, expectedErr
		},
	}

	repo := repository.NewDeckRepository(db)
	err := repo.Delete(context.Background(), 1)

	if err != expectedErr {
		t.Errorf("Delete() error = %v, want %v", err, expectedErr)
	}
}

func TestDeckRepository_UpdateSRSConfig_NotFound(t *testing.T) {
	db := &mockDB{
		execFunc: func(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
			return &mockResult{rowsAffected: 0}, nil
		},
	}

	repo := repository.NewDeckRepository(db)
	config := &model.SRSConfig{
		SM2: &srs.SM2Config{InitialEaseFactor: 2.5},
	}
	err := repo.UpdateSRSConfig(context.Background(), 999, config)

	if !errors.Is(err, model.ErrNotFound) {
		t.Errorf("UpdateSRSConfig() error = %v, want ErrNotFound", err)
	}
}

func TestDeckRepository_UpdateSRSConfig_Success(t *testing.T) {
	db := &mockDB{
		execFunc: func(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
			return &mockResult{rowsAffected: 1}, nil
		},
	}

	repo := repository.NewDeckRepository(db)
	config := &model.SRSConfig{
		SM2: &srs.SM2Config{InitialEaseFactor: 2.5},
	}
	err := repo.UpdateSRSConfig(context.Background(), 1, config)

	if err != nil {
		t.Errorf("UpdateSRSConfig() error = %v, want nil", err)
	}
}

func TestDeckRepository_UpdateSRSConfig_NilConfig(t *testing.T) {
	db := &mockDB{
		execFunc: func(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
			return &mockResult{rowsAffected: 1}, nil
		},
	}

	repo := repository.NewDeckRepository(db)
	err := repo.UpdateSRSConfig(context.Background(), 1, nil)

	if err != nil {
		t.Errorf("UpdateSRSConfig() error = %v, want nil", err)
	}
}

func TestDeckRepositoryInterface(t *testing.T) {
	var _ repository.DeckRepository = (*deckRepoMock)(nil)
}

type deckRepoMock struct{}

func (m *deckRepoMock) Create(ctx context.Context, deck *model.Deck) error {
	return nil
}
func (m *deckRepoMock) GetByID(ctx context.Context, id int64) (*model.Deck, error) {
	return &model.Deck{
		ID:        id,
		UserID:    1,
		Name:      "Test Deck",
		Algorithm: model.AlgorithmSM2,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
func (m *deckRepoMock) GetByUserID(ctx context.Context, userID int64) ([]*model.Deck, error) {
	return []*model.Deck{}, nil
}
func (m *deckRepoMock) Update(ctx context.Context, deck *model.Deck) error {
	return nil
}
func (m *deckRepoMock) UpdateSRSConfig(ctx context.Context, id int64, config *model.SRSConfig) error {
	return nil
}
func (m *deckRepoMock) Delete(ctx context.Context, id int64) error {
	return nil
}

var _ driver.Result = (*mockResult)(nil)
