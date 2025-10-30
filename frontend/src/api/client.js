import { API_BASE_URL } from "../config";

const BASE = API_BASE_URL || "";

async function request(path, { method = "GET", body, headers = {}, credentials = "include" } = {}) {
  const res = await fetch(`${BASE}${path}`, {
    method,
    headers: { "Content-Type": "application/json", ...headers },
    credentials, // supports httpOnly cookie auth
    body: body ? JSON.stringify(body) : undefined
  });

  const ct = res.headers.get("content-type") || "";
  const data = ct.includes("application/json") ? await res.json() : await res.text();

  if (!res.ok) {
    throw new Error(typeof data === "string" ? data : data?.error || "Request failed");
  }
  return data;
}

export default { request };
