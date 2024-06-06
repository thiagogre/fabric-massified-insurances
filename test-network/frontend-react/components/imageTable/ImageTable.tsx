import React from "react";

import type { ImageTableProps } from "./types";

const ImageTable = ({ data }: ImageTableProps) => {
	return (
		<div className="overflow-x-auto max-w-lg">
			<table className="min-w-full bg-white border rounded-lg">
				<thead>
					<tr>
						<th className="px-4 py-2 border-b">Name</th>
						<th className="px-4 py-2 border-b">Size</th>
						<th className="px-4 py-2 border-b">Pixels</th>
						<th className="px-4 py-2 border-b">Dimensions</th>
					</tr>
				</thead>
				<tbody>
					{data.map((item, index) => (
						<tr key={index} className="text-center">
							<td className="px-4 py-2 border-b">{item.name}</td>
							<td className="px-4 py-2 border-b">{item.size}</td>
							<td className="px-4 py-2 border-b">
								{item.pixels}
							</td>
							<td className="px-4 py-2 border-b">
								{item.dimensions}
							</td>
						</tr>
					))}
				</tbody>
			</table>
		</div>
	);
};

export default ImageTable;
