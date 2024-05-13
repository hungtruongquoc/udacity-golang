import type { FunctionComponent } from "../common/types";
const apiUrl: string = import.meta.env.VITE_API_BASE_URL;

export const Home = (): FunctionComponent => {
	return (
		<div className="bg-blue-300  font-bold w-screen h-screen flex flex-col justify-center items-center ">
			<p className="text-white text-6xl">Hello, world!</p>
			<p> API {apiUrl}</p>
		</div>
	);
};
