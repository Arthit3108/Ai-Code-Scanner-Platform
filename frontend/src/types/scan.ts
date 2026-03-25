export type ScanStatus = 'pending' | 'cloning' | 'scanning' | 'done' | 'failed';

export interface AIFix {
  fix_command: string;
  fix_explanation: string;
}

export interface Vulnerability {
  id: number;
  source: string;
  severity: string;
  title: string;
  pkg_name: string;
  description: string;
  file?: string;
  line?: string;
  ai_analysis?: AIFix;
}

export interface ScanResult {
  tool: string;
  target: string;
  findings: Vulnerability[];
}

export interface ScanJob {
  run_id: string;
  repo_url: string;
  status: ScanStatus;
  error?: string;
  results?: ScanResult[];
  created_at: string;
  updated_at: string;
}

export interface ScanResponse {
  run_id: string;
  status: ScanStatus;
  repo_url: string;
}

export interface ScanRequest {
  repo_url: string;
}
