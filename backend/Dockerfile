# syntax=docker/dockerfile:1.7

FROM golang:1.23-alpine@sha256:383395b794dffa5b53012a212365d40c8e37109a626ca30d6151c8348d380b5f AS base

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download && go mod verify

FROM base AS test-runner

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    set -eu; \
    mkdir -p /test-reports; \
    go install github.com/jstemmer/go-junit-report/v2@v2.1.0; \
    set +e; \
    CGO_ENABLED=0 go test -v -coverprofile=/test-reports/coverage.out ./... > /test-reports/test.log 2>&1; \
    test_status=$?; \
    set -e; \
    "$(go env GOPATH)/bin/go-junit-report" < /test-reports/test.log > /test-reports/junit.xml; \
    if [ -f /test-reports/coverage.out ]; then \
      go tool cover -func=/test-reports/coverage.out > /test-reports/coverage.txt || true; \
    fi; \
    printf '%s\n' "$test_status" > /test-reports/test.exit-code; \
    exit 0

FROM scratch AS test
COPY --from=test-runner /test-reports/ /

FROM base AS builder

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/backend \
    ./cmd/api

FROM gcr.io/distroless/static-debian12:nonroot@sha256:d093aa3e30dbadd3efe1310db061a14da60299baff8450a17fe0ccc514a16639 AS runtime

ARG VERSION="0.0.1"

WORKDIR /
COPY --from=builder /out/backend /backend
LABEL org.opencontainers.image.version=${VERSION}

EXPOSE 8081
USER nonroot:nonroot
ENTRYPOINT ["/backend"]
