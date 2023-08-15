import { AuthStore } from "./auth";

export function apiFetch(path: string, r: RequestInit): Promise<Response> {
	const user = AuthStore.user();
	if (!user) {
		AuthStore.logout();
		throw new Error("user not found");
	}

	const req: RequestInit = {
		credentials: "include",
		headers: {
			Authorization: "Bearer " + user.access_token,
		},
	};

	return fetch(path, deepCopy(r, req));
}

function deepCopy(dst: any, src: any): any {
	for (const [k, v] of Object.entries(src)) {
		if (typeof v === "object") {
			deepCopy(dst[k], v);
		} else {
			dst[k] = v;
		}
	}
	return dst;
}
