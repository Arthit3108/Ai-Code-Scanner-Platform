package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"google.golang.org/genai"
)

type GeminiClient struct{}

func NewGeminiClient() *GeminiClient {
	return &GeminiClient{}
}

type VulnerabilityInput struct {
	ID          int    `json:"id"`
	Source      string `json:"source"`
	Severity    string `json:"severity"`
	Title       string `json:"title"`
	PkgName     string `json:"pkg_name"`
	Description string `json:"description"`
	File        string `json:"file"`
	Line        int    `json:"line"`
}

type AiAnalysis struct {
	ID             int    `json:"id"`
	FixCommand     string `json:"fix_command"`
	FixExplanation string `json:"fix_explanation"`
}

func (c *GeminiClient) AnalyzeVulnerabilities(vulns []VulnerabilityInput) ([]AiAnalysis, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  os.Getenv("GEMINI_API_KEY"),
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create gemini client: %w", err)
	}

	var trivyVulns []VulnerabilityInput
	var gitleaksFindings []VulnerabilityInput

	for _, v := range vulns {
		if v.Source == "trivy" {
			trivyVulns = append(trivyVulns, v)
		} else {
			gitleaksFindings = append(gitleaksFindings, v)
		}
	}

	trivyJSON, _ := json.MarshalIndent(trivyVulns, "", "  ")
	gitleaksJSON, _ := json.MarshalIndent(gitleaksFindings, "", "  ")

	prompt := fmt.Sprintf(`
		You are a security analysis assistant. 
		Analyze the following vulnerability scan results and secrets detection and provide a concise security report.

		Vulnerabilities found by Trivy:
		<vulns>
		%s
		</vulns>

		Secrets found by Gitleaks:
		<secrets>
		%s
		</secrets>

		Respond in this exact JSON array format, nothing else (no markdown block):
		[
			{
				"id": 1,
				"fix_command": "command here",
				"fix_explanation": "One sentence why this fixes it"
			}
		]
		Note: For Gitleaks findings, use the ID provided in the list.
		Please provide analysis for each vulnerability and secret.`, string(trivyJSON), string(gitleaksJSON))

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("gemini generate content failed: %w", err)
	}

	cleanText := result.Text()
	cleanText = strings.TrimPrefix(cleanText, "```json")
	cleanText = strings.TrimPrefix(cleanText, "```")
	cleanText = strings.TrimSuffix(cleanText, "```")
	cleanText = strings.TrimSpace(cleanText)

	var analysis []AiAnalysis
	if err := json.Unmarshal([]byte(cleanText), &analysis); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	return analysis, nil
}


