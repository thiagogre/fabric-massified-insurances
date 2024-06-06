import React from "react";
import type { EventProps } from "./types";

const parsePayload = (payload: string): string => {
	try {
		const decodedPayload = atob(payload);
		const parsedPayload = JSON.parse(decodedPayload);
		return JSON.stringify(parsedPayload, null, 2);
	} catch (error) {
		return "Error parsing payload";
	}
};

const Event = (props: EventProps & { isBlockHidden?: boolean }) => {
	const {
		BlockNumber,
		TransactionID,
		ChaincodeName,
		EventName,
		Payload,
		isBlockHidden = false,
	} = props;

	return (
		<div className="flex-1 border border-gray-300 rounded-lg p-4 bg-white shadow-md">
			{!isBlockHidden && (
				<h2 className="text-lg font-semibold mb-2">
					Block Number: {BlockNumber}
				</h2>
			)}
			<p>
				<span className="font-semibold">Transaction ID:</span>{" "}
				{TransactionID}
			</p>
			<p>
				<span className="font-semibold">Chaincode Name:</span>{" "}
				{ChaincodeName}
			</p>
			<p>
				<span className="font-semibold">Event Name:</span> {EventName}
			</p>
			<p>
				<span className="font-semibold">Payload:</span>
			</p>
			<pre className="whitespace-pre-wrap overflow-auto p-2 border border-gray-300 rounded-md bg-gray-100">
				{parsePayload(Payload)}
			</pre>
		</div>
	);
};

export default Event;
