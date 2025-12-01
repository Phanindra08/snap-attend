import { Link } from "react-router-dom";
import { useAuth } from "../auth/AuthContext";

export default function AdminDashboard() {
  const { user } = useAuth() || {};

  return (
    <div className="page">
      <header className="topbar">
        <div>
          <h1>Admin Dashboard</h1>
          <p className="muted">
            {user?.email ? `Signed in as ${user.email}` : "Signed in as Admin"}
          </p>
        </div>
        <Link to="/profile" className="link">
          Update Profile
        </Link>
      </header>

      <div className="card">
        <h2>Management</h2>
        <p className="subtitle">Manage users and courses.</p>
        <div className="grid-menu">
          <Link to="/admin/students" className="btn btn-secondary">
            Manage Students
          </Link>
          <Link to="/admin/professors" className="btn btn-secondary">
            Manage Professors
          </Link>
          <Link to="/admin/courses" className="btn btn-secondary">
            Manage Courses
          </Link>
        </div>
      </div>
    </div>
  );
}
