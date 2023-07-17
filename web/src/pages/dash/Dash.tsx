import { useEffect, useState, ChangeEvent } from "react";
import { Github } from "lib/github/request";
import { Fetch } from "lib/api";
import { Repository } from "./types";
import { Repo } from './Repo'

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
  } = useRepos();

  const user = Github.getUser();

  return (
    <>
      <section>
        <img src={user?.payload.avatar_url} className="w-20 rounded-full" />
        <h1>{user?.payload.name}</h1>
        <p>NOTE: I can display more information as well</p>
      </section>

      <section className="">
        <h2>Repositories:</h2>

        <div className="">
          <div>
            <label htmlFor="sort">Sort by: </label>
            <select
              name="sort"
              className=""
              onChange={handleSort}
              defaultValue={SortDefault}
            >
              {Object.entries(sortFunctions).map(([name, sort]) => (
                <option key={sort} value={sort}>{name}</option>
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
                <option key={i} value={i}>{i}</option>
              ))}
            </select>
          </div>
        </div>

        {loading ? (
          <p>Loading</p>
        ) : error ? (
          <p>{error}</p>
        ) : data ? (
          <>
            <ul>
              {data.length === 0 ? (
                <li>
                  <p>You don't have any repository yet.</p>
                </li>
              ) : (
                data.map((repo) => (
                  <li key={repo.id} className="border border-main">
                    <Repo isOwner={user?.payload.id === repo.owner.id} repo={repo} />
                  </li>
                ))
              )}
            </ul>
          </>
        ) : (
          <> bruh</>
        )}
      </section>

      <div className="">
        <button onClick={handlePrevPage}>prev</button>
        <p>{page}</p>
        <button onClick={handleNextPage}>next</button>
      </div>
    </>
  );
}

function useRepos() {
  const [page, setPage] = useState<number>(1);
  const [sort, setSort] = useState<Sort>(SortDefault);
  const [show, setShow] = useState<number>(ShowDefault);
  const [data, setData] = useState<Repository[]>();
  const [error, setError] = useState<string>();
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    async function f() {
      try {
        setLoading(() => true);
        const param = { page: String(page), sort, show: String(show) };
        const response = await Fetch.GET("/repos/all", param);
        const { status } = response;
        const payload = await response.json();

        switch (status) {
          case 200:
            setData(() => payload as Repository[]);
            return;
          case 401 | 403:
            // TODO: redirect to home page
            return;
          default:
            throw new Error(`status code ${status}`);
        }
      } catch (e) {
        setError(() => "Github is not responding, try again later.");
        console.error(e);
      } finally {
        setLoading(() => false);
      }
    }

    f();
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
  };
}
