import { Product } from "../app/partner/types";

const products: Product[] = [
	{
		id: "0",
		name: "iPhone 13",
		brand: "Apple",
		price: 999,
		image: "/iphone.png",
		insurance: {
			coveredValue: 1200,
			coverageType: "Contra furto e roubo",
			coverageDuration: 12,
			coveredItemDescription: "Smartphone de alta tecnologia da Apple.",
			premiumValue: 60,
		},
	},
	{
		id: "1",
		name: "Galaxy S21",
		brand: "Samsung",
		price: 799,
		image: "/galaxy.png",
		insurance: {
			coveredValue: 1000,
			coverageType: "Contra furto e roubo",
			coveredItemDescription: "Smartphone avançado da Samsung.",
			premiumValue: 55,
			coverageDuration: 12,
		},
	},
	{
		id: "2",
		name: "Pixel 6",
		brand: "Google",
		price: 699,
		image: "/pixel.png",
		insurance: {
			coveredValue: 900,
			coverageType: "Contra furto e roubo",
			coveredItemDescription:
				"Smartphone com recursos avançados da Google.",
			premiumValue: 50,
			coverageDuration: 12,
		},
	},
];

export { products };
