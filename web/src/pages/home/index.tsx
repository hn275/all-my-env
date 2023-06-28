import { GITHUB_SECRET } from "lib/github";

export function Home() {
	return (
		<>
			<a href={getLoginUrl()}>login</a>
		</>
	);
}

function getLoginUrl(): string {
	const githubLogin = "https://github.com/login/oauth/authorize";
	const param = new URLSearchParams({
		client_id: GITHUB_SECRET,
	});
	return githubLogin + "?" + param.toString();
}
