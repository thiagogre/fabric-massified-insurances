"use client";

import React, { useState } from "react";

import Carousel from "../../components/carousel/Carousel";
import ImageTable from "../../components/imageTable/ImageTable";
import Footer from "../../components/footer/Footer";

const Shop = () => {
	const images = [
		"https://via.placeholder.com/600x400",
		"https://via.placeholder.com/600x400/ff7f7f",
		"https://via.placeholder.com/600x400/7f7fff",
	];

	const imageData = [
		{
			name: "Image 1",
			size: "200 KB",
			pixels: "600x400",
			dimensions: "600x400",
		},
		{
			name: "Image 2",
			size: "210 KB",
			pixels: "600x400",
			dimensions: "600x400",
		},
		{
			name: "Image 3",
			size: "220 KB",
			pixels: "600x400",
			dimensions: "600x400",
		},
	];

	return (
		<>
			<h1 className="text-4xl font-bold text-center">Celulares</h1>
			<div className="flex items-center justify-center space-x-8">
				<Carousel images={images} />
				<ImageTable data={imageData} />
			</div>
		</>
	);
};

const Insurance = () => {
	return (
		<>
			<h1 className="text-4xl font-bold text-center">
				Seguro para o celular
			</h1>
			<div className="flex items-center justify-center space-x-8"></div>
		</>
	);
};

const App = () => {
	const [step, setStep] = useState(0);

	const onBack = React.useCallback(() => {
		setStep((prev) => --prev);
	}, []);

	const onNext = React.useCallback(() => {
		setStep((prev) => ++prev);
	}, []);

	return (
		<div className="min-h-screen flex flex-col items-center justify-center space-y-8 p-4">
			{step === 0 && <Shop />}
			{step === 1 && <Insurance />}
			<Footer
				onBack={onBack}
				onNext={onNext}
				currentStep={step}
				totalSteps={2}
			/>
		</div>
	);
};

export default App;
