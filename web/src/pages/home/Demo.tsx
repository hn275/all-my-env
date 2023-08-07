import { AnimatePresence, MotionProps, motion } from "framer-motion";
import { ReactNode } from "react";

export type Variables = {
	key: string;
	value: string;
	encoded: string;
};
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


interface Props {
	encrypted: boolean;
}

export function Demo({ encrypted }: Props) {
	return (
		<ul className="bg-dark-100 flex h-full w-full flex-col gap-2 rounded-lg p-5 lg:w-96">
			{DemoVariables.map(({ key, value, encoded }, i) => (
				<li key={key} className="text-main flex w-96 flex-row font-semibold">
					<p className="font-normal text-white">{key}=</p>
					<AnimatePresence mode="wait">
						{encrypted ? (
							<TextWrapper key={encoded} delay={i}>
								{encoded}
							</TextWrapper>
						) : (
							<TextWrapper key={value} delay={i}>
								{value}
							</TextWrapper>
						)}
					</AnimatePresence>
				</li>
			))}
		</ul>
	);
}

interface TextWrapperProps extends MotionProps {
	children: ReactNode;
	delay: number;
}
function TextWrapper({ delay, children, ...rest }: TextWrapperProps) {
	return (
		<motion.p
			initial={{ y: -10, x: 5, opacity: 0 }}
			animate={{
				y: 0,
				x: 5,
				opacity: 1,
				transition: { delay: delay * 0.01, type: "tween" },
			}}
			exit={{ opacity: 0, transition: { delay: delay * 0.01 } }}
			{...rest}
			className="ml-[-5px] italic"
		>
			"{children}"
		</motion.p>
	);
}
