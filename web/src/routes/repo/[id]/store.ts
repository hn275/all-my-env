import { writable } from "svelte/store";

export type NewVariable = {
	key: string;
	value: string;
};

export interface Variable extends NewVariable {
	id: string;
	created_at: string;
	updated_at: string;
}

type AppState = {
	deleteVariable?: string;
};

export type RepositoryEnv = {
	write_access: boolean;
	variables: Array<Variable>;
	state: AppState;
};

export const store = writable<RepositoryEnv>({
	write_access: false,
	variables: [],
	state: {},
});
