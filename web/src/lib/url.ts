const API_URL = import.meta.env.VITE_ENVHUB_API;
if (!API_URL) throw new Error("`VITE_ENVHUB_API` not set");

export const GITHUB_CLIENT_ID = import.meta.env.VITE_GITHUB_CLIENT_ID;
if (!GITHUB_CLIENT_ID) throw new Error("`VITE_GITHUB_CLIENT_ID` not set");

export function makeUrl(path: string, queries?: object): string {
  if (path[0] !== "/") path = "/" + path;
  let url = API_URL + path;

  if (queries) {
    const p = new URLSearchParams();
    for (const [k, v] of Object.entries(queries)) {
      p.set(k, String(v));
    }
    url += "?";
    url += p.toString();
  }

  return url;
}
export function oauthHref(): string {
  const redirect_uri: string = url(window.location.host+"/auth");

  const client_id: string = GITHUB_CLIENT_ID;
  const scope: string = "repo user read:org";
  const p = new URLSearchParams({ client_id, redirect_uri, scope });

  return "https://github.com/login/oauth/authorize?" + p.toString();
}

function url(r: string): string {
  const mode = import.meta.env.NODE_ENV;
  const production = mode === "production";
  if (!production) {
    console.warn("NODE_ENV:", mode);
    return "http://" + r;
  }
  return "https://" + r;
}
