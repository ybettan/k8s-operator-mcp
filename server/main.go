package main

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/ybettan/k8s-operator-mcp/server/tools"
)

func main() {

	server := mcp.NewServer(&mcp.Implementation{Name: "k8s-operator-manager", Version: "v1.0.0"}, nil)

	mcp.AddTool(server, &mcp.Tool{Name: "create-operator-template", Description: "creates HW operator template"}, tools.CreateOperatorTemplate)

	// Run the server over stdin/stdout, until the client disconnects.
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
