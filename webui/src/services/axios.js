import axios from "axios";
import {getAuthentication} from "./session";


export function getApiUrl(path) {
	if (path.startsWith("/")) {
		return __API_URL__ + path;
	}
	return __API_URL__ + "/" + path;
}

export const api = axios.create({
	baseURL: __API_URL__,
	timeout: 1000 * 5
});


api.interceptors.request.use(
	config => {
		config.headers["Authorization"] = `Bearer ${getAuthentication()}`;
		return config;
	},
	error => {
		return Promise.reject(error);
	}
)
export default api
