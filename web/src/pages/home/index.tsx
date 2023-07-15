import { GITHUB_SECRET } from "lib/github/request";
import { Layout } from "components/Layout";
import { Nav } from "./Nav";
import { AiFillGithub } from "react-icons/ai";
import "./style.css";
import cx from "classnames";

export function Home() {
	return (
		<>
			<Nav oauthUrl={getLoginUrl()} githubUrl="" />

			{/* HERO */}
			<section
				className={cx([
					"relative h-[calc(100vh-64px)] w-full overflow-x-hidden",
					"flex flex-col items-center justify-center gap-6",
					"px-3 text-center md:px-0",
				])}
			>
				<h1 className="hero-text-gradient font-accent text-4xl font-bold uppercase">
					effortless secrets management
				</h1>

				<p className="w-80 font-semibold md:w-full md:text-lg">
					The ultimate <span className="hero-text-underline">open-source</span>
					&nbsp; solution for managing your&nbsp;
					<span className="hero-text-underline">variables</span>.
				</p>

				<div className="flex flex-col items-center justify-center gap-3 md:flex-row md:gap-6">
					<button className="h-10 w-36 rounded-md border border-main font-semibold text-main">
						Get started
					</button>

					<button
						className={cx([
							"flex items-center justify-center border border-main bg-main md:flex-row",
							"h-10 w-36 rounded-md font-semibold",
						])}
					>
						Log In&nbsp;
						<AiFillGithub />
					</button>
				</div>

				<div className="hero-graphic-blur main" />
				<div className="hero-graphic-blur accent" />
			</section>
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
