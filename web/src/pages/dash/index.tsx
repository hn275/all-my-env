import { useEffect, useState, ChangeEvent } from "react";
import { Github } from "lib/github/request";
import cx from "classnames";

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
			</section>

			<section className="relative m-5 flex h-[50vh] flex-col">
				<h2>Repositories:</h2>

				<div className="ml-auto flex w-max items-center gap-3">
					<div>
						<label htmlFor="sort">Sort by: </label>
						<select
							name="sort"
							className="bg-main text-dark"
							onChange={handleSort}
							defaultValue={SortDefault}
						>
							{Object.entries(sortFunctions).map(([name, sort]) => (
								<option value={sort}>{name}</option>
							))}
						</select>
					</div>

					<div>
						<label htmlFor="show">Show: </label>
						<select
							name="show"
							className="bg-main text-dark"
							onChange={handleShow}
							defaultValue={ShowDefault}
						>
							{[10, 20, 30].map((i) => (
								<option value={i}>{i}</option>
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
						<ul
							className={cx([
								"flex-grow overflow-y-scroll",
								"rounded-sm border border-main p-2",
							])}
						>
							{data.length === 0 ? (
								<li className="absolute left-0 top-0 flex h-full w-full items-center justify-center">
									<p>You don't have any repository yet.</p>
								</li>
							) : (
								data.map((repo) => (
									<li className="my-3 flex items-center gap-3">
										<img
											src={repo.owner.avatar_url}
											role="presentation"
											alt={repo.owner.login}
											className="aspect-auto w-8 rounded-full"
										/>
										<a href={`https://www.github.com/${repo.full_name}`}>
											{repo.name}
										</a>
									</li>
								))
							)}
						</ul>
					</>
				) : (
					<> bruh</>
				)}
			</section>

			<div className="mx-auto flex justify-center gap-3">
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
				const param = { page: String(page), sort, per_page: String(show) };
				const response = await Github.GET("/user/repos", param);
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

interface Repository {
	id: number;
	node_id: string;
	name: string;
	full_name: string;
	owner: Owner;
	private: boolean;
	html_url: string;
	description: string;
	fork: boolean;
	url: string;
	archive_url: string;
	assignees_url: string;
	blobs_url: string;
	branches_url: string;
	collaborators_url: string;
	comments_url: string;
	commits_url: string;
	compare_url: string;
	contents_url: string;
	contributors_url: string;
	deployments_url: string;
	downloads_url: string;
	events_url: string;
	forks_url: string;
	git_commits_url: string;
	git_refs_url: string;
	git_tags_url: string;
	git_url: string;
	issue_comment_url: string;
	issue_events_url: string;
	issues_url: string;
	keys_url: string;
	labels_url: string;
	languages_url: string;
	merges_url: string;
	milestones_url: string;
	notifications_url: string;
	pulls_url: string;
	releases_url: string;
	ssh_url: string;
	stargazers_url: string;
	statuses_url: string;
	subscribers_url: string;
	subscription_url: string;
	tags_url: string;
	teams_url: string;
	trees_url: string;
	clone_url: string;
	mirror_url: string;
	hooks_url: string;
	svn_url: string;
	homepage: string;
	language: null;
	forks_count: number;
	stargazers_count: number;
	watchers_count: number;
	size: number;
	default_branch: string;
	open_issues_count: number;
	is_template: boolean;
	topics: string[];
	has_issues: boolean;
	has_projects: boolean;
	has_wiki: boolean;
	has_pages: boolean;
	has_downloads: boolean;
	archived: boolean;
	disabled: boolean;
	visibility: string;
	pushed_at: Date;
	created_at: Date;
	updated_at: Date;
	permissions: Permissions;
	allow_rebase_merge: boolean;
	template_repository: null;
	temp_clone_token: string;
	allow_squash_merge: boolean;
	allow_auto_merge: boolean;
	delete_branch_on_merge: boolean;
	allow_merge_commit: boolean;
	subscribers_count: number;
	network_count: number;
	license: License;
	forks: number;
	open_issues: number;
	watchers: number;
}

interface License {
	key: string;
	name: string;
	url: string;
	spdx_id: string;
	node_id: string;
	html_url: string;
}

interface Owner {
	login: string;
	id: number;
	node_id: string;
	avatar_url: string;
	gravatar_id: string;
	url: string;
	html_url: string;
	followers_url: string;
	following_url: string;
	gists_url: string;
	starred_url: string;
	subscriptions_url: string;
	organizations_url: string;
	repos_url: string;
	events_url: string;
	received_events_url: string;
	type: string;
	site_admin: boolean;
}

interface Permissions {
	admin: boolean;
	push: boolean;
	pull: boolean;
}
