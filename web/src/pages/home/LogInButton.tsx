import { useAuth } from "context/auth";
import { oauthHref } from "lib/url";
import { AiFillGithub } from "react-icons/ai";
import { Link } from "react-router-dom";

interface Props {
  to?: string;
}
export function LogInButton({ to }: Props) {
  const {
    state: { auth },
  } = useAuth();

  return (
    <Link
      className="btn btn-primary flex w-[14ch] items-center justify-center"

      to={auth ? "/dash" : oauthHref(to)}
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
