import api from "./axios";

const API_AUTHENTICATION_KEY = "auth-key"

function setAuthentication(token) {
	if (token) {
		localStorage.setItem(API_AUTHENTICATION_KEY, token)
	} else {
		localStorage.removeItem(API_AUTHENTICATION_KEY)
	}
}

export function getAuthentication() {
	const token = localStorage.getItem(API_AUTHENTICATION_KEY)
	return token
}

export function isAuthed() {
	if (!getAuthentication()) {
		return false
	}
	return true
}


export const SessionService = Object.freeze({
	async doLogin(username){
		const response = await api.post("/session", {username: username})

		if (response.status >= 200 && response.status < 300) {
			setAuthentication(response.data.uuid)
		} else {
			throw new Error(response.statusText)
		}

		return response
	}
})

export default SessionService
