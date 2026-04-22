import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { getScanStatus } from '../services/api';
import type { ScanJob, ScanStatus, Vulnerability } from '../types/scan';
import { ShieldCheck, ShieldAlert, AlertTriangle, TerminalSquare, Loader2 } from 'lucide-react';

const Dashboard: React.FC = () => {
  const { runId } = useParams<{ runId: string }>();
  const navigate = useNavigate();
  const [job, setJob] = useState<ScanJob | null>(null);
  const [error, setError] = useState<string | null>(null);

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

  if (error) {
    return (
      <div className="container mt-8 fade-in text-center">
        <div className="glass-panel p-8" style={{ background: 'var(--danger-glow)', borderColor: 'rgba(239, 68, 68, 0.3)', maxWidth: '600px', margin: '0 auto' }}>
          <AlertTriangle size={48} className="text-error" style={{ margin: '0 auto 16px auto' }} />
          <h2 style={{ fontSize: '1.5rem', marginBottom: '8px', color: 'var(--danger)' }}>Error loading Dashboard</h2>
          <p>{error}</p>
          <button className="btn btn-primary mt-4" onClick={() => navigate('/')}>Back to Home</button>
        </div>
      </div>
    );
  }

  if (!job) {
    return (
      <div className="container mt-8 flex justify-center fade-in">
        <Loader2 className="text-primary spin-animation" size={48} />
      </div>
    );
  }

  const isScanning = job.status === 'pending' || job.status === 'cloning' || job.status === 'scanning';

  return (
    <div className="container mt-8 fade-in">
      {/* Header section */}
      <div className="flex justify-between items-center mb-8">
        <div>
          <h2 style={{ fontSize: '2rem', fontWeight: 700 }}>Scan Results</h2>
          <p className="text-muted flex items-center gap-2 mt-2">
            Target: <a href={job.repo_url} target="_blank" rel="noreferrer" style={{ color: 'var(--primary)', textDecoration: 'none' }}>{job.repo_url}</a>
          </p>
        </div>
        <StatusBadge status={job.status} />
      </div>

      {isScanning && (
        <div className="glass-panel p-8 text-center" style={{ maxWidth: '600px', margin: '40px auto' }}>
          <Loader2 className="text-primary spin-animation" size={64} style={{ margin: '0 auto 24px auto' }} />
          <h3 style={{ fontSize: '1.5rem', marginBottom: '8px' }}>Analysis in Progress</h3>
          <p className="text-muted" style={{ textTransform: 'capitalize' }}>Current phase: {job.status}</p>
          <div style={{ width: '100%', height: '4px', background: 'var(--panel-border)', borderRadius: '2px', marginTop: '24px', overflow: 'hidden' }}>
            <div style={{
              width: job.status === 'pending' ? '10%' : job.status === 'cloning' ? '40%' : '70%',
              height: '100%',
              background: 'var(--primary)',
              transition: 'width 1s ease'
            }} />
          </div>
        </div>
      )}

      {job.status === 'done' && (
        <div className="fade-in">
          <SummaryCards results={job.results || []} />
          <FindingsList results={job.results || []} />
        </div>
      )}

      {job.status === 'failed' && (
        <div className="glass-panel p-8 text-center" style={{ background: 'var(--danger-glow)', borderColor: 'rgba(239, 68, 68, 0.3)' }}>
          <AlertTriangle size={48} className="text-error" style={{ margin: '0 auto 16px auto' }} />
          <h3 style={{ fontSize: '1.5rem', marginBottom: '8px', color: 'var(--danger)' }}>Scan Failed</h3>
          <p>{job.error || 'An unknown error occurred during the scan.'}</p>
        </div>
      )}

      <style>{`
        .spin-animation { animation: spin 2s linear infinite; }
      `}</style>
    </div>
  );
};

const StatusBadge = ({ status }: { status: ScanStatus }) => {
  let color = 'var(--text-muted)';
  let bg = 'var(--panel-border)';
  
  if (status === 'done') { color = 'var(--success)'; bg = 'rgba(16, 185, 129, 0.2)'; }
  if (status === 'failed') { color = 'var(--danger)'; bg = 'var(--danger-glow)'; }
  if (status === 'scanning' || status === 'cloning') { color = 'var(--warning)'; bg = 'rgba(245, 158, 11, 0.2)'; }

  return (
    <div style={{
      display: 'inline-flex',
      alignItems: 'center',
      gap: '8px',
      padding: '8px 16px',
      background: bg,
      color: color,
      borderRadius: '20px',
      fontWeight: 600,
      textTransform: 'uppercase',
      fontSize: '0.875rem'
    }}>
      {status === 'done' ? <ShieldCheck size={18} /> : (status === 'failed' ? <ShieldAlert size={18} /> : <Loader2 size={18} className="spin-animation" />)}
      {status}
    </div>
  );
};

const SummaryCards = ({ results }: { results: any[] }) => {
  let total = 0, critical = 0, high = 0, medium = 0, low = 0;

  results.forEach(r => {
    const findings = r.findings || [];
    total += findings.length;
    findings.forEach((f: Vulnerability) => {
      const sev = (f.severity || '').toLowerCase();
      if (sev === 'critical') critical++;
      else if (sev === 'high') high++;
      else if (sev === 'medium') medium++;
      else low++;
    });
  });

  return (
    <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: '24px', marginBottom: '40px' }}>
      <div className="glass-panel p-6 text-center">
        <div style={{ fontSize: '2.5rem', fontWeight: 700, color: 'var(--text-main)', marginBottom: '8px' }}>{total}</div>
        <div className="text-muted text-sm uppercase" style={{ letterSpacing: '1px' }}>Total Issues</div>
      </div>
      <div className="glass-panel p-6 text-center" style={{ borderColor: critical > 0 ? 'var(--danger)' : '' }}>
        <div style={{ fontSize: '2.5rem', fontWeight: 700, color: 'var(--danger)', marginBottom: '8px' }}>{critical}</div>
        <div className="text-muted text-sm uppercase" style={{ letterSpacing: '1px' }}>Critical</div>
      </div>
      <div className="glass-panel p-6 text-center" style={{ borderColor: high > 0 ? 'var(--warning)' : '' }}>
        <div style={{ fontSize: '2.5rem', fontWeight: 700, color: 'var(--warning)', marginBottom: '8px' }}>{high}</div>
        <div className="text-muted text-sm uppercase" style={{ letterSpacing: '1px' }}>High</div>
      </div>
      <div className="glass-panel p-6 text-center">
        <div style={{ fontSize: '2.5rem', fontWeight: 700, color: 'var(--primary)', marginBottom: '8px' }}>{medium + low}</div>
        <div className="text-muted text-sm uppercase" style={{ letterSpacing: '1px' }}>Medium/Low</div>
      </div>
    </div>
  );
};

const FindingsList = ({ results }: { results: any[] }) => {
  const [filter, setFilter] = useState<string>('ALL');

  const allFindings = results.flatMap(r => r.findings || []).sort((a, b) => {
    const sevScore: Record<string, number> = { 'CRITICAL': 4, 'HIGH': 3, 'MEDIUM': 2, 'LOW': 1 };
    return (sevScore[b.severity?.toUpperCase() || ''] || 0) - (sevScore[a.severity?.toUpperCase() || ''] || 0);
  });

  const filteredFindings = filter === 'ALL'
    ? allFindings
    : allFindings.filter(f => (f.severity?.toUpperCase() || 'UNKNOWN') === filter);

  if (allFindings.length === 0) {
    return (
      <div className="glass-panel p-8 text-center align-center">
        <ShieldCheck size={64} className="text-success" style={{ margin: '0 auto 16px auto' }} />
        <h3 style={{ fontSize: '1.5rem', color: 'var(--success)' }}>Secure!</h3>
        <p className="text-muted">No vulnerabilities found in this repository.</p>
      </div>
    );
  }

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '16px', marginBottom: '60px' }}>
      <div className="flex gap-2 mb-4" style={{ overflowX: 'auto', paddingBottom: '8px' }}>
        {['ALL', 'CRITICAL', 'HIGH', 'MEDIUM', 'LOW'].map(sev => (
          <button
            key={sev}
            onClick={() => setFilter(sev)}
            style={{
              padding: '6px 16px',
              borderRadius: '20px',
              background: filter === sev ? 'var(--primary)' : 'var(--panel-bg)',
              color: filter === sev ? 'white' : 'var(--text-muted)',
              border: `1px solid ${filter === sev ? 'var(--primary)' : 'var(--panel-border)'}`,
              cursor: 'pointer',
              fontWeight: 600,
              fontSize: '0.875rem',
              transition: 'all 0.2s ease'
            }}
          >
            {sev}
          </button>
        ))}
      </div>

      {filteredFindings.length === 0 ? (
        <div className="text-center text-muted p-8 glass-panel">No findings match the selected severity.</div>
      ) : (
        filteredFindings.map((f: Vulnerability, idx) => (
          <div key={idx} className="glass-panel p-6">
            <div className="flex justify-between items-start mb-4">
              <div>
                <div className="flex items-center gap-2 mb-2">
                  <span style={{
                    padding: '4px 8px', borderRadius: '4px', fontSize: '0.75rem', fontWeight: 700,
                    background: f.severity?.toUpperCase() === 'CRITICAL' ? 'var(--danger)' :
                                f.severity?.toUpperCase() === 'HIGH' ? 'var(--warning)' :
                                'var(--primary)', color: 'white'
                  }}>
                    {f.severity || 'UNKNOWN'}
                  </span>
                  <span className="text-muted text-sm">{f.source}</span>
                  <span className="text-accent text-sm ml-2">{f.pkg_name}</span>
                </div>
                <h3 style={{ fontSize: '1.25rem', fontWeight: 600 }}>{f.title || f.id}</h3>
              </div>
              {f.file && (
                <div className="text-muted text-sm bg-panel-bg p-2" style={{ background: 'rgba(0,0,0,0.2)', borderRadius: '4px', border: '1px solid var(--panel-border)' }}>
                  {f.file}{f.line ? `:${f.line}` : ''}
                </div>
              )}
            </div>
            
            {/* Removed description to keep it short as requested */}
            
            {f.ai_analysis && (f.ai_analysis.fix_explanation || f.ai_analysis.fix_command) ? (
              <div style={{ background: 'rgba(139, 92, 246, 0.1)', borderLeft: '4px solid var(--accent)', padding: '16px', borderRadius: '4px' }}>
                <h4 className="flex items-center gap-2 text-accent mb-2"><TerminalSquare size={18} /> AI Remediation</h4>
                
                {f.ai_analysis.fix_explanation && (
                  <p style={{ fontSize: '0.9rem', marginBottom: '12px' }}>{f.ai_analysis.fix_explanation}</p>
                )}
                
                {f.ai_analysis.fix_command && (
                  <div style={{ background: '#000', padding: '12px', borderRadius: '6px', fontFamily: 'monospace', color: '#10b981', overflowX: 'auto', whiteSpace: 'pre-wrap' }}>
                    {f.ai_analysis.fix_command}
                  </div>
                )}
              </div>
            ) : (
              <p className="text-muted text-sm italic">No AI remediation available.</p>
            )}
          </div>
        ))
      )}
    </div>
  );
};

export default Dashboard;
