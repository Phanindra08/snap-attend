import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { useAuth } from "../auth/AuthContext";
import { getSections } from "../api/professor";

export default function ProfessorDashboard() {
  const { user } = useAuth() || {};
  const [sections, setSections] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    getSections()
      .then((data) => {
        setSections(data);
        setLoading(false);
      })
      .catch((err) => {
        console.error(err);
        setError("Failed to load sections.");
        setLoading(false);
      });
  }, []);

  if (loading) return <div className="page p-4">Loading...</div>;
  if (error) return <div className="page p-4 error">{error}</div>;

  return (
    <div className="page">
      <header className="topbar">
        <div>
          <h1>Professor Dashboard</h1>
          <p className="muted">
            {user?.email
              ? `Signed in as ${user.email}`
              : "Signed in as Professor"}
          </p>
        </div>
        <Link to="/profile" className="link">
          Update Profile
        </Link>
      </header>

      {sections.length === 0 ? (
        <div className="card">
          <p>You have no active sections.</p>
        </div>
      ) : (
        sections.map((section) => (
          <div key={section.sectionId} className="card mb-4">
            <h2>{section.courseName}</h2>
            <p className="subtitle">Section {section.sectionNumber}</p>
            <ul className="list">
              <li>
                <strong>Room:</strong> {section.room || "TBD"}
              </li>
              <li>
                <strong>Time:</strong> {section.startTime} – {section.endTime}
              </li>
              <li>
                <strong>Day:</strong> {section.day}
              </li>
            </ul>

            <Link
              to={`/professor/qr?sectionId=${section.sectionId}`}
              className="btn btn-link-button"
            >
              Generate QR for Attendance
            </Link>
            <Link
              to={`/professor/attendance?sectionId=${section.sectionId}`}
              className="btn btn-link-button ml-4"
              style={{ marginLeft: "1rem" }}
            >
              View Report
            </Link>
          </div>
        ))
      )}
    </div>
  );
}
