import { API, Fetch, UnauthorizeError } from "lib/api";
import { oauth } from "lib/auth";
import { WEB } from "lib/routes";
import { Spinner } from "pages/auth/Spinner";
import { useEffect, useState } from "react";
import { useSearchParams, useNavigate, Link } from "react-router-dom";

export function Connect() {
	const { error, loading } = useConnectRepo();
	return (
		<>
			{loading ? (
				<>
					<Spinner />
				</>
			) : error === "Session token expired." ? (
				<>
					<Link to={WEB.home}>Back to Home page</Link>
				</>
			) : (
				<p>{error}</p>
			)}
		</>
	);
}

type ConnectRepoHook = { loading: boolean; error: string | undefined };
function useConnectRepo(): ConnectRepoHook {
	const [param] = useSearchParams();
	const nav = useNavigate();
	const repoNameEncoded = param.get("repo_name");

	const [loading, setLoading] = useState<boolean>(true);
	const [error, setError] = useState<string>();

	useEffect(() => {
		(async () => {
			if (!repoNameEncoded) {
				nav(-1);
				return;
			}

			try {
				const body = { full_name: decodeURIComponent(repoNameEncoded) };
				const queries = null;
				const res = await Fetch.POST(API.repo.link, queries, body);
				switch (res.status) {
					case 201:
						setError(() => undefined);
						nav(`${WEB.repo}/${await res.text()}`);
						return;

					case 403 | 401:
						oauth();
						return;

					default:
						type ErrResponse = { error: string };
						const err = (await res.json()) as ErrResponse;
						setError(() => err.error);
						return;
				}
			} catch (e) {
				setError(() => "Something went wrong, try again later.");
				console.error(e);
			} finally {
				setLoading(() => false);
			}
		})();
	}, []);

	return { loading, error };
}
