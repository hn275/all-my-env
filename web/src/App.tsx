import { Route, Routes } from "react-router-dom";
import { Home } from "pages/home";
import { Connect } from "pages/connect";

export default function App() {
	return (
		<Routes>
			<Route index path="/" element={<Home />} />
			<Route index path="/connect" element={<Connect />} />
		</Routes>
	);
}
