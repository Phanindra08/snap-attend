import client from "./client";

// ========== COURSES ==========
export async function getAllCourses() {
    const res = await client.get("/admin/courses");
    return res.data?.courses || [];
}

export async function createCourse(data) {
    const res = await client.post("/admin/courses", data);
    return res.data;
}

export async function updateCourse(id, data) {
    const res = await client.put(`/admin/courses/${id}`, data);
    return res.data;
}

export async function deleteCourse(id) {
    const res = await client.delete(`/admin/courses/${id}`);
    return res.data;
}

// ========== STUDENTS ==========
export async function getAllStudents() {
    const res = await client.get("/admin/students");
    return res.data?.students || [];
}

export async function createStudent(data) {
    const res = await client.post("/admin/students", data);
    return res.data;
}

export async function getStudentById(id) {
    const res = await client.get(`/admin/students/${id}`);
    return res.data;
}

export async function updateStudent(id, data) {
    const res = await client.put(`/admin/students/${id}`, data);
    return res.data;
}

// ========== PROFESSORS ==========
export async function getAllProfessors() {
    const res = await client.get("/admin/professors");
    return res.data?.professors || [];
}

export async function createProfessor(data) {
    const res = await client.post("/admin/professors", data);
    return res.data;
}

export async function getProfessorById(id) {
    const res = await client.get(`/admin/professors/${id}`);
    return res.data;
}

export async function updateProfessor(id, data) {
    const res = await client.put(`/admin/professors/${id}`, data);
    return res.data;
}

export async function createSection(data) {
    const res = await client.post("/admin/sections", data);
    return res.data;
}

export async function getAllSections() {
    const res = await client.get("/admin/sections");
    return res.data?.sections || [];
}

export async function updateSection(id, data) {
    const res = await client.put(`/admin/sections/${id}`, data);
    return res.data;
}

// ========== LINKING ==========
export async function enrollStudent(studentId, sectionId) {
    const res = await client.post("/admin/enrollments", {
        studentId,
        sectionId,
    });
    return res.data;
}

export async function assignProfessor(professorId, sectionId) {
    const res = await client.post("/admin/sections/assign-professor", {
        professorId,
        sectionId,
    });
    return res.data;
}

// ========== SEARCH ==========
export async function search(params) {
    const res = await client.get("/admin/search", { params });
    return res.data;
}
