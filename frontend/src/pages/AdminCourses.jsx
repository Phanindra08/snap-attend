import { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import Input from "../components/Input";
import Button from "../components/Button";
import {
    getAllCourses,
    createCourse,
    updateCourse,
    deleteCourse,
    getAllStudents,
    getAllProfessors,
    enrollStudent,
    assignProfessor,
    createSection,
    getAllSections,
    updateSection
} from "../api/admin";

export default function AdminCourses() {
    const [courses, setCourses] = useState([]);
    const [sections, setSections] = useState([]);
    const [students, setStudents] = useState([]);
    const [professors, setProfessors] = useState([]);
    const [loading, setLoading] = useState(true);
    const [form, setForm] = useState({
        courseName: "",
    });
    const [editingId, setEditingId] = useState(null);
    const [message, setMessage] = useState("");
    const [error, setError] = useState("");
    const [searchTerm, setSearchTerm] = useState("");

    // Linking state
    const [showEnrollForm, setShowEnrollForm] = useState(false);
    const [showAssignForm, setShowAssignForm] = useState(false);
    const [showCreateSectionForm, setShowCreateSectionForm] = useState(false);
    const [showEditSectionForm, setShowEditSectionForm] = useState(false);
    const [enrollForm, setEnrollForm] = useState({ studentId: "", sectionId: "" });
    const [assignForm, setAssignForm] = useState({ professorId: "", sectionId: "" });
    const [createSectionForm, setCreateSectionForm] = useState({
        courseId: "",
        sectionNumber: "001",
        semesterId: "1",
        roomId: "1",
        latitude: "",
        longitude: ""
    });
    const [editSectionForm, setEditSectionForm] = useState({
        id: "",
        sectionNumber: "",
        latitude: "",
        longitude: "",
        roomId: ""
    });

    useEffect(() => {
        loadData();
    }, []);

    const loadData = async () => {
        try {
            const [coursesData, studentsData, professorsData, sectionsData] = await Promise.all([
                getAllCourses(),
                getAllStudents(),
                getAllProfessors(),
                getAllSections()
            ]);
            setCourses(coursesData);
            setStudents(studentsData);
            setProfessors(professorsData);
            setSections(sectionsData);
            setLoading(false);
        } catch (err) {
            console.error(err);
            setError("Failed to load data");
            setLoading(false);
        }
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setMessage("");
        setError("");

        try {
            if (editingId) {
                await updateCourse(editingId, form);
                setMessage("Course updated successfully!");
                setEditingId(null);
            } else {
                await createCourse(form);
                setMessage("Course created successfully!");
            }
            setForm({ courseName: "" });
            loadData();
        } catch (err) {
            setError("Failed to save course. " + (err.response?.data?.error || ""));
            console.error(err);
        }
    };

    const handleEdit = (course) => {
        setForm({
            courseName: course.courseName,
        });
        setEditingId(course.id);
        window.scrollTo({ top: 0, behavior: "smooth" });
    };

    const handleCancelEdit = () => {
        setEditingId(null);
        setForm({ courseName: "" });
    };

    const handleDelete = async (id) => {
        if (!window.confirm("Are you sure you want to delete this course?")) return;

        try {
            await deleteCourse(id);
            setMessage("Course deleted successfully!");
            loadData();
        } catch (err) {
            setError("Failed to delete course. " + (err.response?.data?.error || ""));
            console.error(err);
        }
    };

    const handleCreateSection = async (e) => {
        e.preventDefault();
        setMessage("");
        setError("");

        try {
            const sectionData = {
                courseId: parseInt(createSectionForm.courseId),
                sectionNumber: createSectionForm.sectionNumber,
                semesterId: parseInt(createSectionForm.semesterId),
                roomId: parseInt(createSectionForm.roomId),
                latitude: parseFloat(createSectionForm.latitude),
                longitude: parseFloat(createSectionForm.longitude)
            };

            const res = await createSection(sectionData);
            setMessage(`Section created successfully! ID: ${res.section.id}`);
            setCreateSectionForm({
                courseId: "",
                sectionNumber: "001",
                semesterId: "1",
                roomId: "1",
                latitude: "",
                longitude: ""
            });
            setShowCreateSectionForm(false);
        } catch (err) {
            setError("Failed to create section. " + (err.response?.data?.error || ""));
            console.error(err);
        }
    };

    const handleEnrollStudent = async (e) => {
        e.preventDefault();
        setMessage("");
        setError("");

        try {
            await enrollStudent(parseInt(enrollForm.studentId), parseInt(enrollForm.sectionId));
            setMessage("Student enrolled successfully!");
            setEnrollForm({ studentId: "", sectionId: "" });
            setShowEnrollForm(false);
        } catch (err) {
            setError("Failed to enroll student. " + (err.response?.data?.error || ""));
            console.error(err);
        }
    };

    const handleAssignProfessor = async (e) => {
        e.preventDefault();
        setMessage("");
        setError("");

        try {
            await assignProfessor(parseInt(assignForm.professorId), parseInt(assignForm.sectionId));
            setMessage("Professor assigned successfully!");
            setAssignForm({ professorId: "", sectionId: "" });
            setShowAssignForm(false);
        } catch (err) {
            setError("Failed to assign professor. " + (err.response?.data?.error || ""));
            console.error(err);
        }
    };

    const handleEditSection = (section) => {
        setEditSectionForm({
            id: section.id,
            sectionNumber: section.sectionNumber,
            latitude: section.latitude.toString(),
            longitude: section.longitude.toString(),
            roomId: section.roomId.toString()
        });
        setShowEditSectionForm(true);
        setShowCreateSectionForm(false);
        setShowEnrollForm(false);
        setShowAssignForm(false);
        window.scrollTo({ top: 0, behavior: "smooth" });
    };

    const handleUpdateSection = async (e) => {
        e.preventDefault();
        setMessage("");
        setError("");

        try {
            const updateData = {
                sectionNumber: editSectionForm.sectionNumber,
                latitude: parseFloat(editSectionForm.latitude),
                longitude: parseFloat(editSectionForm.longitude),
                roomId: parseInt(editSectionForm.roomId)
            };

            await updateSection(editSectionForm.id, updateData);
            setMessage("Section updated successfully!");
            setEditSectionForm({ id: "", sectionNumber: "", latitude: "", longitude: "", roomId: "" });
            setShowEditSectionForm(false);
            loadData();
        } catch (err) {
            setError("Failed to update section. " + (err.response?.data?.error || ""));
            console.error(err);
        }
    };

    const filteredCourses = courses.filter(c =>
        c.courseName.toLowerCase().includes(searchTerm.toLowerCase())
    );

    if (loading) return <div className="page p-4">Loading...</div>;

    return (
        <div className="page">
            <header className="topbar">
                <h1>Manage Courses</h1>
                <Link to="/admin/dashboard" className="link">
                    Back to Dashboard
                </Link>
            </header>

            {message && <div className="card success">{message}</div>}
            {error && <div className="card error">{error}</div>}

            <div className="card">
                <h2>{editingId ? "Edit Course" : "Create New Course"}</h2>
                <form onSubmit={handleSubmit} className="form">
                    <Input
                        label="Course Name"
                        value={form.courseName}
                        onChange={(e) => setForm({ ...form, courseName: e.target.value })}
                        placeholder="e.g. Introduction to Computer Science"
                        required
                    />

                    <div className="flex gap-2">
                        <Button type="submit">{editingId ? "Update" : "Create"} Course</Button>
                        {editingId && (
                            <Button type="button" onClick={handleCancelEdit}>
                                Cancel Edit
                            </Button>
                        )}
                    </div>
                </form>
            </div>

            <div className="card">
                <h2>All Courses ({courses.length})</h2>
                <Input
                    label="Search"
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                    placeholder="Search by course name..."
                />

                {filteredCourses.length === 0 ? (
                    <p className="muted">No courses found.</p>
                ) : (
                    <div className="table-container">
                        <table className="data-table">
                            <thead>
                                <tr>
                                    <th>ID</th>
                                    <th>Course Name</th>
                                    <th>Actions</th>
                                </tr>
                            </thead>
                            <tbody>
                                {filteredCourses.map((course) => (
                                    <tr key={course.id}>
                                        <td>{course.id}</td>
                                        <td>{course.courseName}</td>
                                        <td>
                                            <div className="flex gap-2">
                                                <button
                                                    className="btn btn-sm"
                                                    onClick={() => handleEdit(course)}
                                                >
                                                    Edit
                                                </button>
                                                <button
                                                    className="btn btn-sm btn-danger"
                                                    onClick={() => handleDelete(course.id)}
                                                >
                                                    Delete
                                                </button>
                                            </div>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                )}
            </div>

            {showEditSectionForm && (
                <div className="card">
                    <h2>Edit Section</h2>
                    <form onSubmit={handleUpdateSection} className="form">
                        <Input
                            label="Section Number"
                            value={editSectionForm.sectionNumber}
                            onChange={(e) => setEditSectionForm({ ...editSectionForm, sectionNumber: e.target.value })}
                            required
                        />
                        <div className="grid-2">
                            <Input
                                label="Latitude"
                                type="number"
                                step="any"
                                value={editSectionForm.latitude}
                                onChange={(e) => setEditSectionForm({ ...editSectionForm, latitude: e.target.value })}
                                required
                            />
                            <Input
                                label="Longitude"
                                type="number"
                                step="any"
                                value={editSectionForm.longitude}
                                onChange={(e) => setEditSectionForm({ ...editSectionForm, longitude: e.target.value })}
                                required
                            />
                        </div>
                        <Input
                            label="Room ID"
                            type="number"
                            value={editSectionForm.roomId}
                            onChange={(e) => setEditSectionForm({ ...editSectionForm, roomId: e.target.value })}
                            required
                        />
                        <div className="flex gap-2">
                            <Button type="submit">Update Section</Button>
                            <Button type="button" onClick={() => setShowEditSectionForm(false)}>Cancel</Button>
                        </div>
                    </form>
                </div>
            )}

            <div className="card">
                <h2>All Sections ({sections.length})</h2>
                {sections.length === 0 ? (
                    <p className="muted">No sections found.</p>
                ) : (
                    <div className="table-container">
                        <table className="data-table">
                            <thead>
                                <tr>
                                    <th>ID</th>
                                    <th>Course</th>
                                    <th>Section #</th>
                                    <th>Professor ID</th>
                                    <th>Coordinates</th>
                                    <th>Actions</th>
                                </tr>
                            </thead>
                            <tbody>
                                {sections.map((section) => (
                                    <tr key={section.id}>
                                        <td>{section.id}</td>
                                        <td>{section.courseName}</td>
                                        <td>{section.sectionNumber}</td>
                                        <td>{section.professorId || "Unassigned"}</td>
                                        <td>{section.latitude}, {section.longitude}</td>
                                        <td>
                                            <button
                                                className="btn btn-sm"
                                                onClick={() => handleEditSection(section)}
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
            <div className="card">
                <h2>Linking</h2>
                <p className="subtitle">Create sections, enroll students, or assign professors</p>

                <div className="flex gap-2 mb-4 flex-wrap">
                    <Button onClick={() => {
                        setShowCreateSectionForm(!showCreateSectionForm);
                        setShowEnrollForm(false);
                        setShowAssignForm(false);
                    }}>
                        {showCreateSectionForm ? "Hide" : "Show"} Create Section Form
                    </Button>
                    <Button onClick={() => {
                        setShowEnrollForm(!showEnrollForm);
                        setShowCreateSectionForm(false);
                        setShowAssignForm(false);
                    }}>
                        {showEnrollForm ? "Hide" : "Show"} Enroll Student Form
                    </Button>
                    <Button onClick={() => {
                        setShowAssignForm(!showAssignForm);
                        setShowCreateSectionForm(false);
                        setShowEnrollForm(false);
                    }}>
                        {showAssignForm ? "Hide" : "Show"} Assign Professor Form
                    </Button>
                </div>

                {showCreateSectionForm && (
                    <form onSubmit={handleCreateSection} className="form">
                        <h3>Create New Course Section</h3>
                        <label className="input-label">
                            <span>Course</span>
                            <select
                                className="input-control"
                                value={createSectionForm.courseId}
                                onChange={(e) => setCreateSectionForm({ ...createSectionForm, courseId: e.target.value })}
                                required
                            >
                                <option value="">Select a course</option>
                                {courses.map(c => (
                                    <option key={c.id} value={c.id}>
                                        {c.courseName}
                                    </option>
                                ))}
                            </select>
                        </label>
                        <Input
                            label="Section Number"
                            value={createSectionForm.sectionNumber}
                            onChange={(e) => setCreateSectionForm({ ...createSectionForm, sectionNumber: e.target.value })}
                            placeholder="e.g. 001"
                            required
                        />
                        <div className="grid-2">
                            <Input
                                label="Semester ID"
                                type="number"
                                value={createSectionForm.semesterId}
                                onChange={(e) => setCreateSectionForm({ ...createSectionForm, semesterId: e.target.value })}
                                required
                            />
                            <Input
                                label="Room ID"
                                type="number"
                                value={createSectionForm.roomId}
                                onChange={(e) => setCreateSectionForm({ ...createSectionForm, roomId: e.target.value })}
                                required
                            />
                        </div>
                        <div className="grid-2">
                            <Input
                                label="Latitude"
                                type="number"
                                step="any"
                                value={createSectionForm.latitude}
                                onChange={(e) => setCreateSectionForm({ ...createSectionForm, latitude: e.target.value })}
                                placeholder="e.g. 40.7128"
                                required
                            />
                            <Input
                                label="Longitude"
                                type="number"
                                step="any"
                                value={createSectionForm.longitude}
                                onChange={(e) => setCreateSectionForm({ ...createSectionForm, longitude: e.target.value })}
                                placeholder="e.g. -74.0060"
                                required
                            />
                        </div>
                        <p className="small muted">Note: Default Semester ID is 1 (Fall 2025), Default Room ID is 1 (101)</p>
                        <Button type="submit">Create Section</Button>
                    </form>
                )}

                {showEnrollForm && (
                    <form onSubmit={handleEnrollStudent} className="form">
                        <h3>Enroll Student in Section</h3>
                        <label className="input-label">
                            <span>Student</span>
                            <select
                                className="input-control"
                                value={enrollForm.studentId}
                                onChange={(e) => setEnrollForm({ ...enrollForm, studentId: e.target.value })}
                                required
                            >
                                <option value="">Select a student</option>
                                {students.map(s => (
                                    <option key={s.id} value={s.id}>
                                        {s.firstName} {s.lastName} ({s.email})
                                    </option>
                                ))}
                            </select>
                        </label>
                        <Input
                            label="Section ID"
                            type="number"
                            value={enrollForm.sectionId}
                            onChange={(e) => setEnrollForm({ ...enrollForm, sectionId: e.target.value })}
                            placeholder="Enter section ID"
                            required
                        />
                        <Button type="submit">Enroll Student</Button>
                    </form>
                )}

                {showAssignForm && (
                    <form onSubmit={handleAssignProfessor} className="form">
                        <h3>Assign Professor to Section</h3>
                        <label className="input-label">
                            <span>Professor</span>
                            <select
                                className="input-control"
                                value={assignForm.professorId}
                                onChange={(e) => setAssignForm({ ...assignForm, professorId: e.target.value })}
                                required
                            >
                                <option value="">Select a professor</option>
                                {professors.map(p => (
                                    <option key={p.id} value={p.id}>
                                        {p.firstName} {p.lastName} ({p.email})
                                    </option>
                                ))}
                            </select>
                        </label>
                        <Input
                            label="Section ID"
                            type="number"
                            value={assignForm.sectionId}
                            onChange={(e) => setAssignForm({ ...assignForm, sectionId: e.target.value })}
                            placeholder="Enter section ID"
                            required
                        />
                        <Button type="submit">Assign Professor</Button>
                    </form>
                )}
            </div>
        </div>
    );
}
