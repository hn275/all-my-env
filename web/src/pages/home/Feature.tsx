import { nanoid } from "nanoid";
import { Dispatch, HTMLAttributes, SetStateAction } from "react";
import { Demo, Variables } from "./Demo";
import { AnimatePresence, motion } from "framer-motion";
import cx from "classnames";
import ShareSVG from "./assets/share.svg";

export type FeatureType = {
	title: string;
	content: string;
	svg: any;
};

interface Props {
	features: FeatureType[];
	activeTile: number;
	setActiveTile: Dispatch<SetStateAction<number>>;
}

const DemoVariables: Variables[] = [
	{
		key: "API_KEY",
		value: "api-secret-key",
		encoded: "Ym05M1lYeDV",
	},

	{
		key: "TOKEN",
		value: "my-token",
		encoded: "YXVCZ29ZQm02ZGR2",
	},

	{
		key: "PG_DSN",
		value: "pg://foo:bar@153.531.6.2",
		encoded: "cGc6L2ZvbzpiY",
	},

	{
		key: "AWS_KEY",
		value: "bucket",
		encoded: "c2Q3YWZnZDkwZjg5YThk",
	},

	{
		key: "AWS_SECRET",
		value: "password",
		encoded: "Zjd0c2RmNzk4YTd",
	},

	{
		key: "JWT_SECRET",
		value: "super-secure-secret",
		encoded: "OGZkN2E2c2=",
	},
];

export function Features({ features, activeTile, setActiveTile }: Props) {
	const f = features[activeTile];
	const pages: number[] = [];
	for (let i = 0; i < features.length; ++i) {
		pages.push(i);
	}

	return (
		<section
			className={cx([
				"mx-auto flex w-full flex-col items-center justify-center",
				"gap-8 lg:flex-row lg:gap-24",
			])}
		>
      <div className="mt-5 flex w-3/4 lg:w-1/2 flex-col items-start">
				<div className="h-64 w-full lg:w-[40ch]">
					<h3 className="mb-5 text-xl font-bold">{f.title}</h3>
					<p>{f.content}</p>
				</div>

				<ul className="mt-3 flex items-center justify-start gap-3">
					{pages.map((page) => (
						<li key={nanoid(10)}>
							<PagButton
								active={activeTile === page}
								onClick={() => setActiveTile(() => page)}
							>
								{page + 1}
							</PagButton>
						</li>
					))}
				</ul>
			</div>

			<div className="w-[35ch] lg:w-1/2">
				<AnimatePresence mode="wait">
					{activeTile === 2 ? (
						<motion.div
							className="w-full"
							key="svg"
							initial={{ opacity: 0, x: 50 }}
							animate={{ opacity: 1, x: 0, transition: { type: "tween" } }}
							exit={{ opacity: 0, x: 50, transition: { type: "tween" } }}
						>
							<img
								src={ShareSVG}
								role="presentation"
								alt=""
								className="max-h-52"
							/>
						</motion.div>
					) : (
						<motion.div
							className="h-full"
							key="card"
							initial={{ opacity: 0, x: -50 }}
							animate={{ opacity: 1, x: 0, transition: { type: "tween" } }}
							exit={{ opacity: 0, x: -50, transition: { type: "tween" } }}
						>
							<Demo content={DemoVariables} encrypted={activeTile === 1} />
						</motion.div>
					)}
				</AnimatePresence>
			</div>
		</section>
	);
}

interface PagButtonProps extends HTMLAttributes<HTMLButtonElement> {
	active: boolean;
}

function PagButton({ active, children, ...rest }: PagButtonProps) {
	return (
		<button
			className={cx([
				"h-8 w-8 rounded-lg border transition-all",
				{
					"bg-dark-100 border-dark-100": !active,
					"border-fuchsia-400 bg-fuchsia-400/10 text-fuchsia-400": active,
				},
				"hover:border-fuchsia-400/50 hover:text-fuchsia-400/50",
			])}
			{...rest}
		>
			{children}
		</button>
	);
}
