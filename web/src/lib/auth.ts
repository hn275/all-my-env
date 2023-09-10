import { PUBLIC_GITHUB_CLIENT_ID, PUBLIC_NODE_ENV } from "$env/static/public";
const IS_PROD: boolean = PUBLIC_NODE_ENV === "production";
const HTTP = IS_PROD ? "https://" : "http://";

// store
export type User = {
	access_token: string;
	name: string;
	avatar_url: string;
	login: string;
};

export class Auth {
	public static user(): User | undefined {
		const u = window.localStorage.getItem("envhub:user");
		return u ? JSON.parse(u) : undefined;
	}

	public static login(u: User): void {
		window.localStorage.setItem("envhub:user", JSON.stringify(u));
	}

	public static logout(): void {
		window.localStorage.removeItem("envhub:user");
	}

	public static refresh(): void {
		window.sessionStorage.setItem(
			"envhub:session_refreshed",
			JSON.stringify(true),
		);
	}

	public static isRefreshed(): boolean {
		return (
			window.sessionStorage.getItem("envhub:session_refreshed") !== null
		);
	}
}

export function oauth(redirect?: string): void {
	const client_id = PUBLIC_GITHUB_CLIENT_ID;
	const scope = "read:repo read:user read:org";
	const redirect_uri = HTTP + window.location.host + "/auth";

	let to: string;
	if (redirect) {
		to = HTTP + window.location.host + redirect;
	} else {
		const curLoc: URL = new URL(window.location.href);
		for (const k of curLoc.searchParams.keys()) {
			curLoc.searchParams.delete(k);
		}
		to = curLoc.toString();
	}
	window.sessionStorage.setItem("envhub:redirect_uri", to);

	const p = new URLSearchParams({ client_id, redirect_uri, scope });
	const l = "https://github.com/login/oauth/authorize?" + p.toString();
	window.location.replace(l);
}
