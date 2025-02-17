<script setup>
import {computed, onBeforeMount, ref, watch} from "vue"
import {getAuthentication} from "@/services/sessionService";
import UserService from "@/services/userService";
import {getApiUrl} from "@/services/axios";

const props = defineProps([
	"singleMode", "title", "submitPlaceholder", "closePlaceholder"
])
const emits = defineEmits([
	"close", "addedParticipants", "addParticipant"
]);

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
		return
	}
	addedUsers.value.push(participant)
}

function clearSelections() {
	addedUsers.value = []
	searchedUsername.value = ""

	emits("close")
}
</script>

<template>
	<div class="modal-dialog modal-dialog-centered">
		<div class="modal-content">
			<div class="modal-header">
				<h5 class="modal-title" id="participantsModalLabel">
					{{ title }}
				</h5>
			</div>
			<div class="modal-body">
				<div class="search-box">
					<input v-model="searchedUsername" type="text" class="search-bar" placeholder="Search"/>
				</div>
				<div class="cards-list">
					<div v-for="user in users" :key="user.uuid" class="list-card"
						 :class="{'selected' : (!singleMode && addedUsers.findIndex((p) => p.uuid === user.uuid) !== -1)}"
						 @click="(!singleMode) ? selectParticipant(user) : $emit('addParticipant', user.uuid)"
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
			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-secondary rounded-pill"
						@click="clearSelections">
					{{ closePlaceholder }}
				</button>
				<button type="submit" class="btn btn-primary rounded-pill">
					{{ submitPlaceholder }}
				</button>
			</div>
		</div>
	</div>
</template>

<style scoped>
.selected, .list-card:hover {
	background: var(--SECONDARY-COLOR);
	border: none;
}


.modal-body {
	display: flex;
	flex-flow: column nowrap;
	height: 400px;
	width: 500px;
	padding: 0;
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

.cards-list {
	display: flex;
	flex-flow: column nowrap;
	overflow-y: scroll;
	overflow-x: hidden;
	height: 100%;
	width: 100%;
}

.list-card {
	display: flex;
	height: 4rem;
	width: 100%;
	flex-flow: row nowrap;
	align-items: center;
	border-bottom: 1px solid #e4e4e4;

	&:hover {
		cursor: pointer;
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
