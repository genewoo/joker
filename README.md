# Poker Project

A Go implementation of poker with core game mechanics.

## Features

- Deck management and shuffling
- Card dealing mechanics
- Texas Hold'em specific logic
- Hand evaluation and ranking
- Comprehensive test coverage

## Installation

1. Ensure Go is installed (version 1.20+ recommended)
2. Clone this repository:
   ```bash
   git clone https://github.com/genewoo/poker.git
   cd poker
   ```
3. Install dependencies:
   ```bash
   go mod download
   ```

## Usage

Run the main program:
```bash
go run main.go
```

Run tests:
```bash
go test ./...
```

Build the project:
```bash
make build
```

## Package Structure

```
.
├── internal
│   ├── dealer       # Card dealing logic
│   ├── deck         # Deck management and hand evaluation
│   └── holdem       # Texas Hold'em specific rules
├── bin              # Compiled binaries
├── go.mod           # Go module definition
├── go.sum           # Dependency checksums
└── Makefile         # Build automation
```

## Building

The project includes a Makefile for common tasks:

- `make build`: Build the project
- `make test`: Run all tests
- `make clean`: Remove build artifacts

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a pull request

Please ensure all changes are well-tested and documented.