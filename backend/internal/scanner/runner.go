package scanner

import (
	"fmt"
	"sync"
)

// registry maps tool names to their respective scanner implementations.
// This allows for easy extensibility by simply adding new scanners to the map (Registry Pattern).
var registry = map[string]Scanner{
	"trivy":    &TrivyScanner{},
	"gitleaks": &GitleaksScanner{},
}

// RunAll executes multiple security scanning tools concurrently on a given repository.
// It uses a WaitGroup to ensure all scans are completed before returning the aggregated results.
func RunAll(tools []string, repoPath string) []ScanResult {
	// Initialize slice to hold results for each tool
	result := make([]ScanResult, len(tools))

	// WaitGroup to synchronize concurrent goroutines
	var wg sync.WaitGroup
	
	for i, tool := range tools {
		// Increment WaitGroup counter for each tool
		wg.Add(1)
		
		// Launch each scanner in its own goroutine for parallel execution
		go func(idx int, toolName string) {
			// Decrement counter when the goroutine completes
			defer wg.Done()
			
			// Execute the tool and store the result at the correct index
			result[idx] = runTool(toolName, repoPath)
		}(i, tool) // Pass i and tool as arguments to avoid loop variable capture issues
	}

	// Wait for all goroutines to finish
	wg.Wait()
	return result
}

// runTool is a helper function that fetches a scanner from the registry and executes it.
func runTool(tool, repoPath string) ScanResult {
    s, ok := registry[tool]
    if !ok {
        // Return a failure result if the tool is not found in the registry
        return ScanResult{Tool: tool, Success: false, Error: fmt.Sprintf("tool '%s' not implemented", tool)}
    }
    // Execute the scanner implementation (Polymorphism via Interface)
    return s.Run(repoPath) 
}


