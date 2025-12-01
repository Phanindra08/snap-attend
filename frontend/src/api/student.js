import client from "./client";

// POST /api/student/attendance
export async function markAttendance({
  qrId,
  studentName,
  latitude,
  longitude,
  question,
  answer,
}) {
  const res = await client.post("/student/attendance", {
    qrId,
    studentName,
    latitude,
    longitude,
    question: question || "",
    answer: answer || "",
  });
  return res.data;
}

// GET /api/student/attendance/summary?courseId=1
export async function getAttendanceSummary(courseId = 1) {
  const res = await client.get(
    `/student/attendance/summary?courseId=${courseId}`
  );
  return res.data;
}

// GET /api/student/courses
export async function getCourses() {
  const res = await client.get("/student/courses");
  return res.data?.courses || [];
}
