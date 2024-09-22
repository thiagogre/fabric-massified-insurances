"use client";

import React, { useEffect, useState } from "react";

import { invoke, query } from "../../config/api";
import { wait } from "../../utils/promise";
import { uniqueId } from "../../utils/uniqueId";

const App = () => {
	const [data, setData] = useState<any>(null);

	useEffect(() => {
		const runAssetFlow = async () => {
			const id = String(uniqueId());

			await invoke({
				channelid: "mychannel",
				chaincodeid: "basic",
				function: "CreateAsset",
				args: [id, "Dono", "12", "5000", "0", "Varejista", "300"],
			});

			await wait(1000);

			await invoke({
				channelid: "mychannel",
				chaincodeid: "basic",
				function: "UpdateAsset",
				args: [
					id,
					"Dono",
					"12",
					"5000",
					"0",
					"Varejista",
					"300",
					"Active",
				],
			});

			await wait(1000);

			await invoke({
				channelid: "mychannel",
				chaincodeid: "basic",
				function: "UpdateAsset",
				args: [
					id,
					"Dono",
					"12",
					"5000",
					"0",
					"Varejista",
					"300",
					"Active",
				],
			});

			await wait(1000);

			await invoke({
				channelid: "mychannel",
				chaincodeid: "basic",
				function: "UpdateAsset",
				args: [
					id,
					"Dono",
					"12",
					"5000",
					"0",
					"Varejista",
					"300",
					"Pending",
				],
			});

			await wait(1000);

			await invoke({
				channelid: "mychannel",
				chaincodeid: "basic",
				function: "DeleteAsset",
				args: [id],
			});

			await wait(1000);

			const responseData = await query({
				channelid: "mychannel",
				chaincodeid: "basic",
				function: "GetAssetRecords",
				args: [id],
			});
			setData(responseData);
		};

		const timeout = setTimeout(runAssetFlow, 1);

		// INFO: we can have more than one transaction per block
		// const timeout = setTimeout(() => {
		// 	runAssetFlow();
		// 	runAssetFlow();
		// 	runAssetFlow();
		// 	runAssetFlow();
		// 	runAssetFlow();
		// 	runAssetFlow();
		// 	runAssetFlow();
		// 	runAssetFlow();
		// 	runAssetFlow();
		// 	runAssetFlow();
		// 	runAssetFlow();
		// 	runAssetFlow();
		// 	runAssetFlow();
		// 	runAssetFlow();
		// 	runAssetFlow();
		// 	runAssetFlow();
		// 	runAssetFlow();
		// 	runAssetFlow();
		// 	runAssetFlow();
		// 	runAssetFlow();
		// }, 1);

		return () => clearTimeout(timeout);
	}, []);

	return (
		<div className="min-h-screen flex flex-col items-center justify-center space-y-8 p-4">
			<h1 className="text-4xl font-bold text-center">Asset Records</h1>
			<div className="flex items-center justify-center space-x-8">
				{data && (
					<pre className="whitespace-pre-wrap overflow-auto p-2 border border-gray-300 rounded-md bg-gray-100">
						{JSON.stringify(data, null, 2)}
					</pre>
				)}
			</div>
		</div>
	);
};

export default App;
