FROM golang:1.24-alpine AS build
RUN apk add --no-cache git build-base
WORKDIR /app
COPY ./ ./

ARG BUILD_DATE
ARG GIT_COMMIT
ARG VERSION

RUN go build -o qrquiz \
    -ldflags="-X 'github.com/sekthor/qrquiz/internal/config.Date=${BUILD_DATE}' \
        -X 'github.com/sekthor/qrquiz/internal/config.Commit=${GIT_COMMIT}' \
        -X 'github.com/sekthor/qrquiz/internal/config.Version=${VERSION}'" \
    ./cmd/main.go

FROM alpine:latest AS final
WORKDIR /app
RUN mkdir /app/data
COPY --from=build /app/qrquiz ./
CMD ["./qrquiz"]
