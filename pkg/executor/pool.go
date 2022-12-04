package executor

import (
	"fmt"
	"github.com/linuxsuren/go-ffmpeg/pkg/store"
	"os/exec"
	"path"
	"strings"
	"sync"
)

type Pool struct {
	one   sync.Once
	ch    chan *store.Task
	quite chan struct{}
}

func (p *Pool) Submit(task *store.Task) {
	p.ch <- task
}

func (p *Pool) Close() {
	p.quite <- struct{}{}
}

func (p *Pool) Run() {
	go p.one.Do(func() {
		fmt.Println("task pool started")
		p.ch = make(chan *store.Task, 100)
		p.quite = make(chan struct{})

		for {
			select {
			case task := <-p.ch:
				fmt.Println("start to run task", task.ID)
				if err := p.runTask(task); err != nil {
					task.Retry = task.Retry + 1
					if task.Retry < 3 {
						p.ch <- task
					}
				}
				fmt.Println("task finished", task.ID)
			case <-p.quite:
				fmt.Println("pool has stopped")
				return
			}
		}
	})
}

func (p *Pool) runTask(task *store.Task) (err error) {
	sourceFile := task.Filename
	targetFile := strings.ReplaceAll(sourceFile, path.Ext(sourceFile), task.TargetFormat)

	task.TargetFile = targetFile
	if sourceFile == "" || task.TargetFormat == "" {
		err = fmt.Errorf("source file or target format is missing")
		return
	}

	subCmds := fmt.Sprintf("-i %s -hide_banner -acodec libmp3lame -ab 256k %s -y", sourceFile, targetFile)

	cmd := exec.Command("ffmpeg", strings.Split(subCmds, " ")...)

	task.Command = cmd.String()
	task.Output, err = cmd.CombinedOutput()
	return
}
