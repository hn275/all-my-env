export const GITHUB_SECRET = import.meta.env.VITE_GITHUB_CLIENT_ID;
if (!GITHUB_SECRET) throw new Error("`GITHUB_SECRET` not set");

export function oauth(redirect?: string): void {
	let r = url(window.location.href);
	if (redirect) r = url(window.location.host + redirect);

	const p = new URLSearchParams({
		client_id: GITHUB_SECRET,
		redirect_uri: r,
		scope: "repo user read:org",
	});
	const login = "https://github.com/login/oauth/authorize?" + p.toString();
	window.location.replace(login);
}

function url(r: string): string {
	const mode = import.meta.env.NODE_ENV;
	const production = mode === "production";
	if (!production) {
		console.warn("NODE_ENV:", mode);
		return "http://" + r;
	}
	return "https://" + r;
}
