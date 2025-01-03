package mem0client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
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

type ResponseSingleMemory struct {
	ID        string    `json:"id"`
	Memory    string    `json:"memory"`
	UserID    string    `json:"user_id"`
	AgentID   *string   `json:"agent_id,omitempty"`
	AppID     *string   `json:"app_id,omitempty"`
	RunID     *string   `json:"run_id,omitempty"`
	Hash      string    `json:"hash"`
	Metadata  Metadata  `json:"metadata"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ResponseSearchMemories struct {
	ID        string    `json:"id"`
	Memory    string    `json:"memory"`
	Input     []Message `json:"input"`
	UserID    string    `json:"user_id"`
	Hash      string    `json:"hash"`
	Metadata  *Metadata `json:"metadata,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ResponseGetMemories struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Input         []Message `json:"input"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	TotalMemories int       `json:"total_memories"`
	Owner         string    `json:"owner"`
	Organization  string    `json:"organization"`
	Metadata      Metadata  `json:"metadata"`
	Type          string    `json:"type"`
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
	version        string
}

// Mem0Client is the main client for interacting with memories
type Mem0Client struct {
	config Mem0ClientConfig
}

// NewMem0Client creates a new Mem0 client with default or custom configurations
func NewMem0Client(apiKey string, opts ...func(*Mem0ClientConfig)) *Mem0Client {
	// Generate a default user ID if not provided
	// userID := os.Getenv("MEM0_USER_ID")
	// if userID == "" {
	// 	userID = uuid.New().String()
	// }

	config := Mem0ClientConfig{
		BaseURL: "https://api.mem0.ai/v1",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		APIKey:  apiKey,
		Debug:   true,
		version: "v1.1",
		// UserID:         userID,
		// OrganizationID: os.Getenv("MEM0_ORG_ID"),
		// ProjectID:      os.Getenv("MEM0_PROJECT_ID"),
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
		log.Printf("\033[33m"+format+"\033[0m", v...)
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
	// req.Header.Set("Mem0-User-ID", c.config.UserID)

	// if c.config.OrganizationID != "" {
	// 	req.Header.Set("Mem0-Organization-ID", c.config.OrganizationID)
	// }

	// if c.config.ProjectID != "" {
	// 	req.Header.Set("Mem0-Project-ID", c.config.ProjectID)
	// }
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

// Message represents a single message in the memory
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// StoreOptions represents the full set of options for storing memories
type StoreOptions struct {
	Messages         []Message `json:"messages"`
	UserID           string    `json:"user_id,omitempty"`
	AgentID          string    `json:"agent_id,omitempty"`
	RunID            string    `json:"run_id,omitempty"`
	Metadata         Metadata  `json:"metadata,omitempty"`
	OutputFormat     *string   `json:"output_format,omitempty"`
	AppID            *string   `json:"app_id,omitempty"`
	Includes         *string   `json:"includes,omitempty"`
	Excludes         *string   `json:"excludes,omitempty"`
	Infer            *bool     `json:"infer,omitempty"`
	CustomCategories Metadata  `json:"custom_categories,omitempty"`
	OrganizationName *string   `json:"org_name,omitempty"`
	ProjectName      *string   `json:"project_name,omitempty"`
	OrganizationID   *string   `json:"org_id,omitempty"`
	ProjectID        *string   `json:"project_id,omitempty"`
}

// Store saves memories to the system with full configuration options
func (c *Mem0Client) Store(ctx context.Context, opts *StoreOptions) (*ResponseSingleMemory, error) {
	if opts == nil {
		return nil, fmt.Errorf("store options cannot be nil")
	}

	// Validate that at least one of UserID, AgentID, or RunID is provided
	if opts.UserID == "" && opts.AgentID == "" && opts.RunID == "" {
		return nil, fmt.Errorf("one of the following is required: user_id, agent_id, or run_id")
	}

	// If no messages provided, use the default content as a user message
	if len(opts.Messages) == 0 {
		return nil, fmt.Errorf("at least one message is required")
	}

	// Prepare metadata if not already set
	if opts.Metadata == nil {
		opts.Metadata = make(Metadata)
	}

	// Add UserID, AgentID, RunID to metadata if not already present
	if opts.UserID != "" {
		opts.Metadata["user_id"] = opts.UserID
	}
	if opts.AgentID != "" {
		opts.Metadata["agent_id"] = opts.AgentID
	}
	if opts.RunID != "" {
		opts.Metadata["run_id"] = opts.RunID
	}

	c.debugLog("Storing memory with UserID: %s, AgentID: %s, RunID: %s",
		opts.UserID, opts.AgentID, opts.RunID)

	jsonPayload, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %v", err)
	}

	c.debugLog("Json Payload: %+v", string(jsonPayload))

	req, err := http.NewRequestWithContext(ctx, "POST", c.config.BaseURL+"/memories/", bytes.NewBuffer(jsonPayload))
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

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	c.debugLog("Raw response body: %s", string(body))

	var memory ResponseSingleMemory
	if err := json.NewDecoder(resp.Body).Decode(&memory); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	c.debugLog("Stored memory ID: %s", memory.ID)
	return &memory, nil
}

// GetMemoriesOptions represents optional filters for retrieving memories
type GetMemoriesOptions struct {
	UserID     string            `json:"user_id,omitempty"`
	AgentID    string            `json:"agent_id,omitempty"`
	AppID      string            `json:"app_id,omitempty"`
	RunID      string            `json:"run_id,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
	Categories []string          `json:"categories,omitempty"`
	OrgID      string            `json:"org_id,omitempty"`
	ProjectID  string            `json:"project_id,omitempty"`
	Fields     []string          `json:"fields,omitempty"`
	Keywords   string            `json:"keywords,omitempty"`
	Page       int               `json:"page,omitempty"`
	PageSize   int               `json:"page_size,omitempty"`
}

func (c *Mem0Client) GetMemories(ctx context.Context, opts *GetMemoriesOptions) ([]ResponseGetMemories, error) {
	c.debugLog("Getting memories with options: %+v", opts)

	req, err := http.NewRequestWithContext(ctx, "GET", c.config.BaseURL+"/memories/", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	q := req.URL.Query()
	if opts != nil {
		if opts.UserID != "" {
			q.Add("user_id", opts.UserID)
		}
		if opts.AgentID != "" {
			q.Add("agent_id", opts.AgentID)
		}
		if opts.AppID != "" {
			q.Add("app_id", opts.AppID)
		}
		if opts.RunID != "" {
			q.Add("run_id", opts.RunID)
		}
		if len(opts.Metadata) > 0 {
			metadataJSON, _ := json.Marshal(opts.Metadata)
			q.Add("metadata", string(metadataJSON))
		}
		if len(opts.Categories) > 0 {
			for _, cat := range opts.Categories {
				q.Add("categories", cat)
			}
		}
		if opts.OrgID != "" {
			q.Add("org_id", opts.OrgID)
		}
		if opts.ProjectID != "" {
			q.Add("project_id", opts.ProjectID)
		}
		if len(opts.Fields) > 0 {
			for _, field := range opts.Fields {
				q.Add("fields", field)
			}
		}
		if opts.Keywords != "" {
			q.Add("keywords", opts.Keywords)
		}
		if opts.Page > 0 {
			q.Add("page", fmt.Sprintf("%d", opts.Page))
		}
		if opts.PageSize > 0 {
			q.Add("page_size", fmt.Sprintf("%d", opts.PageSize))
		}
	}
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

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	c.debugLog("Raw response body: %s", string(body))

	// Try to decode as v1.1 API response with 'results' key
	var v11Response struct {
		Results []ResponseGetMemories `json:"results"`
		Count   int                   `json:"count"`
	}
	err = json.Unmarshal(body, &v11Response)
	if err == nil {
		c.debugLog("API Response Count: %d", v11Response.Count)
		if v11Response.Count > 0 {
			c.debugLog("Retrieved %d memories from v1.1 API", len(v11Response.Results))
			return v11Response.Results, nil
		}
		return []ResponseGetMemories{}, nil
	}

	// Try to decode as an array first
	var memories []ResponseGetMemories
	err = json.Unmarshal(body, &memories)
	if err == nil && len(memories) > 0 {
		c.debugLog("Retrieved %d memories", len(memories))
		return memories, nil
	}

	// If array decoding fails, try decoding as a single object
	var singleMemory ResponseGetMemories
	err = json.Unmarshal(body, &singleMemory)
	if err == nil {
		c.debugLog("Retrieved 1 memory")
		return []ResponseGetMemories{singleMemory}, nil
	}

	return nil, fmt.Errorf("failed to decode response: %v. Raw response: %s", err, string(body))
}

// SearchMemoriesOptions represents options for semantic memory search
type SearchMemoriesOptions struct {
	Query                   string            `json:"query"`
	AgentID                 string            `json:"agent_id,omitempty"`
	UserID                  string            `json:"user_id,omitempty"`
	AppID                   string            `json:"app_id,omitempty"`
	RunID                   string            `json:"run_id,omitempty"`
	Metadata                map[string]string `json:"metadata,omitempty"`
	TopK                    int               `json:"top_k,omitempty"`
	Fields                  []string          `json:"fields,omitempty"`
	Rerank                  bool              `json:"rerank,omitempty"`
	OutputFormat            string            `json:"output_format,omitempty"`
	OrgID                   string            `json:"org_id,omitempty"`
	ProjectID               string            `json:"project_id,omitempty"`
	FilterMemories          bool              `json:"filter_memories,omitempty"`
	Categories              []string          `json:"categories,omitempty"`
	OnlyMetadataBasedSearch bool              `json:"only_metadata_based_search,omitempty"`
}

// SearchMemories performs a semantic search on memories
func (c *Mem0Client) SearchMemories(ctx context.Context, opts *SearchMemoriesOptions) ([]ResponseSearchMemories, error) {
	c.debugLog("Searching memories with options: %+v", opts)

	if opts == nil || opts.Query == "" {
		return nil, fmt.Errorf("query is required for searching memories")
	}

	payload, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal search options: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.config.BaseURL+"/memories/search/", bytes.NewBuffer(payload))
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

	var memories []ResponseSearchMemories
	if err := json.NewDecoder(resp.Body).Decode(&memories); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	c.debugLog("Found %d memories in search", len(memories))
	return memories, nil
}

// UpdateMemoryOptions represents the options for updating a memory
type UpdateMemoryOptions struct {
	Text     string            `json:"text"`
	UserID   string            `json:"user_id,omitempty"`
	AgentID  string            `json:"agent_id,omitempty"`
	AppID    string            `json:"app_id,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// UpdateMemory updates a specific memory by its ID
func (c *Mem0Client) UpdateMemory(ctx context.Context, memoryID string, opts *UpdateMemoryOptions) (*Memory, error) {
	c.debugLog("Updating memory %s with options: %+v", memoryID, opts)

	if opts == nil || opts.Text == "" {
		return nil, fmt.Errorf("text is required for updating a memory")
	}

	payload, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal update options: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", fmt.Sprintf("%s/memories/%s/", c.config.BaseURL, memoryID), bytes.NewBuffer(payload))
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

	var updatedMemory Memory
	if err := json.NewDecoder(resp.Body).Decode(&updatedMemory); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	c.debugLog("Updated memory ID: %s", updatedMemory.ID)
	return &updatedMemory, nil
}
