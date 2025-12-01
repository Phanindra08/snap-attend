import axios from "axios";
import { API_BASE_URL, TOKEN_KEY } from "../config";

const client = axios.create({
  baseURL: API_BASE_URL,
});

// Attach JWT token automatically if present
client.interceptors.request.use((config) => {
  const token = localStorage.getItem(TOKEN_KEY);
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export default client;
