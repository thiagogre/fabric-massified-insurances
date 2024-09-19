"use client";
import React, { useMemo, useReducer, useState } from "react";
import { useRouter } from "next/navigation";

import type {
	Product,
	CarouselAction,
	CarouselState,
	CarouselProps,
} from "./types";
import Button from "../../components/button/Button";

import { products } from "../../mocks/products";

const carouselReducer = (
	state: CarouselState,
	action: CarouselAction
): CarouselState => {
	switch (action.type) {
		case "PREV_IMAGE":
			return {
				...state,
				currentIndex:
					state.currentIndex === 0
						? state.images.length - 1
						: state.currentIndex - 1,
			};
		case "NEXT_IMAGE":
			return {
				...state,
				currentIndex:
					state.currentIndex === state.images.length - 1
						? 0
						: state.currentIndex + 1,
			};
		default:
			return state;
	}
};

const Carousel = (props: CarouselProps) => {
	const { products } = props;
	const router = useRouter();

	const initialState: CarouselState = {
		currentIndex: 0,
		images: products.map((product: Product) => product.image),
	};

	const [state, dispatch] = useReducer(carouselReducer, initialState);
	const [showModal, setShowModal] = useState(false);

	const selectedProduct = useMemo(
		() => products[state.currentIndex],
		[products, state.currentIndex]
	);

	const prevImage = () => {
		dispatch({ type: "PREV_IMAGE" });
	};

	const nextImage = () => {
		dispatch({ type: "NEXT_IMAGE" });
	};

	const handlePurchaseClick = () => {
		setShowModal(true);
	};

	const handleInsuranceDecision = (purchaseInsurance: boolean) => {
		setShowModal(false);
		if (purchaseInsurance) {
			router.push(`/partner/insurance/${selectedProduct.id}`);
		}
	};

	return (
		<div className="w-full flex items-center">
			<div className="relative w-full h-full">
				<div className="overflow-hidden rounded-lg">
					<div className="overflow-hidden rounded-lg w-full h-full">
						<div className="relative w-full h-full flex items-center justify-center">
							{state.images.map((image, index) => (
								<img
									key={index}
									src={image}
									alt={`Slide ${index}`}
									className={`w-auto h-full object-contain transition-opacity duration-500 ${
										index === state.currentIndex
											? "opacity-100"
											: "opacity-0 absolute top-0 left-0"
									}`}
									style={{ transition: "opacity 0.5s" }}
								/>
							))}
						</div>
					</div>
				</div>
				<button
					className={`absolute top-1/2 transform -translate-y-1/2 left-0 text-gray-700 text-3xl p-2 transition-transform duration-300`}
					onClick={prevImage}
				>
					<span
						className="text-gray-700"
						style={{ fontWeight: "bold", fontSize: "40px" }}
					>
						&lt;
					</span>
				</button>
				<button
					className={`absolute top-1/2 transform -translate-y-1/2 right-0 text-gray-700 text-3xl p-2 transition-transform duration-300`}
					onClick={nextImage}
				>
					<span
						className="text-gray-700"
						style={{ fontWeight: "bold", fontSize: "40px" }}
					>
						&gt;
					</span>
				</button>
			</div>

			<div className="w-[30%] pl-8">
				<h2 className="text-3xl font-bold">{selectedProduct.name}</h2>
				<p className="text-lg text-gray-600 mt-2">
					Marca: {selectedProduct.brand}
				</p>
				<p className="text-xl font-semibold mt-4">
					Preço: R$ {selectedProduct.price}
				</p>
				<p className="mt-6 text-gray-700">
					Este {selectedProduct.name} é um smartphone de alta
					qualidade, com recursos avançados e desempenho superior,
					produzido pela {selectedProduct.brand}.
				</p>
				<div className="mt-6">
					<Button onClick={handlePurchaseClick}>Comprar Agora</Button>
				</div>
			</div>

			{showModal && (
				<div className="fixed inset-0 flex items-center justify-center z-50 bg-black bg-opacity-50">
					<div className="bg-white p-6 rounded-lg shadow-lg">
						<h3 className="text-xl font-semibold mb-4">
							Deseja adquirir um seguro para este produto?
						</h3>
						<div className="flex justify-end space-x-4">
							<Button
								type="secondary"
								onClick={() => handleInsuranceDecision(false)}
							>
								Não
							</Button>
							<Button
								onClick={() => handleInsuranceDecision(true)}
							>
								Sim
							</Button>
						</div>
					</div>
				</div>
			)}
		</div>
	);
};

const Shop = () => {
	return (
		<div className="flex items-center justify-center space-x-8">
			<Carousel products={products} />
		</div>
	);
};

const App = () => {
	return (
		<div className="min-h-screen flex flex-col items-center justify-center space-y-8 p-4">
			<Shop />
		</div>
	);
};

export default App;
