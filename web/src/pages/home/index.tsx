import { oauth } from "lib/auth";
import { Nav } from "./Nav";
import { AiFillGithub } from "react-icons/ai";
import "./style.css";
import cx from "classnames";
import { type FeatureType, Features } from "./Feature";
import { useState } from "react";

const FEATURES: FeatureType[] = [
	{
		title: "First Thing",
		content:
			"Sit velit molestiae nesciunt aperiam ex. Fuga eius sapiente veritatis perspiciatis officia! Deserunt ut exercitationem accusamus quis iure. Deserunt ducimus ab rem voluptates ad cum, obcaecati consequuntur. Magnam nobis officiis.",
		svg: "some svg 1",
	},
	{
		title: "Second Thing",
		content:
			"Dolor dolorem officiis sed reprehenderit eveniet Iusto doloremque hic sint culpa nam illo ipsa minus culpa Cum eligendi atque delectus",
		svg: "some svg 2",
	},
	{
		title: "Third First",
		content:
			"Lorem cumque veniam tenetur corrupti molestiae? In necessitatibus sequi dolore neque sed? Aut explicabo voluptates laudantium culpa ratione? Fuga corrupti!",
		svg: "some svg 3",
	},
];

export function Home() {
	const [activeTile, setActiveTile] = useState<number>(0);

	return (
		<>
			<Nav handleAuth={oauth} githubUrl="" />

			{/* HERO */}
			<main className="relative mx-auto h-max min-h-[calc(100vh-65px)] max-w-7xl overflow-x-hidden">
				<section className="flex h-[calc(100vh-65px)] flex-col items-center justify-center gap-6">
					<h1 className="hero-text-gradient font-accent text-center text-3xl font-bold uppercase md:text-4xl">
						effortless secrets management
					</h1>

					<p className="w-[25ch] text-center font-semibold md:w-max md:text-lg">
						The ultimate{" "}
						<span className="hero-text-underline">open-source</span>
						&nbsp; solution for managing your&nbsp;
						<span className="hero-text-underline">variables</span>.
					</p>

					<div className="flex flex-col items-center justify-center gap-3 md:flex-row md:gap-6">
						<a className="border-main text-main w-36 rounded-md border py-2 text-center font-semibold">
							Get started
						</a>

						<button
							className={cx([
								"border-main bg-main flex items-center justify-center border md:flex-row",
								"h-10 w-36 rounded-md font-semibold hover:no-underline hover:brightness-90",
							])}
							onClick={() => oauth("/auth")}
						>
							Log In&nbsp;
							<AiFillGithub />
						</button>
					</div>
				</section>

				{/* GET STARTED */}
				<section className="mx-auto w-full max-w-4xl mb-96">
					<h2 className="font-accent hero-text-gradient text-center text-4xl font-bold md:ml-12 md:text-left">
						How it works?
					</h2>
					<Features
						features={FEATURES}
						activeTile={activeTile}
						setActiveTile={setActiveTile}
					/>
				</section>
			</main>
			<div className="hero-graphic-blur main" />
			<div className="hero-graphic-blur accent" />
		</>
	);
}
