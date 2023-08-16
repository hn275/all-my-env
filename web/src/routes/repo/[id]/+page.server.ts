type Params = {
	params: {
		id: string;
	};
};

export type Route = {
	id: number;
};

export function load({ params }: Params): Route {
	const id = parseInt(params.id);
	return { id };
}
