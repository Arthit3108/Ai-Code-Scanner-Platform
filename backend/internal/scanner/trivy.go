package scanner 

import (
	"encoding/json"
	"fmt"
	"os/exec"

)

type TrivyScanner struct{}

func (t *TrivyScanner) Name() string { return "trivy" }
func (t *TrivyScanner) Run(repoPath string) ScanResult {
    return runTrivy(repoPath) 
}

type trivyOutput struct {
	Results []struct {
		Vulnerabilities []struct {
			VulnerabilityID string `json:"VulnerabilityID"`
			PkgName         string `json:"PkgName"`
			Title           string `json:"Title"`
			Severity        string `json:"Severity"`
			Description     string `json:"Description"`
			FixedVersion    string `json:"FixedVersion"`
			PrimaryURL      string `json:"PrimaryURL"`
		} `json:"Vulnerabilities"`
	} `json:"Results"`
}

func runTrivy(repoPath string) ScanResult {
	trivyPath, err := exec.LookPath("trivy")
	if err != nil {
		return ScanResult{
			Tool: "trivy",
			Success: false,
			Error: "trivy not found in Path: install trivy",
		}
	}

	cmd := exec.Command(
		trivyPath, "fs",
		"--format", "json",
		"--quiet",     
		repoPath,
	)

	output, err := cmd.Output()
	if err != nil {
		return ScanResult{
			Tool: "trivy",
			Success: false,
			Error: fmt.Sprintf("trivy failed: %v", err),
		}
	}

	var trivyOut trivyOutput
	if err := json.Unmarshal(output, &trivyOut); err != nil {
		return ScanResult{
			Tool: "trivy",
			Success: false,
			Error: fmt.Sprintf("parse failed: %v", err),
		}
	}

	var findings []Vulnerability
	for _, r := range trivyOut.Results {
		for _, v := range r.Vulnerabilities {
			findings = append(findings, Vulnerability{
				ID:           len(findings) + 1, // local sequence; global ID re-assigned by scan_service
				Source:       "trivy",
				Severity:     v.Severity,
				Title:        v.Title,
				PkgName:      v.PkgName,
				Description:  v.Description,
				FixedVersion: v.FixedVersion,
				PrimaryURL:   v.PrimaryURL,
			})
		}
	}

	return ScanResult{
		Tool:     "trivy",
		Success:  true,
		Findings: findings,
	}
}