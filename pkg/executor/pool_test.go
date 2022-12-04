package executor

import (
	"github.com/linuxsuren/go-ffmpeg/pkg/store"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	pool := &Pool{}
	defer pool.Close()
	pool.Run()
	pool.Run()
	pool.Run()

	time.Sleep(2 * time.Second)
	pool.Submit(&store.Task{})

	time.Sleep(10 * time.Second)
}
