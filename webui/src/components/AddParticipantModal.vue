<script setup>
import {computed, ref, watch} from "vue"
import {getAuthentication} from "@/services/session";
import UserService from "@/services/userService";
import ConversationService from "@/services/conversationService";

const props = defineProps(["conversation", "participants", "singleMode"])

const searchedUsername = ref("")
const users = ref([])
const addedUsers = ref([])

const currentParticipants = computed(() => {
	return props.participants.map((user) => user.uuid)
})

const searchQueryParams = computed(() => {
	return {
		filter: `username like '${searchedUsername.value}%' and uuid ne '${getAuthentication()}'`
	}
})

await searchForUsers()

watch(searchedUsername, async () => {
	await searchForUsers()
})

async function addParticipants() {
	try {
		const response = await ConversationService.addToGroup(props.conversation.id, addedUsers.value)
		Object.assign(props.conversation, response.data)
		clearSelections()
	} catch (error) {
		console.error(error.toString())
	}
}

async function searchForUsers() {
	try {
		const response = await UserService.getUsers(searchQueryParams.value)
		users.value = response.data.content
	} catch (error) {
		console.error(error.toString())
	}
}

function selectParticipant(participantUuid) {
	let index = addedUsers.value.indexOf(participantUuid)
	if (index !== -1) {
		addedUsers.value.splice(index, 1)
		return
	}
	addedUsers.value.push(participantUuid)
}

function clearSelections() {
	addedUsers.value = []
	searchedUsername.value = ""
}
</script>

<template>
	<div
		class="modal fade"
		ref="participants-modal"
		id="participants-modal"
		tabindex="-1"
		aria-labelledby="participantsModalLabel"
		aria-hidden="true"
		data-bs-backdrop="static"
		data-bs-keyboard="false"
	>
		<div class="modal-dialog modal-dialog-centered">
			<div class="modal-content">
				<div class="modal-header">
					<h5 class="modal-title" id="participantsModalLabel">Add group participants</h5>
				</div>
				<div class="modal-body">
					<div class="search-box">
						<input v-model="searchedUsername" type="text" class="search-bar" placeholder="Search"/>
					</div>
					<div class="user-list">
						<template v-for="user in users" :key="user.uuid">
							<div v-if="currentParticipants.includes(user.uuid)"
								 class="user-item selected pe-none"
							>
								<span class="user-username">{{ user.username }}</span>
							</div>
							<div v-else class="user-item"
								 :class="(!singleMode && addedUsers.includes(user.uuid)) ? 'selected' : ''"
								 @click="(!singleMode) ? selectParticipant(user.uuid) : $emit('addParticipant', user.uuid)"
							>
								<span class="user-username">{{ user.username }}</span>
							</div>

						</template>
					</div>
				</div>
				<div class="modal-footer">
					<button type="button" class="btn btn-secondary rounded-pill" data-bs-dismiss="modal" @click="clearSelections">Close
					</button>
					<button type="submit" class="btn btn-primary rounded-pill" data-bs-dismiss="modal"
							:disabled="addedUsers.length === 0" @click="addParticipants">
						Add participants
					</button>
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped>
.selected, .user-item:hover {
	background: var(--SECONDARY-COLOR);
	border: none !important;
}


.modal-body {
	display: flex;
	flex-flow: column nowrap;
	height: 400px;
	width: 500px;
	padding: 0px;
	overflow: hidden;
}

.search-box {
	height: 4.5rem;
	padding: 1rem;
	border-bottom: 1px solid #e4e4e4;
	flex-shrink: 0;
}

.search-bar {
	height: 100%;
	width: 100%;
	border: 1px solid #aaaaaa;
	background: #e4e4e4;
	border-radius: 1rem;
	padding: 0 .5rem;
}

.user-list {
	display: flex;
	flex-flow: column nowrap;
	overflow-y: scroll;
	overflow-x: hidden;
	height: 100%;
	width: 100%;
}

.user-item {
	flex-shrink: 0;
	height: 3rem;
	width: 100%;
	display: flex;
	align-items: center;
	padding: 0 1rem;
	border-bottom: 1px solid #e4e4e4;
}

.user-item:hover {
	cursor: pointer;
	border: none;
}
</style>
