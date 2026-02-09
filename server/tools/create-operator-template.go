package tools

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Input struct {
	OperatorName string `json:"operatorName" jsonschema:"the name of the operator to create"`
}

type Output struct {
	DirName string `json:"dirName" jsonschema:"the directory containing the new operator code base"`
}

func runShellCommand(dir, command string, args []string) (string, string, error) {

	var stdout, stderr bytes.Buffer

	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", "", fmt.Errorf("failed to run command %v: %v", cmd, err)
	}

	return stdout.String(), stderr.String(), nil
}

func buildTemplaterIfNeeded() (string, string, error) {

	templaterPath := "gpu-operator-templater/templater"
	if _, err := os.Stat(templaterPath); os.IsNotExist(err) {
		return runShellCommand("gpu-operator-templater", "make", []string{"templater"})
	}

	return "", "", nil
}

func CreateOperatorTemplate(ctx context.Context, req *mcp.CallToolRequest, input Input) (*mcp.CallToolResult, Output, error) {

	stdout, stderr, err := buildTemplaterIfNeeded()
	if err != nil {
		result := &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to build templater: %v\nStdout: %s\nStderr: %s", err, stdout, stderr),
				},
			},
		}
		return result, Output{}, fmt.Errorf("failed to build tempalter: %v", err)
	}

	operatorPath := filepath.Join("generated-operators", input.OperatorName)
	if err := os.MkdirAll(operatorPath, 0755); err != nil {
		return nil, Output{}, fmt.Errorf("couldn't create %s directory for the new operator: %v", operatorPath, err)
	}
	return nil, Output{DirName: operatorPath}, nil
}
