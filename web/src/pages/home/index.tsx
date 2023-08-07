import { oauth } from "lib/auth";
import { Nav, StarUsOnGitHub } from "./Nav";
import { AiFillGithub } from "react-icons/ai";
import "./style.css";
import cx from "classnames";
import { Features } from "./Feature";

export function Home() {
  return (
    <>
      <Nav handleAuth={oauth} />

      {/* HERO */}
      <main className="relative mx-auto h-max min-h-[calc(100vh-65px)] max-w-7xl overflow-x-hidden">
        <section className="flex h-[calc(100vh-65px)] flex-col items-center justify-center gap-6">
          <h1 className="hero-text-gradient font-accent text-center text-3xl font-bold uppercase md:text-4xl">
            effortless secrets management
          </h1>

          <p className="w-[25ch] text-center font-semibold md:w-max md:text-lg">
            The ultimate{" "}
            <span className="hero-text-underline">open-source</span>
            &nbsp; solution for managing your&nbsp;
            <span className="hero-text-underline">variables</span>.
          </p>

          <div className="flex flex-col items-center justify-center gap-3 md:flex-row md:gap-6">
            <a className="border-main text-main w-36 rounded-md border py-2 text-center font-semibold">
              Get started
            </a>

            <button
              className={cx([
                "border-main bg-main flex items-center justify-center border md:flex-row",
                "h-10 w-36 rounded-md font-semibold hover:no-underline hover:brightness-90",
              ])}
              onClick={() => oauth("/auth")}
            >
              Log In&nbsp;
              <AiFillGithub />
            </button>
          </div>
        </section>

        {/* GET STARTED */}
        <section className="mx-auto w-full max-w-4xl mb-40">
          <h2 className="font-accent hero-text-gradient text-center text-4xl font-bold md:ml-12 md:text-left">
            How it works?
          </h2>
          <Features />
        </section>
      </main>
      <footer
        className={cx([
          "relative grid grid-cols-3 place-items-center max-w-screen-2xl mx-auto",
          "text-sm text-light/40 py-5",
        ])}
      >
        <div>
          <p>EnvHub</p>
        </div>

        <div className="">
          <StarUsOnGitHub />
        </div>

        <div className="flex gap-5">
          <a href="" className="hover:text-light">
            Terms & Conditions
          </a>
          <a href="" className="hover:text-light">
            Privacy Agreement
          </a>
        </div>
      </footer>
      <div className="hero-graphic-blur main" />
      <div className="hero-graphic-blur accent" />
    </>
  );
}
