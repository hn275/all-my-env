import { HTMLAttributes, ReactNode, useEffect, useState } from "react";
import { useSearchParams, useNavigate } from "react-router-dom";
import cx from "classnames";
import { Spinner } from "./Spinner";
import * as jwt from "jose";
import { WEB } from "lib/routes";

const API = import.meta.env.VITE_ENVHUB_API;
const JWT_SECRET: string = import.meta.env.VITE_JWT_SECRET;

interface Response {
  code: string;
}

export function Connect() {
  const { loading, err } = useLogin();
  return (
    <section className="flex justify-center items-center h-[100vh]">
      {err ? (
        <Container>
          <h1 className="text-2xl font-semibold text-main mb-2">Whoops!</h1>
          <p>{err}</p>
        </Container>
      ) : (
        loading && (
          <div className="flex flex-col items-center gap-4">
            <Spinner />
            <p className="text-main text-lg font-semibold">Authenticating...</p>
          </div>
        )
      )}
    </section>
  );
}

interface ContainerProps extends HTMLAttributes<HTMLDivElement> {
  children: ReactNode;
}
function Container({ children, className, ...rest }: ContainerProps) {
  return (
    <div
      {...rest}
      className={cx([className, "p-10 border border-main text-center"])}
    >
      {children}
    </div>
  );
}

function useLogin() {
  const [param, _] = useSearchParams("code");
  const [loading, setLoading] = useState<boolean>(true);
  const [err, setErr] = useState<string>();
  const nav = useNavigate();

  useEffect(() => {
    const githubCode = param.get("code");
    if (!githubCode || githubCode === "") return;

    (async () => {
      setLoading(() => true);
      setErr(() => undefined);
      try {
        const response = await fetch(`${API}/auth`, {
          method: "POST",
          headers: {
            "content-type": "application/json",
          },
          body: JSON.stringify({ code: githubCode }),
        });

        const { status } = response;
        switch (true) {
          case status === 200:
            const { code } = (await response.json()) as Response;
            // TODO: parse header
            // const header = jwt.decodeProtectedHeader(code);
            const secret = new Uint8Array(
              JWT_SECRET.split("").map((x) => x.charCodeAt(0)),
            );
            const claims = await jwt.jwtVerify(code, secret);
            await setSession("users", JSON.stringify(claims));
            nav(WEB.dash);
            return;

          case status === 401:
            setErr(() => "Authentication failed. Try again.");
            return;

          default:
            const err = await response.json();
            const msg = err["error"] as string;
            setErr(() => msg.replace(/_+/g, " "));
            return;
        }
      } catch (e) {
        setErr(() => "Authentication failed, try again later.");
        console.error(e);
      } finally {
        setLoading(() => false);
      }
    })();
  }, []);

  return { loading, err };
}

function setSession(k: string, v: string): Promise<void> {
  return new Promise((resolve) => {
    window.sessionStorage.setItem(k, v);
    resolve();
  });
}
