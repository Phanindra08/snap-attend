import { useState } from "react";
import Input from "../components/Input";
import Button from "../components/Button";
import FormError from "../components/FormError";
import { register } from "../api/auth";

export default function Signup() {
  const [f, setF] = useState({
    first_name: "",
    last_name: "",
    email: "",
    password: "",
    role: "Student"
  });
  const [err, setErr] = useState("");

  const onChange = (e) => setF({ ...f, [e.target.name]: e.target.value });

  const onSubmit = async (e) => {
    e.preventDefault();
    setErr("");
    try {
      await register(f);
      window.location.href = "/login";
    } catch (e) {
      setErr(e.message);
    }
  };

  return (
    <div className="max-w-md mx-auto mt-16 p-6 border rounded">
      <h1 className="text-2xl font-bold mb-4">Create your account</h1>
      <form onSubmit={onSubmit}>
        <Input label="First Name" name="first_name" value={f.first_name} onChange={onChange} required />
        <Input label="Last Name" name="last_name" value={f.last_name} onChange={onChange} required />
        <Input label="Email" name="email" type="email" value={f.email} onChange={onChange} required />
        <Input label="Password" name="password" type="password" value={f.password} onChange={onChange} required />

        <label className="block mb-3">
          <span className="block mb-1 text-sm font-medium">Role</span>
          <select
            name="role"
            value={f.role}
            onChange={onChange}
            className="w-full border rounded px-3 py-2"
          >
            <option>Student</option>
            <option>Professor</option>
            <option>Admin</option>
          </select>
        </label>

        <FormError message={err} />
        <Button type="submit" className="mt-2">Sign Up</Button>
      </form>
      <p className="mt-4 text-sm">
        Already have an account? <a className="underline" href="/login">Login</a>
      </p>
    </div>
  );
}
