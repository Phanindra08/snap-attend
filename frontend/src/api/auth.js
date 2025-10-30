import client from "./client";

export const login = (email, password) =>
  client.request("/auth/login", { method: "POST", body: { email, password } });

export const register = (payload) =>
  client.request("/auth/register", { method: "POST", body: payload });

export const me = () => client.request("/auth/me", { method: "GET" });

export const logout = () => client.request("/auth/logout", { method: "POST" });
