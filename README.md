# Go OpenAI CLI

A simple command-line interface for interacting with OpenAI's API using Go.

## Features

- Configurable model selection
- Adjustable temperature and max tokens
- Secure API key management
- Comprehensive error handling
- Detailed logging

## Prerequisites

- Go 
- OpenAI API key

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/go-openai.git
cd go-openai
```

2. Install dependencies:
```bash
go mod download
```

3. Set up your OpenAI API key:
```bash
export OPENAI_API_KEY=your_api_key_here
```

## Usage

### Basic Usage

```bash
go run main.go "What is the capital of France?"
```

### Advanced Usage

```bash
go run main.go \
  --model=gpt-4 \
  --temperature=0.7 \
  --max-tokens=100 \
  "Explain quantum computing"
```

### Command Line Options

- `--api-key`: OpenAI API key (default: uses OPENAI_API_KEY environment variable)
- `--model`: Model to use (default: "gpt-3.5-turbo")
- `--temperature`: Temperature setting (default: 0.5)
- `--max-tokens`: Maximum number of tokens to generate (default: 100)

## Configuration

The application can be configured through:
- Command line flags
- Environment variables


## Development

### Running Tests

```bash
go test -v ./...
```

### Test Coverage

```bash
go test -cover ./...
```

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. 