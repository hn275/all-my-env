export interface Repository {
	id: number;
	created_at: Date;
	full_name: string;
	url: string;
	variables: Variable[]; // TODO: type this
}

export interface Variable {
	id: string;
	created_at: Date;
	updated_at: Date;
	key: string;
	value: string;
}
