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
  private static Session = "users";

  static saveUser(u: User) {
    window.sessionStorage.setItem(this.Session, JSON.stringify(u));
  }

  static getUser(): User | null {
    const b = window.sessionStorage.getItem(this.Session);
    if (!b) return null;

    return JSON.parse(b) as User;
  }
}
