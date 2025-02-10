<script setup>
import {ref, watch, watchEffect, useTemplateRef, nextTick, onBeforeMount, reactive, computed} from "vue"
import UserService from "@/services/userService";
import {getAuthentication} from "@/services/session";
import ConversationService from "@/services/conversationService";
import {useProfileStore} from "@/stores/profileStore";
import TheConversationList from "@/components/TheConversationList.vue";
import router from "@/router";

const emits = defineEmits(["switch"])

const profileStore = useProfileStore()
const newImageField = useTemplateRef("file-upload")

const newGroup = reactive({
	name: null,
	image: null,
	participants: [],
})

const searchedUsername = ref("")
const users = ref([])

const searchQueryParams = computed(() => {
	return {
		filter: `username like '${searchedUsername.value}%' and uuid ne '${getAuthentication()}'`
	}
})
const newImagePreviewUrl = computed(() => {
	if (!newGroup.image) {
		return "/src/assets/images/default_group_image.jpg"
	}
	return URL.createObjectURL(newGroup.image)
})

await searchForUsers()

watch(searchedUsername, async (newValue) => {
	await searchForUsers()
})

async function searchForUsers() {
	let params = searchQueryParams
	try {
		const response = await UserService.getUsers(params.value)
		users.value = response.data.content
	} catch (error) {
		console.error(error.toString())
		return error
	}
}

async function createGroup() {
	try {
		const response = await ConversationService.createGroup(newGroup)
		await profileStore.getConversations()
		emits("switch", TheConversationList.__name)
		router.push({name: "conversation", params: {convId: response.data.id}})
	} catch (error) {
		console.error(error.toString())
		return error
	}
}

function selectParticipant(participantUuid) {
	let index = newGroup.participants.indexOf(participantUuid)
	if (index !== -1) {
		newGroup.participants.splice(index, 1)
		return
	}
	newGroup.participants.push(participantUuid)
	return
}

function fileUploaded() {
	newGroup.image = newImageField.value.files.item(0)
}

</script>

<template>
	<div class="sidebar-content">
		<div class="sidebar-header">
			<h3 class="header-title">New Group</h3>
		</div>

		<form @submit.prevent="createGroup" class="group-form">
			<div class="sidebar-body">
				<label for="file-upload" class="avatar-box flex-shrink-0">
					<img :src="newImagePreviewUrl" class="avatar"/>
					<img class="form-icon image-edit" src="@/assets/images/edit-pencil.svg"/>
				</label>
				<input type="file"
					   id="file-upload"
					   ref="file-upload"
					   class="d-none"
					   @change="fileUploaded"/>
				<div class="profile-username-box">
					<input
						v-model="newGroup.name"
						ref="profile-username"
						class="profile-username-input d-inline mb-1"
						minlength="3"
						maxlength="16"
						placeholder="Group name"
						required
					/>
				</div>

				<div class="sidebar-group h-100">
					<div class="sidebar-group-header">
				<span class="sidebar-group-title">
					Search user
				</span>
						<input v-model="searchedUsername" class="search-bar" placeholder="Search"/>
					</div>
					<div class="users-container">
						<div v-if="users.length !== 0" v-for="user in users" :key="user.uuid"
							 class="sidebar-item user-item"
							 :class="(newGroup.participants.includes(user.uuid)) ? 'selected' : ''"
							 @click="selectParticipant(user.uuid)">
							<span class="username">{{ user.username }}</span>
						</div>
						<div v-else class="sidebar-item">
							<span class="sidebar-item-title">No users found</span>
						</div>
					</div>
				</div>
			</div>
			<div class="form-footer">
				<button type="submit" class="btn btn-primary">Crea gruppo</button>
			</div>
		</form>
	</div>
</template>

<style scoped>
.selected {
	background: var(--SECONDARY-COLOR);
	border: none !important;
}

.form-icon {
	width: 1.5rem;
	height: 1.5rem;
	object-fit: contain;
	flex-shrink: 0;
}

.group-form {
	display: flex;
	flex-flow: column nowrap;
	height: 100%;
	width: 100%;
	overflow: hidden;
}

.sidebar-body {
	display: flex;
	flex-flow: column nowrap;
	height: 100%;
	overflow: hidden;
	align-items: center;
	gap: 1rem;
}

.profile-username-box {
	display: flex;
	flex-flow: row nowrap;
	align-items: center;
	gap: 1rem;
	overflow: hidden;
	border-bottom: 5px var(--PRIMARY-COLOR) solid;
	margin: 0 1rem;
	flex-shrink: 0;
}

.profile-username-input {
	font-size: 2rem;
	color: #000;
	width: 100%;
	background-color: transparent;
	border: none;
	box-shadow: none;
	white-space: nowrap;
	text-overflow: ellipsis;
	overflow: hidden;
}

.profile-username-input:focus {
	box-shadow: none;
	border: none;
	outline: none;
}

.avatar-box {
	position: relative;
	width: 250px;
	height: 250px;
}

.avatar-box:hover {
	filter: brightness(.8);

	.image-edit {
		display: block;
	}
}

.image-edit {
	display: none;
	position: absolute;
	top: calc(50% - 11.5px);
	left: calc(50% - 7px);
}

.avatar {
	width: 100%;
	height: 100%;
	border-radius: 50%;
	object-fit: cover;
}

.sidebar-group {
	display: flex;
	width: 100%;
	flex-shrink: 0;
	flex-flow: column nowrap;
	border-bottom: 1px solid #aaaaaa;
	max-height: 100%;
}

.sidebar-group-header {
	width: 100%;
	padding: 1rem;
	display: flex;
	flex-flow: column nowrap;
	gap: 1rem;
	border-top: 1px solid #aaaaaa;
	border-bottom: 1px solid #aaaaaa;
}

.search-bar {
	height: 2.5rem;
	width: 100%;
	border: 1px solid #aaaaaa;
	background: #e4e4e4;
	border-radius: 1rem;
	padding: 0 .5rem;
}

.sidebar-group-title {
	font-size: 1.3rem;
	font-weight: 600;
}

.sidebar-item {
	padding: 0 1rem;
	display: flex;
	flex-shrink: 0;
	align-items: center;
	width: 100%;
	height: 5rem;
	border-bottom: 1px solid #e4e4e4;
}

.sidebar-item:hover {
	background: var(--SECONDARY-COLOR);
	border: none;
	cursor: pointer;
}

.sidebar-item-title {
	font-size: 1.2rem;
	font-weight: 600;
}

.users-container {
	overflow-y: scroll;
	overflow-x: hidden;
	height: 100%;
	width: 100%;
}

.user-item {
	height: 3rem;
}

.username {
	font-size: 1rem;
	font-weight: 500;
}

.form-footer {
	display: flex;
	height: 4rem;
	border-top: 1px solid #aaaaaa;
	align-items: center;
	justify-content: center;
	flex-shrink: 0;
}
</style>
