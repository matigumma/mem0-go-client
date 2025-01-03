package main

import (
	"context"
	"fmt"
	"log"
	"os"

	mem "github.com/matigumma/mem0-go-client/mem0client"
)

func main() {
	apiKey := os.Getenv("MEM0_API_KEY")
	if apiKey == "" {
		log.Fatal("MEM0_API_KEY environment variable is required")
	}

	// Create client with debug logging and custom configuration
	client := mem.NewMem0Client(apiKey,
		mem.WithDebug(true),
		mem.WithBaseURL("https://api.mem0.ai/v1"), // Optional custom base URL
		// mem.WithUserID("custom-user-123"),
		// mem.WithOrganizationID("org-123"),
		// mem.WithProjectID("project-456"),
	)

	ctx := context.Background()

	// Store a memory
	memory, err := client.Store(ctx, "This is a test memory", mem.Metadata{
		"source": "test",
		"tags":   []string{"example", "demo"},
	})
	if err != nil {
		log.Fatalf("Failed to store memory: %v", err)
	}
	fmt.Printf("Stored Memory ID: %s\n", memory.ID)

	// Query memories
	memories, err := client.Query(ctx, "test memory", 5)
	if err != nil {
		log.Fatalf("Failed to query memories: %v", err)
	}
	fmt.Printf("Found %d memories\n", len(memories))
}
