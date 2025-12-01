import { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import Input from "../components/Input";
import Button from "../components/Button";
import { getAllProfessors, createProfessor, updateProfessor } from "../api/admin";

export default function AdminProfessors() {
    const [professors, setProfessors] = useState([]);
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
        loadProfessors();
    }, []);

    const loadProfessors = async () => {
        try {
            const data = await getAllProfessors();
            setProfessors(data);
            setLoading(false);
        } catch (err) {
            console.error(err);
            setError("Failed to load professors");
            setLoading(false);
        }
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setMessage("");
        setError("");

        try {
            if (editingId) {
                await updateProfessor(editingId, form);
                setMessage("Professor updated successfully!");
                setEditingId(null);
            } else {
                await createProfessor(form);
                setMessage("Professor created successfully!");
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
            loadProfessors();
        } catch (err) {
            setError("Failed to save professor. " + (err.response?.data?.error || ""));
            console.error(err);
        }
    };

    const handleEdit = (professor) => {
        setForm({
            firstName: professor.firstName,
            lastName: professor.lastName,
            email: professor.email,
            password: "",
            address1: "123 Campus Dr",
            city: "Charlotte",
            state: "NC",
            zipCode: "28262",
            country: "USA",
        });
        setEditingId(professor.id);
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

    const filteredProfessors = professors.filter(p =>
        `${p.firstName} ${p.lastName} ${p.email}`.toLowerCase().includes(searchTerm.toLowerCase())
    );

    if (loading) return <div className="page p-4">Loading...</div>;

    return (
        <div className="page">
            <header className="topbar">
                <h1>Manage Professors</h1>
                <Link to="/admin/dashboard" className="link">
                    Back to Dashboard
                </Link>
            </header>

            <div className="card">
                <h2>{editingId ? "Edit Professor" : "Create New Professor"}</h2>
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
                        <Button type="submit">{editingId ? "Update" : "Create"} Professor</Button>
                        {editingId && (
                            <Button type="button" onClick={handleCancelEdit}>
                                Cancel Edit
                            </Button>
                        )}
                    </div>
                </form>
            </div>

            <div className="card">
                <h2>All Professors ({professors.length})</h2>
                <Input
                    label="Search"
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                    placeholder="Search by name or email..."
                />

                {filteredProfessors.length === 0 ? (
                    <p className="muted">No professors found.</p>
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
                                {filteredProfessors.map((professor) => (
                                    <tr key={professor.id}>
                                        <td>{professor.id}</td>
                                        <td>{professor.firstName} {professor.lastName}</td>
                                        <td>{professor.email}</td>
                                        <td>
                                            <button
                                                className="btn btn-sm"
                                                onClick={() => handleEdit(professor)}
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
