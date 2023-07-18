import { Route, Routes } from "react-router-dom";
import { Home } from "pages/home";
import { Auth } from "pages/auth";
import { NotFound } from "pages/notfound";
import { Dash } from "pages/dash";
import { Repo } from "pages/repos";
import { WEB } from "lib/routes";
import { Connect } from "pages/connect";

export default function App() {
	return (
		<Routes>
			<Route index path={WEB.home} element={<Home />} />
			<Route path={WEB.auth} element={<Auth />} />
			<Route path={WEB.dash} element={<Dash />} />
			<Route path={WEB.repo}>
				<Route path={`${WEB.repo}/:id`} element={<Repo />} />
				<Route path={`${WEB.repo}/connect`} element={<Connect />} />
			</Route>
			<Route index path="*" element={<NotFound />} />
		</Routes>
	);
}
