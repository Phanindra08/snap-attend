import AppRoutes from "./routes";
import "./styles/App.css";
import { useAuth } from "./auth/AuthContext";
import { Link } from "react-router-dom";

export default function App() {
  const { user, logout } = useAuth() || {};

  return (
    <div className="app-root">
      <nav className="navbar">
        <div className="navbar-left">
          <span className="logo">SnapAttend</span>
        </div>
        <div className="navbar-right">
          {user && (
            <>
              <span className="navbar-user">
                {user.email} ({user.profile || user.role})
              </span>
              <button className="btn btn-small" onClick={logout}>
                Logout
              </button>
            </>
          )}
          {!user && (
            <Link to="/" className="link">
              Login
            </Link>
          )}
        </div>
      </nav>
      <main className="app-main">
        <AppRoutes />
      </main>
    </div>
  );
}
