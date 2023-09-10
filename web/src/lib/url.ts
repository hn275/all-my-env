import { PUBLIC_ENVHUB_API } from "$env/static/public";

export function makeUrl(path: string, param?: object): string {
	let url: string = PUBLIC_ENVHUB_API;
	url += path;
	if (param) {
		const p = new URLSearchParams();
		for (const [k, v] of Object.entries(param)) {
			p.set(k, String(v));
		}
		url += "?" + p.toString();
	}
	return url;
}

export function urlParam(k: string): string | null {
	const url = new URL(window.location.href);
	return url.searchParams.get(k);
}
