import { useEffect, useState, ChangeEvent } from "react";
import { Api } from "lib/api";
import { Repository } from "./types";
import { Repo } from "./Repo";
import { useAuth } from "context/auth";

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
		<>
			<section>
				<img src={user?.avatar_url} className="w-20 rounded-full" />
				<h1>{user?.name}</h1>
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

					<div>
						<label htmlFor="repo-search">Search bar</label>
						<input name="repo-search" value={search} onChange={handleSearch} />
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
								data
									.filter((d) => d.full_name.includes(search ?? ""))
									.map((repo) => (
										<li key={repo.id} className="border border-main">
											<Repo {...repo} />
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

/* HOOKS */
function useRepos() {
	const { getToken, dispatch } = useAuth();
	// HANDLE QUERIES
	const [page, setPage] = useState<number>(1);
	const [sort, setSort] = useState<Sort>(SortDefault);
	const [show, setShow] = useState<number>(ShowDefault);
	const [data, setData] = useState<Repository[]>();
	const [error, setError] = useState<string>();
	const [loading, setLoading] = useState<boolean>(true);
	useEffect(() => {
		(async () => {
			try {
				setLoading(() => true);
				const url = Api.makeUrl("/repos", { page, sort, show });
				const tok = getToken();
				if (!tok) {
					console.log("tok not found");
					dispatch({ type: "logout" });
					return;
				}
				const response = await fetch(url, {
					method: "GET",
					headers: {
						Accept: "application/json",
						Authorization: `Bearer ${tok}`,
					},
					credentials: "include",
				});
				const { status } = response;
				const payload = await response.json();

				switch (status) {
					case 200:
						setData(() => payload as Repository[]);
						return;
					case 401 | 403:
						return;
					default:
						throw new Error(`status code ${status} `);
				}
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
