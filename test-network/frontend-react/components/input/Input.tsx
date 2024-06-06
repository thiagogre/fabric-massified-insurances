import React from "react";

import type { InputProps } from "./types";

const Input = ({
	type = "text",
	placeholder,
	value,
	disabled = false,
	onChange,
}: InputProps) => {
	return (
		<input
			type={type}
			placeholder={placeholder}
			value={value}
			disabled={disabled}
			onChange={onChange}
			className="border border-gray-300 p-2 rounded w-full"
		/>
	);
};

export default Input;
