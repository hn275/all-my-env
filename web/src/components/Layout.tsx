import { HTMLAttributes, ReactNode, useEffect } from "react";
import { useAuth } from "context/auth";
import { AnimatePresence, motion } from "framer-motion";
import cx from "classnames";

interface Props extends HTMLAttributes<HTMLDivElement> {
  children: ReactNode;
}
export function Layout({ children, className, ...rest }: Props) {
  const {
    state: { auth },
  } = useAuth();
  return (
    <>
      <div
        className={cx(["mx-auto w-full max-w-screen-2xl", className])}
        {...rest}
          >
        {children}
      </div>
      <TokenExpired auth={auth} />
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
            "fixed left-0 top-0",
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
