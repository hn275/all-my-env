import { useEffect, useState, ChangeEvent } from "react";
import { Repository } from "./types";
import { Repo } from "./Repo";
import { useAuth } from "context/auth";
import { makeUrl } from "lib/url";
import { Layout } from "components/Layout";
import { Spinner } from "components/Spinner";

// sort
type Sort = "created" | "updated" | "pushed" | "full_name";
const SortDefault: Sort = "pushed";
const sortFunctions: Record<string, Sort> = {
  Created: "created",
  Pushed: "pushed",
  Updated: "updated",
  Name: "full_name",
};
// page limit
const ShowDefault: number = 30;

export function Dash() {
  const {
    data,
    loading,
    handlePrevPage,
    handleNextPage,
    error,
    page,
    handleSort,
    handleShow,
    search,
    handleSearch,
  } = useRepos();

  const {
    state: { user },
  } = useAuth();

  return (
    <Layout className="px-4 md:px-10">
      <section className="mt-5 px-4">
        <div className="flex flex-col items-center justify-center gap-2 md:w-max">
          <img src={user?.avatar_url} className="w-20 rounded-full" />
          <div>
            <h2 className="inline-block text-lg font-semibold">{user?.name}</h2>&nbsp;
            <span className="text-xs">({user?.login})</span>
          </div>
        </div>
      </section>

      <section className="">
        <h1 className="text-gradient text-2xl font-semibold">Repositories</h1>

        <div className="">
          <div className="grid grid-cols-2">
            <div>
              <label htmlFor="sort">Sort by: </label>
              <select
                name="sort"
                className=""
                onChange={handleSort}
                defaultValue={SortDefault}
              >
                {Object.entries(sortFunctions).map(([name, sort]) => (
                  <option key={sort} value={sort}>
                    {name}
                  </option>
                ))}
              </select>
            </div>

            <div>
              <label htmlFor="show">Show: </label>
              <select
                name="show"
                className=""
                onChange={handleShow}
                defaultValue={ShowDefault}
              >
                {[10, 20, 30].map((i) => (
                  <option key={i} value={i}>
                    {i}
                  </option>
                ))}
              </select>
            </div>
          </div>

          <div className="flex flex-col">
            <label htmlFor="repo-search">Search bar</label>
            <input name="repo-search" value={search} onChange={handleSearch} />
          </div>
        </div>


        <hr className="text-main border-main my-4 rounded-lg border" />


        {loading ? (
          <div className="flex h-64 w-full flex-col items-center justify-center gap-3">
            <Spinner className="stroke-main" />
            <p>Fetching data...</p>
          </div>
        ) : error ? (
          <p>{error}</p>
        ) : (
          <>
            {data.length === 0 ? (
              <p>You don't have any repository yet.</p>
            ) : (
              <>
                <ul className="mt-2 flex flex-col gap-3">
                  {data
                    .filter((d) => d.full_name.includes(search ?? ""))
                    .map((repo) => (
                      <Repo key={repo.id} {...repo} />
                    ))}
                </ul>
                <div className="">
                  <button onClick={handlePrevPage}>prev</button>
                  <p>{page}</p>
                  <button onClick={handleNextPage}>next</button>
                </div>
              </>
            )}
          </>
        )}
      </section>

    </Layout>
  );
}

/* HOOKS */
function useRepos() {
  const { getToken, dispatch } = useAuth();
  // HANDLE QUERIES
  const [page, setPage] = useState<number>(1);
  const [sort, setSort] = useState<Sort>(SortDefault);
  const [show, setShow] = useState<number>(ShowDefault);
  const [data, setData] = useState<Repository[]>([]);
  const [error, setError] = useState<string>();
  const [loading, setLoading] = useState<boolean>(true);
  useEffect(() => {
    const tok = getToken();
    if (!tok) {
      dispatch({ type: "logout" });
      return;
    }

    const url = makeUrl("/repos", { page, sort, show });

    const headers = new Headers({
      Accept: "application/json",
      Authorization: `Bearer ${tok}`,
    });
    (async () => {
      try {
        setLoading(() => true);
        const response = await fetch(url, {
          method: "GET",
          headers,
          credentials: "include",
        });
        const { status } = response;

        if (status === 401 || status === 403) {
          dispatch({ type: "logout" });
          return;
        }

        const payload = await response.json();
        if (status === 200) setData(() => payload as Repository[]);
        else setError(() => payload["message"]);
      } catch (e) {
        setError(() => "Server is not responding, try again later.");
        console.error(e);
      } finally {
        setLoading(() => false);
      }
    })();
  }, [page, show, sort]);

  function handleNextPage() {
    const len = data?.length ?? 0;
    if (len < show) return;
    setPage((p) => p + 1);
  }

  function handlePrevPage() {
    if (page <= show) {
      setPage(() => 1);
      return;
    }
    setPage((p) => p - 1);
  }

  function handleSort(e: ChangeEvent<HTMLSelectElement>) {
    setSort(() => e.target.value as Sort);
  }

  function handleShow(e: ChangeEvent<HTMLSelectElement>) {
    setShow(() => Number(e.target.value as Sort));
  }

  // HANDLE SEARCH
  const [search, setSearch] = useState<string>();
  function handleSearch(e: ChangeEvent<HTMLInputElement>) {
    setSearch(() => e.target.value);
  }

  return {
    data,
    error,
    handleNextPage,
    handlePrevPage,
    page,
    loading,
    show,
    handleSort,
    handleShow,
    search,
    handleSearch,
  };
}
