package main

import (
	"fmt"
	restful "github.com/emicklei/go-restful/v3"
	"github.com/linuxsuren/go-ffmpeg/pkg/executor"
	"github.com/linuxsuren/go-ffmpeg/pkg/memory_store"
	"github.com/linuxsuren/go-ffmpeg/pkg/store"
	"io"
	"net/http"
	"os"
	"strconv"
)

var pool *executor.Pool
var simpleStore store.Store

func main() {
	ws := new(restful.WebService)

	pool = &executor.Pool{}
	defer pool.Close()
	pool.Run()

	simpleStore = &memory_store.SimpleStore{}

	ws.Route(ws.GET("/").To(indexFile))
	ws.Route(ws.POST("/upload").To(upload))
	ws.Route(ws.GET("/download").To(download))
	ws.Route(ws.GET("/queryTask").To(queryTask))

	restful.Add(ws)
	fmt.Println("start in port 8080")
	http.ListenAndServe(":8080", nil)
}

func indexFile(req *restful.Request, resp *restful.Response) {
	http.ServeFile(
		resp.ResponseWriter,
		req.Request, "index.html")
}

func download(req *restful.Request, resp *restful.Response) {
	f := req.QueryParameter("file")

	w := resp.ResponseWriter
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", f))
	w.Header().Set("Content-Type", "application/octet-stream")

	ff, _ := os.Open(f)

	io.Copy(w, ff)
}

func queryTask(req *restful.Request, resp *restful.Response) {
	id := req.QueryParameter("id")
	task := simpleStore.Get(id)
	resp.WriteAsJson(task)
}

func upload(req *restful.Request, resp *restful.Response) {
	r := req.Request
	format := r.FormValue("format")
	file, header, _ := r.FormFile("file")

	// source and target file
	sourceFile := header.Filename

	f, _ := os.OpenFile(sourceFile, os.O_WRONLY|os.O_CREATE, 0666)
	io.Copy(f, file)

	task := &store.Task{
		Filename:     sourceFile,
		TargetFormat: format,
		BeginTime:    r.FormValue("beginTime"),
		EndTime:      r.FormValue("endTime"),
	}
	task.TargetWidth, _ = strconv.Atoi(r.FormValue("width"))
	task.TargetHeight, _ = strconv.Atoi(r.FormValue("height"))

	id := simpleStore.Save(task)
	pool.Submit(task)

	fmt.Printf("task [%s] has created\n", id)
	resp.Write([]byte(id))
}
