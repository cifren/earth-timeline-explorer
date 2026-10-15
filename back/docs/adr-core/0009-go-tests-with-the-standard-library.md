# 9. Go tests with the standard library

Date: 2026-10-15
Status: Accepted

## Context

ADR common 0009 requires unit and integration tests, fast and isolated from external
services.

## Decision

- The standard `testing` package only: no assertion or mocking framework.
- Tests live in `_test.go` files beside the code; table-driven, with named cases run as
  subtests (`t.Run`).
- **Unit tests** fake external HTTP services with `net/http/httptest` servers (the
  gateway under test receives the fake server's address) and Redis with **miniredis**, an
  in-process Redis server: the standard library has no equivalent, and the code keeps
  using the real Redis client.
- **Integration tests** are files named `*_integration_test.go` with the build
  constraint `//go:build integration`. They start the real HTTP handler of the
  application (through `httptest`) on a real Redis, the one of the development stack,
  with external services faked by `httptest` servers.
- `make test` runs `go test -race ./...`; `make test-integration` adds
  `-tags integration`.

## Consequences

- No dependency to learn for testing besides the Redis fake.
- Integration tests never run by accident in the fast suite.
