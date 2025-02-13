import {getAuthentication} from "./sessionService";
import api from "../services/axios";
import {useConversationsStore} from "@/stores/conversationsStore";

export const MessageService = Object.freeze({
	get store() {
		return useConversationsStore()
	},
	get authedUserUUID() {
		return getAuthentication()
	},
	async setSeen(conversation) {
		const response = await api.put(`/users/${this.authedUserUUID}/conversations/${conversation.id}/messages`)

		if (response.status !== 200) {
			throw new Error(response.statusText)
		}
		if (response.data.content[0]) {
			conversation.read = true
			this.store.updateConversation(conversation)
		}
		this.store.updateLatestMessage(conversation, response.data.content[0])
		return response.data
	},
	async getMessages(conversation, params) {
		const response = await api.get(
			`/users/${this.authedUserUUID}/conversations/${conversation.id}/messages`,
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
	async getMessage(conversationId, messageId) {
		const response = await api.get(`/users/${this.authedUserUUID}/conversations/${conversationId}/messages/${messageId}`)

		if (response.status !== 200) {
			throw new Error(response.statusText)
		}

		return response.data
	},
	async sendMessage(conversation, newMessage) {
		let cleanedForm = Object.fromEntries(Object.entries(newMessage).filter(([key, v]) => v != null && key !== "repliedMessage"))

		if (newMessage.repliedMessage) {
			cleanedForm.repliedMessageId = newMessage.repliedMessage.id
		}

		const response = await api.post(
			`/users/${this.authedUserUUID}/conversations/${conversation.id}/messages`,
			cleanedForm,
			{
				headers: {"Content-Type": "multipart/form-data"},
			}
		)

		if (response.status !== 200) {
			throw new Error(response.statusText)
		}

		return response.data
	},
	async deleteMessage(message) {
		const response = await api.delete(`/users/${this.authedUserUUID}/conversations/${message.conversationId}/messages/${message.id}`)

		if (response.status !== 204) {
			throw new Error(response.statusText)
		}

		return response.data
	},
	async getComments(message) {
		const response = await api.get(
			`/users/${this.authedUserUUID}/conversations/${message.conversationId}/messages/${message.id}/comments`,
		)

		if (response.status !== 200) {
			throw new Error(response.statusText)
		}

		return response.data
	},
	async commentMessage(message, newComment) {
		const response = await api.put(
			`/users/${this.authedUserUUID}/conversations/${message.conversationId}/messages/${message.id}/comments`,
			{comment: newComment}
		)

		if (response.status !== 200) {
			throw new Error(response.statusText)
		}

		return response.data
	},
	async uncommentMessage(message) {
		const response = await api.delete(
			`/users/${this.authedUserUUID}/conversations/${message.conversationId}/messages/${message.id}/comments`,
		)

		if (response.status !== 204) {
			throw new Error(response.statusText)
		}

		return response.data
	},
	async forwardMessage(message, destination) {
		const response = await api.post(
			`/users/${this.authedUserUUID}/conversations/${message.conversationId}/messages/${message.id}/forward`,
			{
				destConversationId : destination.id,
			}
		)

		if (response.status !== 200) {
			throw new Error(response.statusText)
		}

		return response.data
	}

})

export default MessageService
