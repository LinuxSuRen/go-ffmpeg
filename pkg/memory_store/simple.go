package memory_store

import (
	"github.com/linuxsuren/go-ffmpeg/pkg/store"
	"math/rand"
	"time"
)

type SimpleStore struct {
	data map[string]*store.Task
}

func (s *SimpleStore) Save(task *store.Task) (id string) {
	id = RandStringRunes(5)
	task.ID = id
	if s.data == nil {
		s.data = make(map[string]*store.Task)
	}
	s.data[id] = task
	return
}
func (s *SimpleStore) Get(id string) (task *store.Task) {
	task = s.data[id]
	return
}
func (s *SimpleStore) GetByName(name string) (task *store.Task) {
	return
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
