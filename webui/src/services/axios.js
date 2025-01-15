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

api.defaults.headers.common['Authorization'] = 'Bearer ' + getAuthentication();

export default api
