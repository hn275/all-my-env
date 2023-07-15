import { useState } from "react";
import Logo from "assets/logo.svg";
import cx from "classnames";
import { MdMenu } from "react-icons/md";
import { BsChevronCompactUp } from "react-icons/bs";
import { AiFillGithub, AiFillStar } from "react-icons/ai";

interface Props {
	oauthUrl: string;
	githubUrl: string;
}
export function Nav({ oauthUrl, githubUrl }: Props) {
	const { open, toggleOpen } = useMenu();

	return (
		<nav className="justify-between md:flex md:backdrop-blur">
			<div className="flex h-16 items-center justify-between px-5">
				<img src={Logo} alt="logo" />
				<button
					onClick={toggleOpen}
					className="rounded-md bg-inherit p-2 transition-all hover:bg-[#3a3a3a] md:hidden"
				>
					<MdMenu className="text-xl text-main" />
				</button>
			</div>

			<div
				className={cx(["absolute left-0 top-16 w-full", "md:static md:w-max"])}
			>
				<ul
					className={cx([
						`w-full ${open ? "h-[350px] py-5" : "h-0 py-0"}`,
						"bg-[#1e1e1e] font-semibold text-light",
						"flex flex-col items-center justify-between",
						"transition-all",
						"relative overflow-clip",
						"z-50 gap-10 md:h-full md:flex-row md:bg-inherit",
					])}
				>
					<li>
						<a
							href={oauthUrl}
							className={cx([
								"flex items-center hover:cursor-pointer hover:no-underline",
								"text-main md:rounded-md md:bg-main md:px-3 md:py-2 md:text-dark",
								"transition-all hover:brightness-95",
							])}
						>
							<span>
								<AiFillGithub />
							</span>
							&nbsp; Sign in
						</a>
					</li>
					<li>
						<a>Pricing</a>
					</li>
					<li>
						<a>Docs</a>
					</li>
					<li>
						<a>FAQ</a>
					</li>
					<li>
						<a
							href={githubUrl}
							target="_blank"
							className="flex items-center md:text-main"
						>
							<span className="inline-block">
								<AiFillStar />
							</span>
							&nbsp; us on GitHub
						</a>
					</li>
					<li>
						<button
							onClick={toggleOpen}
							className="absolute bottom-2 left-1/2 -translate-x-1/2 md:hidden"
						>
							<BsChevronCompactUp className="text-main" />
						</button>
					</li>
				</ul>
			</div>
		</nav>
	);
}

function useMenu() {
	const [open, setOpen] = useState<boolean>(false);
	const toggleOpen = () => setOpen((o) => !o);
	return { open, toggleOpen };
}
