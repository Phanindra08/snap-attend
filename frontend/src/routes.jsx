import { Routes, Route } from "react-router-dom";
import Login from "./pages/Login";
import Signup from "./pages/Signup";
import Profile from "./pages/Profile";
import StudentDashboard from "./pages/StudentDashboard";
import StudentScan from "./pages/StudentScan";
import ProfessorDashboard from "./pages/ProfessorDashboard";
import QrGenerator from "./pages/QrGenerator";
import AdminDashboard from "./pages/AdminDashboard";
import AdminCourses from "./pages/AdminCourses";
import AdminStudents from "./pages/AdminStudents";
import AdminProfessors from "./pages/AdminProfessors";
import NotFound from "./pages/NotFound";
import ProfessorAttendance from "./pages/ProfessorAttendance";

export default function AppRoutes() {
  return (
    <Routes>
      <Route path="/" element={<Login />} />
      <Route path="/login" element={<Login />} />
      <Route path="/signup" element={<Signup />} />

      <Route path="/profile" element={<Profile />} />

      <Route path="/student/dashboard" element={<StudentDashboard />} />
      <Route path="/student/scan" element={<StudentScan />} />

      <Route path="/professor/dashboard" element={<ProfessorDashboard />} />
      <Route path="/professor/qr" element={<QrGenerator />} />
      <Route path="/professor/attendance" element={<ProfessorAttendance />} />

      <Route path="/admin/dashboard" element={<AdminDashboard />} />
      <Route path="/admin/courses" element={<AdminCourses />} />
      <Route path="/admin/students" element={<AdminStudents />} />
      <Route path="/admin/professors" element={<AdminProfessors />} />

      <Route path="*" element={<NotFound />} />
    </Routes>
  );
}
