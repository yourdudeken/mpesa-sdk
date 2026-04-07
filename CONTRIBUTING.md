# Contributing to Mpesa SDK

Thank you for your interest in contributing to the Mpesa Daraja SDK!

## Code of Conduct

Please be respectful and professional when contributing. We aim to foster an inclusive community.

## How to Contribute

### Reporting Bugs

1. Check if the bug has already been reported
2. Create a detailed issue with:
   - Clear title
   - Steps to reproduce
   - Expected vs actual behavior
   - Environment details

### Suggesting Features

1. Check existing issues and PRs
2. Open an issue with:
   - Feature description
   - Use case
   - Proposed implementation (optional)

### Pull Requests

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Commit with clear messages
7. Push to your fork
8. Submit a PR

## Development Setup

### PHP (Laravel)
```bash
cd packages/php
composer install
composer test
```

### Node.js
```bash
cd packages/node
npm install
npm test
```

### Python
```bash
cd packages/python
pip install -e .
pytest
```

### Java
```bash
cd packages/java
mvn test
```

### C#
```bash
cd packages/dotnet
dotnet test
```

### Go
```bash
cd packages/go
go test ./...
```

## Style Guidelines

- Follow existing code style in each package
- Use meaningful variable/function names
- Add comments for complex logic
- Keep functions focused and small

## Review Process

- PRs require review before merging
- Address feedback promptly
- Ensure CI passes

## Questions?

Open an issue for questions about contributing.