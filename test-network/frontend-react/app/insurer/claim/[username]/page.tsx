"use client";
import React, { useState, useCallback } from "react";
import { useDropzone } from "react-dropzone";
import { fetchAPI } from "../../../../config/api";
import Button from "../../../../components/button/Button";
import { useRouter } from "next/navigation";
import SpinLoading from "../../../../components/loading/Loading";

const App = ({ params }: { params: { username: string } }) => {
	const { username } = params;
	const router = useRouter();

	const [description, setDescription] = useState("");
	const [selectedFiles, setSelectedFiles] = useState<File[]>([]);
	const [btnLoading, setBtnLoading] = useState(false);

	const onDrop = useCallback((acceptedFiles: File[]) => {
		setSelectedFiles((prevFiles) => [...prevFiles, ...acceptedFiles]);
	}, []);

	const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
		const files = e.target.files;
		if (files) {
			setSelectedFiles((prevFiles) => [
				...prevFiles,
				...Array.from(files),
			]);
		}
	};

	const handleRemoveFile = (index: number) => {
		setSelectedFiles((prevFiles) =>
			prevFiles.filter((_, i) => i !== index)
		);
	};

	const { getRootProps, getInputProps, isDragActive } = useDropzone({
		onDrop,
		accept: { "application/pdf": [] },
		multiple: true,
	});

	const handleFormSubmit = async () => {
		setBtnLoading(true);
		if (selectedFiles.length === 0) {
			alert("Selecione pelo menos um arquivo.");
			return;
		}

		const formData = new FormData();
		selectedFiles.forEach((file) => {
			formData.append("files", file);
		});
		formData.append("username", username);

		const response = await fetchAPI({
			method: "POST",
			endpoint: "/smartcontract/claim",
			bodyData: formData,
			headers: {},
		});

		if (response?.success) {
			alert(response?.data);

			router.back();
		} else {
			alert(response?.message);
		}

		setBtnLoading(false);
	};

	return (
		<div className="max-w-lg mx-auto p-6">
			<h1 className="text-2xl font-bold mb-4">Envio de Evidências</h1>

			<form>
				<div className="mb-4">
					<label
						htmlFor="description"
						className="block text-lg font-semibold mb-2"
					>
						Descreva o evento
					</label>
					<textarea
						id="description"
						rows={10}
						name="description"
						className="w-full p-2 border rounded"
						value={description}
						onChange={(e) => setDescription(e.target.value)}
						required
					/>
				</div>

				<div
					{...getRootProps()}
					className={`p-6 border-2 border-dashed rounded-md cursor-pointer transition-colors ${
						isDragActive
							? "border-blue-500 bg-blue-50"
							: "border-gray-300 bg-gray-100"
					}`}
				>
					<input {...getInputProps()} onChange={handleFileChange} />
					{isDragActive ? (
						<p className="text-blue-500">
							Solte os arquivos aqui...
						</p>
					) : (
						<p className="text-gray-600">
							Arraste ou clique e selecione os arquivos.
						</p>
					)}
				</div>

				{selectedFiles.length > 0 && (
					<div className="mt-4">
						<h3 className="font-medium text-lg mb-2">
							Files Selected:
						</h3>
						<ul className="list-none">
							{selectedFiles.map((file, index) => (
								<li
									key={index}
									className="flex items-center justify-between mb-2"
								>
									<div className="flex items-center">
										<span className="text-gray-700">
											{file.name}
										</span>
									</div>
									<button
										type="button"
										className="text-red-500 hover:text-red-700 ml-2"
										onClick={() => handleRemoveFile(index)}
									>
										Remover
									</button>
								</li>
							))}
						</ul>
					</div>
				)}

				<div className="mt-6">
					<Button onClick={handleFormSubmit} disabled={btnLoading}>
						<span className="flex items-center">
							Enviar {btnLoading && <SpinLoading />}
						</span>
					</Button>
				</div>
			</form>
		</div>
	);
};

export default App;
