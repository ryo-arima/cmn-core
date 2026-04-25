# cmn-core Documentation

Built with [mdBook](https://rust-lang.github.io/mdBook/).

## Structure

```
docs/
├── architecture/        # Architecture diagrams (.mmd)
├── books/               # mdBook source
│   └── src/             # Markdown pages
├── godoc/               # Go package docs
└── swagger/             # OpenAPI spec (swagger.yaml)
```

## Local Preview

```bash
# mdBook
cargo install mdbook
mdbook serve docs/books --open

# GoDoc
go install golang.org/x/tools/cmd/godoc@latest
godoc -http=:6060
open http://localhost:6060/pkg/github.com/ryo-arima/cmn-core/

# Swagger UI (server must be running on :8000)
open http://localhost:8000/swagger/index.html
```

## Build

```bash
make docs        # build all docs
make docs-serve  # serve locally
```

Online: **https://ryo-arima.github.io/cmn-core/**
