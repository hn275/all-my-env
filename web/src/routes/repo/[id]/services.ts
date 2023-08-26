import { apiFetch } from "@lib/requests";
import { store } from "./store";
import type { RepositoryEnv, Variable } from "./store";
import { makeUrl } from "@lib/url";

export async function handleDelete(repoID: number, vID: string): Promise<void> {
    if (!repoID) {
        throw new Error("repository id not found.");
    }
    const url = makeUrl(`/repos/${repoID}/variables/delete`, { id: vID });
    const rsp = await apiFetch(url, { method: "DELETE" });
    if (rsp.status === 204) {
        store.update((s) => {
            const newVariables: Array<Variable> = s.variables.filter(
                (i) => i.id !== vID,
            );
            return {
                ...s,
                variables: newVariables,
                deleteVariable: undefined,
            };
        });
        return;
    }
    const err: EnvHub.Error = await rsp.json();
    throw new Error(err.message);
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
