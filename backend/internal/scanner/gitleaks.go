package scanner

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

type GitleaksScanner struct{}

func (g *GitleaksScanner) Name() string { return "gitleaks" }
func (g *GitleaksScanner) Run(repoPath string) ScanResult {
	return runGitleaks(repoPath)
}

type gitleaksFinding struct {
	Description string `json:"Description"`
	StartLine   int    `json:"StartLine"`
	File        string `json:"File"`
	RuleID      string `json:"RuleID"`
	Secret      string `json:"Secret"`
	Match       string `json:"Match"`
}

func runGitleaks(repoPath string) ScanResult {
	gitleaksPath, err := exec.LookPath("gitleaks")
	if err != nil {
		return ScanResult{
			Tool:    "gitleaks",
			Success: false,
			Error:   "gitleaks not found in PATH: install gitleaks",
		}
	}

	// Use a temp file per scan to avoid concurrent write conflicts
	tmpFile, err := os.CreateTemp("", "gitleaks-report-*.json")
	if err != nil {
		return ScanResult{
			Tool:    "gitleaks",
			Success: false,
			Error:   fmt.Sprintf("failed to create temp report file: %v", err),
		}
	}
	reportPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(reportPath)

	cmd := exec.Command(
		gitleaksPath, "detect",
		"--source", repoPath,
		"--report-format", "json",
		"--report-path", reportPath,
		"--exit-code", "0", // never return non-zero on findings; handle manually
		"--no-git",         // avoid git-ownership issues
	)

	if output, err := cmd.CombinedOutput(); err != nil {
		return ScanResult{
			Tool:    "gitleaks",
			Success: false,
			Error:   fmt.Sprintf("gitleaks error: %v, output: %s", err, string(output)),
		}
	}

	data, err := os.ReadFile(reportPath)
	if err != nil || len(data) == 0 {
		// No findings — return empty success
		return ScanResult{Tool: "gitleaks", Success: true, Findings: []Vulnerability{}}
	}

	var raw []gitleaksFinding
	if err := json.Unmarshal(data, &raw); err != nil {
		return ScanResult{
			Tool:    "gitleaks",
			Success: false,
			Error:   fmt.Sprintf("failed to parse gitleaks report: %v", err),
		}
	}

	var findings []Vulnerability
	for i, f := range raw {
		findings = append(findings, Vulnerability{
			ID:       i + 1, // local sequence; global ID assigned by scan_service
			Source:   "gitleaks",
			Severity: "HIGH", // secrets are always treated as HIGH
			Title:    f.Description,
			File:     f.File,
			Line:     f.StartLine,
			PkgName:  f.RuleID,
		})
	}

	return ScanResult{
		Tool:     "gitleaks",
		Success:  true,
		Findings: findings,
	}
}
