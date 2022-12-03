package main

import (
	"fmt"
	restful "github.com/emicklei/go-restful/v3"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
)

func main() {
	ws := new(restful.WebService)

	ws.Route(ws.GET("/").To(indexFile))
	ws.Route(ws.POST("/upload").To(upload))
	ws.Route(ws.GET("/download").To(download))

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

func upload(req *restful.Request, resp *restful.Response) {
	r := req.Request
	format := "." + r.FormValue("format")
	file, header, _ := r.FormFile("file")

	// source and target file
	sourceFile := header.Filename
	targetFile := strings.ReplaceAll(sourceFile, path.Ext(sourceFile), format)

	f, _ := os.OpenFile(sourceFile, os.O_WRONLY|os.O_CREATE, 0666)
	io.Copy(f, file)

	fmt.Println("start to convert")
	out, _ := exec.Command("ffmpeg",
		strings.Split(fmt.Sprintf("-i %s -acodec libmp3lame -ab 256k %s -y", sourceFile, targetFile), " ")...).CombinedOutput()
	fmt.Println(string(out))

	fmt.Println("done with convert")
	resp.Write([]byte(fmt.Sprintf("/download?file=%s", targetFile)))
}
