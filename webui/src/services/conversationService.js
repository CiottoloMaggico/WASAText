import {getAuthentication} from "./sessionService";
import api from "../services/axios";
import {useConversationsStore} from "@/stores/conversationsStore";
import router from "@/router";

export const ConversationService = Object.freeze({
	get store() {
		return useConversationsStore()
	},
	get authedUserUUID() {
		return getAuthentication()
	},
	async refresh() {
		const data = await this.setDelivered()
		this.store.update(data)
	},
	async createChat(recipient) {
		const response = await api.post(
			`/users/${this.authedUserUUID}/chats`,
			{
				recipient: recipient.uuid,
			}
		)

		if (response.status === 409) {
			let conversation = this.store.getChatByRecipient(recipient.username);
			if (conversation) {
				return conversation
			}
			let missingConversationData = await this.getConversations({
				filter: `name eq '${recipient.username}' and type eq 'chat'`
			})
			conversation = missingConversationData.content[0]
			this.store.updateConversation(conversation)
			return conversation
		} else if (response.status !== 200) {
			throw new Error(response.statusText)
		}

		this.store.updateConversation(response.data)
		return response.data
	},
	async createGroup(newGroup) {
		let groupData = Object.fromEntries(Object.entries(newGroup).filter(([key, v]) => v != null && key !== 'participants'))

		let response = await api.post(
			`/users/${this.authedUserUUID}/groups`,
			groupData,
			{
				headers: {"Content-Type": "multipart/form-data"},
			}
		)

		if (response.status !== 200) {
			throw new Error(response.statusText)
		}

		this.store.updateConversation(response.data)
		if (newGroup.participants.length === 0) {
			return response.data
		}
		let participants = Object.fromEntries(Object.entries(newGroup).filter(([key, v]) => key === 'participants'))

		response = await api.put(
			`/users/${this.authedUserUUID}/groups/${response.data.id}`,
			participants,
		)

		if (response.status !== 200) {
			throw {
				status: 201,
				created: response.data,
				message: "The group was created without participants because an error occurred while adding them"
			}
		}
		return response.data
	},
	async addToGroup(group, participants) {
		const response = await api.put(
			`/users/${this.authedUserUUID}/groups/${group.id}`,
			{
				participants: participants,
			},
		)

		if (response.status !== 200) {
			throw new Error("An error occurred while adding the participants to the group")
		}
		this.store.updateConversation(response.data)
		return response.data
	},
	async setGroupName(group, newName) {
		const response = await api.put(
			`/users/${this.authedUserUUID}/groups/${group.id}/name`,
			{
				name: newName,
			}
		)

		if (response.status !== 200) {
			throw new Error(response.statusText)
		}
		this.store.updateConversation(response.data)
		return response.data
	},
	async setGroupPhoto(group, newPhoto) {
		const response = await api.put(
			`/users/${this.authedUserUUID}/groups/${group.id}/photo`,
			{
				image: newPhoto,
			},
			{
				headers: {"Content-Type": "multipart/form-data"},
			}
		)

		if (response.status !== 200) {
			throw new Error(response.statusText)
		}
		this.store.updateConversation(response.data)
		return response.data
	},
	async leaveGroup(group) {
		const response = await api.delete(`/users/${this.authedUserUUID}/groups/${group.id}`)

		if (response.status !== 204) {
			throw new Error(response.statusText)
		}

		this.store.removeConversation(group)
		return response.data
	},
	async getConversations(params) {
		const response = await api.get(
			`/users/${this.authedUserUUID}/conversations`,
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
	async getNextPage() {
		let nextPage = this.store.pagination.nextPage
		if (!nextPage) return

		const response = await api.get(
			nextPage,
		)

		if (response.status !== 200) {
			throw new Error(response.statusText)
		}

		this.store.addPage(response.data)
		return response.data
	},
	async setDelivered() {
		const response = await api.put(
			`/users/${this.authedUserUUID}/conversations`
		)

		if (response.status !== 200) {
			throw new Error(response.statusText)
		}

		return response.data
	},
	async getConversation(conversationId) {
		const response = await api.get(`/users/${this.authedUserUUID}/conversations/${conversationId}`)

		if (response.status !== 200) {
			throw new Error(response.statusText)
		}

		this.store.updateConversation(response.data)
		return response.data
	},
})

export default ConversationService
