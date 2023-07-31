import { API, Fetch } from "lib/api";
import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { Repository } from "./type";
import { WEB } from "lib/routes";
import { AiFillGithub } from "react-icons/ai";
import { formatUtcTime } from "lib/time";

export function Repo() {
	const nav = useNavigate();
	const { id } = useParams();
	if (!id) {
		nav(-1);
		return;
	}
	const { data, loading, error, fetchVariables } = useFetchVariables(id);
	const newVarHook = useNewVariable(id, fetchVariables);

	return (
		<>
			<section>
				<div>
					<p>Repository:</p>
					{data?.full_name && <h1>{data.full_name}</h1>}
				</div>
				<p>
					Linked at:&nbsp;
					{data?.created_at && formatUtcTime(new Date(data.created_at))}
				</p>

				<a href={data?.url} target="_blank">
					Link to Github&nbsp;
					<span className="inline-block">
						<AiFillGithub />
					</span>
				</a>
			</section>

			{loading ? (
				<p>Loading</p>
			) : error ? (
				<p>{error}</p>
			) : (
				data && (
					<section>
						<h2>Variables:</h2>

						<button className="bg-main" onClick={newVarHook.toggleNewVarModal}>
							Add new variable
						</button>

						{data.variables.length === 0 ? (
							<p>No variables stored</p>
						) : (
							<ul>
								{data.variables.map(
									({ id, created_at, updated_at, key, value }) => (
										<li key={id} className="m-5">
											<p>key: {key}</p>
											<p>value: {value}</p>
											<p>created at: {formatUtcTime(new Date(created_at))}</p>
											<p>updated at: {formatUtcTime(new Date(updated_at))}</p>
										</li>
									),
								)}
							</ul>
						)}
					</section>
				)
			)}

			{newVarHook.newVarModal && (
				<>
					<section>
						<h2>New Variable</h2>
						<form onSubmit={newVarHook.onSubmit}>
							<div>
								<label htmlFor="key">Key</label>
								<input
									id="key"
									name="key"
									type="text"
									value={newVarHook.key}
									onChange={(e) => newVarHook.setKey(() => e.target.value)}
								/>
							</div>

							<div>
								<label htmlFor="value">Value</label>
								<input
									id="value"
									name="value"
									type="text"
									value={newVarHook.value}
									onChange={(e) => newVarHook.setValue(() => e.target.value)}
								/>
							</div>
							<button type="submit">Submit</button>
						</form>
					</section>
				</>
			)}
		</>
	);
}

type NewVarProps = {
	newVarModal: boolean;
	toggleNewVarModal: () => void;
	key: string;
	setKey: React.Dispatch<React.SetStateAction<string>>;
	value: string;
	setValue: React.Dispatch<React.SetStateAction<string>>;
	onSubmit: React.FormEventHandler;
	loading: boolean;
	error: string | undefined;
};

function useNewVariable(
	repoID: string,
	refresh: () => Promise<void>,
): NewVarProps {
	const [newVarModal, setNewVarModal] = useState<boolean>(false);

	const [key, setKey] = useState<string>("");
	const [value, setValue] = useState<string>("");

	const [error, setError] = useState<string | undefined>(undefined);
	const [loading, setLoading] = useState<boolean>(false);

	function resetFields() {
		setValue(() => "");
		setKey(() => "");
	}

	function toggleNewVarModal() {
		setNewVarModal((o) => !o);
		if (newVarModal) resetFields();
	}

	async function onSubmit(e: React.ChangeEvent<HTMLFormElement>) {
		e.preventDefault();
		try {
			setLoading(() => true);
			const url = `${API.repo.index}/${repoID}/variables/new`;
			const query = null;
			const body = { key, value };
			const res = await Fetch.POST(url, query, body);
			if (res.status !== 201) {
				const data = (await res.json()) as ErrorResponse;
				setError(() => data.error);
				return;
			}
			refresh();
			setNewVarModal(() => false);
			resetFields();
		} catch (e) {
			console.error(e);
			setError(() => "Server failed to save a variable, try again later.");
		} finally {
			setLoading(() => false);
		}
	}

	return {
		newVarModal,
		toggleNewVarModal,
		key,
		setKey,
		value,
		setValue,
		error,
		loading,
		onSubmit,
	};
}

type FetchVariableProps = {
	data: Repository | undefined;
	loading: boolean;
	error: string | undefined;
	fetchVariables: () => Promise<void>;
};
function useFetchVariables(repoID: string): FetchVariableProps {
	const [data, setData] = useState<Repository | undefined>();
	const [loading, setLoading] = useState<boolean>(true);
	const [error, setError] = useState<string | undefined>(undefined);
	const nav = useNavigate();

	async function fetchVariables() {
		try {
			const res = await Fetch.GET(`/repo/${repoID}/variables`);
			const payload = await res.json();
			switch (res.status) {
				case 200:
					setData(() => payload as Repository);
					return;

				case 403 | 401:
					nav(WEB.home);
					return;

				default:
					setError(() => (payload as ErrorResponse).error);
					return;
			}
		} catch (e) {
			console.error(e);
			setError(() => "Something went wrong, try again later.");
		} finally {
			setLoading(() => false);
		}
	}
	useEffect(() => {
		fetchVariables();
	}, []);

	return {
		data,
		loading,
		error,
		fetchVariables,
	};
}
