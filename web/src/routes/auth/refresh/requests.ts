import { AuthStore, type User } from "@lib/auth";
import { makeUrl } from "@lib/url";

export async function refreshSession(): Promise<User> {
	const url: string = makeUrl("/auth/refresh");
	const headers: Headers = new Headers({
		Accept: "application/json",
	});
	const rsp: Response = await fetch(url, {
		method: "GET",
		headers,
		credentials: "include",
	});
	const payload = await rsp.json();
	switch (rsp.status) {
		case 200:
			AuthStore.login(payload as User);
			return payload as User;
		case 401 | 403:
			AuthStore.logout();
			throw new Error(payload["message"]);
		default:
			throw new Error(payload["message"]);
	}
}
