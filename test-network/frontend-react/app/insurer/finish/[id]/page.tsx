"use client";

import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import ProductInsuranceCard from "../../../../components/productInsuranceCard/ProductInsuranceCard";
import { fetchAPI, query } from "../../../../config/api";
import { products } from "../../../../mocks/products";
import Button from "../../../../components/button/Button";

const EvidencesPDFs = ({ insuredAsset }: any) => {
	const [urls, setUrls] = useState<any>(null);

	useEffect(() => {
		const reduceUrlByFilename = (urls: string[]) => {
			return urls.reduce((acc: any, url: string) => {
				const splited = url.split("/");
				acc[splited[splited.length - 1]] = url;
				return acc;
			}, {});
		};

		(async () => {
			const response = await fetchAPI({
				method: "GET",
				endpoint: `/claim/evidence/${insuredAsset.Insured}`,
			});
			if (response?.success && response?.data) {
				setUrls(reduceUrlByFilename(response.data));
			} else {
				alert(response?.message);
			}
		})();
	}, [insuredAsset]);

	return (
		<div>
			<h1 className="text-xl font-semibold mb-2">Evidências</h1>
			<div className="mb-4">
				{urls && (
					<ul>
						{Object.entries(urls).map(([filename, url]) => (
							<li key={url as string}>
								<a
									href={("http://" + url) as string}
									className="text-blue-500 hover:underline"
									target="_blank"
									rel="noopener noreferrer"
								>
									{filename}
								</a>
							</li>
						))}
					</ul>
				)}
			</div>
		</div>
	);
};

const App = ({ params }: { params: { id: string } }) => {
	const { id } = params;

	const router = useRouter();

	const [insuredAsset, setInsuredAsset] = useState<any>(null);
	const product = products.find(
		(product) => product.id === String(insuredAsset?.CoverageType)
	);

	const finish = async (result: boolean) => {
		const response = await fetchAPI({
			method: "POST",
			endpoint: "/claim/finish",
			bodyData: JSON.stringify({
				username: insuredAsset.Insured,
				isApproved: result,
			}),
		});
		if (response?.success) {
			router.back();
		} else {
			alert(response?.message);
		}
	};

	useEffect(() => {
		const timeout = setTimeout(() => {
			(async () => {
				const response = await query({
					channelid: "mychannel",
					chaincodeid: "basic",
					function: "ReadAsset",
					args: [id],
				});
				if (response?.success && response?.data) {
					setInsuredAsset(response.data);
				} else {
					alert(response?.message);
				}
			})();
		}, 1);

		return () => clearTimeout(timeout);
	}, []);

	if (!product || !product.insurance) {
		return <div>Produto não encontrado ou sem informações de seguro.</div>;
	}

	return (
		<div className="flex min-h-screen items-center">
			<div className="flex max-w-screen-lg mx-auto">
				{!!product && (
					<ProductInsuranceCard
						{...{
							...product,
							claimStatus: insuredAsset?.ClaimStatus,
						}}
					/>
				)}
				<div className="p-6">
					{!!insuredAsset && (
						<EvidencesPDFs insuredAsset={insuredAsset} />
					)}
					<div className="mt-6 flex">
						<div className="flex">
							<div className="mr-6">
								<Button onClick={() => finish(true)}>
									<span className="flex items-center">
										Aprovar
									</span>
								</Button>
							</div>
							<Button
								type="secondary"
								onClick={() => finish(false)}
							>
								<span className="flex items-center">
									Reprovar
								</span>
							</Button>
						</div>
					</div>
				</div>
			</div>
		</div>
	);
};

export default App;
