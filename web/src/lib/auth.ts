import { PUBLIC_GITHUB_CLIENT_ID } from "$env/static/public";
import { map } from "nanostores";

// store
export type User = {
	access_token: string;
	name: string;
	avatar_url: string;
	login: string;
};

export const user = map<User>(undefined);

// fn
export function oauth(redirect?: string): void {
	const client_id = PUBLIC_GITHUB_CLIENT_ID;
	const scope = "repo user read:org";

	let redirect_uri: string;
	if (redirect) {
		redirect_uri = window.location.host + redirect;
	} else {
		const curLoc: URL = new URL(window.location.href);
		curLoc.searchParams.delete("code");
		redirect_uri = curLoc.toString();
	}

	const p = new URLSearchParams({ client_id, redirect_uri, scope });
	const l = "https://github.com/login/oauth/authorize?" + p.toString();
	window.location.replace(l);
}
