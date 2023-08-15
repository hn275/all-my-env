import { makeUrl } from "@lib/url";
import type { Repository } from "./types";
import { AuthStore } from "@lib/auth";
import { apiFetch } from "@lib/requests";

export type Sort = "created" | "updated" | "pushed" | "full_name";

export async function fetchRepos(
	page: number,
	sort: Sort,
	show: string,
): Promise<Array<Repository>> {
	const url = makeUrl("/repos", { page, sort, show });
	const headers = new Headers({
		Accept: "application/json",
	});

	const rsp = await apiFetch(url, {
		method: "GET",
		headers,
	});
	const payload = await rsp.json();

	switch (rsp.status) {
		case 200:
			return payload as Array<Repository>;
		case 401 | 403:
			// AuthStore.refreshSession();
			throw new Error(payload["message"]);
		default:
			throw new Error(payload["message"]);
	}
}

export async function linkRepo(repo: Repository): Promise<void> {
  const url = makeUrl("/repos/link");
  const body = JSON.stringify({
    id: repo.id,
    full_name: repo.full_name
  });
  
  const rsp = await apiFetch(url, {
    method: "POST",
    body
  });
  if (rsp.status === 201) return;
  const payload: EnvHub.Error = await rsp.json();
  throw new Error(payload.message);
}
