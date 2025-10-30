import { useState } from "react";
import Input from "../components/Input";
import Button from "../components/Button";
import FormError from "../components/FormError";
import { login } from "../api/auth";

export default function Login() {
  const [f, setF] = useState({ email: "", password: "" });
  const [err, setErr] = useState("");

  const onChange = (e) => setF({ ...f, [e.target.name]: e.target.value });

  const onSubmit = async (e) => {
    e.preventDefault();
    setErr("");
    try {
      await login(f.email, f.password);
      window.location.href = "/";
    } catch (e) {
      setErr(e.message);
    }
  };

  return (
    <div className="max-w-md mx-auto mt-16 p-6 border rounded">
      <h1 className="text-2xl font-bold mb-4">Login to SnapAttend</h1>
      <form onSubmit={onSubmit}>
        <Input label="Email" name="email" type="email" value={f.email} onChange={onChange} required />
        <Input label="Password" name="password" type="password" value={f.password} onChange={onChange} required />
        <FormError message={err} />
        <Button type="submit" className="mt-4">Login</Button>
      </form>
      <p className="mt-4 text-sm">
        Don’t have an account? <a className="underline" href="/signup">Sign up</a>
      </p>
    </div>
  );
}
