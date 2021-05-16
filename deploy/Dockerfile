FROM golang:1.16-alpine AS builder
WORKDIR /src
COPY . .
WORKDIR /src/cmd
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o main .
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /src/cmd/main /usr/bin/
WORKDIR /src/cmd
CMD ["main"]
