import {getAuthentication} from "./session";
import api from "../services/axios";

export const ConversationService = Object.freeze({

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
