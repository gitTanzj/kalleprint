# syntax=docker/dockerfile:1

FROM golang:1.25.8 AS build-stage
  WORKDIR /app

  COPY go.mod go.sum ./
  RUN go mod download

  COPY . .

  RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/main.go

FROM build-stage AS run-test-stage
  RUN go test -v ./...

FROM debian:bookworm-slim AS prusa-stage
  WORKDIR /tmp
  ARG PRUSA_APPIMAGE_URL
  RUN apt-get update && apt-get install -y --no-install-recommends wget ca-certificates \
      && rm -rf /var/lib/apt/lists/* \
      && wget -q "${PRUSA_APPIMAGE_URL}" -O prusa.AppImage \
      && chmod +x prusa.AppImage \
      && ./prusa.AppImage --appimage-extract \
      && rm prusa.AppImage

FROM debian:bookworm-slim AS build-release-stage
  WORKDIR /

  RUN apt-get update && apt-get install -y --no-install-recommends \
      libglib2.0-0 libgomp1 libgl1 ca-certificates \
      && rm -rf /var/lib/apt/lists/*

  COPY --from=build-stage /api /api
  COPY --from=prusa-stage /tmp/squashfs-root /prusa-slicer

  ENV PRUSA_SLICER_PATH=/prusa-slicer/AppRun
  ENV PRUSA_SLICER_CONFIG=/prusa-slicer/usr/bin/resources/profiles/Creality.ini

  EXPOSE 8080

  ENTRYPOINT ["/api"]
