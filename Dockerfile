
FROM golang:1.22-alpine AS builder


RUN apk add --no-cache git ca-certificates tzdata

RUN adduser -D -g '' appuser

WORKDIR /build


COPY go.mod go.sum ./
RUN go mod download

COPY . .


RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o main ./cmd/


FROM scratch


COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/


COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo


COPY --from=builder /etc/passwd /etc/passwd


COPY --from=builder /build/main /main

COPY --from=builder /build/config.yaml /config.yaml


USER appuser


EXPOSE 8000


HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
    CMD ["/main", "--health"] || exit 1


ENTRYPOINT ["/main"]
