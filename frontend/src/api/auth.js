import client from "./client";

// POST /api/auth/sign-up
export async function signup(payload) {
  const res = await client.post("/auth/sign-up", payload);
  return res.data;
}

// POST /api/auth/login
export async function login({ email, password, role }) {
  const res = await client.post("/auth/login", {
    email,
    password,
    profile: role, // backend expects "profile"
  });
  return res.data; // expect { token, user: {...} } or similar
}

// PUT /api/user/update-profile
export async function updateProfile(payload) {
  const res = await client.put("/user/update-profile", payload);
  return res.data;
}
