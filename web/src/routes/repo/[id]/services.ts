import { apiFetch } from "@lib/requests";
import { store } from "./store";
import type { RepositoryEnv, Variable } from "./store";
import { makeUrl } from "@lib/url";

export async function handlePermission(
	repoID: number,
	userIDs: number[],
): Promise<void> {
	const url: string = makeUrl(`/repos/${repoID}/permissions`);
	const headers: Headers = new Headers({
		"Content-Type": "application/json",
	});
	const body: BodyInit = JSON.stringify({ userIDs });
	const rsp = await apiFetch(url, {
		method: "POST",
		headers,
		body,
	});
	if (rsp.status === 200) {
		return;
	}

	const err = (await rsp.json()) as EnvHub.Error;
	throw new Error(err.message);
}

export async function handleEdit(
	repoID: number,
	variableID: string,
	newKey: string,
	newValue: string,
): Promise<void> {
	const url = makeUrl(`/repos/${repoID}/variables/edit`);
	const body: BodyInit = JSON.stringify({
		id: variableID,
		key: newKey,
		value: newValue,
	});
	const headers: Headers = new Headers({
		"Content-Type": "application/json",
		Accept: "application/json",
	});

	const rsp = await apiFetch(url, { method: "PUT", body, headers });

	type UpdateResult = {
		updated_at: string;
	};
	const payload: EnvHub.Response<UpdateResult> = await rsp.json();

	if (rsp.status !== 200) {
		throw new Error((payload as EnvHub.Error).message);
	}

	store.update((state) => {
		const variables: Array<Variable> = state.variables.map((v) => {
			if (v.id === variableID) {
				v.key = newKey;
				v.value = newValue;

				const { updated_at } = payload as UpdateResult;
				v.updated_at = formatTime(new Date(updated_at));
			}
			return v;
		});
		return { ...state, variables };
	});
}

export async function handleDelete(repoID: number, vID: string): Promise<void> {
	const url = makeUrl(`/repos/${repoID}/variables/delete`, { id: vID });
	const rsp = await apiFetch(url, { method: "DELETE" });
	if (rsp.status !== 204) {
		const err: EnvHub.Error = await rsp.json();
		throw new Error(err.message);
	}

	store.update((s) => {
		const variables: Array<Variable> = s.variables.filter(
			(i) => i.id !== vID,
		);
		return { ...s, variables };
	});
}

export async function getVariables(repoID: number): Promise<void> {
	const url = makeUrl(`/repos/${repoID}`);
	const headers = new Headers({
		Accept: "application/json",
	});
	const rsp = await apiFetch(url, { method: "GET", headers });
	const payload: EnvHub.Response<RepositoryEnv> = await rsp.json();
	if (rsp.status !== 200) {
		throw new Error((payload as EnvHub.Error).message);
	}
	store.set({ ...(payload as RepositoryEnv), repoID });
}

export function formatTime(d: Date): string {
	let dt = d.toLocaleDateString() + " ";
	dt += d.getHours().toString() + ":";
	dt += d.getMinutes().toString().padStart(2, "0");
	return dt;
}
