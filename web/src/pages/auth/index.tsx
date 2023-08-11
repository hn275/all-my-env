import { HTMLAttributes, ReactNode, useEffect, useState } from "react";
import { useSearchParams, useNavigate } from "react-router-dom";
import cx from "classnames";
import { Spinner } from "./Spinner";
import { WEB } from "lib/routes";
import { useAuth, User } from "context/auth";
import { Entry } from "lib/storage";
import { Api } from "lib/api";

export function Auth() {
	const { loading, err } = useLogin();
	return (
		<section className="flex h-[100vh] items-center justify-center">
			{err ? (
				<Container>
					<h1 className="mb-2 text-2xl font-semibold text-main">Error:</h1>
					<p>{err}</p>
				</Container>
			) : (
				loading && (
					<div className="flex flex-col items-center gap-4">
						<Spinner />
						<p className="text-lg font-semibold text-main">Authenticating...</p>
					</div>
				)
			)}
		</section>
	);
}

interface ContainerProps extends HTMLAttributes<HTMLDivElement> {
	children: ReactNode;
}
function Container({ children, className, ...rest }: ContainerProps) {
	return (
		<div
			{...rest}
			className={cx([className, "border border-main p-10 text-center"])}
		>
			{children}
		</div>
	);
}

export interface AuthResponse extends User {
	access_token: string;
}
function useLogin() {
	const [param, _] = useSearchParams("code");
	const [loading, setLoading] = useState<boolean>(true);
	const [err, setErr] = useState<string>();
	const nav = useNavigate();
	const { dispatch } = useAuth();

	useEffect(() => {
		(async () => {
			setLoading(() => true);
			setErr(() => undefined);
			try {
				const code = param.get("code");
				if (!code || code === "") throw new Error("`code` param not set");

				const response = await fetch(Api.makeUrl("/auth/github"), {
					method: "POST",
					headers: {
						"content-type": "application/json",
					},
					body: JSON.stringify({ code }),
					credentials: "include",
				});

				const { status } = response;
				switch (true) {
					case status === 200:
						const payload: AuthResponse = await response.json();
						console.log(payload);
						window.localStorage.setItem(Entry.token, payload.access_token);
						dispatch({ type: "login", payload });
						nav(WEB.dash);
						return;

					default:
						const err = await response.json();
						const msg = err["error"] as string;
						setErr(() => msg.replace(/_+/g, " "));
						return;
				}
			} catch (e: any) {
				setErr(() => (e as Error).message);
				console.error(e);
			} finally {
				setLoading(() => false);
			}
		})();
	}, []);

	return { loading, err };
}
