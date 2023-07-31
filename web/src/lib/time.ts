export function formatUtcTime(t: Date): string {
	const [date, time] = t.toLocaleString().split(",");
	return `${date} at ${time}`;
}
