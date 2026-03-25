import { useEffect, useState } from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import Home from './pages/Home';
import Dashboard from './pages/Dashboard';
import Login from './pages/Login';
import Repositories from './pages/Repositories';
import { getUser, logout } from './services/api';
import { LogOut } from 'lucide-react';

function App() {
  const [user, setUser] = useState<any>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    getUser()
      .then(data => setUser(data))
      .catch(() => setUser(null))
      .finally(() => setLoading(false));
  }, []);

  const handleLogout = async () => {
    await logout();
    setUser(null);
  };

  if (loading) return null; // or standard loader

  return (
    <BrowserRouter>
      <div className="app-wrapper">
        <header style={{
          padding: '24px',
          borderBottom: '1px solid var(--panel-border)',
          background: 'rgba(11, 15, 25, 0.8)',
          backdropFilter: 'blur(10px)',
          position: 'sticky',
          top: 0,
          zIndex: 100
        }}>
          <div className="container flex items-center justify-between">
            <h1 style={{ fontSize: '1.5rem', fontWeight: 700, margin: 0, display: 'flex', alignItems: 'center', gap: '8px' }}>
              <span className="text-primary">AI</span>
              <span>DevSecOps</span>
            </h1>
            <nav className="flex items-center gap-4">
              {user ? (
                <>
                  <a href="/" style={{ color: 'var(--text-main)', textDecoration: 'none', fontWeight: 500, marginRight: '16px' }}>Scanner</a>
                  <a href="/repos" style={{ color: 'var(--text-main)', textDecoration: 'none', fontWeight: 500, marginRight: '16px' }}>My Repositories</a>
                  <span className="text-muted text-sm mr-2">{user.name}</span>
                  <button onClick={handleLogout} className="btn" style={{ padding: '6px 12px', background: 'transparent', border: '1px solid var(--panel-border)', color: 'var(--text-main)', fontSize: '0.875rem' }}>
                    <LogOut size={16} /> Logout
                  </button>
                </>
              ) : (
                <a href="/login" style={{ color: 'var(--text-main)', textDecoration: 'none', fontWeight: 500 }}>Login</a>
              )}
            </nav>
          </div>
        </header>

        <main style={{ minHeight: 'calc(100vh - 85px)' }}>
          <Routes>
            <Route path="/login" element={user ? <Navigate to="/" /> : <Login />} />
            <Route path="/" element={user ? <Home /> : <Navigate to="/login" />} />
            <Route path="/dashboard/:runId" element={user ? <Dashboard /> : <Navigate to="/login" />} />
            <Route path="/repos" element={user ? <Repositories /> : <Navigate to="/login" />} />
          </Routes>
        </main>
      </div>
    </BrowserRouter>
  );
}

export default App;
