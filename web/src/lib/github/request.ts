import { Session } from "../sessionStorage.ts";

export const GITHUB_SECRET = import.meta.env.VITE_GITHUB_CLIENT_ID;
if (!GITHUB_SECRET) throw new Error("`GITHUB_SECRET` not set");

export interface User {
  payload: Payload;
  protectedHeader: ProtectedHeader;
}

export interface Payload {
  token: string;
  id: number;
  login: string;
  avatar_url: string;
  name: string;
  email: string;
  iss: string;
  sub: string;
}

export interface ProtectedHeader {
  alg: string;
  typ: string;
}

export class Github {
  static saveUser(u: User) {
    window.sessionStorage.setItem(Session.user, JSON.stringify(u));
  }

  static getUser(): User | null {
    const b = window.sessionStorage.getItem(Session.user);
    if (!b) return null;

    return JSON.parse(b) as User;
  }

  static saveToken(token: string) {
    window.sessionStorage.setItem(Session.token, token);
  }

  static POST(
    path: string,
    param: Record<string, string> | null,
    body: object | null,
  ): Promise<Response> {
    return fetch(this.url(path, param), {
      method: "POST",
      headers: this.headers(),
      body: JSON.stringify(body)
    })

  }

  static GET(
    path: string,
    params: Record<string, string> | null,
  ): Promise<Response> {
    return fetch(this.url(path, params), {
      method: "GET",
      headers: this.headers()
    });
  }

  private static url(
    path: string,
    params: Record<string, string> | null,
  ): string {
    if (path[0] != "/") path = "/" + path;
    let url = "https://api.github.com" + path;
    if (params) {
      const p = new URLSearchParams(params);
      url += "?" + p.toString();
    }

    return url
  }

  private static headers(h?: Record<string, string>): Headers {
    const user = this.getUser();
    if (!user) throw new Error("user not found")
    return new Headers({
      Accept: "application/vnd.github+json",
      Authorization: `Bearer ${user.payload.token}`,
      "X-GitHub-Api-Version": "2022-11-28",
      ...h
    })
  }
}
