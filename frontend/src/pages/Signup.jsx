import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import Input from "../components/Input";
import Button from "../components/Button";
import FormError from "../components/FormError";
import { signup } from "../api/auth";

export default function Signup() {
  const navigate = useNavigate();
  const [form, setForm] = useState({
    firstName: "",
    lastName: "",
    email: "",
    password: "",
    confirmPassword: "",
    role: "Student",
    address1: "",
    address2: "",
    city: "",
    state: "",
    zipCode: "",
    country: "",
  });
  const [error, setError] = useState("");
  const [busy, setBusy] = useState(false);
  const [success, setSuccess] = useState("");

  const validate = () => {
    if (!form.firstName || !form.lastName) return "Name is required.";
    if (!form.email) return "Email is required.";
    if (form.password.length < 8)
      return "Password must be at least 8 characters long.";
    if (form.password !== form.confirmPassword)
      return "Password and Confirm Password must match.";
    return "";
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setSuccess("");
    const msg = validate();
    if (msg) {
      setError(msg);
      return;
    }

    setBusy(true);
    try {
      await signup({
        firstName: form.firstName,
        lastName: form.lastName,
        email: form.email,
        password: form.password,
        profile: form.role,
        address1: form.address1,
        address2: form.address2,
        city: form.city,
        state: form.state,
        zipCode: form.zipCode,
        country: form.country,
      });
      setSuccess("Signup successful! Please login.");
      setTimeout(() => navigate("/"), 1000);
    } catch (err) {
      setError(
        err?.response?.data?.message || err.message || "Signup failed."
      );
    } finally {
      setBusy(false);
    }
  };

  return (
    <div className="page-center">
      <div className="card">
        <h1>Create account</h1>
        <p className="subtitle">
          Sign up as a Student or Professor to use SnapAttend.
        </p>

        <form onSubmit={handleSubmit} className="form form-grid">
          <Input
            label="First Name"
            value={form.firstName}
            onChange={(e) => setForm({ ...form, firstName: e.target.value })}
          />
          <Input
            label="Last Name"
            value={form.lastName}
            onChange={(e) => setForm({ ...form, lastName: e.target.value })}
          />
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
            </select>
          </label>

          <Input
            label="Password (min 8 chars)"
            type="password"
            value={form.password}
            onChange={(e) => setForm({ ...form, password: e.target.value })}
          />
          <Input
            label="Confirm Password"
            type="password"
            value={form.confirmPassword}
            onChange={(e) =>
              setForm({ ...form, confirmPassword: e.target.value })
            }
          />

          <Input
            label="Address 1"
            value={form.address1}
            onChange={(e) => setForm({ ...form, address1: e.target.value })}
          />
          <Input
            label="Address 2"
            value={form.address2}
            onChange={(e) => setForm({ ...form, address2: e.target.value })}
          />
          <Input
            label="City"
            value={form.city}
            onChange={(e) => setForm({ ...form, city: e.target.value })}
          />
          <Input
            label="State"
            value={form.state}
            onChange={(e) => setForm({ ...form, state: e.target.value })}
          />
          <Input
            label="Zip Code"
            value={form.zipCode}
            onChange={(e) => setForm({ ...form, zipCode: e.target.value })}
          />
          <Input
            label="Country"
            value={form.country}
            onChange={(e) => setForm({ ...form, country: e.target.value })}
          />
        </form>

        <FormError message={error} />
        {success && <div className="form-success">{success}</div>}

        <Button onClick={handleSubmit} disabled={busy}>
          {busy ? "Signing up..." : "Sign up"}
        </Button>

        <p className="muted">
          Already have an account? <Link to="/">Login</Link>
        </p>
      </div>
    </div>
  );
}
