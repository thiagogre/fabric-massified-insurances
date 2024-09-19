"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";

import { products } from "../../../../mocks/products";
import Button from "../../../../components/button/Button";
import Input from "../../../../components/input/Input";
import SpinLoading from "../../../../components/loading/Loading";
import { fetchAPI, invoke } from "../../../../config/api";
import { uniqueId } from "../../../../utils/uniqueId";

const App = ({ params }: { params: { id: string } }) => {
	const { id } = params;
	const product = products.find((product) => product.id === id);
	if (!product || !product.insurance) {
		return <div>Produto não encontrado ou sem informações de seguro.</div>;
	}

	const { insurance } = product;
	const router = useRouter();

	const defaultFormDataState = { insured: "" };

	const [formData, setFormData] = useState(defaultFormDataState);
	const [showModal, setShowModal] = useState(false);
	const [btnLoading, setBtnLoading] = useState(false);
	const [identity, setIdentity] = useState({ username: "", password: "" });

	const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
		setFormData({ ...formData, [e.target.name]: e.target.value });
	};

	const confirm = async () => {
		setBtnLoading(true);

		try {
			const response = await fetchAPI({
				method: "POST",
				endpoint: "/identity",
				bodyData: {},
			});
			if (response?.success && response?.data) {
				const uniqueID = String(uniqueId());
				await invoke({
					channelid: "mychannel",
					chaincodeid: "basic",
					function: "CreateAsset",
					args: [
						uniqueID,
						response.data.username,
						String(product.insurance?.coverageDuration),
						String(product.insurance?.coveredValue),
						String(0),
						"Varejista",
						String(product.insurance?.premiumValue),
					],
				});

				setIdentity(response.data);
				setShowModal(true);
			}
		} catch (err) {
			console.error(err);
		}

		setBtnLoading(false);
	};

	const finish = () => {
		setShowModal(false);
		setFormData(defaultFormDataState);
		router.push("/");
	};

	return (
		<div className="max-w-screen-lg mx-auto p-6">
			<h1 className="text-3xl font-bold mb-6">Contratar Seguro</h1>
			<div className="flex gap-8">
				<div className="flex-1 bg-white shadow-md rounded-lg p-6">
					<h2 className="text-2xl font-semibold mb-2">
						{product.name}
					</h2>
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
								<span>
									R$ {insurance.coveredValue.toFixed(2)}
								</span>
							</div>
							<div className="flex justify-between">
								<span className="font-medium">
									Tipo de Cobertura:
								</span>
								<span>{insurance.coverageType}</span>
							</div>
							<div className="flex justify-between">
								<span className="font-medium">
									Valor do Prêmio:
								</span>
								<span>
									R$ {insurance.premiumValue.toFixed(2)} por
									mês
								</span>
							</div>
							<div className="flex justify-between">
								<span className="font-medium">
									Prazo do Seguro:
								</span>
								<span>{insurance.coverageDuration} meses</span>
							</div>
						</div>
					</div>
				</div>

				<div className="flex-1 p-6">
					<h2 className="text-2xl font-semibold mb-2">
						Informações pessoais
					</h2>
					<form>
						<div className="mb-4">
							<label
								htmlFor="name"
								className="block text-lg font-semibold"
							>
								Documento de identificação
							</label>
							<Input
								placeholder=""
								type="text"
								value={formData.insured}
								name="insured"
								onChange={handleInputChange}
								required
							/>
						</div>

						<div className="mt-6">
							<Button onClick={confirm} disabled={btnLoading}>
								<span className="flex items-center">
									Finalizar {btnLoading && <SpinLoading />}
								</span>
							</Button>
						</div>
					</form>
				</div>
			</div>
			{showModal && (
				<div className="fixed inset-0 flex items-center justify-center z-50 bg-black bg-opacity-50">
					<div className="bg-white p-6 rounded-lg shadow-lg">
						<h3 className="text-xl font-semibold mb-4">
							Usuário: {identity.username}
							<br />
							Senha: {identity.password}
						</h3>
						<div className="flex justify-end space-x-4">
							<Button onClick={finish}>
								Já salvei minhas credenciais
							</Button>
						</div>
					</div>
				</div>
			)}
		</div>
	);
};

export default App;
