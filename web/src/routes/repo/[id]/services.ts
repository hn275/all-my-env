import { apiFetch } from "@lib/requests";
import { store } from "./store";
import type { RepositoryEnv, Variable } from "./store";
import { makeUrl } from "@lib/url";

export async function handleEdit(
	repoID: number,
	variableID: string,
	newKey: string,
	newValue: string,
): Promise<void> {
	let url: string = `/repos/${repoID}/variables/edit`;
	url = makeUrl(url, { id: variableID });
	const body: BodyInit = JSON.stringify({ key: newKey, value: newValue });
	const headers: Headers = new Headers({
		"Content-Type": "application/json",
		Accept: "application/json",
	});

	const rsp = await apiFetch(url, { method: "PUT", body, headers });
	const payload: EnvHub.Response<Variable> = await rsp.json();
	if (rsp.status !== 200) {
		throw new Error((payload as EnvHub.Error).message);
	}

	store.update((state) => {
		const variables: Array<Variable> = state.variables.map((v) => {
			if (v.id === variableID) {
				return payload as Variable;
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
