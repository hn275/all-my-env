import { useAuth } from "context/auth";
import { Fetch, API } from "lib/api";
import { Entry } from "lib/storage";
import { useEffect } from "react";
import { AuthResponse } from ".";

export function useRefreshToken() {
	const refreshed = window.sessionStorage.getItem(Entry.tokenRefresh);
	const { dispatch } = useAuth();

	useEffect(() => {
		if (refreshed) return;
		(async () => {
			try {
				const res = await Fetch.POST(API.auth.refresh);
				const { status } = res;
				if (status === 401 || status === 403) {
					console.log("auth failed");
					// dispatch({ type: "logout" });
					return;
				}
				const payload = await res.json();
				if (status !== 200) throw new Error(payload["message"]);
				const user: AuthResponse = payload;

				window.localStorage.setItem(Entry.token, user.access_token);
				dispatch({ type: "login", payload: user });
			} catch (e) {
				console.error(e);
			}
		})();
	}, []);
}
