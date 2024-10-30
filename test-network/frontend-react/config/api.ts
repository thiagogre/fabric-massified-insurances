const SERVER_PORT = process.env.NEXT_PUBLIC_SERVER_PORT;
const API_BASE_URL = `http://localhost:${SERVER_PORT}`;

const query = async (params: Record<string, any>): Promise<any> => {
	const queryString = new URLSearchParams(params).toString();
	const url = `${API_BASE_URL}/smartcontract/query?${queryString}`;

	try {
		const response = await fetch(url);
		return response.json();
	} catch (error) {
		return error;
	}
};

const invoke = async (body: Record<string, any>): Promise<any> => {
	const url = `${API_BASE_URL}/smartcontract/invoke`;

	try {
		const response = await fetch(url, {
			method: "POST",
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify(body),
		});
		return response.json();
	} catch (error) {
		return error;
	}
};

type RequestMethod = "GET" | "POST" | "PUT";

type FetchOptions = {
	method: RequestMethod;
	endpoint: string;
	queryParams?: Record<string, string>;
	bodyData?: any;
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
		config.body = bodyData;
	}

	try {
		const response = await fetch(url, config);
		return response.json();
	} catch (error) {
		return error;
	}
};

export { query, invoke, fetchAPI };
