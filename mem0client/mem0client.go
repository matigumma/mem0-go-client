package mem0client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Mem0Error represents specific errors from the Mem0 API
type Mem0Error struct {
	Detail   string         `json:"detail"`
	Code     string         `json:"code"`
	Messages []ErrorMessage `json:"messages"`
}

type ErrorMessage struct {
	TokenClass string `json:"token_class"`
	TokenType  string `json:"token_type"`
	Message    string `json:"message"`
}

func (e *Mem0Error) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("Mem0 API Error: %s (Code: %s)", e.Detail, e.Code)
	}

	if len(e.Messages) > 0 {
		return fmt.Sprintf("Mem0 API Error: %s (Token Type: %s, Token Class: %s)",
			e.Messages[0].Message,
			e.Messages[0].TokenType,
			e.Messages[0].TokenClass)
	}

	return "Unknown Mem0 API Error"
}

// Memory represents a single memory entry
type Memory struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Metadata  Metadata  `json:"metadata"`
	Timestamp time.Time `json:"timestamp"`
}

// Metadata allows for additional context about a memory
type Metadata map[string]interface{}

// Mem0ClientConfig allows customization of the Mem0 client
type Mem0ClientConfig struct {
	BaseURL        string
	HTTPClient     *http.Client
	APIKey         string
	Debug          bool
	UserID         string
	OrganizationID string
	ProjectID      string
}

// Mem0Client is the main client for interacting with memories
type Mem0Client struct {
	config Mem0ClientConfig
}

// NewMem0Client creates a new Mem0 client with default or custom configurations
func NewMem0Client(apiKey string, opts ...func(*Mem0ClientConfig)) *Mem0Client {
	// Generate a default user ID if not provided
	userID := os.Getenv("MEM0_USER_ID")
	if userID == "" {
		userID = uuid.New().String()
	}

	config := Mem0ClientConfig{
		BaseURL: "https://api.mem0.ai/v1",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		APIKey:         apiKey,
		Debug:          false,
		UserID:         userID,
		OrganizationID: os.Getenv("MEM0_ORG_ID"),
		ProjectID:      os.Getenv("MEM0_PROJECT_ID"),
	}

	// Apply optional configurations
	for _, opt := range opts {
		opt(&config)
	}

	return &Mem0Client{config: config}
}

// WithBaseURL allows customizing the base API URL
func WithBaseURL(url string) func(*Mem0ClientConfig) {
	return func(c *Mem0ClientConfig) {
		c.BaseURL = url
	}
}

// WithHTTPClient allows providing a custom HTTP client
func WithHTTPClient(client *http.Client) func(*Mem0ClientConfig) {
	return func(c *Mem0ClientConfig) {
		c.HTTPClient = client
	}
}

// WithDebug enables debug logging
func WithDebug(debug bool) func(*Mem0ClientConfig) {
	return func(c *Mem0ClientConfig) {
		c.Debug = debug
	}
}

// WithUserID sets a custom user ID
func WithUserID(userID string) func(*Mem0ClientConfig) {
	return func(c *Mem0ClientConfig) {
		c.UserID = userID
	}
}

// WithOrganizationID sets a custom organization ID
func WithOrganizationID(orgID string) func(*Mem0ClientConfig) {
	return func(c *Mem0ClientConfig) {
		c.OrganizationID = orgID
	}
}

// WithProjectID sets a custom project ID
func WithProjectID(projectID string) func(*Mem0ClientConfig) {
	return func(c *Mem0ClientConfig) {
		c.ProjectID = projectID
	}
}

// debugLog prints debug information if debug mode is enabled
func (c *Mem0Client) debugLog(format string, v ...interface{}) {
	if c.config.Debug {
		log.Printf(format, v...)
	}
}

// prepareRequest adds common headers and parameters
func (c *Mem0Client) prepareRequest(req *http.Request) {
	// Ensure the API key has the 'Token ' prefix
	authHeader := c.config.APIKey
	if !strings.HasPrefix(authHeader, "Token ") {
		authHeader = "Token " + authHeader
	}

	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Mem0-User-ID", c.config.UserID)

	if c.config.OrganizationID != "" {
		req.Header.Set("Mem0-Organization-ID", c.config.OrganizationID)
	}

	if c.config.ProjectID != "" {
		req.Header.Set("Mem0-Project-ID", c.config.ProjectID)
	}
}

// parseErrorResponse attempts to parse a Mem0 API error response
func (c *Mem0Client) parseErrorResponse(body io.Reader) error {
	var apiError Mem0Error
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("failed to read error response: %v", err)
	}

	// Log the raw error response for debugging
	c.debugLog("Error response body: %s", string(bodyBytes))

	// Try to parse the error
	if err := json.Unmarshal(bodyBytes, &apiError); err != nil {
		// If parsing fails, return the raw response as an error
		return fmt.Errorf("API error: %s", string(bodyBytes))
	}

	return &apiError
}

// Store saves a new memory
func (c *Mem0Client) Store(ctx context.Context, content string, metadata Metadata) (*Memory, error) {
	payload := map[string]interface{}{
		"content":  content,
		"metadata": metadata,
	}

	c.debugLog("Storing memory: %+v", payload)

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.config.BaseURL+"/memories", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	c.prepareRequest(req)

	resp, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.parseErrorResponse(resp.Body)
	}

	var memory Memory
	if err := json.NewDecoder(resp.Body).Decode(&memory); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	c.debugLog("Stored memory ID: %s", memory.ID)
	return &memory, nil
}

// Query searches for memories based on a query string
func (c *Mem0Client) Query(ctx context.Context, query string, limit int) ([]Memory, error) {
	c.debugLog("Querying memories: query='%s', limit=%d", query, limit)

	req, err := http.NewRequestWithContext(ctx, "GET", c.config.BaseURL+"/memories/search", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	q := req.URL.Query()
	q.Add("query", query)
	q.Add("limit", fmt.Sprintf("%d", limit))
	req.URL.RawQuery = q.Encode()

	c.prepareRequest(req)

	resp, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.parseErrorResponse(resp.Body)
	}

	var memories []Memory
	if err := json.NewDecoder(resp.Body).Decode(&memories); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	c.debugLog("Found %d memories", len(memories))
	return memories, nil
}
