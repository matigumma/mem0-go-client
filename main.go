package main

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/matigumma/mem0-go-client/pyexample"
)

func main() {
	apiKey := os.Getenv("MEM0_API_KEY")
	if apiKey == "" {
		// Try getting from .env file
		if file, err := os.Open(".env"); err == nil {
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "MEM0_API_KEY=") {
					apiKey = strings.TrimSpace(strings.TrimPrefix(line, "MEM0_API_KEY="))
					break
				}
			}
		}
		if apiKey == "" {
			log.Fatal("MEM0_API_KEY environment variable is required")
		}
	}

	messages := []pyexample.Message{
		{Role: "user", Content: "Hi, I'm Alex. I'm a vegetarian and I'm allergic to nuts."},
		{Role: "assistant", Content: "Hello Alex! I've noted that you're a vegetarian and have a nut allergy."},
	}

	err := pyexample.StoreMemory(messages, "alex", apiKey)
	if err != nil {
		log.Fatalf("Failed to store memory: %v", err)
	}

	/*
		// Create client with debug logging and custom configuration
		client := mem.NewMem0Client(apiKey,
			mem.WithDebug(true),
			mem.WithBaseURL("https://api.mem0.ai/v1"), // Optional custom base URL
			// mem.WithUserID("custom-user-123"),
			// mem.WithOrganizationID("org-123"),
			// mem.WithProjectID("project-456"),
		)

		ctx := context.Background()

		timestamp := time.Now().UTC().Format(time.RFC3339)

		UserID := "test-user"
		// AgentID := "test-agent"
		// AppID := "test-app"
		// RunID := "test-run"

		// Store a memory with comprehensive options
		storeOpts := &mem.StoreOptions{
			Messages: []mem.Message{
				{
					Role:    "user",
					Content: fmt.Sprintf("Test memory created at %s", timestamp),
				},
			},
			OutputFormat: ptr.String("v1.1"),
			UserID:       UserID,
			// Infer:        ptr.Bool(false),
			// AgentID:      AgentID,
			// AppID:        AppID,
			// RunID:        RunID,
			Metadata: mem.Metadata{
				"source":    "test_client",
				"timestamp": timestamp,
				"test_run":  "memory_retrieval_test",
			},
		}

		memory, err := client.Store(ctx, storeOpts)
		if err != nil {
			log.Fatalf("Failed to store memory: %v", err)
		}
		fmt.Printf("Stored Memory ID: %s\n", memory.ID)

		// Demonstrate GetMemories
		getOpts := &mem.GetMemoriesOptions{
			UserID:   UserID,
			PageSize: 10,
			Page:     1,
			Keywords: "test",
			Metadata: map[string]string{
				"test_run": "memory_retrieval_test",
			},
		}
		memories, err := client.GetMemories(ctx, getOpts)
		if err != nil {
			log.Fatalf("Failed to get memories: %v", err)
		}
		fmt.Printf("Retrieved %d memories\n", len(memories))

		for i, m := range memories {
			fmt.Printf("Memory %d: %+v\n", i+1, m)
		} */
	/*

		// Demonstrate SearchMemories
		searchOpts := &mem.SearchMemoriesOptions{
			Query:  "test memory",
			TopK:   5,
			UserID: "test-user",
			// AgentID: "test-agent",
			Rerank: true,
		}
		searchResults, err := client.SearchMemories(ctx, searchOpts)
		if err != nil {
			log.Fatalf("Failed to search memories: %v", err)
		}
		fmt.Printf("Found %d memories in search\n", len(searchResults))

		// Demonstrate UpdateMemory (if a memory ID is available)
		if memory != nil && memory.ID != "" {
			updateOpts := &mem.UpdateMemoryOptions{
				Text: "Updated test memory content",
				Metadata: map[string]string{
					"updated": "true",
				},
			}
			updatedMemory, err := client.UpdateMemory(ctx, memory.ID, updateOpts)
			if err != nil {
				log.Fatalf("Failed to update memory: %v", err)
			}
			fmt.Printf("Updated Memory ID: %s\n", updatedMemory.ID)
		}
	*/
}
