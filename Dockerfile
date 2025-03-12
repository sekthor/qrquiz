FROM golang:1.22-alpine AS build
RUN apk add --no-cache git build-base
WORKDIR /app
COPY ./ ./

ARG BUILD_DATE="$(date -u +'%Y-%m-%dT%H:%M:%SZ')"
ARG GIT_COMMIT
ARG VERSION

RUN go build -o qrquiz \
    -ldflags="-X 'github.com/sekthor/qrquiz/internal/server.date=${BUILD_DATE}' \
        -X 'github.com/sekthor/qrquiz/internal/server.commit=${GIT_COMMIT}' \
        -X 'github.com/sekthor/qrquiz/internal/server.tag=${VERSION}'" \
    ./cmd/main.go

FROM alpine:latest AS final
WORKDIR /app
RUN mkdir /app/data
COPY --from=build /app/qrquiz ./
CMD ["./qrquiz"]
