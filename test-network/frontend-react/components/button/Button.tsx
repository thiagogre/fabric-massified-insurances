import React from "react";

import type { ButtonProps } from "./types";

const Button = ({ children, onClick, disabled = false }: ButtonProps) => {
	return (
		<button
			disabled={disabled}
			onClick={onClick}
			className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
		>
			{children}
		</button>
	);
};

export default Button;
