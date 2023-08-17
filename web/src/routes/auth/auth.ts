import { makeUrl, urlRedirect } from "@lib/url";
import { AuthStore, type User } from "@lib/auth";
import { apiFetch } from "@lib/requests";

export async function signIn(code: string): Promise<User> {
  const res = await fetch(makeUrl("/auth"), {
    method: "POST",
    headers: {
      "content-type": "application/json",
      accept: "application/json",
    },
    body: JSON.stringify({ code }),
    credentials: "include",
  });
  const payload: EnvHub.Response<User> = await res.json();
  if (res.status !== 200) {
    throw new Error((payload as EnvHub.Error)["message"]);
  }
  AuthStore.login(payload as User);
  urlRedirect("/dash");
  return payload as User;
}

export async function refresh(): Promise<User> {
  const res = await apiFetch(makeUrl("/auth/refresh"), {
    method: "GET",
    headers: {
      Accept: "application/json",
    },
  });
  const payload = await res.json();
  if (res.status !== 200) throw new Error(payload["message"]);
  return payload as User;
}
