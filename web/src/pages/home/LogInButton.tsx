import { useAuth } from "context/auth";
import { oauthHref } from "lib/url";
import { AiFillGithub } from "react-icons/ai";
import { Link } from "react-router-dom";

export function LogInButton() {
  const {
    state: { auth },
  } = useAuth();

  return (
    <Link
      className="btn btn-primary flex w-[14ch] items-center justify-center"
      to={auth ? "/dash" : oauthHref()}
    >
      {auth ? (
        "Dashboard"
      ) : (
        <>
          <span>
            <AiFillGithub />
          </span>
          &nbsp; Sign in
        </>
      )}
    </Link>
  );
}
