import type { Product } from "../../app/partner/types";
import type { BadgeProps, ClaimStatus } from "./types";

const Badge = ({ backgroundColor, textColor, text }: BadgeProps) => {
	return (
		<div
			className={`flex px-3 rounded-full text-${textColor} font-semibold items-center mb-4`}
			style={{ backgroundColor }}
		>
			<span>{text}</span>
		</div>
	);
};

const ProductInsuranceCard = (
	product: Product & { claimStatus?: ClaimStatus }
) => {
	const { insurance, claimStatus } = product;
	const claimStatusColor = {
		Active: { bgColor: "yellow", textColor: "black", text: "Ativo" },
		Pending: { bgColor: "orange", textColor: "black", text: "Em Análise" },
		Rejected: { bgColor: "red", textColor: "white", text: "Rejeitado" },
		Approved: { bgColor: "green", textColor: "white", text: "Aprovado" },
	};

	return (
		<div className="flex-1 bg-white shadow-md rounded-lg p-6">
			<div className="flex flex-row justify-between">
				<h2 className="text-2xl font-semibold mb-2">{product.name}</h2>
				{!!claimStatus && (
					<Badge
						backgroundColor={claimStatusColor[claimStatus].bgColor}
						textColor={claimStatusColor[claimStatus].textColor}
						text={claimStatusColor[claimStatus].text}
					/>
				)}
			</div>
			<div className="w-full mb-6">
				<img
					src={product.image}
					alt={product.name}
					className="w-full h-auto object-cover rounded-md"
				/>
			</div>
			<div className="space-y-4">
				<p className="text-lg text-gray-700 mb-4">
					<span className="font-medium">
						{insurance.coveredItemDescription}
					</span>
				</p>
				<div className="space-y-2">
					<div className="flex justify-between">
						<span className="font-medium">
							Valor Coberto pelo Seguro:
						</span>
						<span>R$ {insurance.coveredValue.toFixed(2)}</span>
					</div>
					<div className="flex justify-between">
						<span className="font-medium">Tipo de Cobertura:</span>
						<span>{insurance.coverageType}</span>
					</div>
					<div className="flex justify-between">
						<span className="font-medium">Valor do Prêmio:</span>
						<span>
							R$ {insurance.premiumValue.toFixed(2)} por mês
						</span>
					</div>
					<div className="flex justify-between">
						<span className="font-medium">Prazo do Seguro:</span>
						<span>{insurance.coverageDuration} meses</span>
					</div>
				</div>
			</div>
		</div>
	);
};

export default ProductInsuranceCard;
