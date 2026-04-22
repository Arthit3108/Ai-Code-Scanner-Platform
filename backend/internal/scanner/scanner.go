package scanner

// Scanner is an interface that every security tool must implement.
// This allows the runner to execute different tools interchangeably (Polymorphism).
// Each tool (e.g., Trivy, Gitleaks) defines its own logic for scanning a repository
// by implementing the Name() and Run() methods.
type Scanner interface {
	// Name returns the unique identifier for the security tool.
	Name() string
	// Run executes the security tool against the provided repository path.
	// It handles command execution (exec.Command) and output parsing.
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
