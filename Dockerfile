FROM golang:1.18 AS builder
WORKDIR /go/src/github.com/restlesswhy/video-merger/

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/app/main.go

FROM alpine:latest AS app
RUN apk --no-cache add ca-certificates
RUN apk add --no-cache ffmpeg
WORKDIR /app
COPY --from=builder /go/src/github.com/restlesswhy/video-merger/app ./
ENTRYPOINT [ "/app/app" ]  