import { useState } from "react";
import Logo from "assets/logo.svg";
import cx from "classnames";
import { MdMenu } from "react-icons/md";
import { BsChevronCompactUp } from "react-icons/bs";

interface Props {
  url: string;
}
export function Nav({ url }: Props) {
  const { open, toggleOpen } = useMenu();

  return (
    <>
      <div className="flex justify-between items-center h-16 px-5">
        <img src={Logo} alt="logo" />
        <button onClick={toggleOpen}>
          <MdMenu className="text-main text-xl" />
        </button>
      </div>

      <nav className="absolute top-16 left-0 w-full">
        <ul
          className={cx([
            `w-full ${open ? "h-48 py-5" : "h-0 py-0"}`,
            "bg-[#1e1e1e] text-light",
            "flex flex-col justify-between items-center",
            "transition-all",
            "relative overflow-clip",
          ])}
        >
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
            <a href={url} target="_blank">
              start us on GitHub
            </a>
          </li>
          <li>
            <button
              onClick={toggleOpen}
              className="absolute bottom-2 left-1/2 -translate-x-1/2"
            >
              <BsChevronCompactUp className="text-main" />
            </button>
          </li>
        </ul>
      </nav>
    </>
  );
}

function useMenu() {
  const [open, setOpen] = useState<boolean>(false);
  const toggleOpen = () => setOpen((o) => !o);
  return { open, toggleOpen };
}
