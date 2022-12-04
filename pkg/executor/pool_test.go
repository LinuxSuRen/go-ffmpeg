package executor

import (
	"github.com/linuxsuren/go-ffmpeg/pkg/store"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	pool := &Pool{}
	defer pool.Close()
	for i := 0; i < 3; i++ {
		pool.Run()
	}

	time.Sleep(1 * time.Second)
	pool.Submit(&store.Task{})
	time.Sleep(2 * time.Second)
}
