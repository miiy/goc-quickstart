import { Configuration } from "./generated/runtime";

const DEVELOPMENT_API_URL = "http://localhost:8080/api/v1";

function apiBaseURL() {
  const fromEnv = import.meta.env.VITE_API_URL?.trim();
  const baseURL = fromEnv || (import.meta.env.DEV ? DEVELOPMENT_API_URL : "");
  if (!baseURL) {
    throw new Error("VITE_API_URL is required for production builds");
  }
  return baseURL.replace(/\/+$/, "");
}

function apiConfiguration() {
  return new Configuration({
    basePath: apiBaseURL(),
    credentials: "include"
  });
}

export { apiBaseURL, apiConfiguration };
