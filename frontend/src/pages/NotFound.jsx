import { Link } from "react-router-dom";

export default function NotFound() {
  return (
    <div className="page-center">
      <div className="card">
        <h2>Page not found</h2>
        <p className="muted">
          The page you&apos;re looking for doesn&apos;t exist.
        </p>
        <Link to="/" className="btn btn-link-button">
          Go to Login
        </Link>
      </div>
    </div>
  );
}

