import { useState, useEffect } from "react";
import { useSearchParams, useNavigate, useLocation } from "react-router-dom";
import Button from "../components/Button";
import FormError from "../components/FormError";
import Input from "../components/Input";
import { markAttendance } from "../api/student";
import { useAuth } from "../auth/AuthContext";

export default function StudentScan() {
  const [searchParams] = useSearchParams();
  const qrIdParam = searchParams.get("qr"); // Get QR ID from URL
  const { user } = useAuth();

  const [qrId, setQrId] = useState(qrIdParam || "");
  const [name, setName] = useState(user?.firstName ? `${user.firstName} ${user.lastName}` : "");
  const [question, setQuestion] = useState("");
  const [answer, setAnswer] = useState("");
  const [error, setError] = useState("");
  const [message, setMessage] = useState("");
  const [busy, setBusy] = useState(false);

  useEffect(() => {
    if (user) {
      setName(user.firstName ? `${user.firstName} ${user.lastName}` : "");
    }
  }, [user]);

  useEffect(() => {
    if (qrIdParam) {
      setQrId(qrIdParam);
    }
  }, [qrIdParam]);

  const handleSubmit = (e) => {
    e.preventDefault();
    setError("");
    setMessage("");

    if (!qrId) {
      setError("Please enter QR Id from the scanned code.");
      return;
    }

    if (!name.trim()) {
      setError("Please enter your name.");
      return;
    }

    setBusy(true);

    // Use browser geolocation
    navigator.geolocation.getCurrentPosition(
      async (position) => {
        try {
          const latitude = position.coords.latitude;
          const longitude = position.coords.longitude;

          await markAttendance({
            qrId: Number(qrId),
            studentName: name.trim(),
            latitude,
            longitude,
            question: question.trim(),
            answer: answer.trim(),
          });

          setMessage("Attendance submitted successfully!");
          setName("");
          setQuestion("");
          setAnswer("");
        } catch (err) {
          setError(
            err?.response?.data?.error ||
            err.message ||
            "Failed to submit attendance."
          );
        } finally {
          setBusy(false);
        }
      },
      (geoError) => {
        setBusy(false);
        setError(
          `Error getting location: ${geoError.message}. Please enable location access.`
        );
      }
    );
  };

  return (
    <div className="page-center">
      <div className="card">
        <h2>Scan QR & Mark Attendance</h2>

        {!qrIdParam && (
          <div className="alert alert-warning">
            No QR Code detected. Please scan the code provided by your professor.
          </div>
        )}

        <form onSubmit={handleSubmit} className="form">
          {/* QR ID is handled automatically from URL */}
          {qrIdParam && <input type="hidden" value={qrId} />}

          <Input
            label="Your Name"
            placeholder="Enter your full name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            required
          />

          <Input
            label="Answer to Class Question (Optional)"
            placeholder="Your answer..."
            value={answer}
            onChange={(e) => setAnswer(e.target.value)}
          />

          <Input
            label="Your Question (Optional)"
            placeholder="Ask a question..."
            value={question}
            onChange={(e) => setQuestion(e.target.value)}
          />

          <FormError message={error} />
          {message && <div className="form-success">{message}</div>}

          <Button type="submit" disabled={busy}>
            {busy ? "Submitting..." : "Submit Attendance"}
          </Button>
        </form>

        <p className="muted small">
          This page uses{" "}
          <code>navigator.geolocation.getCurrentPosition(...)</code> to send
          your current latitude and longitude along with the attendance
          submission.
        </p>
      </div>
    </div>
  );
}
