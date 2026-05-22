FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG APP_VERSION=dev
ARG GIT_COMMIT=unknown
ARG BUILD_DATE

RUN BUILD_DATE=${BUILD_DATE:-$(date -u +%Y-%m-%dT%H:%M:%SZ)} && \
    CGO_ENABLED=0 go build \
      -ldflags="-X github.com/mikefero/deslacked/internal/version.AppVersion=${APP_VERSION} \
                -X github.com/mikefero/deslacked/internal/version.AppCommit=${GIT_COMMIT} \
                -X github.com/mikefero/deslacked/internal/version.BuildDate=${BUILD_DATE}" \
      -o /app/bin/deslacked ./main.go

FROM cgr.dev/chainguard/static:latest

ARG APP_VERSION=dev
ARG GIT_COMMIT=unknown
ARG BUILD_DATE

LABEL org.opencontainers.image.title="deslacked" \
      org.opencontainers.image.description="deslacked — turning your coworkers into ghosts, one profile pic at a time" \
      org.opencontainers.image.version="${APP_VERSION}" \
      org.opencontainers.image.revision="${GIT_COMMIT}" \
      org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.source="https://github.com/mikefero/slack-deactivated" \
      org.opencontainers.image.url="https://github.com/mikefero/slack-deactivated" \
      org.opencontainers.image.documentation="https://github.com/mikefero/slack-deactivated#readme" \
      org.opencontainers.image.licenses="MIT" \
      org.opencontainers.image.authors="Michael Fero"

COPY --from=builder /app/bin/deslacked /usr/local/bin/deslacked

EXPOSE 8080

ENTRYPOINT ["deslacked", "serve", "--no-browser"]
