FROM golang:1.19 as builder

WORKDIR /ws
COPY go.mod go.mod
COPY go.sum go.sum
COPY main.go main.go
COPY index.html index.html
COPY sources.list sources.list
RUN go build -o go-ffmpeg .

FROM ubuntu:20.04

COPY --from=builder /ws/sources.list /etc/apt/sources.list

RUN apt update -y && \
    apt install ffmpeg -y

COPY --from=builder /ws/go-ffmpeg /go-ffmpeg
COPY --from=builder /ws/index.html /index.html

ENTRYPOINT ["/go-ffmpeg"]
