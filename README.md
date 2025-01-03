# Mem0 Go Client

A Golang client for interacting with the Mem0 Memory API.

## Features

- Store memories
- Query memories
- Flexible configuration
- Supports user, organization, and project IDs
- Debug logging

## Installation

```bash
go get github.com/matigumma/mem0-go-client/mem0client
```

### Set environment variables for configuration:

- `MEM0_API_KEY`: Required API key
- `MEM0_USER_ID`: Optional custom user ID
- `MEM0_ORG_ID`: Optional organization ID
- `MEM0_PROJECT_ID`: Optional project ID

### When initializing the client, you can:

- Use default configurations
- Provide custom configurations via option functions
- Override environment variables programmatically


```go
// Using environment variables
client := NewMem0Client(os.Getenv("MEM0_API_KEY"))

// Custom configuration
client := NewMem0Client(apiKey, 
    WithUserID("custom-user"),
    WithOrganizationID("org-123"),
    WithProjectID("project-456")
)
```

### Option functions:
- ```WithBaseURL(url string)```: Customize the base API URL
- ```WithHTTPClient(client *http.Client)```: Use a custom HTTP client
- ```WithDebug(debug bool)```: Enable debug logging
- ```WithUserID(userID string)```: Set a custom user ID
- ```WithOrganizationID(orgID string)```: Set a custom organization ID
- ```WithProjectID(projectID string)```: Set a custom project ID