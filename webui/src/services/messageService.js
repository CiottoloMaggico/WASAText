import {getAuthentication} from "./sessionService";
import api from "../services/axios";

export const MessageService = Object.freeze({
	async setSeen(conversation) {
		const response = await api.put(`/users/${getAuthentication()}/conversations/${conversation.id}/messages`)

		if (response.status !== 200) {
			throw new Error(response.statusText)
		}

		return response
	},
	async getMessages(conversation) {
		const response = await api.get(`/users/${getAuthentication()}/conversations/${conversation.id}/messages`)

		if (response.status !== 200) {
			throw new Error(response.statusText)
		}

		return response
	},
	async sendMessage(conversation, newMessage) {
		let cleanedForm = Object.fromEntries(Object.entries(newMessage).filter(([key, v]) => v != null && key !== "repliedMessage"))

		if (newMessage.repliedMessage) {
			cleanedForm.repliedMessageId = newMessage.repliedMessage.id
		}

		const response = await api.post(
			`/users/${getAuthentication()}/conversations/${conversation.id}/messages`,
			cleanedForm,
			{
				headers: {"Content-Type": "multipart/form-data"},
			}
		)

		if (response.status !== 200) {
			throw new Error(response.statusText)
		}

		return response
	},
	async deleteMessage(message) {
		const response = await api.delete(`/users/${getAuthentication()}/conversations/${message.conversationId}/messages/${message.id}`)

		if (response.status !== 204) {
			throw new Error(response.statusText)
		}

		return response
	},
	async getComments(message) {
		const response = await api.get(
			`/users/${getAuthentication()}/conversations/${message.conversationId}/messages/${message.id}/comments`,
		)

		if (response.status !== 200) {
			throw new Error(response.statusText)
		}

		return response

	},
	async commentMessage(message, newComment) {
		const response = await api.put(
			`/users/${getAuthentication()}/conversations/${message.conversationId}/messages/${message.id}/comments`,
			{comment: newComment}
		)

		if (response.status !== 200) {
			throw new Error(response.statusText)
		}

		return response
	},
	async uncommentMessage(message) {
		const response = await api.delete(
			`/users/${getAuthentication()}/conversations/${message.conversationId}/messages/${message.id}/comments`,
		)

		if (response.status !== 204) {
			throw new Error(response.statusText)
		}

		return response
	},
	async forwardMessage(message, destination) {
		const response = await api.post(
			`/users/${getAuthentication()}/conversations/${message.conversationId}/messages/${message.id}/forward`,
			{
				destConversationId : destination.id,
			}
		)

		if (response.status !== 200) {
			throw new Error(response.statusText)
		}

		return response
	}

})

export default MessageService
