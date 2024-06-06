import React from "react";

import type { FooterProps } from "./types";

import Button from "../button/Button";

const Footer = ({ onBack, onNext, currentStep, totalSteps }: FooterProps) => {
	return (
		<div className="fixed bottom-0 left-0 right-0 bg-white shadow-lg py-4 flex justify-between items-center px-8 border-t border-gray-300 h-24">
			<div className="flex justify-start items-center">
				{currentStep > 0 && <Button onClick={onBack}>Voltar</Button>}
			</div>
			<div className="flex justify-end items-center">
				{currentStep < totalSteps - 1 && (
					<Button onClick={onNext}>Avançar</Button>
				)}
			</div>
		</div>
	);
};

export default Footer;
