package scanner

import (
	"fmt"
	"sync"
)

var registry = map[string]Scanner{
	"trivy":    &TrivyScanner{},
	"gitleaks": &GitleaksScanner{},
}

func RunAll(tools []string, repoPath string) []ScanResult {
	result := make([]ScanResult, len(tools))

	var wg sync.WaitGroup
	
	for i, tool := range tools {
		wg.Add(1)
		go func(idx int, toolName string) {
			defer wg.Done()
			result[idx] = runTool(toolName, repoPath)
		}(i, tool)
	}

	wg.Wait()
	return result
}

func runTool(tool, repoPath string) ScanResult {
    s, ok := registry[tool]
    if !ok {
        return ScanResult{Tool: tool, Success: false, Error: fmt.Sprintf("tool '%s' not implemented", tool)}
    }
    return s.Run(repoPath) 
}


