import { PUBLIC_GITHUB_CLIENT_ID, PUBLIC_NODE_ENV } from "$env/static/public";

// store
export type User = {
	access_token: string;
	name: string;
	avatar_url: string;
	login: string;
};

export type AuthStoreType = {
	user?: User;
	tokenRefreshed: boolean;
};

export class AuthStore {
	private static userEntry: string = "envhub:user";
	private static tokenEntry: string = "envhub:refreshed";
	private static redirectEntry: string = "envhub:redirect";

	public static init() {}

	public static login(u: User): AuthStore {
		window.localStorage.setItem(AuthStore.userEntry, JSON.stringify(u));
		return AuthStore;
	}

	public static logout(): AuthStore {
		window.localStorage.removeItem(AuthStore.userEntry);
		return AuthStore;
	}

	public static refreshed(): AuthStore {
		window.sessionStorage.setItem(
			AuthStore.tokenEntry,
			JSON.stringify(true),
		);
		return AuthStore;
	}

	public static user(): User | undefined {
		const u = window.localStorage.getItem(AuthStore.userEntry);
		if (!u) return;
		return JSON.parse(u!) as User;
	}

	public static sessionRefreshed(): boolean {
		const r = window.sessionStorage.getItem(AuthStore.tokenEntry);
		return r !== null;
	}

	public static refreshSession(redirect?: string): void {
		const r = redirect ?? window.location.href;
		window.sessionStorage.setItem(this.redirectEntry, r);
		window.location.replace("/auth/refresh");
	}

	public static redirectUrl(): string {
		return window.sessionStorage.getItem(this.redirectEntry) ?? "/";
	}
}

export function oauth(redirect?: string): void {
	const client_id = PUBLIC_GITHUB_CLIENT_ID;
	const scope = "repo user read:org";

	let redirect_uri: string;
	if (redirect) {
		const isProd: boolean = PUBLIC_NODE_ENV === "production";
		const http = isProd ? "https://" : "http://";
		redirect_uri = http + window.location.host + redirect;
	} else {
		const curLoc: URL = new URL(window.location.href);
		for (const k of curLoc.searchParams.keys()) {
			console.log(k);
			curLoc.searchParams.delete(k);
		}
		redirect_uri = curLoc.toString();
	}

	const p = new URLSearchParams({ client_id, redirect_uri, scope });
	const l = "https://github.com/login/oauth/authorize?" + p.toString();
	window.location.replace(l);
}
