import cx from "classnames";
import { useAuth } from "context/auth";
import { HTMLAttributes } from "react";
import { AiFillGithub } from "react-icons/ai";
import { Link } from "react-router-dom";

interface Props extends HTMLAttributes<HTMLButtonElement> {}
export function LogInButton({ className, ...rest }: Props) {
	const {
		state: { auth },
	} = useAuth();
	return (
		<button
			className={cx([
				"flex items-center justify-center hover:cursor-pointer hover:no-underline",
				"md:bg-main font-bold md:text-dark md:rounded-md md:px-3 md:py-2",
				"transition-all hover:brightness-95",
				className,
			])}
			{...rest}
		>
			{auth ? (
				<>
          <Link to="/dash">Dashboard</Link>
        </>
			) : (
				<>
					<span>
						<AiFillGithub />
					</span>
					&nbsp; Sign in
				</>
			)}
		</button>
	);
}
