import { AnimatePresence, MotionProps, motion } from "framer-motion";
import { ReactNode } from "react";

export type Variables = {
	key: string;
	value: string;
	encoded: string;
};

interface Props {
	content: Variables[];
	encrypted: boolean;
}

export function Demo({ content, encrypted }: Props) {
	return (
		<ul className="bg-dark-100 flex h-full w-full flex-col gap-2 rounded-lg p-5 lg:w-96">
			{content.map(({ key, value, encoded }, i) => (
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
