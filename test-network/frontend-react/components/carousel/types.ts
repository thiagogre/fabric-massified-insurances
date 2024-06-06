type CarouselProps = {
	images: string[];
};

type CarouselState = {
	currentIndex: number;
	images: string[];
};

type CarouselAction = { type: "PREV_IMAGE" } | { type: "NEXT_IMAGE" };

export type { CarouselProps, CarouselAction, CarouselState };
