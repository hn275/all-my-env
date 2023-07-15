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
    <nav className="md:flex justify-between md:backdrop-blur">
      <div className="flex justify-between items-center h-16 px-5">
        <img src={Logo} alt="logo" />
        <button
          onClick={toggleOpen}
          className="md:hidden bg-inherit hover:bg-[#3a3a3a] p-2 rounded-md transition-all"
        >
          <MdMenu className="text-main text-xl" />
        </button>
      </div>

      <div
        className={cx(["absolute top-16 left-0 w-full", "md:static md:w-max"])}
      >
        <ul
          className={cx([
            `w-full ${open ? "h-[350px] py-5" : "h-0 py-0"}`,
            "bg-[#1e1e1e] text-light font-semibold",
            "flex flex-col justify-between items-center",
            "transition-all",
            "relative overflow-clip",
            "md:h-full md:flex-row gap-10 md:bg-inherit",
          ])}
        >
          <li>
            <a
              href={oauthUrl}
              className={cx([
                "flex items-center hover:cursor-pointer hover:no-underline",
                "text-main md:bg-main md:text-dark md:py-2 md:px-3 md:rounded-md",
                "hover:brightness-95 transition-all",
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
              className="absolute md:hidden bottom-2 left-1/2 -translate-x-1/2"
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
