"use client";

import React, { useState } from "react";
import { useRouter } from "next/navigation";

import Button from "../../components/button/Button";
import Input from "../../components/input/Input";
import { fetchAPI } from "../../config/api";
import SpinLoading from "../../components/loading/Loading";

const Login = () => {
	const router = useRouter();
	const [username, setUsername] = useState("");
	const [password, setPassword] = useState("");
	const [btnLoading, setBtnLoading] = useState(false);

	const handleLogin = async () => {
		if ([username, password].some((v) => !v)) {
			return;
		}

		setBtnLoading(true);

		const response = await fetchAPI({
			method: "POST",
			endpoint: "/auth",
			bodyData: { username, password },
		});

		if (response.success) {
			router.push(`/insurer/${username}`);
		} else {
			alert(response.message);
		}

		setBtnLoading(false);
	};

	return (
		<div className="min-h-screen flex items-center justify-center bg-gray-100">
			<div className="bg-white shadow-md rounded p-8 max-w-md w-full">
				<h2 className="text-2xl font-bold mb-6 text-center">Login</h2>
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
					<Button onClick={handleLogin} disabled={btnLoading}>
						<span className="flex items-center">
							Login {btnLoading && <SpinLoading />}
						</span>
					</Button>
				</div>
			</div>
		</div>
	);
};

export default Login;
