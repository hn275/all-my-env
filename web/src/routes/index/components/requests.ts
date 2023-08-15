import { makeUrl } from "@lib/url";

export async function logout(): Promise<void> {
	const url: string = makeUrl("/auth/logout");
	const rsp: Response = await fetch(url, {
		method: "GET",
		credentials: "include",
	});
	if (rsp.status !== 200) {
		const payload: EnvHub.Error = await rsp.json();
		throw new Error(payload.message);
	}
}
