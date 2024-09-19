type ButtonProps = {
	children: React.ReactNode;
	onClick: () => void;
	disabled?: boolean;
	type?: "primary" | "secondary";
};

export type { ButtonProps };
