const SESSION_ENTRY = "users";
const SESSION_TOKEN = "token";

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
		window.sessionStorage.setItem(SESSION_ENTRY, JSON.stringify(u));
	}

	static getUser(): User | null {
		const b = window.sessionStorage.getItem(SESSION_ENTRY);
		if (!b) return null;

		return JSON.parse(b) as User;
	}

	static saveToken(token: string) {
		window.sessionStorage.setItem(SESSION_TOKEN, token);
	}

	static GET(path: string, params?: Record<string, string>) {
		if (path[0] != "/") path = "/" + path;
		const user = this.getUser();

		let url = "https://api.github.com" + path;
		if (params) {
			const p = new URLSearchParams(params);
			url += "?" + p.toString();
		}

		return fetch(url, {
			method: "GET",
			headers: {
				Accept: "application/vnd.github+json",
				Authorization: `Bearer ${user?.payload.token}`,
				"X-GitHub-Api-Version": "2022-11-28",
			},
		});
	}
}
