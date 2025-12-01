import { useEffect, useState } from "react";
import { useSearchParams } from "react-router-dom";
import QRCode from "react-qr-code";
import sampleCourses from "../data/sampleCourses";
import Button from "../components/Button";
import FormError from "../components/FormError";
import { generateAttendanceQr } from "../api/professor";

export default function QrGenerator() {
  const [searchParams] = useSearchParams();
  const sectionId = Number(searchParams.get("sectionId") || 1);
  const course = sampleCourses.find((c) => c.id === sectionId) || sampleCourses[0];

  const [qrValue, setQrValue] = useState("");
  const [error, setError] = useState("");
  const [busy, setBusy] = useState(false);

  const handleGenerate = async () => {
    setBusy(true);
    setError("");
    try {
      const data = await generateAttendanceQr(sectionId);
      // Encode full URL for scanning
      const qrId = data?.qr?.id || data?.qrId || data?.id;
      if (!qrId) {
        throw new Error("Failed to get QR ID from server");
      }
      const url = `${window.location.origin}/student/scan?qr=${qrId}`;
      setQrValue(url);
    } catch (err) {
      setError(
        err?.response?.data?.message ||
        err.message ||
        "Failed to generate QR."
      );
    } finally {
      setBusy(false);
    }
  };

  // Optionally auto-generate once
  useEffect(() => {
    handleGenerate();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [sectionId]);

  return (
    <div className="page-center">
      <div className="card">
        <h2>QR Code for Attendance</h2>
        <p className="subtitle">
          {course.courseName} ({course.code}) - Room {course.room}
        </p>

        <FormError message={error} />

        {qrValue ? (
          <div className="qr-wrapper">
            <QRCode value={qrValue} size={220} />
            <p className="muted small" style={{ wordBreak: "break-all", marginTop: "1rem" }}>
              Scan to visit: <a href={qrValue} target="_blank" rel="noreferrer">{qrValue}</a>
            </p>
            <p className="muted small">
              Students will scan this code in class.
            </p>
          </div>
        ) : (
          <p className="muted">No QR generated yet.</p>
        )}

        <Button onClick={handleGenerate} disabled={busy}>
          {busy ? "Generating..." : "Regenerate QR"}
        </Button>
      </div>
    </div>
  );
}
