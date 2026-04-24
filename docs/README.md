# cmn-core Documentation

This directory contains the comprehensive documentation for cmn-core, built with [mdBook](https://rust-lang.github.io/mdBook/).

## Documentation Structure

```
docs/
├── architecture/           # Architecture diagrams
│   ├── high-level-architecture.mmd
│   └── high-level-architecture.svg
├── books/                  # mdBook documentation
│   ├── book/              # Generated HTML (gitignored)
│   ├── src/               # Markdown source files
│   │   ├── SUMMARY.md     # Table of contents
│   │   ├── introduction.md
│   │   ├── architecture/  # Architecture docs
│   │   ├── api/          # API reference
│   │   ├── configuration/# Configuration guides
│   │   ├── development/  # Development guides
│   │   └── appendix/     # Additional resources
│   └── book.toml         # mdBook configuration
├── godoc/                 # Go documentation
└── swagger/               # OpenAPI/Swagger specs
    └── swagger.yaml
```

## Viewing Documentation

### Online (Recommended for Users)

When published, documentation will be available at:
- **GitHub Pages**: https://ryo-arima.github.io/cmn-core/
- **pkg.go.dev**: https://pkg.go.dev/github.com/ryo-arima/cmn-core

### Local Development

#### mdBook Documentation

```bash
# Install mdBook (if not already installed)
cargo install mdbook

# Serve with live reload (recommended for development)
mdbook serve docs/books --open

# Or build static HTML
mdbook build docs/books

# Output will be in docs/books/book/
open docs/books/book/index.html
```

#### GoDoc

```bash
# Install godoc
go install golang.org/x/tools/cmd/godoc@latest

# Start local server
godoc -http=:6060

# View documentation
open http://localhost:6060/pkg/github.com/ryo-arima/cmn-core/
```

#### Swagger UI

```bash
# Start cmn-core server
./.bin/cmn-server

# View Swagger UI
open http://localhost:8080/swagger/index.html
```

## Building Documentation

### Generate All Documentation

```bash
# Generate architecture diagram
./scripts/main.sh docs architecture

# Build mdBook
cd docs/books && mdbook build

# Generate godoc HTML
godoc -html github.com/ryo-arima/cmn-core > docs/godoc/index.html
```

### Using Makefile

```bash
# Generate all documentation
make docs

# Serve documentation locally
make docs-serve

# Clean generated files
make docs-clean
```

## Documentation Content

### Architecture

- **High-Level Architecture**: System overview with component diagram (SVG)
- **Component Details**: Deep dive into each layer
- **Data Flow**: Request/response lifecycle
- **Security Model**: Authentication and authorization

### API Reference

- **Endpoints**: Complete REST API documentation
- **Authentication**: JWT token usage
- **Authorization**: Casbin RBAC policies
- **Request/Response Examples**: Real-world usage

### Configuration

- **Setup Guide**: Getting started with configuration
- **Environment Variables**: Alternative configuration method
- **Casbin Policies**: RBAC policy management
- **Security Best Practices**: Production deployment

### Development

- **Getting Started**: Local development setup
- **Building**: Compilation and build process
- **Testing**: Test execution and coverage
- **Contributing**: Contribution guidelines

### Appendix

- **Swagger/OpenAPI**: Interactive API documentation
- **GoDoc**: Go package documentation
- **Troubleshooting**: Common issues and solutions

## Contributing to Documentation

### Adding New Pages

1. Create markdown file in appropriate directory:
   ```bash
   touch docs/books/src/category/new-page.md
   ```

2. Add entry to `docs/books/src/SUMMARY.md`:
   ```markdown
   - [New Page](./category/new-page.md)
   ```

3. Write content using markdown

4. Build and verify:
   ```bash
   mdbook serve docs/books
   ```

### Documentation Style Guide

- **Headers**: Use sentence case for headers
- **Code Blocks**: Always specify language for syntax highlighting
- **Links**: Use relative links within documentation
- **Examples**: Provide runnable code examples
- **Images**: Store in appropriate subdirectories

### Markdown Features

mdBook supports:
- Standard Markdown
- GitHub Flavored Markdown
- Syntax highlighting
- Table of contents generation
- Search functionality
- Responsive design

## Continuous Documentation

### Pre-commit Checks

```bash
# Validate markdown
markdownlint docs/**/*.md

# Check links
markdown-link-check docs/**/*.md

# Build test
mdbook test docs/books
```

### CI/CD Integration

Documentation is automatically:
1. Built on every commit
2. Validated for broken links
3. Deployed to GitHub Pages (main branch)
4. Published to pkg.go.dev (tagged releases)

## Documentation Tools

### Installed Tools

- **mdBook**: Documentation framework
- **godoc**: Go documentation generator
- **Swagger/OpenAPI**: API documentation
- **Mermaid**: Diagram generation

### Optional Tools

```bash
# Markdown linting
npm install -g markdownlint-cli

# Link checking
npm install -g markdown-link-check

# Diagram editing
npm install -g @mermaid-js/mermaid-cli
```

## Documentation Maintenance

### Regular Tasks

- [ ] Update API documentation when endpoints change
- [ ] Regenerate architecture diagrams for major changes
- [ ] Update godoc comments in code
- [ ] Validate all links quarterly
- [ ] Review and update examples
- [ ] Check for outdated information

### Version-Specific Documentation

Each major version should maintain its own documentation:

```
docs/
├── v1/  # Version 1.x documentation
├── v2/  # Version 2.x documentation
└── latest/ -> v2/  # Symlink to latest
```

## Getting Help

- **Issues**: https://github.com/ryo-arima/cmn-core/issues
- **Discussions**: https://github.com/ryo-arima/cmn-core/discussions
- **Email**: ryo.arima@example.com

## License

Documentation is licensed under [MIT License](../../LICENSE), same as the project code.
