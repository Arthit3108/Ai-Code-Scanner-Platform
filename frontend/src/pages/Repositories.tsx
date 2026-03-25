import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { getUserRepos } from '../services/api';
import { Database, Clock, RefreshCw, ShieldAlert, ShieldCheck } from 'lucide-react';
import type { ScanJob } from '../types/scan';

const Repositories: React.FC = () => {
  const [repos, setRepos] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchRepos = async () => {
      try {
        const data = await getUserRepos();
        setRepos(data);
      } catch (err: any) {
        setError(err.message || 'Failed to fetch repositories');
      } finally {
        setLoading(false);
      }
    };
    fetchRepos();
  }, []);

  if (loading) {
    return (
      <div className="container mt-8 flex justify-center fade-in">
        <RefreshCw className="text-primary spin-animation" size={48} />
      </div>
    );
  }

  if (error) {
    return (
      <div className="container mt-8 text-center fade-in">
        <div className="glass-panel p-8" style={{ background: 'var(--danger-glow)', borderColor: 'rgba(239, 68, 68, 0.3)', maxWidth: '600px', margin: '0 auto' }}>
          <ShieldAlert size={48} className="text-error" style={{ margin: '0 auto 16px auto' }} />
          <h2 style={{ fontSize: '1.5rem', marginBottom: '8px', color: 'var(--danger)' }}>Error loading repositories</h2>
          <p>{error}</p>
        </div>
      </div>
    );
  }

  return (
    <div className="container mt-8 fade-in">
      <div className="flex justify-between items-center mb-8">
        <div>
          <h2 style={{ fontSize: '2rem', fontWeight: 700 }} className="flex items-center gap-3">
            <Database size={32} className="text-primary" /> My Repositories
          </h2>
          <p className="text-muted mt-2">View all your previously scanned templates and code repositories.</p>
        </div>
      </div>

      {repos.length === 0 ? (
        <div className="glass-panel p-12 text-center flex flex-col items-center justify-center">
          <Database size={64} className="text-muted mb-4" />
          <h3 style={{ fontSize: '1.5rem', marginBottom: '8px' }}>No repositories yet</h3>
          <p className="text-muted mb-6">Scan your first repository to see it listed here.</p>
          <button className="btn btn-primary" onClick={() => navigate('/')}>Scan a Repository</button>
        </div>
      ) : (
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(320px, 1fr))', gap: '24px' }}>
          {repos.map((repo, idx) => {
            const latestScan = repo.scan_jobs && repo.scan_jobs.length > 0 ? repo.scan_jobs[0] as ScanJob : null;
            
            let totalIssues = 0;
            if (latestScan && latestScan.results) {
              latestScan.results.forEach(r => {
                totalIssues += (r.findings || []).length;
              });
            }

            return (
              <div key={repo.id || idx} className="glass-panel p-6" style={{ transition: 'all 0.3s ease' }}>
                <h3 style={{ fontSize: '1.2rem', fontWeight: 600, wordBreak: 'break-all', marginBottom: '12px' }}>
                  <a href={repo.url} target="_blank" rel="noreferrer" className="text-accent" style={{ textDecoration: 'none' }}>
                    {repo.url.replace('https://github.com/', '')}
                  </a>
                </h3>
                
                <div className="flex justify-between items-center mb-6 text-sm text-muted">
                  <span className="flex items-center gap-1"><Clock size={16} /> {new Date(repo.updated_at).toLocaleDateString()}</span>
                  <span>{repo.scan_jobs?.length || 0} Scans</span>
                </div>

                {latestScan ? (
                  <div style={{ background: 'var(--panel-bg)', padding: '16px', borderRadius: '8px', border: '1px solid var(--panel-border)' }}>
                    <div className="flex justify-between items-center mb-4">
                      <span style={{ fontSize: '0.875rem', fontWeight: 600, color: 'var(--text-muted)' }}>LATEST SCAN</span>
                      <span style={{ 
                        fontSize: '0.75rem', padding: '4px 8px', borderRadius: '20px', fontWeight: 700, textTransform: 'uppercase',
                        background: latestScan.status === 'done' ? 'rgba(16, 185, 129, 0.2)' : 'var(--danger-glow)',
                        color: latestScan.status === 'done' ? 'var(--success)' : 'var(--danger)'
                      }}>
                        {latestScan.status}
                      </span>
                    </div>
                    
                    {latestScan.status === 'done' && (
                      <div className="flex items-center gap-3 mb-4">
                        {totalIssues === 0 ? (
                          <span className="text-success flex items-center gap-1 font-semibold"><ShieldCheck size={18} /> Secure (0 Issues)</span>
                        ) : (
                          <span className="text-danger flex items-center gap-1 font-semibold"><ShieldAlert size={18} /> {totalIssues} Issues Found</span>
                        )}
                      </div>
                    )}

                    <button 
                      className="btn btn-secondary w-full"
                      onClick={() => navigate(`/dashboard/${latestScan.run_id}`)}
                    >
                      View Report
                    </button>
                  </div>
                ) : (
                  <p className="text-sm text-muted">No completed scans yet.</p>
                )}
              </div>
            );
          })}
        </div>
      )}

      <style>{`
        .spin-animation { animation: spin 2s linear infinite; }
        @keyframes spin { 100% { transform: rotate(360deg); } }
      `}</style>
    </div>
  );
};

export default Repositories;
