image:
	docker build . -t ghcr.io/linuxsuren/go-ffmpeg

image-run: image
	docker run --rm -p 8080:8080 ghcr.io/linuxsuren/go-ffmpeg
