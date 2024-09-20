"use client";

import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import ProductInsuranceCard from "../../../components/productInsuranceCard/ProductInsuranceCard";
import { query } from "../../../config/api";
import { products } from "../../../mocks/products";
import Button from "../../../components/button/Button";

const App = ({ params }: { params: { username: string } }) => {
	const { username } = params;

	const router = useRouter();

	const [insuredAsset, setInsuredAsset] = useState<any>(null);
	const product = products.find(
		(product) => product.id === String(insuredAsset?.CoverageType)
	);

	const claim = () => {
		router.push(`/insurer/claim/${username}`);
	};

	useEffect(() => {
		const timeout = setTimeout(() => {
			(async () => {
				const response = await query({
					channelid: "mychannel",
					chaincodeid: "basic",
					function: "GetAssetsByRichQuery",
					args: [`{"selector":{"Insured":"${username}"}}`],
				});
				if (response?.success && response?.data?.docs?.length) {
					setInsuredAsset(response.data.docs[0]);
				}
			})();
		}, 1);

		return () => clearTimeout(timeout);
	}, []);

	if (!product || !product.insurance) {
		return <div>Produto não encontrado ou sem informações de seguro.</div>;
	}

	return (
		<div className="max-w-screen-sm mx-auto p-6">
			<h2 className="text-2xl font-semibold mb-2"></h2>
			{!!product && (
				<ProductInsuranceCard
					{...{ ...product, claimStatus: insuredAsset?.ClaimStatus }}
				/>
			)}
			{insuredAsset?.ClaimStatus === "Active" && (
				<div className="mt-6 flex justify-end">
					<Button onClick={claim}>
						<span className="flex items-center">
							Acionar seguro
						</span>
					</Button>
				</div>
			)}
		</div>
	);
};

export default App;
