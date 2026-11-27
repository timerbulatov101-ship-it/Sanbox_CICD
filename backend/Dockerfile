FROM golang:1.23-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/backend \
    ./cmd/api

FROM builder AS test
RUN go test ./...

FROM gcr.io/distroless/static-debian12:nonroot AS runtime

ARG VERSION="0.0.1"

WORKDIR /
COPY --from=builder /out/backend /backend
LABEL org.opencontainers.image.version=${VERSION}

EXPOSE 8081
USER nonroot:nonroot
ENTRYPOINT ["/backend"]
