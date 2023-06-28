export const GITHUB_SECRET = import.meta.env.VITE_GITHUB_CLIENT_ID;
if (!GITHUB_SECRET) throw new Error("`GITHUB_SECRET` not set");
