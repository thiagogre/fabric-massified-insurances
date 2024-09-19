type InputProps = {
	type: React.HTMLInputTypeAttribute;
	placeholder: string;
	value: string;
	onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
	disabled?: boolean;
	required?: boolean;
	name?: string;
};

export type { InputProps };
