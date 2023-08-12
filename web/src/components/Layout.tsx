import { ReactNode, useEffect } from "react";
import { useAuth } from "context/auth";
import { AnimatePresence, motion } from "framer-motion";
import cx from "classnames";

interface Props {
	children: ReactNode;
}
export function Layout({ children }: Props) {
	const {
		state: { auth },
	} = useAuth();
	return (
		<>
			<TokenExpired auth={auth} />

			<div className="max-w-screen-2xl w-full mx-auto">{children}</div>
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
						<a href="">Log in to continue</a>
					</motion.section>
				</motion.div>
			) : (
				<></>
			)}
		</AnimatePresence>
	);
}
