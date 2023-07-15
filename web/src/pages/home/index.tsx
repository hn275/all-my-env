import { GITHUB_SECRET } from "lib/github/request";
import { Layout } from "components/Layout";
import { Nav } from "./Nav";

export function Home() {
	return (
		<Layout>
			<Nav url={getLoginUrl()} />
		</Layout>
	);
}

function getLoginUrl(): string {
	const githubLogin = "https://github.com/login/oauth/authorize";
	const param = new URLSearchParams({
		client_id: GITHUB_SECRET,
	});
	return githubLogin + "?" + param.toString();
}
