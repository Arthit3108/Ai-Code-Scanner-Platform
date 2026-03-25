package models

type ScanRequest struct {
	RepoURL string `json:"repo_url"`
}

type ScanResponse struct {
	RunID   string `json:"run_id"`
	Status  string `json:"status"`
	RepoURL string `json:"repo_url"`
}