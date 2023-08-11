import { Route, Routes } from "react-router-dom";
import { Home } from "pages/home";
import { Auth } from "pages/auth";
import { NotFound } from "pages/notfound";
import { Dash } from "pages/dash";
import { Repo, Connect } from "pages/repos";
import { WEB } from "lib/routes";
import { useAuth } from "context/auth";
import { AnimatePresence, motion } from "framer-motion";
import { oauth } from "lib/auth";
import cx from "classnames";
import { LogInButton } from "pages/home/LogInButton";
import { useEffect } from "react";
import { useRefreshToken } from "pages/auth/useRefreshToken";

export default function App() {
	const {
		state: { auth },
	} = useAuth();
	useRefreshToken();
	return (
		<>
			<TokenExpired auth={auth} />
			<Routes>
				<Route index path={WEB.home} element={<Home />} />
				<Route path={WEB.auth} element={<Auth />} />
				<Route path={WEB.dash} element={<Dash />} />
				<Route path={WEB.repo}>
					<Route path={`${WEB.repo}/:id/`} element={<Repo />} />
					<Route path={`${WEB.repo}/:id/connect/`} element={<Connect />} />
				</Route>
				<Route index path="*" element={<NotFound />} />
			</Routes>
		</>
	);
}

type TokenExpiredProps = {
	auth: boolean;
};
function TokenExpired({ auth }: TokenExpiredProps) {
	useEffect(() => {
		document.body.style.overflowY = auth ? "scroll" : "hidden";
	}, [auth]);

	return (
		<AnimatePresence>
			{!auth ? (
				<motion.div
					key="login-modal"
					initial={{ opacity: 0, backdropFilter: "blur(0px)" }}
					animate={{ opacity: 1, backdropFilter: "blur(8px)" }}
					exit={{ opacity: 0, backdropFilter: "blur(0px)" }}
					className={cx([
						"fixed top-0 left-0",
						"bg-dark-200/80 z-50 h-[100vh] w-[100vw]",
						"grid place-items-center",
					])}
				>
					<motion.section>
						<h2 className="hero-text-gradient text-lg font-bold">
							Session Expired, sign in again to continue.
						</h2>
						<LogInButton onClick={() => oauth()} />
					</motion.section>
				</motion.div>
			) : (
				<></>
			)}
		</AnimatePresence>
	);
}
