import { useEffect, useState } from "react";
import Logo from "assets/logo.svg";
import cx from "classnames";
import { MdMenu } from "react-icons/md";
import { BsChevronCompactUp } from "react-icons/bs";
import { AiFillStar } from "react-icons/ai";
import { LogInButton } from "./LogInButton";
import { oauth } from "lib/auth";

export function Nav() {
	const { open, toggleOpen } = useMenu();
	const show = useNavScroll();

	return (
		<nav
			className={cx([
				"sticky left-0 top-0 z-[49] transition-all",
				"bg-dark -translate-y-full justify-between bg-transparent backdrop-blur md:flex",
				{ "translate-y-0": show },
			])}
		>
			<div className="flex h-16 items-center justify-between px-5">
				<img src={Logo} alt="logo" />
				<button
					onClick={toggleOpen}
					className="rounded-md bg-inherit p-2 transition-all hover:bg-[#3a3a3a] md:hidden"
				>
					<MdMenu className="text-main text-xl" />
				</button>
			</div>

			<div
				className={cx(["absolute left-0 top-16 w-full", "md:static md:w-max"])}
			>
				<ul
					className={cx([
						`w-full ${open ? "h-[350px] py-5" : "h-0 py-0"}`,
						"text-light bg-[#1e1e1e] font-semibold",
						"flex flex-col items-center justify-between",
						"transition-all",
						"relative overflow-clip",
						"z-50 gap-10 md:h-full md:flex-row md:bg-inherit",
					])}
				>
					<li>
						<LogInButton onClick={() => oauth("/auth")} />
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
						<StarUsOnGitHub />
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

export function StarUsOnGitHub() {
	return (
		<a href="" target="_blank" className="md:text-main flex items-center">
			<span className="inline-block">
				<AiFillStar />
			</span>
			&nbsp; us on GitHub
		</a>
	);
}

function useMenu() {
	const [open, setOpen] = useState<boolean>(false);
	const toggleOpen = () => setOpen((o) => !o);
	return { open, toggleOpen };
}

function useNavScroll() {
	const [show, setShow] = useState<boolean>(false);
	useEffect(() => {
		const fn = () => setShow(() => window.scrollY > window.innerHeight / 4);

		window.addEventListener("scroll", fn);
		return () => window.removeEventListener("scroll", fn);
	}, [window.scrollY, window.innerHeight]);

	return show;
}
