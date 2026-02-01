package main

import (
	"context"
	"log"
	"os/exec"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	ctx := context.Background()

	// Create a new client, with no features.
	client := mcp.NewClient(&mcp.Implementation{Name: "mcp-client", Version: "v1.0.0"}, nil)

	// Connect to a server over stdin/stdout.
	transport := &mcp.CommandTransport{Command: exec.Command("server/myserver")}
	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	listRes, err := session.ListTools(ctx, nil)
	if err != nil {
		log.Fatalf("ListTool failed: %v", err)
	}
	for _, t := range listRes.Tools {
		log.Printf("Tool: %s - %s", t.Name, t.Description)
	}

	// Call a tool on the server.
	params := &mcp.CallToolParams{
		Name:      "greet",
		Arguments: map[string]any{"name": "Yoni"},
	}
	callRes, err := session.CallTool(ctx, params)
	if err != nil {
		log.Fatalf("CallTool failed: %v", err)
	}
	if callRes.IsError {
		log.Fatal("tool failed")
	}
	for _, c := range callRes.Content {
		log.Print(c.(*mcp.TextContent).Text)
	}
}
