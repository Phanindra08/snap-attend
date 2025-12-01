import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import Input from "../components/Input";
import Button from "../components/Button";
import FormError from "../components/FormError";
import { login } from "../api/auth";
import { useAuth } from "../auth/AuthContext";
import { TOKEN_KEY } from "../config";

export default function Login() {
  const navigate = useNavigate();
  const { setUser } = useAuth() || { setUser: () => { } };

  const [form, setForm] = useState({
    email: "",
    password: "",
    role: "Student",
  });
  const [error, setError] = useState("");
  const [busy, setBusy] = useState(false);

  const validate = () => {
    if (!form.email) return "Email is required.";
    if (form.password.length < 8)
      return "Password must be at least 8 characters long.";
    return "";
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    const msg = validate();
    if (msg) {
      setError(msg);
      return;
    }
    setError("");
    setBusy(true);

    try {
      const res = await login(form);
      const token = res?.token;
      const user = res?.user || {
        profile: form.role,
        email: form.email,
      };

      if (token) {
        localStorage.setItem(TOKEN_KEY, token);
      }
      if (setUser) {
        setUser(user);
      }

      if (user.profile === "Student") {
        navigate("/student/dashboard");
      } else if (user.profile === "Professor") {
        navigate("/professor/dashboard");
      } else if (user.profile === "Admin") {
        navigate("/admin/dashboard");
      } else {
        navigate("/"); // unknown
      }
    } catch (err) {
      setError(err?.response?.data?.message || err.message || "Login failed.");
    } finally {
      setBusy(false);
    }
  };

  return (
    <div className="page-center">
      <div className="card">
        <h1>SnapAttend</h1>
        <p className="subtitle">
          Welcome back. Sign in to mark attendance and view your classes.
        </p>

        <form onSubmit={handleSubmit} className="form">
          <Input
            label="Email"
            type="email"
            value={form.email}
            onChange={(e) => setForm({ ...form, email: e.target.value })}
          />

          <label className="input-label">
            <span>Role</span>
            <select
              className="input-control"
              value={form.role}
              onChange={(e) => setForm({ ...form, role: e.target.value })}
            >
              <option value="Student">Student</option>
              <option value="Professor">Professor</option>
              <option value="Admin">Admin</option>
            </select>
          </label>

          <Input
            label="Password (min 8 chars)"
            type="password"
            value={form.password}
            onChange={(e) => setForm({ ...form, password: e.target.value })}
          />

          <FormError message={error} />

          <Button type="submit" disabled={busy}>
            {busy ? "Logging in..." : "Login"}
          </Button>
        </form>

        <p className="muted">
          Don&apos;t have an account? <Link to="/signup">Sign up</Link>
        </p>
      </div>
    </div>
  );
}
