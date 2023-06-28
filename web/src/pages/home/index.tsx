import { GITHUB_SECRET } from "lib/github";
import Logo from "assets/logo.svg";
import { Layout } from "components/Layout";

export function Home() {
	return (
		<Layout>
			<img src={Logo} alt="logo" />
			<a href={getLoginUrl()}>login</a>
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
