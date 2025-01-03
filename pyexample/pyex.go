package pyexample

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

// Message represents a single message in the conversation
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// StoreMemory executes a Python script to store a memory using mem0
func StoreMemory(messages []Message, userID string, apiKey string) error {
	// Create a temporary JSON file with messages
	tempFile, err := os.CreateTemp("", "messages*.json")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Encode messages to JSON
	encoder := json.NewEncoder(tempFile)
	if err := encoder.Encode(messages); err != nil {
		return fmt.Errorf("failed to encode messages: %v", err)
	}
	tempFile.Close()

	// Get API key from environment variable
	if apiKey == "" {
		return fmt.Errorf("MEM0_API_KEY environment variable is not set")
	}

	// Prepare Python script execution
	cmd := exec.Command("python3", "-c", `
import json
import os
from mem0 import MemoryClient

# Read messages from temp file
with open(os.environ['MESSAGES_FILE'], 'r') as f:
	messages = json.load(f)

# Initialize MemoryClient
client = MemoryClient(api_key=os.environ['MEM0_API_KEY'])

# Store messages
result = client.add(messages, user_id=os.environ['USER_ID'])

print(result)
`)

	// Set environment variables
	cmd.Env = append(os.Environ(),
		"MESSAGES_FILE="+tempFile.Name(),
		"MEM0_API_KEY="+apiKey,
		"USER_ID="+userID,
	)

	// Run the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to store memory: %v\nOutput: %s", err, string(output))
	}

	fmt.Println(string(output))

	return nil
}
