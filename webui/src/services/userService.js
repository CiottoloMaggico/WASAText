import api from "./axios";
import qs from "qs";
import {getAuthentication} from "./session";

export const UserService = Object.freeze({
	async getProfile() {
		const authedUserUUID = getAuthentication()
		const response = await api.get(`/users/${authedUserUUID}`)

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
		const authedUserUUID = getAuthentication()
		const response = await api.put(
			`/users/${authedUserUUID}/username`,
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
		const authedUserUUID = getAuthentication()
		const response = await api.put(
			`/users/${authedUserUUID}/avatar`,
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
