enum WEB {
  home = "/",
  dash = "/dash",
  auth = "/auth",
  repo = "/repo"
}

const API = import.meta.env.VITE_ENVHUB_API;
if (!API) throw new Error("`VITE_ENVHUB_API` not set");

export { API, WEB };
