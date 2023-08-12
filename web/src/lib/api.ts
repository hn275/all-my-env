import { Entry } from "./storage";

const API_URL = import.meta.env.VITE_ENVHUB_API;
if (!API_URL) throw new Error("`VITE_ENVHUB_API` not set");

/**
 * @deprecated too much boilerplate
 * */
export const API = {
	repo: {
		index: "/repo",
		link: "/repo/link",
	},
	auth: {
		refresh: "/auth/refresh",
	},
};

/**
 * @deprecated too much boilerplate
 * */
export class Api {
	public static makeUrl(path: string, param?: object): string {
		if (path[0] !== "/") {
			path = "/" + path;
		}

		let url = API_URL + path;
		if (param) {
			const q = new URLSearchParams();
			for (const [k, v] of Object.entries(param)) {
				q.set(k, String(v));
			}
			url += "?" + q.toString();
		}

		return url;
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
		return fetch(url, { method: "GET", headers, credentials: "include" });
	}

	static POST(
		path: string,
		query?: Record<string, string>,
		body?: object,
	): Promise<Response> {
		return fetch(this.url(path, query ?? undefined), {
			method: "POST",
			credentials: "include",
			headers: this.headers(),
			body: body ? JSON.stringify(body) : undefined,
		});
	}

	private static headers(headers?: Record<string, string>): Headers {
		const t = window.localStorage.getItem(Entry.token);
		const h = new Headers({
			Authorization: `Bearer ${t}`,
		});
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
