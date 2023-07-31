import { Session } from "./sessionStorage";

const API_URL = import.meta.env.VITE_ENVHUB_API;
if (!API_URL) throw new Error("`VITE_ENVHUB_API` not set");

export const API = {
	repo: {
		index: "/repo",
		link: "/repo/link",
	},
};

export class UnauthorizeError extends Error {
	constructor() {
		super("Session token not found.");
		this.name = "UnauthorizeError";
	}
}

export class Fetch {
	static async GET(
		path: string,
		query?: Record<string, string>,
	): Promise<Response> {
		let url = API_URL + path;
		if (query) url += "?" + new URLSearchParams(query).toString();
		const headers = this.headers();
		return fetch(url, { method: "GET", headers });
	}

	static POST(
		path: string,
		query: Record<string, string> | null,
		body?: object,
	): Promise<Response> {
		return fetch(this.url(path, query ?? undefined), {
			method: "POST",
			headers: this.headers(),
			body: body ? JSON.stringify(body) : undefined,
		});
	}

	private static headers(headers?: Record<string, string>): Headers {
		const t = window.sessionStorage.getItem(Session.token);
		if (!t) throw new UnauthorizeError();
		const h = new Headers({ Authorization: `Bearer ${t}` });
		if (headers) Object.assign(h, headers);
		return h;
	}

	private static url(path: string, query?: Record<string, string>): string {
		if (path[0] != "/") path = "/" + path;
		const url = API_URL + path;
		if (!query) return url;
		return url + "?" + new URLSearchParams(query).toString();
	}
}
