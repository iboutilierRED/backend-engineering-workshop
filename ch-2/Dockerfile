# applied standard go image in order to compile application binary
FROM golang:1.25 AS go

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod cmd ./
RUN go mod download

# Disable CGO_ENABLED in order to make compiled binary static because scratch image does not contain
# C libraries for dynamic linking
# GOOS=linux to compile for linux kernal instead of MacOS
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /ch2-wikiapp ./...

# Ran into a problem where a TLS error was occurring because the scratch image lacks a 
# certificate store. I added a small Linux distribution in order to get access to a 
# package manager so I could download the certificates needed to fix the TLS error when
# receiving messages from Wikipedia and then copied those to my container image
FROM alpine:latest AS builder
# adds the certificates
RUN apk add --no-cache ca-certificates

FROM scratch

# Copy binary from go image
COPY --from=go /ch2-wikiapp /ch2-wikiapp

# Copy configs to a directory in the container so app can read it
COPY ./configs ./configs

# Copy the certificates from the builder stage
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 7000

CMD ["/ch2-wikiapp"]