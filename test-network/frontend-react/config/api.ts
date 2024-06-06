const API_BASE_URL = "http://localhost:3001";

const query = async (params: Record<string, any>): Promise<any> => {
	const queryString = new URLSearchParams(params).toString();
	const url = `${API_BASE_URL}/query?${queryString}`;

	const response = await fetch(url);
	if (!response.ok) {
		throw new Error(
			`Failed to fetch data from ${url}: ${response.statusText}`
		);
	}

	return response.json();
};

const invoke = async (body: Record<string, any>): Promise<any> => {
	const url = `${API_BASE_URL}/invoke`;

	const response = await fetch(url, {
		method: "POST",
		headers: {
			"Content-Type": "application/json",
		},
		body: JSON.stringify(body),
	});
	if (!response.ok) {
		throw new Error(`Failed to invoke ${url}: ${response.statusText}`);
	}

	return response.json();
};

type RequestMethod = "GET" | "POST" | "PUT";

type FetchOptions = {
	method: RequestMethod;
	endpoint: string;
	queryParams?: Record<string, string>;
	bodyData?: Record<string, any>;
	headers?: Record<string, string>;
};

const fetchAPI = async ({
	method,
	endpoint,
	queryParams,
	bodyData,
	headers = {
		"Content-Type": "application/json",
	},
}: FetchOptions): Promise<any> => {
	const queryString = queryParams
		? "?" + new URLSearchParams(queryParams).toString()
		: "";

	const url = API_BASE_URL + endpoint + queryString;

	const config: RequestInit = {
		method,
		headers,
	};

	if (method !== "GET" && bodyData) {
		config.body = JSON.stringify(bodyData);
	}

	try {
		const response = await fetch(url, config);

		if (!response.ok) {
			throw new Error(`HTTP error! status: ${response.status}`);
		}

		return await response.json();
	} catch (error) {
		console.error("Fetch API error:", error);
		throw error;
	}
};

export { query, invoke, fetchAPI };
