import { Route, Routes } from "react-router-dom";
import { Home } from "pages/home";
import { Connect } from "pages/connect";
import { NotFound } from "pages/notfound";
import { Dash } from "pages/dash";

export default function App() {
	return (
		<Routes>
			<Route index path="/" element={<Home />} />
			<Route index path="/connect" element={<Connect />} />
			<Route index path="/dash" element={<Dash />} />
			<Route index path="*" element={<NotFound />} />
		</Routes>
	);
}
