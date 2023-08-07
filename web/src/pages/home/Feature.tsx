import { nanoid } from "nanoid";
import {  MouseEventHandler, ReactNode,  useState } from "react";
import { Demo} from "./Demo";
import { AnimatePresence, Variants, motion } from "framer-motion";
import cx from "classnames";
import ShareSVG from "./assets/share.svg";

type FeatureType = {
	title: string;
	content: string;
	svg: any;
};
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

export function Features() {
  const [activeTile, setActiveTile] = useState<number>(0);
	const f = FEATURES[activeTile];
	const pages: number[] = [];
	for (let i = 0; i < FEATURES.length; ++i) {
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
				<AnimatePresence mode="wait">
					<FeatureWrapper key={f.title} title={f.title} content={f.content} />
				</AnimatePresence>

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
							<Demo encrypted={activeTile === 1} />
						</motion.div>
					)}
				</AnimatePresence>
			</div>
		</section>
	);
}

interface FeatureWrapperProps {
	key: any;
	content: string;
	title: string;
}
function FeatureWrapper({ content, title, key }: FeatureWrapperProps) {
	const variants: Variants = {
		init: { opacity: 0 },
		animate: { opacity: 1 },
		exit: { opacity: 0 },
	};

	return (
		<motion.div className="h-64 w-full lg:h-48 lg:w-[40ch]">
			<motion.h3
				key={key}
				variants={variants}
				initial="init"
				animate="animate"
				exit="exit"
				className="font-bold text-xl mb-5"
			>
				{title}
			</motion.h3>
			<motion.p
				key={key + 1}
				variants={variants}
				initial="init"
				animate="animate"
				exit="exit"
			>
				{content}
			</motion.p>
		</motion.div>
	);
}

interface PagButtonProps {
	active: boolean;
	children: ReactNode;
	onClick: MouseEventHandler<HTMLButtonElement>;
}

function PagButton({ active, children, onClick }: PagButtonProps) {
	return (
		<button
			onClick={onClick}
			className={cx([
				"h-8 w-8 rounded-lg border transition-all",
				{
					"bg-dark-100 border-dark-100": !active,
					"border-fuchsia-400 bg-fuchsia-400/10 text-fuchsia-400": active,
					"hover:border-fuchsia-400/50 hover:text-fuchsia-400/50": !active,
				},
			])}
		>
			{children}
		</button>
	);
}
