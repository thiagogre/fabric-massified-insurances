import React from "react";

import type { InputProps } from "./types";

const Input = ({
	type = "text",
	placeholder,
	value,
	disabled = false,
	required = false,
	onChange,
	name = undefined,
}: InputProps) => {
	return (
		<input
			type={type}
			placeholder={placeholder}
			value={value}
			disabled={disabled}
			required={required}
			onChange={onChange}
			name={name}
			className="border border-gray-300 p-2 rounded w-full"
		/>
	);
};

export default Input;
