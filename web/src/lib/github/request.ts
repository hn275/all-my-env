import { useEffect, useState } from "react";

const SESSION_ENTRY = "users";

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
}

export function useGithubFetch<T>(path: string) {
	const [data, setData] = useState<T | null>();
	const [error, setError] = useState<string>();

	useEffect(() => {
		const session = window.sessionStorage.getItem(SESSION_ENTRY);
		if (!session) {
			// redirect to login page
			return;
		}
		const { payload } = JSON.parse(session) as User;

		(async () => {
			try {
				const { token } = payload;
				const req = await fetch(path, {
					method: "GET",
					headers: {
						Accept: "application/vnd.github+json",
						Authorization: `Bearer ${token}`,
						"X-GitHub-Api-Version": "2022-11-28",
					},
				});

				const { status } = req;
				const res = await req.json();
				switch (status) {
					case 200:
						setData(() => payload as T);
						return;
					default:
						setError(() => res["error"]);
						return;
				}
			} catch (e) {
				setError(() => "Server not responding");
				console.error(e);
			}
		})();
	}, []);

	return { data, error };
}
