import React from "react";

import type { ButtonProps } from "./types";

const Button = ({
	children,
	onClick,
	disabled = false,
	type = "primary",
}: ButtonProps) => {
	return type === "primary" ? (
		<button
			type="button"
			disabled={disabled}
			onClick={onClick}
			className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-md"
		>
			{children}
		</button>
	) : (
		<button
			type="button"
			disabled={disabled}
			onClick={onClick}
			className="px-4 py-2 bg-gray-200 hover:bg-gray-400 text-black rounded-md"
		>
			{children}
		</button>
	);
};

export default Button;
