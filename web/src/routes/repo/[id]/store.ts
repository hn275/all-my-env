import { writable } from "svelte/store";

export type Variable = {
	id: string;
	created_at: Date;
	updated_at: Date;
	key: string;
	value: string;
};

export const store = writable<Array<Variable>>([]);
