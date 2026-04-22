import type { ScanResponse, ScanJob } from '../types/scan';

const API_URL = import.meta.env.VITE_API_URL || '/api/v1';

// startScan triggers a new security scan for the given repository URL.
// Credentials are included to associate the scan request with the authenticated user session.
export const startScan = async (repoUrl: string): Promise<ScanResponse> => {
  const response = await fetch(`${API_URL}/scan`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ repo_url: repoUrl }),
    // Include credentials to send JWT cookies automatically for session-based auth
    credentials: 'include',
  });

  if (!response.ok) {
    const err = await response.json();
    throw new Error(err.error || 'Failed to start scan');
  }

  return response.json();
};

// ScanJob represents a single security scanning execution for a repository (Backend Model).
/*
type ScanJob struct {
	RunID        string               `gorm:"primaryKey" json:"run_id"`
	UserID       *uint                `gorm:"index" json:"user_id,omitempty"` 
	RepositoryID *uint                `gorm:"index" json:"repository_id,omitempty"` 
	RepoURL      string               `json:"repo_url"`
	Status       ScanStatus           `json:"status"`
	Error        string               `json:"error,omitempty"`
	// Results are stored as a JSON blob in the database for flexibility.
	Results      []scanner.ScanResult `gorm:"serializer:json" json:"results,omitempty"`
	CreatedAt    time.Time            `json:"created_at"`
	UpdatedAt    time.Time            `json:"updated_at"`
}
*/

// getScanStatus fetches the current status and results of a specific scan job.
// Uses polling logic on the frontend to monitor async scan progression.
export const getScanStatus = async (runId: string): Promise<ScanJob> => {
  const response = await fetch(`${API_URL}/scan/${runId}`, {
    credentials: 'include' // Must include credentials to access protected routes
  });

  if (!response.ok) {
    const err = await response.json();
    throw new Error(err.error || 'Failed to get scan status');
  }

  return response.json();
};

/*
// Frontend Polling Logic Example:
useEffect(() => {
  if (!runId) return;

  // fetchStatus polls the backend to check if the scan is finished.
  const fetchStatus = async () => {
    try {
      const data = await getScanStatus(runId);
      setJob(data);
      // If status is 'done' or 'failed', we stop the polling process.
      if (data.status === 'done' || data.status === 'failed') {
        return true; // Stop polling
      }
      return false;
    } catch (err: any) {
      setError(err.message || 'Failed to fetch scan status');
      return true; // Stop polling on error
    }
  };

  fetchStatus();
  // Set up a 3-second interval to check for status updates (Polling Pattern)
  const interval = setInterval(async () => {
    const shouldStop = await fetchStatus();
    if (shouldStop) clearInterval(interval);
  }, 3000);

  return () => clearInterval(interval);
}, [runId]);
*/

export const getUser = async (): Promise<any> => {
  const response = await fetch(`${API_URL}/auth/me`, {
    credentials: 'include'
  });
  if (!response.ok) {
    throw new Error('Unauthorized');
  }
  return response.json();
};

export const getUserRepos = async (): Promise<any[]> => {
  const response = await fetch(`${API_URL}/repos`, {
    method: 'GET',
    headers: {
      'Accept': 'application/json'
    },
    credentials: 'include'
  });
  if (!response.ok) {
    throw new Error('Failed to fetch user repositories');
  }
  return response.json();
};

export const logout = async () => {
  await fetch(`${API_URL}/auth/logout`, {
    method: 'POST',
    credentials: 'include'
  });
};
