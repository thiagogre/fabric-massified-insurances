"use client";

import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { query } from "../../config/api";
import ProductInsuranceCard from "../../components/productInsuranceCard/ProductInsuranceCard";
import { products } from "../../mocks/products";
import { Product } from "../partner/types";

const App = () => {
	const router = useRouter();
	const [pendingClaims, setPendingClaims] = useState<any>([]);

	const getProductDataByAsset = (pendingClaim: any): Product | undefined => {
		return products.find(
			(p) => p.id === String(pendingClaim?.CoverageType)
		);
	};

	const handleAnalysis = (id: string) =>
		router.push(`/evidence-analysis/${id}`);

	useEffect(() => {
		(async () => {
			const response = await query({
				channelid: "mychannel",
				chaincodeid: "basic",
				function: "GetAssetsByRichQuery",
				args: [`{"selector":{"ClaimStatus":"Pending"}}`],
			});
			if (response?.success && response?.data?.docs?.length) {
				setPendingClaims(response.data.docs);
			}
		})();
	}, []);

	return (
		<div className="min-h-screen flex flex-col items-center justify-center p-4 bg-gray-50">
			<h1 className="text-3xl font-bold mb-8">Pedidos para análise</h1>
			<div className="w-full max-w-md flex flex-col justify-center gap-4">
				{!!pendingClaims?.length &&
					pendingClaims.map(
						(pc: any, i: number) =>
							!!getProductDataByAsset(pc) && (
								<ProductInsuranceCard
									key={i}
									{...{
										...(getProductDataByAsset(
											pc
										) as Product),
										claimStatus: pc?.ClaimStatus,
										btn: {
											onClick: () =>
												handleAnalysis(pc.ID),
											title: "Analisar",
										},
									}}
								/>
							)
					)}
			</div>
		</div>
	);
};

export default App;
