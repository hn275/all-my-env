import { HTMLAttributes } from "react";
import "./spinner.css";
interface SpinnerProps extends HTMLAttributes<HTMLOrSVGElement> {

}
export function Spinner(props: SpinnerProps) {
  const {id: _} = props
  return (
    <svg viewBox="25 25 50 50" id="spinner" {...props}>
      <circle r="20" cy="50" cx="50" />
    </svg>
  );
}
