import { HTMLAttributes, ReactNode, useEffect, useState } from "react";
import { useSearchParams, useNavigate } from "react-router-dom";
import cx from "classnames";
import { Spinner } from "./Spinner";
import * as jwt from "jose";
import { WEB } from "lib/routes";
import { Github, User } from "lib/github/request";

const API = import.meta.env.VITE_ENVHUB_API;
const JWT_SECRET: string = import.meta.env.VITE_JWT_SECRET;
const JWT_ALGO = import.meta.env.VITE_JWT_HEADER_ALGO;
const JWT_TYP = import.meta.env.VITE_JWT_HEADER_TYPE;

export function Connect() {
	const { loading, err } = useLogin();
	return (
		<section className="flex h-[100vh] items-center justify-center">
			{err ? (
				<Container>
					<h1 className="mb-2 text-2xl font-semibold text-main">Whoops!</h1>
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

function useLogin() {
	const [param, _] = useSearchParams("code");
	const [loading, setLoading] = useState<boolean>(true);
	const [err, setErr] = useState<string>();
	const nav = useNavigate();

	useEffect(() => {
		(async () => {
			setLoading(() => true);
			setErr(() => undefined);
			try {
				const code = param.get("code");
				if (!code || code === "") throw new Error("`code` param not set");

				const response = await fetch(`${API}/auth/github`, {
					method: "POST",
					headers: {
						"content-type": "application/json",
					},
					body: JSON.stringify({ code }),
				});

				const { status } = response;
				switch (true) {
					case status === 200:
						const code = await response.text();

						const { alg, typ } = jwt.decodeProtectedHeader(code);
						const invalidHeader = alg !== JWT_ALGO || typ !== JWT_TYP;
						if (invalidHeader) {
							setErr(() => "Authentication failed.");
							return;
						}

						const secret = new Uint8Array(
							JWT_SECRET.split("").map((x) => x.charCodeAt(0)),
						);
						const claims = await jwt.jwtVerify(code, secret);
						Github.saveUser(claims as unknown as User);
						Github.saveToken(code);
						nav(WEB.dash);
						return;

					case status === 401 || status === 403:
						setErr(() => "Authentication failed.");
						return;

					default:
						const err = await response.json();
						const msg = err["error"] as string;
						setErr(() => msg.replace(/_+/g, " "));
						return;
				}
			} catch (e) {
				setErr(() => "Server is not responding, try again later.");
				console.error(e);
			} finally {
				setLoading(() => false);
			}
		})();
	}, []);

	return { loading, err };
}
