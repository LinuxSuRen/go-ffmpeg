package memory_store

import (
	"github.com/linuxsuren/go-ffmpeg/pkg/store"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSimpleStore(t *testing.T) {
	simpleStore := &SimpleStore{}

	id := simpleStore.Save(&store.Task{Name: "good"})
	assert.NotEmpty(t, id)
	assert.NotNil(t, simpleStore.Get(id))
}
