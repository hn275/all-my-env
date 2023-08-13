import { makeUrl } from "@lib/url";
import type { Repository } from "./types";
import { AuthStore } from "@lib/auth";

export type Sort = "created" | "updated" | "pushed" | "full_name";

export async function fetchRepos(
	page: number,
	sort: Sort,
	show: number,
): Promise<Repository[]> {
	const url = makeUrl("/repos", { page, sort, show });
	const headers = new Headers({
		Accept: "application/json",
	});

	const rsp = await fetch(url, {
		method: "GET",
		credentials: "include",
		headers,
	});
	const payload = await rsp.json();

	switch (rsp.status) {
		case 200:
			return payload as Repository[];
		case 401 | 403:
			AuthStore.logout();
			throw new Error(payload["message"]);
		default:
			throw new Error(payload["message"]);
	}
}
