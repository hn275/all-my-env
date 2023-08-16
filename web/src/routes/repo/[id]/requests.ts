import { apiFetch } from "@lib/requests";
import { makeUrl } from "@lib/url";
import { store, type Variable } from "./store";

export async function getVariables(repoID: number): Promise<void> {
  const url = makeUrl(`/repos/${repoID}`);
  const headers = new Headers({
    Accept: "application/json",
  });
  const rsp = await apiFetch(url, { method: "GET", headers });
  const payload: EnvHub.Response<Array<Variable>> = await rsp.json();
  if (rsp.status !== 200) {
    throw new Error((payload as EnvHub.Error).message);
  }
  store.set(payload as Array<Variable>);
}
