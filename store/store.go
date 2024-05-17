package store

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/qualidafial/pomo"
	"gopkg.in/yaml.v3"
)

const (
	currentPomo = "current"
	historyKey  = "history"
)

const (
	timeKeyFormat = "2006-01-02_150405"
)

func New(path string) (*Store, error) {
	storeDir, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("getting store directory absolute path: %w", err)
	}

	err = os.MkdirAll(storeDir, 0700)
	if err != nil {
		return nil, fmt.Errorf("creating pomo store directory: %w", err)
	}

	historyDir := filepath.Join(path, historyKey)
	err = os.MkdirAll(historyDir, 0700)
	if err != nil {
		return nil, fmt.Errorf("creating pomo history directory: %w", err)
	}

	return &Store{
		path: path,
	}, nil
}

type Store struct {
	path string
}

func (s *Store) ClearCurrent() error {
	err := s.Delete(currentPomo)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("clearing current pomo: %w", err)
	}
	return nil
}

func (s *Store) GetCurrent() (pomo.Pomo, error) {
	p, err := s.Read(currentPomo)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return p, fmt.Errorf("reading current pomo: %w", err)
	}
	return p, nil
}

func (s *Store) SaveCurrent(p pomo.Pomo) error {
	err := s.Save(currentPomo, p)
	if err != nil {
		return err
	}
	key := filepath.Join(currentPomo, s.formatTimeKey(time.Now()))
	return s.Save(key, p)
}

func (s *Store) List(fromTo ...time.Time) ([]pomo.Pomo, error) {
	keys, err := s.ListKeys(fromTo...)
	if err != nil {
		return nil, fmt.Errorf("listing keys: %w", err)
	}

	slices.Sort(keys)

	var pomos []pomo.Pomo
	for _, key := range keys {
		pomo, err := s.Read(key)
		if err != nil {
			return pomos, fmt.Errorf("reading key %s: %w", key, err)
		}
		pomos = append(pomos, pomo)
	}
	return pomos, nil
}

func (s *Store) SavePomo(p pomo.Pomo) error {
	key := filepath.Join(historyKey, s.formatTimeKey(p.End))
	return s.Save(key, p)
}

func (s *Store) Read(key string) (pomo.Pomo, error) {
	var p pomo.Pomo

	path := s.pomoFile(key)
	f, err := os.Open(path)
	if err != nil {
		return p, fmt.Errorf("opening file: %w", err)
	}

	err = yaml.NewDecoder(f).Decode(&p)
	return p, errors.Join(err, f.Close())
}

func (s *Store) Save(key string, p pomo.Pomo) (err error) {
	defer func() {
		if err != nil {
			log.Error("saving pomo", "key", key, "err", err)
		} else {
			log.Info("saved pomo", "key", key)
		}
	}()

	path := s.pomoFile(key)

	data, err := yaml.Marshal(p)
	if err != nil {
		return fmt.Errorf("encoding pomodoro to file: %w", err)
	}

	dir := filepath.Dir(path)
	err = os.MkdirAll(dir, 0o700)
	if err != nil {
		return fmt.Errorf("creating parent directory for file: %w", err)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating pomodoro file: %w", err)
	}

	_, err = f.Write(data)
	if err != nil {
		return fmt.Errorf("writing pomodoro to file: %w", err)
	}

	closeErr := f.Close()
	if closeErr != nil {
		closeErr = fmt.Errorf("closing pomodoro file: %w", closeErr)
	}

	return errors.Join(err, closeErr)
}

func (s *Store) Delete(key string) error {
	return os.Remove(s.pomoFile(key))
}

func (s *Store) formatTimeKey(t time.Time) string {
	return t.UTC().Format(timeKeyFormat)
}

func (s *Store) ListKeys(fromTo ...time.Time) ([]string, error) {
	var from, to string
	if len(fromTo) > 0 {
		from = fromTo[0].UTC().Format(timeKeyFormat)
	}
	if len(fromTo) > 1 {
		to = fromTo[1].UTC().Format(timeKeyFormat)
	}

	historyDir := filepath.Join(s.path, historyKey)

	entries, err := os.ReadDir(historyDir)
	if err != nil {
		return nil, fmt.Errorf("reading pomodoro directory: %w", err)
	}

	var keys []string

	for _, entry := range entries {
		name := strings.TrimSuffix(entry.Name(), ".yaml")
		if from != "" && name < from {
			continue
		}
		if to != "" && name >= to {
			continue
		}

		keys = append(keys, filepath.Join(historyKey, name))
	}

	return keys, nil
}

func (s *Store) pomoFile(key string) string {
	return filepath.Join(s.path, key+".yaml")
}
