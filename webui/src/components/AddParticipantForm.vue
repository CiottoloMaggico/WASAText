<script setup>
import {computed, onBeforeMount, ref, watch} from "vue"
import {getAuthentication} from "@/services/sessionService";
import UserService from "@/services/userService";

const props = defineProps(["participants", "singleMode"])
const emits = defineEmits(["addParticipant"])

const searchedUsername = ref("")
const users = ref([])

const searchQueryParams = computed(() => {
	return {
		filter: `username like '${searchedUsername.value}%' and uuid ne '${getAuthentication()}'`
	}
})

onBeforeMount(async () => {
	await searchForUsers()
})

watch(searchedUsername, async () => {
	await searchForUsers()
})

async function searchForUsers() {
	try {
		const data = await UserService.getUsers(searchQueryParams.value)
		users.value = data.content
	} catch (error) {
		console.error(error.toString())
	}
}

function selectParticipant(participantUuid) {
	let index = props.participants.indexOf(participantUuid)
	if (index !== -1) {
		props.participants.splice(index, 1)
		return
	}
	props.participants.push(participantUuid)
}
</script>

<template>
	<div class="sidebar-group h-100">
		<div class="sidebar-group-header">
				<span class="sidebar-group-title">
					Search user
				</span>
			<input v-model="searchedUsername" class="search-bar" placeholder="Search"/>
		</div>
		<div class="users-container">
			<div v-for="user in users" v-if="users.length !== 0" :key="user.uuid"
				 :class="(!singleMode && participants.includes(user.uuid)) ? 'selected' : ''"
				 class="sidebar-item user-item"
				 @click="(!singleMode) ? selectParticipant(user.uuid) : $emit('addParticipant', user)">
				<span class="username">{{ user.username }}</span>
			</div>
			<div v-else class="sidebar-item">
				<span class="sidebar-item-title">No users found</span>
			</div>
		</div>
	</div>
</template>

<style scoped>
.sidebar-group-header, .sidebar-item {
	border-bottom: 1px solid #e4e4e4;
}

.sidebar-item {
	padding: 0 1rem;
}

.search-bar {
	height: 2.5rem;
	width: 100%;
	border: 1px solid #aaaaaa;
	background: #e4e4e4;
	border-radius: 1rem;
	padding: 0 .5rem;
}

.users-container {
	overflow-y: scroll;
	overflow-x: hidden;
	height: 100%;
	width: 100%;
}

.user-item {
	height: 3rem !important;
}

.username {
	font-size: 1rem;
	font-weight: 500;
}

</style>
