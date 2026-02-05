package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Input struct {
	OperatorName string `json:"operatorName" jsonschema:"the name of the operator to create"`
}

type Output struct {
	DirName string `json:"dirName" jsonschema:"the directory containing the new operator code base"`
}

func CreateOperatorTemplate(ctx context.Context, req *mcp.CallToolRequest, input Input) (*mcp.CallToolResult, Output, error) {
	operatorPath := filepath.Join("generated-operators", input.OperatorName)
	if err := os.MkdirAll(operatorPath, 0755); err != nil {
		return nil, Output{}, fmt.Errorf("couldn't create %s directory for the new operator: %v", operatorPath, err)
	}
	return nil, Output{DirName: operatorPath}, nil
}
