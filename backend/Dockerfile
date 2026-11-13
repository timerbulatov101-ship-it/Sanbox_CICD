ARG VERSION="0.0.1"

FROM golang:alpine AS builder

ARG VERSION

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o "./build/backend-${VERSION}/app" \
    ./cmd/api/main.go ./cmd/api/router.go

FROM golang:alpine AS runtime

ARG VERSION

WORKDIR /app

COPY --from=builder /app/build/backend-${VERSION}/app ./app

EXPOSE 8081

CMD ["./app"]