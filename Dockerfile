FROM golang:1.22-alpine AS build
RUN apk add --no-cache git build-base
WORKDIR /app
COPY ./ ./
RUN export BUILD_DATE="$(date -u +'%Y-%m-%dT%H:%M:%SZ')" && \
    export GIT_COMMIT="$(git rev-parse HEAD)" && \
    export VERSION="$(git describe --tags --abbrev=0 | tr -d '\n')" && \
    go build -o qrquiz \
    -ldflags="-X 'github.com/sekthor/qrquiz/internal/server.date=${BUILD_DATE}' -X 'github.com/sekthor/qrquiz/internal/server.commit=${GIT_COMMIT}' -X 'github.com/sekthor/qrquiz/internal/server.tag=${VERSION}'" \
    ./cmd/main.go

FROM alpine:latest AS final
WORKDIR /app
RUN mkdir /app/data
COPY --from=build /app/qrquiz ./
CMD ["./qrquiz"]
