package scanner

type Scanner interface {
	Name() string
	Run(repoPath string) ScanResult
}

type ScanResult struct {
	Tool     string          `json:"tool"`
	Success  bool            `json:"success"`
	Error    string          `json:"error,omitempty"`
	Findings []Vulnerability `json:"findings"`
}

type AIFix struct {
	FixCommand     string `json:"fix_command"`
	FixExplanation string `json:"fix_explanation"`
}

type Vulnerability struct {
	ID           int      `json:"id"`            // auto-incremented across all tools
	Source       string   `json:"source"`        // "trivy" | "gitleaks"
	Severity     string   `json:"severity"`
	Title        string   `json:"title"`
	PkgName      string   `json:"pkg_name"`
	Description  string   `json:"description"`
	File         string   `json:"file"`
	Line         int      `json:"line"`
	Confidence   string   `json:"confidence"`
	FixedVersion string   `json:"fixed_version"`
	PrimaryURL   string   `json:"primary_url"`
	References   []string `json:"references"`
	AIFix        *AIFix   `json:"ai_analysis,omitempty"`
}
