import React from 'react';
import { ShieldCheck } from 'lucide-react';

const API_URL = import.meta.env.VITE_API_URL || '/api/v1';

const Login: React.FC = () => {
  const handleGithubLogin = () => {
    window.location.href = `${API_URL}/auth/github/login`;
  };

  return (
    <div className="container flex items-center justify-center fade-in" style={{ minHeight: '80vh' }}>
      <div className="glass-panel" style={{ padding: '48px', maxWidth: '400px', width: '100%', textAlign: 'center' }}>
        <div style={{
          width: '64px',
          height: '64px',
          background: 'var(--primary-glow)',
          borderRadius: '50%',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          margin: '0 auto 24px auto',
          boxShadow: '0 0 20px rgba(59, 130, 246, 0.4)'
        }}>
          <ShieldCheck size={32} className="text-primary" />
        </div>
        
        <h2 style={{ fontSize: '1.75rem', fontWeight: 700, marginBottom: '8px' }}>Welcome</h2>
        <p className="text-muted" style={{ marginBottom: '32px' }}>
          Sign in to access your secure code analysis platform.
        </p>

        <div className="flex flex-col gap-4 justify-center">
          <button 
            onClick={handleGithubLogin} 
            className="btn" 
            style={{ background: '#24292e', color: 'white', padding: '14px' }}
          >
            😎 Continue with GitHub
          </button>
        </div>
      </div>
    </div>
  );
};

export default Login;
