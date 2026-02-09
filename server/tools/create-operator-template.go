package tools

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Input struct {
	OperatorName string `json:"operatorName" jsonschema:"the name of the operator to create"`
}

type Output struct {
	DirName string `json:"dirName" jsonschema:"the directory containing the new operator code base"`
}

func runShellCommand(fromDir, command string, args []string) error {

	var stdout, stderr bytes.Buffer

	cmd := exec.Command(command, args...)
	cmd.Dir = fromDir
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	cmdStr := fmt.Sprintf("%s", strings.Join(append([]string{command}, args...), " "))

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run command '%s': %v", cmdStr, err)
	}

	return nil
}

func CreateOperatorTemplate(ctx context.Context, req *mcp.CallToolRequest, input Input) (*mcp.CallToolResult, Output, error) {

	// Check if templater binary exists, if not build it
	if _, err := os.Stat("gpu-operator-templater/templater"); os.IsNotExist(err) {
		err := runShellCommand("gpu-operator-templater", "make", []string{"templater"})
		if err != nil {
			return nil, Output{}, fmt.Errorf("failed to build tempalter: %v", err)
		}
	}

	// Create a target directory for the new operator
	operatorPath := filepath.Join("generated-operators", input.OperatorName)
	if err := os.MkdirAll(operatorPath, 0755); err != nil {
		return nil, Output{}, fmt.Errorf("couldn't create %s directory for the new operator: %v", operatorPath, err)
	}

	// Run the templater binary to generate the new operator codebase
	err := runShellCommand(operatorPath, "../../gpu-operator-templater/templater",
		[]string{"-f", "../../gpu-operator-templater/examples/config-kmm-only.yaml"})
	if err != nil {
		return nil, Output{}, fmt.Errorf("failed to generate operator's code for %s: %v", input.OperatorName, err)
	}

	return nil, Output{DirName: operatorPath}, nil
}
