import { useEffect } from "react";
import { useParams } from "react-router-dom";

export function Repo() {
  const { id } = useParams();
  useEffect(() => console.log(id), [id]);
  return <>repo page</>;
}
