import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Search, ShieldAlert } from 'lucide-react';
import { startScan } from '../services/api';

const Home: React.FC = () => {
  const [repoUrl, setRepoUrl] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  const handleScan = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!repoUrl) return;

    setLoading(true);
    setError(null);

    try {
      const resp = await startScan(repoUrl);
      navigate(`/dashboard/${resp.run_id}`);
    } catch (err: any) {
      setError(err.message || 'An error occurred while starting the scan');
      setLoading(false);
    }
  };

  return (
    <div className="container flex items-center justify-center fade-in" style={{ minHeight: '80vh' }}>
      <div className="glass-panel" style={{ padding: '48px', maxWidth: '600px', width: '100%', textAlign: 'center' }}>
        <div style={{
          width: '80px',
          height: '80px',
          background: 'var(--primary-glow)',
          borderRadius: '50%',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          margin: '0 auto 24px auto',
          boxShadow: '0 0 30px rgba(59, 130, 246, 0.4)'
        }}>
          <ShieldAlert size={40} className="text-primary" />
        </div>
        
        <h2 style={{ fontSize: '2.5rem', fontWeight: 700, marginBottom: '16px' }}>Secure Your Code</h2>
        <p className="text-muted" style={{ fontSize: '1.1rem', marginBottom: '40px' }}>
          Enter your GitHub repository URL below to trigger an AI-enhanced security scan detecting vulnerabilities, secrets, and misconfigurations.
        </p>

        <form onSubmit={handleScan} className="flex flex-col gap-4">
          <div style={{ position: 'relative' }}>
            <div style={{ position: 'absolute', left: '16px', top: '50%', transform: 'translateY(-50%)', color: 'var(--text-muted)' }}>
              <Search size={20} />
            </div>
            <input
              type="url"
              className="input-field"
              placeholder="https://github.com/username/repository"
              style={{ paddingLeft: '48px', paddingRight: '120px' }}
              value={repoUrl}
              onChange={(e) => setRepoUrl(e.target.value)}
              required
            />
            <button 
              type="submit" 
              className="btn btn-primary" 
              style={{ position: 'absolute', right: '6px', top: '6px', bottom: '6px', padding: '0 20px' }}
              disabled={loading || !repoUrl}
            >
              {loading ? (
                <span style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                  <div className="spinner" style={{ width: '16px', height: '16px', border: '2px solid rgba(255,255,255,0.3)', borderTopColor: 'white', borderRadius: '50%', animation: 'spin 1s linear infinite' }} />
                  Scanning...
                </span>
              ) : 'Scan Now'}
            </button>
          </div>
          
          {error && (
            <div className="glass-panel" style={{ padding: '12px', background: 'var(--danger-glow)', borderColor: 'rgba(239, 68, 68, 0.3)', color: 'var(--danger)', marginTop: '8px', fontSize: '0.9rem' }}>
              {error}
            </div>
          )}
        </form>

        <style>{`
          @keyframes spin {
            to { transform: rotate(360deg); }
          }
        `}</style>
      </div>
    </div>
  );
};

export default Home;
