package store

import (
	"fmt"
	"github.com/qualidafial/pomo"
	"os"
	"path/filepath"
	"time"
)

func New(path string) (*Store, error) {
	err := os.MkdirAll(path, 0700)
	if err != nil {
		return nil, fmt.Errorf("creating store directory: %w", err)
	}

	return &Store{
		path: path,
	}, nil
}

type Store struct {
	path string
}

func (s *Store) Load(key string) (pomo.Pomodoro, error) {
	filename := filepath.Join(s.path, "pomodoro."+key)
	return pomo.Pomodoro{}, nil
}

func (s *Store) List(fromTo ...time.Time) ([]string, error) {

}

func (s *Store) Store(key string, p pomo.Pomodoro) error {
	return nil
}
