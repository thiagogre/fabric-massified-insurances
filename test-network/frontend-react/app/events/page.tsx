"use client";

import { useEffect, useState } from "react";

import Event from "../../components/event/Event";
import { EventProps } from "../../components/event/types";
import { fetchAPI } from "../../config/api";

const App = () => {
	const [events, setEvents] = useState<Record<number, EventProps[]>>({});

	useEffect(() => {
		// Group events by block number
		const groupedEvents = (docs: EventProps[] | []) =>
			docs.reduce((acc, event) => {
				if (!acc[event.BlockNumber]) {
					acc[event.BlockNumber] = [];
				}
				acc[event.BlockNumber].push(event);
				return acc;
			}, {} as Record<number, EventProps[]>);

		const getEvents = async () => {
			try {
				const response = await fetchAPI({
					method: "GET",
					endpoint: "/event",
				});
				const docs: EventProps[] | [] = response?.docs ?? [];
				setEvents(groupedEvents(docs));
			} catch (error) {
				console.error("Error fetching events: ", error);
			}
		};

		getEvents();
		const interval = setInterval(getEvents, 5000);

		return () => clearInterval(interval);
	}, []);

	return (
		<div className="min-h-screen flex flex-col items-center justify-center p-4 bg-gray-50">
			<h1 className="text-3xl font-bold mb-8">Events</h1>
			<div className="w-full max-w-6xl flex flex-col justify-center gap-4">
				{Object.entries(events)
					.reverse()
					.map(([blockNumber, blockEvents]) => (
						<div
							key={blockNumber}
							className="w-full border border-gray-300 bg-gray-200 shadow-lg rounded-lg p-4 mb-2"
						>
							<h2 className="text-2xl font-semibold mb-2">
								Block #{blockNumber}
							</h2>
							<div className="w-full flex flex-wrap justify-center gap-2">
								{blockEvents.map((event, index) => (
									<Event
										key={index}
										{...event}
										isBlockHidden
									/>
								))}
							</div>
						</div>
					))}
			</div>
		</div>
	);
};

export default App;
