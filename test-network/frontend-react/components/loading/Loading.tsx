import React from "react";

const SpinLoading = ({
	type = "primary",
}: {
	type?: "primary" | "secondary";
}) => {
	return (
		<div className="flex justify-center items-center">
			<div
				className={`animate-spin rounded-full h-4 w-4 ml-4 border-t-2 ${
					type === "primary" ? "border-white" : "border-black"
				} border-solid border-opacity-50`}
			></div>
		</div>
	);
};

export default SpinLoading;
