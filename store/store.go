package store

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/qualidafial/pomo"
	"gopkg.in/yaml.v3"
)

const (
	keyFormat = "2006-01-02_150405"
)

func New(path string) (*Store, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("getting store directory absolute path: %w", err)
	}

	err = os.MkdirAll(path, 0700)
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

func (s *Store) Read(key string) (pomo.Pomodoro, error) {
	var p pomo.Pomodoro

	path := s.pomoFile(key)
	f, err := os.Open(path)
	if err != nil {
		return p, fmt.Errorf("opening file: %w", err)
	}

	err = yaml.NewDecoder(f).Decode(&p)
	return p, errors.Join(err, f.Close())
}

func (s *Store) Save(p pomo.Pomodoro) (string, error) {
	key := s.Key(p)
	path := s.pomoFile(key)

	data, err := yaml.Marshal(p)
	if err != nil {
		return "", fmt.Errorf("encoding pomodoro to file: %w", err)
	}

	f, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("creating pomodoro file: %w", err)
	}

	_, err = f.Write(data)
	if err != nil {
		return "", fmt.Errorf("writing pomodoro to file: %w", err)
	}

	closeErr := f.Close()
	if closeErr != nil {
		closeErr = fmt.Errorf("closing pomodoro file: %w", closeErr)
	}

	return key, errors.Join(err, closeErr)
}

func (s *Store) Key(p pomo.Pomodoro) string {
	return time.Time(p.Start).UTC().Format(keyFormat)
}

func (s *Store) List(fromTo ...time.Time) ([]string, error) {
	var from, to string
	if len(fromTo) > 0 {
		from = fromTo[0].Format(time.DateOnly)
	}
	if len(fromTo) > 1 {
		to = fromTo[1].Format(time.DateOnly)
	}

	entries, err := os.ReadDir(s.path)
	if err != nil {
		return nil, fmt.Errorf("reading pomodoro directory: %w", err)
	}

	var keys []string

	for _, entry := range entries {
		name := entry.Name()
		if from != "" && name < from {
			continue
		}
		if to != "" && name >= to {
			continue
		}
		keys = append(keys, name)
	}

	return keys, nil
}

func (s *Store) pomoFile(key string) string {
	return filepath.Join(s.path, key+".yaml")
}

type pomodoro struct {
	State    string `yaml:"state"`
	Start    string `yaml:"start,omitempty"`
	End      string `yaml:"end,omitempty"`
	Duration string `yaml:"duration"`
	Tasks    map[string][]task
}

type task struct {
	Name  string `yaml:"name"`
	Notes string `yaml:"notes,omitempty"`
}
