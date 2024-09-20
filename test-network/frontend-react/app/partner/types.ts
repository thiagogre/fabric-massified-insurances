type InsuranceData = {
	coveredValue: number;
	coverageType: string;
	coveredItemDescription: string;
	premiumValue: number;
	coverageDuration: number;
};

type Product = {
	id: string;
	name: string;
	brand: string;
	price: number;
	image: string;
	insurance: InsuranceData;
};

type CarouselProps = {
	products: Product[];
};

type CarouselState = {
	currentIndex: number;
	images: string[];
};

type CarouselAction = { type: "PREV_IMAGE" } | { type: "NEXT_IMAGE" };

export type { Product, CarouselProps, CarouselAction, CarouselState };
