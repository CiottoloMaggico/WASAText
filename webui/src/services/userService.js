import api from "./axios";
import qs from "qs";
import {getAuthentication} from "./sessionService";

export const UserService = Object.freeze({
	async getProfile() {
		const response = await api.get(`/users/${getAuthentication()}`)

		if (response.status !== 200) {
			throw new Error(response.statusText)
		}

		return response
	},
	async getUsers(params) {
		const response = await api.get(
			"/users",
			{
				params: params,
				paramsSerializer: (params) => {
					return qs.stringify(params, {arrayFormat: "repeat"})
				}
			},
		)
		if (response.status !== 200) {
			throw new Error(response.statusText)
		}
		return response
	},
	async setMyUsername(newUsername) {
		const response = await api.put(
			`/users/${getAuthentication()}/username`,
			{
				username: newUsername,
			}
		)

		if (response.status !== 200) {
			throw new Error(response.statusText)
		}

		return response
	},
	async setMyPhoto(newPhoto) {
		const response = await api.put(
			`/users/${getAuthentication()}/avatar`,
			{
				photo: newPhoto,
			},
			{
				headers: {"Content-Type": "multipart/form-data"},
			}
		)

		if (response.status !== 200) {
			throw new Error(response.statusText)
		}

		return response
	},
})

export default UserService
