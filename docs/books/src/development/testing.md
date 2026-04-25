# Testing

## Unit Tests

```bash
make test

# With coverage
go test -v -cover ./...
```

## E2E Tests

Requires dev environment running (`make dev-up`).

```bash
go test -v -timeout 15m ./test/e2e/testcase/
```

## CI/CD

Unit and E2E tests run automatically on GitHub Actions for every pull request and push to `dev` / `main`.
