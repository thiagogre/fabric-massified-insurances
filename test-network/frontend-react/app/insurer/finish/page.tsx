"use client";

import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { query } from "../../../config/api";
import ProductInsuranceCard from "../../../components/productInsuranceCard/ProductInsuranceCard";
import { products } from "../../../mocks/products";
import { Product } from "../../partner/types";

const App = () => {
	const router = useRouter();
	const [evidencesApprovedClaims, setEvidencesApprovedClaims] = useState<any>(
		[]
	);

	const getProductDataByAsset = (
		evidencesApprovedClaim: any
	): Product | undefined => {
		return products.find(
			(p) => p.id === String(evidencesApprovedClaim?.CoverageType)
		);
	};

	const handleAnalysis = (id: string) => router.push(`/insurer/finish/${id}`);

	useEffect(() => {
		(async () => {
			const response = await query({
				channelid: "mychannel",
				chaincodeid: "basic",
				function: "GetAssetsByRichQuery",
				args: [`{"selector":{"ClaimStatus":"EvidencesApproved"}}`],
			});
			if (response?.success && response?.data?.docs?.length) {
				setEvidencesApprovedClaims(response.data.docs);
			}
		})();
	}, []);

	return (
		<div className="min-h-screen flex flex-col items-center justify-center p-4 bg-gray-50">
			<h1 className="text-3xl font-bold mb-8">
				Pedidos com evidências aprovadas
			</h1>
			<div className="w-full max-w-md flex flex-col justify-center gap-4">
				{!!evidencesApprovedClaims?.length &&
					evidencesApprovedClaims.map(
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
