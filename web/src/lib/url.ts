import { PUBLIC_ENVHUB_API } from "$env/static/public";

export function makeUrl(path: string, param?: Record<string, string>): string {
  let url: string = PUBLIC_ENVHUB_API;
  url += path;
  if (param) {
    const p = new URLSearchParams(param);
    url += "?";
    url += p.toString();
  }
  return url;
}

export function urlParam(k: string): string | null {
  const url = new URL(window.location.href);
  return url.searchParams.get(k);
}

export function urlRedirect(k: string = "/"): void {
  window.location.replace(k);
}