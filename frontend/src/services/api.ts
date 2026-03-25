import type { ScanResponse, ScanJob } from '../types/scan';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:3000/api/v1';

export const startScan = async (repoUrl: string): Promise<ScanResponse> => {
  const response = await fetch(`${API_URL}/scan`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ repo_url: repoUrl }),
    credentials: 'include',
  });

  if (!response.ok) {
    const err = await response.json();
    throw new Error(err.error || 'Failed to start scan');
  }

  return response.json();
};

export const getScanStatus = async (runId: string): Promise<ScanJob> => {
  const response = await fetch(`${API_URL}/scan/${runId}`, {
    credentials: 'include'
  });

  if (!response.ok) {
    const err = await response.json();
    throw new Error(err.error || 'Failed to get scan status');
  }

  return response.json();
};

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
