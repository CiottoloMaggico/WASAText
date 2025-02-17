<script setup>
import {ref, toRefs, watch} from "vue"
import {getAuthentication} from "@/services/sessionService";
import {getApiUrl} from "@/services/axios";
import MessageService from "@/services/messageService";
import UserService from "@/services/userService";
import ConversationService from "@/services/conversationService";
import router from "@/router";

const props = defineProps(["message", "show"])
const emits = defineEmits(["close"]);


const conversations = ref([])
const users = ref([])
const forwardToChat = ref(true)
const newConversationMode = ref(false)
const searchedUsername = ref("")

const {message, show} = toRefs(props)

watch(show, async (newVal, oldVal) => {
	if (!oldVal && newVal) {
		await getConversations()
	}
})

watch(forwardToChat, async () => {
	await getConversations()
})

watch(newConversationMode, async (newMode) => {
	if (newMode) {
		await searchForUsers()
	}
})

watch(searchedUsername, async () => {
	await searchForUsers()
})

async function getConversations() {
	let filterQuery = `id ne ${message.value.conversationId}`
	filterQuery += (forwardToChat.value) ? " and type eq 'chat'" : " and type eq 'group'"

	try {
		const data = await ConversationService.getConversations(
			{filter: filterQuery},
		)
		conversations.value = data.content
	} catch (e) {
		console.error(e)
	}
}

async function searchForUsers() {
	let filterQuery = `username like '${searchedUsername.value}%' and uuid ne '${getAuthentication()}'`

	try {
		const data = await UserService.getUsers({filter: filterQuery})
		users.value = data.content
	} catch (error) {
		console.error(error.toString())
	}
}

async function forwardMessage(destination, newChat) {
	if (newChat) {
		try {
			destination = await ConversationService.createChat(destination)
		} catch (e) {
			console.error(e)
		}
	}

	try {
		await MessageService.forwardMessage(message.value, destination)
		await router.push({name: "conversation", params: {convId: destination.id}})
	} catch (e) {
		console.error(e)
	}
	closeModal()
}

function newConversationModeHandler(activate) {
	searchedUsername.value = ""
	newConversationMode.value = activate
}

function closeModal() {
	emits("close")
	newConversationMode.value = false
	forwardToChat.value = true
}
</script>

<template>
	<div
		class="modal fade"
		:ref="`forward-modal-${message.id}`"
		:id="`forward-modal-${message.id}`"
		tabindex="-1"
		aria-hidden="true"
		data-bs-backdrop="static"
		data-bs-keyboard="false"
	>
		<div class="modal-dialog modal-dialog-centered">
			<div class="modal-content">
				<div class="modal-header">
					<h5 class="modal-title">Forward to</h5>
				</div>
				<div class="modal-body">
					<template v-if="!newConversationMode">
						<div class="conversation-switch">
							<div class="switch">
								<div class="switch-btn" :class="{'switch-btn-selected' : forwardToChat}"
									 @click="forwardToChat = true">
									Your chats
								</div>
								<div class="switch-btn" :class="{'switch-btn-selected' : !forwardToChat}"
									 @click="forwardToChat = false">
									Your groups
								</div>
							</div>
						</div>
						<div class="cards-list">
							<div v-if="forwardToChat" class="list-card" @click="newConversationModeHandler(true)">
								<div class="card-content">
									<img class="new-chat-icon" src="@/assets/images/add-chat.svg" alt="start new chat"/>
									<span>New chat</span>
								</div>
							</div>
							<div v-for="conversation in conversations" :key="conversation.id" class="list-card"
								 data-bs-dismiss="modal" @click="forwardMessage(conversation, false)">
								<div class="card-content">
									<div class="img-box">
										<img class="img" :src="getApiUrl(conversation.image.fullUrl)"
											 :width="conversation.image.width"
											 :height="conversation.image.height" :alt="`${conversation.name}`"/>
									</div>
									<span class="card-title">
										{{ conversation.name }}
								</span>
								</div>
							</div>
						</div>
					</template>
					<template v-else>
						<div class="search-box">
							<input v-model="searchedUsername" type="text" class="search-bar" placeholder="Search"/>
						</div>
						<div class="cards-list">
							<div v-for="user in users" :key="user.uuid" class="list-card"
								 @click="forwardMessage(user, true)" data-bs-dismiss="modal"
							>
								<div class="card-content">
									<div class="img-box">
										<img class="img" :src="getApiUrl(user.photo.fullUrl)"
											 :width="user.photo.width"
											 :height="user.photo.height" :alt="`${user.username}`"/>
									</div>
									<span class="card-title">
										{{ user.username }}
								</span>
								</div>
							</div>
						</div>
					</template>
				</div>
				<div class="modal-footer">
					<button v-if="!newConversationMode" type="button" class="btn btn-secondary rounded-pill"
							data-bs-dismiss="modal"
							@click="closeModal">
						Close
					</button>
					<button v-else type="button" class="btn btn-secondary rounded-pill"
							@click="newConversationModeHandler(false)">
						Indietro
					</button>
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped>
.modal-body {
	width: 100%;
	height: 400px;
	padding: 0;
	display: flex;
	flex-flow: column nowrap;
}

.conversation-switch {
	width: 100%;
	height: 4rem;
	padding: 1rem 3rem;
	border-bottom: 1px solid #e4e4e4;
}

.switch {
	width: 100%;
	height: 100%;
	background-color: #e4e4e4;
	display: flex;
	flex-flow: row nowrap;
	border-radius: 5px;
	padding: 3px;
}

.switch-btn {
	display: flex;
	justify-content: center;
	align-items: center;
	width: 50%;
	height: 100%;
	user-select: none;
}

.switch-btn:hover {
	cursor: pointer;
}

.switch-btn-selected {
	background-color: var(--PRIMARY-COLOR);
	border-radius: 3px;
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

.cards-list {
	width: 100%;
	height: 100%;
	overflow-y: scroll;
	overflow-x: hidden;
}

.new-chat-icon {
	width: 1.5rem;
	height: 1.5rem;
}

.list-card {
	display: flex;
	height: 4rem;
	width: 100%;
	flex-flow: row nowrap;
	align-items: center;
	border-bottom: 1px solid #e4e4e4;

	&:hover{
		cursor: pointer;
		background: var(--SECONDARY-COLOR);
		border: none;
	}
}

.card-content {
	padding: .5rem;
	display: flex;
	width: 100%;
	height: 100%;
	flex-flow: row nowrap;
	align-items: center;
	gap: .5rem;
}

.img-box {
	display: flex;
	flex-shrink: 0;
	width: 3rem;
	height: 100%;
	flex-flow: row nowrap;
	align-items: center;
	gap: .5rem;
}

.img {
	width: 100%;
	height: 100%;
	border-radius: 50%;
	object-fit: cover;
}

.card-title {
	font-size: 1.2rem;
}
</style>

