import { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import Input from "../components/Input";
import Button from "../components/Button";
import { getAllStudents, createStudent, updateStudent } from "../api/admin";

export default function AdminStudents() {
    const [students, setStudents] = useState([]);
    const [loading, setLoading] = useState(true);
    const [form, setForm] = useState({
        firstName: "",
        lastName: "",
        email: "",
        password: "",
        address1: "123 Campus Dr",
        city: "Charlotte",
        state: "NC",
        zipCode: "28262",
        country: "USA",
    });
    const [editingId, setEditingId] = useState(null);
    const [message, setMessage] = useState("");
    const [error, setError] = useState("");
    const [searchTerm, setSearchTerm] = useState("");

    useEffect(() => {
        loadStudents();
    }, []);

    const loadStudents = async () => {
        try {
            const data = await getAllStudents();
            setStudents(data);
            setLoading(false);
        } catch (err) {
            console.error(err);
            setError("Failed to load students");
            setLoading(false);
        }
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setMessage("");
        setError("");

        try {
            if (editingId) {
                await updateStudent(editingId, form);
                setMessage("Student updated successfully!");
                setEditingId(null);
            } else {
                await createStudent(form);
                setMessage("Student created successfully!");
            }
            setForm({
                firstName: "",
                lastName: "",
                email: "",
                password: "",
                address1: "123 Campus Dr",
                city: "Charlotte",
                state: "NC",
                zipCode: "28262",
                country: "USA",
            });
            loadStudents();
        } catch (err) {
            setError("Failed to save student. " + (err.response?.data?.error || ""));
            console.error(err);
        }
    };

    const handleEdit = (student) => {
        setForm({
            firstName: student.firstName,
            lastName: student.lastName,
            email: student.email,
            password: "",
            address1: "123 Campus Dr",
            city: "Charlotte",
            state: "NC",
            zipCode: "28262",
            country: "USA",
        });
        setEditingId(student.id);
        window.scrollTo({ top: 0, behavior: "smooth" });
    };

    const handleCancelEdit = () => {
        setEditingId(null);
        setForm({
            firstName: "",
            lastName: "",
            email: "",
            password: "",
            address1: "123 Campus Dr",
            city: "Charlotte",
            state: "NC",
            zipCode: "28262",
            country: "USA",
        });
    };

    const filteredStudents = students.filter(s =>
        `${s.firstName} ${s.lastName} ${s.email}`.toLowerCase().includes(searchTerm.toLowerCase())
    );

    if (loading) return <div className="page p-4">Loading...</div>;

    return (
        <div className="page">
            <header className="topbar">
                <h1>Manage Students</h1>
                <Link to="/admin/dashboard" className="link">
                    Back to Dashboard
                </Link>
            </header>

            <div className="card">
                <h2>{editingId ? "Edit Student" : "Create New Student"}</h2>
                <form onSubmit={handleSubmit} className="form">
                    <div className="grid-2">
                        <Input
                            label="First Name"
                            value={form.firstName}
                            onChange={(e) => setForm({ ...form, firstName: e.target.value })}
                            required
                        />
                        <Input
                            label="Last Name"
                            value={form.lastName}
                            onChange={(e) => setForm({ ...form, lastName: e.target.value })}
                            required
                        />
                    </div>
                    <Input
                        label="Email"
                        type="email"
                        value={form.email}
                        onChange={(e) => setForm({ ...form, email: e.target.value })}
                        required
                    />
                    <Input
                        label="Password"
                        type="password"
                        value={form.password}
                        onChange={(e) => setForm({ ...form, password: e.target.value })}
                        required={!editingId}
                        placeholder={editingId ? "Leave blank to keep current password" : ""}
                    />

                    {message && <div className="success">{message}</div>}
                    {error && <div className="error">{error}</div>}

                    <div className="flex gap-2">
                        <Button type="submit">{editingId ? "Update" : "Create"} Student</Button>
                        {editingId && (
                            <Button type="button" onClick={handleCancelEdit}>
                                Cancel Edit
                            </Button>
                        )}
                    </div>
                </form>
            </div>

            <div className="card">
                <h2>All Students ({students.length})</h2>
                <Input
                    label="Search"
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                    placeholder="Search by name or email..."
                />

                {filteredStudents.length === 0 ? (
                    <p className="muted">No students found.</p>
                ) : (
                    <div className="table-container">
                        <table className="data-table">
                            <thead>
                                <tr>
                                    <th>ID</th>
                                    <th>Name</th>
                                    <th>Email</th>
                                    <th>Actions</th>
                                </tr>
                            </thead>
                            <tbody>
                                {filteredStudents.map((student) => (
                                    <tr key={student.id}>
                                        <td>{student.id}</td>
                                        <td>{student.firstName} {student.lastName}</td>
                                        <td>{student.email}</td>
                                        <td>
                                            <button
                                                className="btn btn-sm"
                                                onClick={() => handleEdit(student)}
                                            >
                                                Edit
                                            </button>
                                        </td>
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
