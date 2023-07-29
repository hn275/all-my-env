import { Github } from "lib/github/request";
import { WEB } from "lib/routes";
import { useEffect, useState } from "react";
import { useSearchParams, useNavigate } from "react-router-dom";

export function Connect() {
  const { error, loading } = useConnectRepo()
  return <>
    {
      loading ? (
        <>
          <p>loading</p>
        </>
      ) : error ? (
        <>{error}</>
      ) : (
        <>
          <p>redirect</p>
        </>
      )
    }
  </>;
}

type ConnectRepoHook = { loading: boolean, error: string | undefined }
function useConnectRepo(): ConnectRepoHook {
  const [param] = useSearchParams()
  const repoNameEncoded = param.get("repo_name") as string
  const nav = useNavigate()

  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string>();

  useEffect(() => {
    (async () => {
      try {
        if (!repoNameEncoded) {
          setError(() => "Repository name not found.")
          return
        }

        const body = { "full_name": decodeURIComponent(repoNameEncoded) }
        const res = await Github.POST("/repos/link", null, body)

        switch (res.status) {
          case 201:
            setError(() => undefined)
            return

          case 403 | 401:
            nav(WEB.home)
            return

          default:
            type ErrResponse = { error: string }
            const err = await res.json() as ErrResponse
            setError(() => err.error)
            return
        }
      } catch (e) {
        setError(() => "Something went wrong, try again later.")
        console.error(e)
      } finally {
        setLoading(() => false)
      }
      //
    })()

  }, [])

  return { loading, error }
}
