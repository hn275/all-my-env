import { Route, Routes } from "react-router-dom";
import { Home } from "pages/home";
import { Connect } from "pages/connect";
import { NotFound } from "pages/notfound";
import { Dash } from "pages/dash";
import { Repo } from "pages/repos";

export default function App() {
	return (
		<Routes>
			<Route index path="/" element={<Home />} />
			<Route path="/connect" element={<Connect />} />
			<Route path="/dash" element={<Dash />} />
			<Route path="/repos">
				<Route path="/repos/:id" element={<Repo />} />
			</Route>
			<Route index path="*" element={<NotFound />} />
		</Routes>
	);
}
