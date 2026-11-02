FROM golang:1.27-bookworm AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 go build -trimpath -ldflags="-s -w" -o /out/hmachronos-api .

FROM golang:1.27-bookworm AS development

WORKDIR /src
RUN go install github.com/air-verse/air@v1.61.7

COPY go.mod go.sum ./
RUN go mod download
COPY . .

CMD ["air", "-c", ".air.toml"]

FROM debian:bookworm-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/* \
    && useradd --system --uid 10001 --create-home app \
    && mkdir -p /data \
    && chown app:app /data

COPY --from=build /out/hmachronos-api /usr/local/bin/hmachronos-api

ENV PORT=8080 \
    GIN_MODE=release \
    DB_DRIVER=postgres

USER app
EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/hmachronos-api"]
