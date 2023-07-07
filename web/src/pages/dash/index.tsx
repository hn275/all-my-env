import { Repository } from "lib/github/repositories";
import { useState } from "react";

export function Dash() {
	return <></>;
}

function useRepo() {
	const [data, setData] = useState<Repository[]>([]);
	const [error, setError] = useState<string>();
}
