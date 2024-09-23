"use client";

import React from "react";
import { useRouter } from "next/navigation";

import { uniqueId } from "../utils/uniqueId";
import Button from "../components/button/Button";

const App = () => {
	const router = useRouter();

	const routes = [
		{ route: "/partner", title: "Loja do Parceiro de distribuição" },
		{ route: "/login", title: "Login" },
		{ route: "/evidence-analysis", title: "Análise de Evidências" },
		{ route: "/insurer/finish", title: "Aprovar/reprovar claim" },
		{ route: "/events", title: "Events" },
		{ route: "/asset-records", title: "Asset Records" },
	];

	return (
		<div className="min-h-screen flex flex-col items-center justify-center space-y-8 p-4">
			<h1 className="text-4xl font-bold text-center">Home</h1>
			<div className="flex flex-col md:flex-row items-start justify-center space-y-8 space-x-8 mt-8 w-full">
				{routes.map((r) => (
					<div
						key={uniqueId()}
						className="mt-8 flex flex-col items-center space-y-4"
					>
						<Button onClick={() => router.push(r.route)}>
							{r.title}
						</Button>
					</div>
				))}
			</div>
		</div>
	);
};

export default App;
