package main

import (
	"flag"
	"fmt"
	"os"
	"testing"
)

// Test BuildPrompt
func TestBuildPrompt(t *testing.T) {
	tests := []struct {
		args []string
		want string
	}{
		{[]string{"Hello", "world!"}, "Hello world!"},
		{[]string{"What", "is", "the", "capital", "of", "France?"}, "What is the capital of France?"},
		{[]string{}, ""},
	}

	for _, tt := range tests {
		got := BuildPrompt(tt.args)
		if got != tt.want {
			t.Errorf("BuildPrompt(%v) = %q, want %q", tt.args, got, tt.want)
		}
	}
}

// Test ValidateAPIKey
func TestValidateAPIKey(t *testing.T) {
	tests := []struct {
		key     string
		wantErr bool
	}{
		{"sk-123456", false},
		{"", true},
	}

	for _, tt := range tests {
		err := ValidateAPIKey(tt.key)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidateAPIKey(%q) error = %v, wantErr %v", tt.key, err, tt.wantErr)
		}
	}
}

func TestParseConfig_Success(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Setenv("OPENAI_API_KEY", "sk-testkey")
	os.Args = []string{"cmd", "--model=gpt-4", "--temperature=0.7", "--max-tokens=50", "What", "is", "AI?"}

	model, temperature, maxTokens, apiKey, prompt, err := ParseConfig()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if model != "gpt-4" {
		t.Errorf("expected model gpt-4, got %s", model)
	}
	if temperature != 0.7 {
		t.Errorf("expected temperature 0.7, got %f", temperature)
	}
	if maxTokens != 50 {
		t.Errorf("expected maxTokens 50, got %d", maxTokens)
	}
	if apiKey != "sk-testkey" {
		t.Errorf("expected apiKey sk-testkey, got %s", apiKey)
	}
	if prompt != "What is AI?" {
		t.Errorf("expected prompt 'What is AI?', got %q", prompt)
	}
}

func TestParseConfig_MissingAPIKey(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Unsetenv("OPENAI_API_KEY")
	os.Args = []string{"cmd", "Hello"}

	_, _, _, _, _, err := ParseConfig()
	if err == nil || err.Error() != "no API key provided. set --api-key or OPENAI_API_KEY in .env or environment" {
		t.Errorf("expected missing API key error, got %v", err)
	}
}

func TestParseConfig_MissingPrompt(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Setenv("OPENAI_API_KEY", "sk-testkey")
	os.Args = []string{"cmd"}

	_, _, _, _, _, err := ParseConfig()
	if err == nil || err.Error() == "" {
		t.Errorf("expected error for missing prompt, got %v", err)
	}
}

type MockChatCompleter struct {
	Response string
	Err      error
}

func (m *MockChatCompleter) CreateChatCompletion(model string, temperature float64, maxTokens int, prompt string) (string, error) {
	return m.Response, m.Err
}

func TestRunChatCompletion(t *testing.T) {
	mock := &MockChatCompleter{Response: "Hello, world!"}
	resp, err := RunChatCompletion(mock, "gpt-3.5-turbo", 0.5, 100, "Say hi")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != "Hello, world!" {
		t.Errorf("expected 'Hello, world!', got %q", resp)
	}
}

func TestRunChatCompletion_Error(t *testing.T) {
	mock := &MockChatCompleter{
		Response: "",
		Err:      fmt.Errorf("API error"),
	}
	_, err := RunChatCompletion(mock, "gpt-3.5-turbo", 0.5, 100, "Say hi")
	if err == nil {
		t.Error("expected error, got nil")
	}
	if err.Error() != "API error" {
		t.Errorf("expected 'API error', got %v", err)
	}
}

func TestRunChatCompletion_DifferentInputs(t *testing.T) {
	mock := &MockChatCompleter{Response: "Different response"}
	resp, err := RunChatCompletion(mock, "gpt-4", 1.0, 200, "Another prompt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != "Different response" {
		t.Errorf("expected 'Different response', got %q", resp)
	}
}

func TestRunChatCompletion_EmptyPrompt(t *testing.T) {
	mock := &MockChatCompleter{Response: "No prompt"}
	resp, err := RunChatCompletion(mock, "gpt-3.5-turbo", 0.5, 100, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != "No prompt" {
		t.Errorf("expected 'No prompt', got %q", resp)
	}
}
