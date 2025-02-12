import api from "./axios";
import qs from "qs";
import {getAuthentication} from "./sessionService";
import {useProfileStore} from "@/stores/profileStore";

export const UserService = Object.freeze({
	get store() {
		return useProfileStore()
	},
	get authedUserUUID() {
		return getAuthentication()
	},
	async refresh() {
		let data = await this.getProfile()

		this.store.update(data)
		return data
	},
	async getProfile() {
		const response = await api.get(`/users/${this.authedUserUUID}`)

		if (response.status !== 200) {
			throw new Error(response.statusText)
		}

		this.store.update(response.data)
		return response.data
	},
	async setMyUsername(newUsername) {
		const response = await api.put(
			`/users/${this.authedUserUUID}/username`,
			{
				username: newUsername,
			}
		)

		if (response.status === 409) {
			throw new Error("This username already exists, please choose another username.")
		} else if (response.status !== 200) {
			throw new Error(response.statusText)
		}
		this.store.update(response.data)
		return response.data
	},
	async setMyPhoto(newPhoto) {
		const response = await api.put(
			`/users/${this.authedUserUUID}/avatar`,
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

		this.store.update(response.data)
		return response.data
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
		return response.data
	},
})

export default UserService
