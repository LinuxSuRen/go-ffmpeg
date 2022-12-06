image:
	docker build . -t ghcr.io/linuxsuren/go-ffmpeg
image-with-proxy:
	docker build --build-arg GOPROXY=https://goproxy.io,direct . -t ghcr.io/linuxsuren/go-ffmpeg

image-run: image
	docker run --rm -p 8080:8080 ghcr.io/linuxsuren/go-ffmpeg
image-run-with-proxy: image-with-proxy
	docker run --rm -p 8080:8080 ghcr.io/linuxsuren/go-ffmpeg
pre-commit:
	go test ./...
