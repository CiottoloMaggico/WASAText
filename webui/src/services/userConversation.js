import {getAuthentication} from "./session";
import api from "../services/axios";

export const UserConversationService = Object.freeze({

	async getConversations() {
		const authedUserUUID = getAuthentication()
		const response = await api.put(`/users/${authedUserUUID}/conversations`)

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
