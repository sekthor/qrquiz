FROM golang:1.22-alpine AS build
RUN apk add --no-cache git build-base
WORKDIR /app
COPY ./ ./
RUN go build -o qrquiz ./cmd/main.go

FROM alpine:latest AS final
WORKDIR /app
COPY --from=build /app/qrquiz ./
CMD ["./qrquiz"]
