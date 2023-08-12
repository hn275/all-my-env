import { useContext } from "react";
import { AuthCtxProps, AuthCtx, User } from "./auth";
import { Entry } from "lib/storage";

export const GITHUB_SECRET = import.meta.env.VITE_GITHUB_CLIENT_ID;
if (!GITHUB_SECRET) throw new Error("`VITE_GITHUB_CLIENT_ID` not set");

export interface UseAuthProps extends AuthCtxProps {
	getToken: () => string | null;
}
export function useAuth(): UseAuthProps {
	const { state, dispatch } = useContext(AuthCtx);

	function getToken(): string | null {
		if (state.user) return state.user.access_token;
		const u = window.localStorage.getItem(Entry.user);
		if (!u) return null;
		const user: User = JSON.parse(u);
		return user.access_token;
	}

	return { state, dispatch, getToken };
}

export function oauth(redirect?: string): void {
	let r = url(window.location.href);
	if (redirect) r = url(window.location.host + redirect);

	const p = new URLSearchParams({
		client_id: GITHUB_SECRET,
		redirect_uri: r,
		scope: "repo user read:org",
	});
	const login = "https://github.com/login/oauth/authorize?" + p.toString();
	window.location.replace(login);
}

function url(r: string): string {
	const mode = import.meta.env.NODE_ENV;
	const production = mode === "production";
	if (!production) {
		console.warn("NODE_ENV:", mode);
		return "http://" + r;
	}
	return "https://" + r;
}
