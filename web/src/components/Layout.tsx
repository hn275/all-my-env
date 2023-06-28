import { ReactNode } from "react";

interface Props {
  children: ReactNode;
}
// TODO: global styles go here
export function Layout({ children }: Props) {
  return <div className="bg-dark">{children}</div>;
}
