import { Session } from "./sessionStorage";

const API = import.meta.env.VITE_ENVHUB_API;
if (!API) throw new Error("`VITE_ENVHUB_API` not set");

export class Fetch {
  static async GET(
    path: string,
    query?: Record<string, string>,
  ): Promise<Response> {
    if (path[0] != "/") path = "/" + path;
    let url = API + path;
    if (query) {
      url += "?" + new URLSearchParams(query).toString();
    }
    const headers = new Headers({
      Authorization: `Bearer ${this.getToken()}`,
      "Content-Type": "application/json",
    });

    return fetch(url, { method: "GET", headers });
  }

  static POST(
    path: string,
    query: Record<string, string> | null,
    body: object | null
  ): Promise<Response> {
    return fetch(this.url(path, query))

  }


  private static url(path: string, query: Record<string, string> | null): string {
    if (path[0] != "/") path = "/" + path;
    let url = API + path;
    if (query) {
      url += "?" + new URLSearchParams(query).toString();
    }

    return url
  }

  private static getToken(): string {
    const t = window.sessionStorage.getItem(Session.token);
    if (!t) throw new Error("session token not found");
    return t;
  }
}
