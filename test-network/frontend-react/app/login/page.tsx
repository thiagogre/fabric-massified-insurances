"use client";

import React, { useState } from "react";

import Button from "../../components/button/Button";
import Input from "../../components/input/Input";
import { fetchAPI } from "../../config/api";

const Login = () => {
	const [username, setUsername] = useState("");
	const [password, setPassword] = useState("");
	const [error, setError] = useState("");

	const handleLogin = async () => {
		if ([username, password].some((v) => !v)) {
			setError("Please fill in all fields.");
			return;
		}

		try {
			const response = await fetchAPI({
				method: "POST",
				endpoint: "/auth",
				bodyData: { username, password },
			});

			// Handle successful login here
			console.log("Login successful:", response);
		} catch (error) {
			console.error("Error during login:", error);
			setError("Login failed. Please try again.");
		}
	};

	return (
		<div className="min-h-screen flex items-center justify-center bg-gray-100">
			<div className="bg-white shadow-md rounded p-8 max-w-md w-full">
				<h2 className="text-2xl font-bold mb-6 text-center">Login</h2>
				{error && <p className="text-red-500 mb-4">{error}</p>}
				<div className="mb-4">
					<Input
						type="text"
						placeholder="Username"
						value={username}
						onChange={(e) => setUsername(e.target.value)}
					/>
				</div>
				<div className="mb-6">
					<Input
						type="password"
						placeholder="Password"
						value={password}
						onChange={(e) => setPassword(e.target.value)}
					/>
				</div>
				<div className="flex justify-center">
					<Button onClick={handleLogin}>Login</Button>
				</div>
			</div>
		</div>
	);
};

export default Login;
