import {getAuthentication} from "./session";
import api from "../services/axios";
import qs from "qs";

export const UserConversationService = Object.freeze({
	async setDelivered() {
		const authedUserUUID = getAuthentication()
		const response = await api.put(
			`/users/${authedUserUUID}/conversations`
		)

		if (response.status !== 200) {
			throw new Error(response.statusText)
		}

		return response
	},
	async getConversations(params) {
		const authedUserUUID = getAuthentication()
		const response = await api.get(
			`/users/${authedUserUUID}/conversations`,
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
	async getConversation(id) {
		const authedUserUUID = getAuthentication()
		const response = await api.get(`/users/${authedUserUUID}/conversations/${id}`)

		if (response.status !== 200) {
			throw new Error(response.statusText)
		}

		return response
	}

})

export default UserConversationService
