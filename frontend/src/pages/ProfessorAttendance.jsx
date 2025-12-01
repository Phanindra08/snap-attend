import { useState, useEffect } from "react";
import { useSearchParams, Link } from "react-router-dom";
import client from "../api/client";
import Button from "../components/Button";
import Input from "../components/Input";

export default function ProfessorAttendance() {
    const [searchParams] = useSearchParams();
    const sectionId = searchParams.get("sectionId");
    const [attendance, setAttendance] = useState([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState("");
    const [date, setDate] = useState(new Date().toISOString().split('T')[0]);

    useEffect(() => {
        if (sectionId && date) {
            loadAttendance();
        }
    }, [sectionId, date]);

    const loadAttendance = async () => {
        setLoading(true);
        setError("");
        try {
            const res = await client.get(`/professor/sections/${sectionId}/attendance/report?date=${date}`);
            setAttendance(res.data.students || []);
        } catch (err) {
            console.error(err);
            setError("Failed to load attendance report.");
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="page">
            <header className="topbar">
                <h1>Attendance Report</h1>
                <Link to="/professor/dashboard" className="link">
                    Back to Dashboard
                </Link>
            </header>

            <div className="card mb-4">
                <div className="flex gap-4 items-end">
                    <Input
                        label="Date"
                        type="date"
                        value={date}
                        onChange={(e) => setDate(e.target.value)}
                    />
                    <Button onClick={loadAttendance}>Refresh</Button>
                </div>
            </div>

            {error && <div className="card error">{error}</div>}

            <div className="card">
                <h2>Attendance for {date}</h2>
                {loading ? (
                    <p>Loading...</p>
                ) : attendance.length === 0 ? (
                    <p className="muted">No records found for this date.</p>
                ) : (
                    <div className="table-container">
                        <table className="data-table">
                            <thead>
                                <tr>
                                    <th>Student Name</th>
                                    <th>Email</th>
                                    <th>Status</th>
                                    <th>Time</th>
                                    <th>Question</th>
                                    <th>Answer</th>
                                </tr>
                            </thead>
                            <tbody>
                                {attendance.map((record) => (
                                    <tr key={record.studentId}>
                                        <td>{record.studentName || `${record.firstName} ${record.lastName}`}</td>
                                        <td>{record.email}</td>
                                        <td>
                                            <span className={`badge ${record.present ? "success" : "error"}`}>
                                                {record.present ? "Present" : "Absent"}
                                            </span>
                                        </td>
                                        <td>
                                            {record.attendedAt
                                                ? new Date(record.attendedAt).toLocaleTimeString()
                                                : "-"}
                                        </td>
                                        <td>{record.question || "-"}</td>
                                        <td>{record.answer || "-"}</td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                )}
            </div>
        </div>
    );
}
