import { ReactNode } from "react";

interface Props {
	children: ReactNode;
}
export function Layout({ children }: Props) {
	return <div className="max-w-screen-2xl w-full mx-auto">{children}</div>;
}
