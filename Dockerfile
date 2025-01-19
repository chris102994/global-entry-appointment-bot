ARG BRANCH
ARG BUILD_TIMESTAMP
ARG COMMIT_HASH
ARG VERSION

FROM golang:1.23 AS downloader

WORKDIR /work

COPY ./ /work

RUN apt install -y \
        ca-certificates

RUN --mount=type=cache,target=/root/.local/share/golang \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

FROM downloader AS builder

RUN --mount=type=cache,target=/root/.local/share/golang \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -a -o global-entry-appointment-bot -ldflags \
      "-X 'main.Branch=${BRANCH:-N/A}' \
      -X 'main.BuildTimestamp=${BUILD_TIMESTAMP:-N/A}' \
      -X 'main.CommitHash=${COMMIT_HASH:-N/A}' \
      -X 'main.Version=${VERSION:-N/A}'" \
    main.go



FROM scratch
COPY --from=builder /work/global-entry-appointment-bot .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
USER 10000

ENTRYPOINT ["./global-entry-appointment-bot", "run", "--config", "/config/config.yaml"]