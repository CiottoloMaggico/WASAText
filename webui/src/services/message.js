import {getAuthentication} from "./session";
import api from "../services/axios";

export const MessageService = Object.freeze({
	authedUserUUID: getAuthentication(),

	async getMessages(conversationId) {
		const response = await api.get(`/users/${this.authedUserUUID}/conversations/${conversationId}/messages`)

		if (response.status !== 200) {
			throw new Error(response.statusText)
		}

		return response
	},

	async sendMessage(conversationId, message) {
		let cleanedForm = Object.fromEntries(Object.entries(message).filter(([key, v]) => v != null))

		const response = await api.post(
			`/users/${this.authedUserUUID}/conversations/${conversationId}/messages`,
			cleanedForm,
			{
				headers: {"Content-Type": "multipart/form-data"},
			}
		)
		// TODO: change to 201 in backend
		if (response.status !== 200 && response.status !== 201) {
			throw new Error(response.statusText)
		}

		return response
	},

	async deleteMessage(conversationId, messageId) {
		const response = await api.delete(`/users/${this.authedUserUUID}/conversations/${conversationId}/messages/${messageId}`)

		if (response.status !== 204) {
			throw new Error(response.statusText)
		}

		return response
	}
})

export default MessageService
