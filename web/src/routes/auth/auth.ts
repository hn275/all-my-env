import { makeUrl } from "@lib/url";
import { Auth } from "@lib/auth";
import type { User } from "@lib/auth";

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
	Auth.login(payload as User);
	return payload as User;
}
