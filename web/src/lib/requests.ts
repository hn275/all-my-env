import { AuthStore } from "./auth";

export function apiFetch(path: string, r: RequestInit): Promise<Response> {
	const user = AuthStore.user();
	if (!user) {
		AuthStore.logout();
		throw new Error("user not found");
	}

	const headers = new Headers({
		...r.headers,
		Authorization: "Bearer " + user.access_token,
	});
	r.credentials = "include";
	r.headers = headers;

	return fetch(path, r);
}
