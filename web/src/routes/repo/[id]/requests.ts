import { apiFetch } from "@lib/requests";
import { makeUrl } from "@lib/url";
import { store } from "./store";
import type { NewVariable, Variable } from "./store";

export async function writeNewVariable(
	repoID: number,
	v: NewVariable,
): Promise<void> {
	const url: string = makeUrl(`/repos/${repoID}/variables/new`);
	const headers: Headers = new Headers({
		Accept: "application/json",
		"Content-type": "application/json",
	});

	const body: BodyInit = JSON.stringify(v);
	const rsp = await apiFetch(url, { method: "POST", headers, body });

	switch (rsp.status) {
		case 201:
			const newVar: Variable = await rsp.json();
			newVar.value = v.value;
			store.update((s) => ({
				...s,
				variables: [newVar, ...s.variables],
			}));
			return;

		case 409:
			throw new Error("Variable exists.");

		default:
			const payload: EnvHub.Error = await rsp.json();
			throw new Error(payload.message);
	}
}
