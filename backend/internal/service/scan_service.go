package services

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/Arthit3108/devsecops-platform/internal/ai"
	"github.com/Arthit3108/devsecops-platform/internal/db"
	"github.com/Arthit3108/devsecops-platform/internal/models"
	"github.com/Arthit3108/devsecops-platform/internal/scanner"
	"github.com/google/uuid"
)

type ScanService struct {}

func NewScanService() *ScanService {
	return &ScanService{}
}

func (s *ScanService) InitScan(repoURL string, user *models.User) (string, error) {
	runID := uuid.New().String()

	var userID *uint
	var repositoryID *uint

	if user != nil {
		userID = &user.ID
		var repo models.Repository
		if err := db.DB.Where("user_id = ? AND url = ?", user.ID, repoURL).FirstOrCreate(&repo, models.Repository{
			UserID: user.ID,
			URL:    repoURL,
		}).Error; err == nil {
			repositoryID = &repo.ID
		}
	}

	job := models.ScanJob{
		RunID:        runID,
		UserID:       userID,
		RepositoryID: repositoryID,
		RepoURL:      repoURL,
		Status:       models.StatusPending,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := db.DB.Create(&job).Error; err != nil {
		return "", err
	}

	go s.runPipeline(job)

	return runID, nil
}

func (s *ScanService) runPipeline(job models.ScanJob) {
	// ── Step 1: Clone ──────────────────────────────────────────────────────
	s.updateStatus(&job, models.StatusCloning, "")

	tmpDir, err := os.MkdirTemp("", "devsecops-scan-*")
	if err != nil {
		s.updateStatus(&job, models.StatusFailed, err.Error())
		return
	}
	defer os.RemoveAll(tmpDir) // always clean up after scan

	if err := s.cloneRepo(job.RepoURL, tmpDir); err != nil {
		s.updateStatus(&job, models.StatusFailed, fmt.Sprintf("clone failed: %v", err))
		return
	}

	// ── Step 2: Scan (all tools concurrently) ──────────────────────────────
	s.updateStatus(&job, models.StatusScanning, "")

	tools := []string{"trivy", "gitleaks"}
	results := scanner.RunAll(tools, tmpDir)

	// ── Step 3: Re-assign global IDs across all tools ──────────────────────
	globalID := 1
	for i := range results {
		for j := range results[i].Findings {
			results[i].Findings[j].ID = globalID
			globalID++
		}
	}

	// ── Step 4: Collect ALL findings, send to AI in one call ───────────────
	var allFindings []*scanner.Vulnerability
	for i := range results {
		for j := range results[i].Findings {
			allFindings = append(allFindings, &results[i].Findings[j])
		}
	}

	if len(allFindings) > 0 {
		aiClient := ai.NewGeminiClient()

		aiInputs := make([]ai.VulnerabilityInput, len(allFindings))
		for i, f := range allFindings {
			aiInputs[i] = ai.VulnerabilityInput{
				ID:          f.ID,
				Source:      f.Source,
				Severity:    f.Severity,
				Title:       f.Title,
				PkgName:     f.PkgName,
				Description: f.Description,
				File:        f.File,
				Line:        f.Line,
			}
		}

		aiResp, err := aiClient.AnalyzeVulnerabilities(aiInputs)
		if err == nil {
			fixMap := make(map[int]scanner.AIFix)
			for _, fix := range aiResp {
				fixMap[fix.ID] = scanner.AIFix{
					FixCommand:     fix.FixCommand,
					FixExplanation: fix.FixExplanation,
				}
			}
			for _, f := range allFindings {
				if fix, ok := fixMap[f.ID]; ok {
					f.AIFix = &fix
				}
			}
		} else {
			fmt.Printf("AI analysis failed: %v\n", err)
		}
	}

	// ── Step 5: Store results ──────────────────────────────────────────────
	job.Results = results
	job.Status = models.StatusDone
	job.UpdatedAt = time.Now()
	db.DB.Save(&job)
}

func (s *ScanService) updateStatus(job *models.ScanJob, status models.ScanStatus, errMsg string) {
	job.Status = status
	job.Error = errMsg
	job.UpdatedAt = time.Now()
	db.DB.Model(job).Updates(map[string]interface{}{
		"status":     status,
		"error":      errMsg,
		"updated_at": time.Now(),
	})
}

func (s *ScanService) GetJob(runID string) (*models.ScanJob, bool) {
	var job models.ScanJob
	if err := db.DB.Where("run_id = ?", runID).First(&job).Error; err != nil {
		return nil, false
	}
	return &job, true
}

func (s *ScanService) cloneRepo(url, dest string) error {
	fmt.Printf("Cloning repository %s into: %s\n", url, dest)
	cmd := exec.Command("git", "clone", "--depth=1", url, dest)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, string(output))
	}
	return nil
}
