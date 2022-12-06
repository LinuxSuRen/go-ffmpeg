package executor

import (
	"github.com/linuxsuren/go-ffmpeg/pkg/store"
	"github.com/stretchr/testify/assert"
	"strings"
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
	// an invalid task
	pool.Submit(&store.Task{})

	// audio and videos
	tasks := []struct {
		format       string
		expectSubCmd string
	}{{
		format:       "mp3",
		expectSubCmd: "ffmpeg",
	}, {
		format:       "mp4",
		expectSubCmd: "ffmpeg",
	}, {
		format:       "wav",
		expectSubCmd: "ffmpeg",
	}, {
		format:       "mkv",
		expectSubCmd: "ffmpeg",
	}, {
		format:       "unknown",
		expectSubCmd: "",
	}}
	for i := range tasks {
		tt := tasks[i]
		pool.On(tt.format, func(task *store.Task) {
			assert.True(t, strings.Contains(task.Command, tt.expectSubCmd),
				"format is %s, should contains %s, but is %s",
				tt.format, tt.expectSubCmd, task.Command)
		})
		pool.Submit(&store.Task{
			DryRun:       true,
			Filename:     tt.format,
			TargetFormat: tt.format,
			TargetWidth:  -1,
			TargetHeight: -1,
			BeginTime:    "00:00:00",
			EndTime:      "00:00:00",
		})
	}
	time.Sleep(2 * time.Second)
}
