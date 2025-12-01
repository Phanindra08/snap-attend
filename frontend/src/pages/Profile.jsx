import { useState } from "react";
import Input from "../components/Input";
import Button from "../components/Button";
import FormError from "../components/FormError";
import { updateProfile } from "../api/auth";
import { useAuth } from "../auth/AuthContext";

export default function Profile() {
  const { user } = useAuth() || {};
  const [form, setForm] = useState({
    firstName: user?.firstName || "",
    lastName: user?.lastName || "",
    oldPassword: "",
    newPassword: "",
    confirmNewPassword: "",
    address1: "",
    address2: "",
    city: "",
    state: "",
    zipCode: "",
    country: "",
  });

  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  const [busy, setBusy] = useState(false);

  const validate = () => {
    if (!form.oldPassword) return "Old password is required.";
    if (form.newPassword.length < 8)
      return "New password must be at least 8 characters long.";
    if (form.newPassword === form.oldPassword)
      return "New password must be different from old password.";
    if (form.newPassword !== form.confirmNewPassword)
      return "New password and confirm password must match.";
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
      await updateProfile({
        firstName: form.firstName,
        lastName: form.lastName,
        oldPassword: form.oldPassword,
        newPassword: form.newPassword,
        address1: form.address1,
        address2: form.address2,
        city: form.city,
        state: form.state,
        zipCode: form.zipCode,
        country: form.country,
      });
      setSuccess("Profile updated successfully.");
    } catch (err) {
      setError(
        err?.response?.data?.message ||
          err.message ||
          "Failed to update profile."
      );
    } finally {
      setBusy(false);
    }
  };

  return (
    <div className="page-center">
      <div className="card">
        <h2>Update Profile</h2>
        <form onSubmit={handleSubmit} className="form form-grid">
          <Input
            label="First Name"
            value={form.firstName}
            onChange={(e) =>
              setForm((f) => ({ ...f, firstName: e.target.value }))
            }
          />
          <Input
            label="Last Name"
            value={form.lastName}
            onChange={(e) =>
              setForm((f) => ({ ...f, lastName: e.target.value }))
            }
          />
          <Input
            label="Old Password"
            type="password"
            value={form.oldPassword}
            onChange={(e) =>
              setForm((f) => ({ ...f, oldPassword: e.target.value }))
            }
          />
          <Input
            label="New Password (min 8 chars)"
            type="password"
            value={form.newPassword}
            onChange={(e) =>
              setForm((f) => ({ ...f, newPassword: e.target.value }))
            }
          />
          <Input
            label="Confirm New Password"
            type="password"
            value={form.confirmNewPassword}
            onChange={(e) =>
              setForm((f) => ({
                ...f,
                confirmNewPassword: e.target.value,
              }))
            }
          />
          <Input
            label="Address 1"
            value={form.address1}
            onChange={(e) =>
              setForm((f) => ({ ...f, address1: e.target.value }))
            }
          />
          <Input
            label="Address 2"
            value={form.address2}
            onChange={(e) =>
              setForm((f) => ({ ...f, address2: e.target.value }))
            }
          />
          <Input
            label="City"
            value={form.city}
            onChange={(e) =>
              setForm((f) => ({ ...f, city: e.target.value }))
            }
          />
          <Input
            label="State"
            value={form.state}
            onChange={(e) =>
              setForm((f) => ({ ...f, state: e.target.value }))
            }
          />
          <Input
            label="Zip Code"
            value={form.zipCode}
            onChange={(e) =>
              setForm((f) => ({ ...f, zipCode: e.target.value }))
            }
          />
          <Input
            label="Country"
            value={form.country}
            onChange={(e) =>
              setForm((f) => ({ ...f, country: e.target.value }))
            }
          />
        </form>

        <FormError message={error} />
        {success && <div className="form-success">{success}</div>}

        <Button onClick={handleSubmit} disabled={busy}>
          {busy ? "Saving..." : "Save changes"}
        </Button>
      </div>
    </div>
  );
}
