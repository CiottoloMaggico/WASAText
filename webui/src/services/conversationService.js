import {getAuthentication} from "./session";
import api from "../services/axios";

export const ConversationService = Object.freeze({

	async createChat(recipientUuid) {
		const authedUserUUID = getAuthentication()
		const response = await api.post(
			`/users/${authedUserUUID}/chats`,
			{
				recipient: recipientUuid,
			}
		)

		if (response.status !== 200) {
			throw new Error(response.statusText)
		}

		return response
	},
	async createGroup(newGroupData) {
		const authedUserUUID = getAuthentication()
		let groupData = Object.fromEntries(Object.entries(newGroupData).filter(([key, v]) => v != null && key !== 'participants'))
		let participants = Object.fromEntries(Object.entries(newGroupData).filter(([key, v]) => key === 'participants'))

		let response = await api.post(
			`/users/${authedUserUUID}/groups`,
			groupData,
			{
				headers: {"Content-Type": "multipart/form-data"},
			}
		)

		if (response.status !== 200) {
			throw new Error(response.statusText)
		}
		if (participants.participants.length === 0) {
			return response
		}

		response = this.addToGroup(response.data.id, participants.participants)

		if (response.status !== 200) {
			throw new Error(response.statusText)
		}

		return response
	},
	async addToGroup(groupId, participants) {
		const authedUserUUID = getAuthentication()
		const response = await api.put(
			`/users/${authedUserUUID}/groups/${groupId}`,
			{
				participants: participants,
			},
		)

		if (response.status !== 200) {
			throw new Error(response.statusText)
		}
		return response
	},
	async setGroupName(groupId, newName) {
		const authedUserUUID = getAuthentication()
		const response = await api.put(
			`/users/${authedUserUUID}/groups/${groupId}/name`,
			{
				name: newName,
			}
		)

		if (response.status !== 200) {
			throw new Error(response.statusText)
		}

		return response
	},

	async setGroupPhoto(groupId, newPhoto) {
		const authedUserUUID = getAuthentication()
		const response = await api.put(
			`/users/${authedUserUUID}/groups/${groupId}/photo`,
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

		return response
	},

	async leaveGroup(groupId) {
		const authedUserUUID = getAuthentication()
		const response = await api.delete(`/users/${authedUserUUID}/groups/${groupId}`)

		if (response.status !== 204) {
			throw new Error(response.statusText)
		}
		return response
	},
})

export default ConversationService
