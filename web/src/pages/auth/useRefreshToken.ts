import { useAuth } from "context/auth";
import { Entry } from "lib/storage";
import { useEffect } from "react";
import { AuthResponse } from ".";
import { makeUrl } from "lib/url";

export function useRefreshToken() {
	const refreshed = window.sessionStorage.getItem(Entry.tokenRefresh);
	const { dispatch, getToken } = useAuth();

	useEffect(() => {
		if (refreshed) return;
		(async () => {
			const r = window.sessionStorage.getItem(Entry.tokenRefresh);
			if (r) return;

			const tok = getToken();
			if (!tok) {
				dispatch({ type: "logout" });
				return;
			}

			const url = makeUrl("/auth/refresh");
			const headers = new Headers({
				Accept: "application/json",
				Authorization: "Bearer " + tok,
			});
			try {
				const res = await fetch(url, {
					method: "GET",
					credentials: "include",
					headers,
				});
				const { status } = res;
				if (status === 401 || status === 403) {
					dispatch({ type: "logout" });
					return;
				}
				const payload = await res.json();
				if (status !== 200) throw new Error(payload["message"]);
				const user: AuthResponse = payload;
				window.localStorage.setItem(Entry.token, user.access_token);
				window.sessionStorage.setItem(Entry.tokenRefresh, JSON.stringify(true));
				dispatch({ type: "login", payload: user });
			} catch (e) {
				console.error(e);
			}
		})();
	}, []);
}
