package executor

import (
	"fmt"
	"github.com/linuxsuren/go-ffmpeg/pkg/store"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
)

type Pool struct {
	one      sync.Once
	ch       chan *store.Task
	quite    chan struct{}
	listener map[string]func(task *store.Task)
}

func (p *Pool) Submit(task *store.Task) {
	p.ch <- task
}

func (p *Pool) Close() {
	p.quite <- struct{}{}
}

func (p *Pool) On(filename string, callback func(task *store.Task)) {
	p.listener[filename] = callback
}

func (p *Pool) Run() {
	go p.one.Do(func() {
		fmt.Println("task pool started")
		p.ch = make(chan *store.Task, 100)
		p.quite = make(chan struct{})
		p.listener = map[string]func(task *store.Task){}

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

				if callback := p.listener[task.Filename]; callback != nil {
					go callback(task)
				}
			case <-p.quite:
				fmt.Println("pool has stopped")
				return
			}
		}
	})
}

// read more from https://github.com/Onelinerhub/onelinerhub/tree/main/ffmpeg
func (p *Pool) runTask(task *store.Task) (err error) {
	sourceFile := task.Filename
	targetFile := path.Join("output", strings.ReplaceAll(sourceFile, path.Ext(sourceFile), "."+task.TargetFormat))

	// make sure the directory exists
	if err = os.MkdirAll(path.Dir(targetFile), 0644); err != nil {
		return
	}

	task.TargetFile = targetFile
	if sourceFile == "" || task.TargetFormat == "" {
		err = fmt.Errorf("source file or target format is missing")
		return
	}

	var cmd *exec.Cmd
	switch task.TargetFormat {
	case "mp3", "mp4", "mkv", "wav":
		flags := strings.Split(fmt.Sprintf("-i %s -hide_banner -y", sourceFile), " ")
		if task.BeginTime != "" && task.EndTime != "" {
			flags = append(flags, []string{"-ss", task.BeginTime}...)
			flags = append(flags, []string{"-to", task.EndTime}...)
		}
		if task.TargetFormat == "mp3" {
			flags = append(flags, []string{"-acodec", "libmp3lame"}...)

			if task.AudioBitrate > 0 {
				flags = append(flags, []string{"-b:a", fmt.Sprintf("%dk", task.AudioBitrate)}...)
			}
		}
		if task.TargetFormat == "mp4" || task.TargetFormat == "mkv" {
			//flags = append(flags, []string{"-c:v", "copy"}...)

			//flags = append(flags, []string{"-vf", `drawtext=text='this':x=100:y=50:fontsize=64:box=1`}...)

			//-i sample.png -filter_complex [0:v][1:v]overlay=0:0
			//flags = append(flags, []string{"-vf", `"drawtext=timecode='00:00:00:00':r=30:x=10:y=10:fontsize=24:fontcolor=white"`}...)

			if task.TargetWidth != 0 && task.TargetHeight != 0 {
				flags = append(flags, []string{"-vf", fmt.Sprintf("scale=%d:%d,setsar=1:1", task.TargetWidth, task.TargetHeight)}...)
			}
		}
		flags = append(flags, targetFile)
		cmd = exec.Command("ffmpeg", flags...)

		infoData, _ := exec.Command("ffmpeg", "-i", sourceFile, "-hide_banner").Output()
		task.Info = string(infoData)
	default:
		return
	}

	task.Command = cmd.String()
	if !task.DryRun {
		var data []byte
		if data, err = cmd.CombinedOutput(); err != nil {
			task.ErrOutput = err.Error()
		} else {
			task.Output = string(data)
		}
	}
	return
}
