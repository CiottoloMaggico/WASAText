<script setup>
import {computed, ref, watch} from "vue"
import UserService from "../services/userService";
import ConversationService from "@/services/conversationService";
import router from "../router";
import {useProfileStore} from "@/stores/profileStore";
import TheConversationList from "@/components/TheConversationList.vue";
import TheNewGroup from "@/components/TheNewGroup.vue";
import TheProfile from "@/components/TheProfile.vue";
import {getAuthentication} from "@/services/session";

const emits = defineEmits(["switch"])

const profileStore = useProfileStore()
const searchedUsername = ref("")
const searchQueryParams = computed(() => {
	return {
		filter: `username like '${searchedUsername.value}%' and uuid ne '${getAuthentication()}'`
	}
})
const users = ref([])

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

async function createChat(recipientUuid) {
	try {
		const response = await ConversationService.createChat(recipientUuid)
		await profileStore.getConversations()
		emits("switch", TheConversationList.__name)
		router.push({name: "conversation", params: {convId: response.data.id}})
	} catch (error) {
		// TODO: redirect to the existing conversation
		console.error(error.toString())
		return error
	}
}

</script>

<template>
	<div class="sidebar-content">
		<div class="sidebar-header">
			<h3 class="header-title">Add conversation</h3>
		</div>
		<div class="sidebar-body">
			<div class="sidebar-group" @click="emits('switch', TheNewGroup.__name)">
				<div class="sidebar-item">
					<span class="sidebar-item-title">New group</span>
				</div>
			</div>
			<div class="sidebar-group h-100">
				<div class="sidebar-group-header">
				<span class="sidebar-group-title">
					Search user
				</span>
					<input v-model="searchedUsername" class="search-bar" placeholder="Search"/>
				</div>
				<div class="users-container">
					<div v-if="users.length !== 0" v-for="user in users" :key="user.uuid" class="sidebar-item user-item" @click="createChat(user.uuid)">
						<span class="username">{{user.username}}</span>
					</div>
					<div v-else class="sidebar-item">
						<span class="sidebar-item-title">No users found</span>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped>
.sidebar-body {
	height: 100%;
	max-height: 100%;
	display: flex;
	flex-flow: column nowrap;
	border-top: 1px solid #aaaaaa;
	overflow: hidden;
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

.sidebar-item:hover{
	background: var(--SECONDARY-COLOR);
	border: none ;
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


</style>
