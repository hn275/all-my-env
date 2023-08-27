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

export type Contributor = {
    id: number;
    name: string;
    avatar_url: string;
}

export type RepositoryEnv = {
	write_access: boolean;
    is_owner: boolean;
	variables: Array<Variable>;
	repoID?: number;
    contributors: Array<Contributor>;
};

export const store = writable<RepositoryEnv>({
	write_access: false,
    is_owner: false,
	variables: [],
    contributors: []
});
