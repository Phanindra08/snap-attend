import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import Button from "../components/Button";
import { useAuth } from "../auth/AuthContext";
import { getCourses } from "../api/student";

export default function StudentDashboard() {
  const { user } = useAuth() || {};
  const [courses, setCourses] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    getCourses()
      .then((data) => {
        setCourses(data);
        setLoading(false);
      })
      .catch((err) => {
        console.error(err);
        setError("Failed to load courses.");
        setLoading(false);
      });
  }, []);

  if (loading) return <div className="page p-4">Loading...</div>;
  if (error) return <div className="page p-4 error">{error}</div>;

  return (
    <div className="page">
      <header className="topbar">
        <div>
          <h1>Student Dashboard</h1>
          <p className="muted">
            {user?.email
              ? `Signed in as ${user.email}`
              : "Signed in as Student"}
          </p>
        </div>
        <Link to="/profile" className="link">
          Update Profile
        </Link>
      </header>

      {courses.length === 0 ? (
        <div className="card">
          <p>You are not enrolled in any courses for this semester.</p>
        </div>
      ) : (
        courses.map((course) => (
          <div key={course.sectionId} className="card mb-4">
            <h2>{course.courseName}</h2>
            <p className="subtitle">Section {course.sectionNumber}</p>
            <ul className="list">
              <li>
                <strong>Room:</strong> {course.room || "TBD"}
              </li>
              {course.schedules && course.schedules.map((sched, idx) => (
                <li key={idx}>
                  <strong>{sched.day}:</strong> {sched.startTime} – {sched.endTime}
                </li>
              ))}
            </ul>

            <div className="flex gap-2 mt-4">
              <Link to="/student/scan" className="btn btn-primary">
                Mark Attendance
              </Link>
              {/* Add View Attendance button if needed */}
            </div>
          </div>
        ))
      )}
    </div>
  );
}
