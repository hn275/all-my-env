import { PUBLIC_ENVHUB_API } from "$env/static/public";

export function makeUrl(path: string, param?: Record<string, string>): string {
	let url: string = PUBLIC_ENVHUB_API;
	url += path;
	if (param) {
		const p = new URLSearchParams(param);
		url += "?";
		url += p.toString();
	}
	return url;
}
