<script setup>
import {computed, onBeforeMount, ref, watch} from "vue"
import {getAuthentication} from "@/services/sessionService";
import UserService from "@/services/userService";

const props = defineProps(["participants", "singleMode"])
const emits = defineEmits(["addedParticipants", "addParticipant"]);

const searchedUsername = ref("")
const users = ref([])
const addedUsers = ref([])

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

function selectParticipant(participant) {
	let index = addedUsers.value.findIndex((p) => p.uuid === participant.uuid)
	if (index !== -1) {
		addedUsers.value.splice(index, 1)
	} else {
		addedUsers.value.push(participant)
	}
	emits("addedParticipants", addedUsers)
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
			<template v-for="user in users" v-if="users.length !== 0" :key="user.uuid">
				<div v-if="participants.findIndex((p) => p.uuid === user.uuid) !== -1"
					 class="sidebar-item user-item selected pe-none">
					<span class="username">{{ user.username }}</span>
				</div>
				<div v-else :class="{'selected' : (!singleMode && addedUsers.findIndex((p) => p.uuid === user.uuid) !== -1)}"
					 class="sidebar-item user-item"
					 @click="(!singleMode) ? selectParticipant(user) : $emit('addParticipant', user)">
					<span class="username">{{ user.username }}</span>
				</div>
			</template>
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
