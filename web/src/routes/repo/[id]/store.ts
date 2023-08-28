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
    login: string;
    avatar_url: string;
}

export type RepositoryEnv = {
    write_access: boolean;
    owner_id: number;
    is_owner: boolean;
    variables: Array<Variable>;
    repoID?: number;
    contributors: Array<Contributor>;
};

export const store = writable<RepositoryEnv>({
    write_access: false,
    owner_id: 0,
    variables: [],
    contributors: [],
    is_owner: false
});
