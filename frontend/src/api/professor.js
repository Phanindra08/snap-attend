import client from "./client";

// POST /api/professor/sections/:sectionId/attendance-qr
export async function generateAttendanceQr(sectionId) {
  const res = await client.post(
    `/professor/sections/${sectionId}/attendance-qr`,
    {}
  );
  return res.data; // whatever backend returns, we'll show as QR
}

export async function getSections() {
  const res = await client.get("/professor/sections");
  return res.data?.courses || [];
}
