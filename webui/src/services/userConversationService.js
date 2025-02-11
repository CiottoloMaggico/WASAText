import {getAuthentication} from "./sessionService";
import api from "../services/axios";
import qs from "qs";

export const UserConversationService = Object.freeze({
	async setDelivered() {
		const response = await api.put(
			`/users/${getAuthentication()}/conversations`
		)

		if (response.status !== 200) {
			throw new Error(response.statusText)
		}

		return response
	},
	async getConversations(params) {
		const response = await api.get(
			`/users/${getAuthentication()}/conversations`,
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
	async getConversation(conversationId) {
		const response = await api.get(`/users/${getAuthentication()}/conversations/${conversationId}`)

		if (response.status !== 200) {
			throw new Error(response.statusText)
		}

		return response
	}
})

export default UserConversationService
