import React, { useReducer } from "react";

import type { CarouselProps, CarouselAction, CarouselState } from "./types";

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

const Carousel = ({ images }: CarouselProps) => {
	const initialState: CarouselState = {
		currentIndex: 0,
		images: images,
	};

	const [state, dispatch] = useReducer(carouselReducer, initialState);

	const prevImage = () => {
		dispatch({ type: "PREV_IMAGE" });
	};

	const nextImage = () => {
		dispatch({ type: "NEXT_IMAGE" });
	};

	return (
		<div className="relative w-full max-w-lg mx-auto">
			<div className="overflow-hidden rounded-lg">
				<div className="relative">
					{state.images.map((image, index) => (
						<img
							key={index}
							src={image}
							alt={`Slide ${index}`}
							className={`w-full transition-opacity duration-500 ${
								index === state.currentIndex
									? "opacity-100"
									: "opacity-0 absolute top-0 left-0"
							}`}
							style={{ transition: "opacity 0.5s" }}
						/>
					))}
				</div>
			</div>
			<button
				className={`absolute top-1/2 transform -translate-y-1/2 left-4 text-gray-700 text-3xl p-2 transition-transform duration-300`}
				onClick={prevImage}
			>
				&lt;
			</button>
			<button
				className={`absolute top-1/2 transform -translate-y-1/2 right-4 text-gray-700 text-3xl p-2 transition-transform duration-300`}
				onClick={nextImage}
			>
				&gt;
			</button>
		</div>
	);
};

export default Carousel;
