type InputProps = {
	type: React.HTMLInputTypeAttribute;
	placeholder: string;
	value: string;
	onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
	disabled?: boolean;
};

export type { InputProps };
