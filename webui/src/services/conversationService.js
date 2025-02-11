import {getAuthentication} from "./sessionService";
import api from "../services/axios";

export const ConversationService = Object.freeze({
	async createChat(recipient) {
		const response = await api.post(
			`/users/${getAuthentication()}/chats`,
			{
				recipient: recipient.uuid,
			}
		)

		if (response.status !== 200) {
			throw new Error(response.statusText)
		}

		return response
	},
	async createGroup(newGroup) {
		let groupData = Object.fromEntries(Object.entries(newGroup).filter(([key, v]) => v != null && key !== 'participants'))

		let response = await api.post(
			`/users/${getAuthentication()}/groups`,
			groupData,
			{
				headers: {"Content-Type": "multipart/form-data"},
			}
		)

		if (response.status !== 200) {
			throw new Error(response.statusText)
		} else if (newGroup.participants.length === 0) {
			return response
		}
		let participants = Object.fromEntries(Object.entries(newGroup).filter(([key, v]) => key === 'participants'))

		response = await api.put(
			`/users/${getAuthentication()}/groups/${response.data.id}`,
			participants,
		)

		if (response.status !== 200) {
			throw new Error(response.statusText)
		}

		return response
	},
	async addToGroup(group, participants) {
		const response = await api.put(
			`/users/${getAuthentication()}/groups/${group.id}`,
			{
				participants: participants,
			},
		)

		if (response.status !== 200) {
			throw new Error(response.statusText)
		}
		return response
	},
	async setGroupName(group, newName) {
		const response = await api.put(
			`/users/${getAuthentication()}/groups/${group.id}/name`,
			{
				name: newName,
			}
		)

		if (response.status !== 200) {
			throw new Error(response.statusText)
		}

		return response
	},
	async setGroupPhoto(group, newPhoto) {
		const response = await api.put(
			`/users/${getAuthentication()}/groups/${group.id}/photo`,
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
	async leaveGroup(group) {
		const response = await api.delete(`/users/${getAuthentication()}/groups/${group.id}`)

		if (response.status !== 204) {
			throw new Error(response.statusText)
		}
		return response
	},
})

export default ConversationService
