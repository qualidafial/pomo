package store_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/qualidafial/pomo"
	"github.com/qualidafial/pomo/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	p := pomo.Pomo{
		State:    pomo.StateActive,
		Start:    now,
		Duration: 25 * time.Minute,
		Tasks: []pomo.Task{
			{
				Status: pomo.Todo,
				Name:   "Paint the fence",
				Notes:  "Up, down, up down",
			},
			{
				Status: pomo.Doing,
				Name:   "Wax the car",
				Notes:  "Wax on, wax off",
			},
			{
				Status: pomo.Done,
				Name:   "Sand the floor",
				Notes:  "Use little circles\nVery important",
			},
		},
	}

	tempDir := t.TempDir()
	t.Logf("temp dir for test: %s", tempDir)

	storePath := filepath.Join(tempDir, ".pomo")
	t.Logf("store path: %s", storePath)
	s, err := store.New(storePath)
	require.NoError(t, err)

	err = s.Save("test", p)
	require.NoError(t, err)

	loaded, err := s.Read("test")
	assert.Equal(t, p, loaded)
}
