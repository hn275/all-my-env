import { makeUrl, urlRedirect } from "@lib/url";
import { AuthStore, type User } from "@lib/auth";

export async function signIn(code: string): Promise<User> {
  const res = await fetch(makeUrl("/auth/github"), {
    method: "POST",
    headers: {
      "content-type": "application/json",
      accept: "application/json",
    },
    body: JSON.stringify({ code }),
    credentials: "include",
  });
  const payload: EnvHub.Response<User> = await res.json();
  if (res.status !== 200) throw new Error((payload as EnvHub.Error)["message"]);
  AuthStore.login(payload as User);
  urlRedirect("/dash");
  return payload as User;
}

export async function refresh(token: string): Promise<User> {
  const res = await fetch(makeUrl("/auth/refresh"), {
    method: "GET",
    credentials: "include",
    headers: {
      Accept: "application/json",
      Authorization: `Bearer ${token}`,
    },
  });
  const payload = await res.json();
  if (res.status !== 200) throw new Error(payload["message"]);
  return payload as User;
}
