# Development: Go toolchain, live reload and linter; the source is bind-mounted.
FROM golang:1.27 AS dev
RUN GOBIN=/usr/local/bin go install github.com/air-verse/air@latest \
 && curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b /usr/local/bin \
 && rm -rf /root/go /root/.cache
CMD ["air"]

# Build: static binary, no cgo.
FROM golang:1.27 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

# Production: distroless static base, non-root user.
FROM gcr.io/distroless/static-debian12:nonroot AS prod
COPY --from=build /out/server /server
EXPOSE 8080
ENTRYPOINT ["/server"]
