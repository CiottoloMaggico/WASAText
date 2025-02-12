<script setup>
import {computed, nextTick, reactive, ref, useTemplateRef, watchEffect} from "vue"
import ConversationService from "../services/conversationService";
import {getApiUrl} from "../services/axios";
import router from "../router";
import AddParticipantModal from "@/components/AddParticipantModal.vue";
import {storeToRefs} from "pinia";
import {useConversationsStore} from "@/stores/conversationsStore";
import conversationService from "../services/conversationService";

const {activeConversation} = storeToRefs(useConversationsStore());

const newImageField = useTemplateRef("file-upload")
const newNameField = useTemplateRef("group-name")

const loading = ref(false);
const newGroupName = ref("")
const editName = ref(false);
const newGroupImage = ref(null)

const conversation = reactive({
	id: Number,
	name: String,
	image: Object,
	type: String,
	read: Boolean,
	participants: [],
})

const newImagePreviewUrl = computed(() => {
	if (!newGroupImage.value) {
		return getApiUrl(conversation.image.fullUrl)
	}
	return URL.createObjectURL(newGroupImage.value)
})

watchEffect(async (onCleanup) => {
	if (activeConversation.value) await getConversation()

	onCleanup(() => {
		initializePage()
	})
})

async function getConversation() {
	loading.value = true

	try {
		const data = await ConversationService.getConversation(activeConversation.value.id)
		Object.assign(conversation, data)
		newGroupName.value = conversation.name
	} catch (err) {
		console.log(err.toString())
	}

	loading.value = false
}

async function leaveGroup() {
	try {
		await ConversationService.leaveGroup(conversation)
	} catch (err) {
		console.log(err.toString())
		return
	}
	await router.push({name: "homepage"})
}

async function changeName() {
	loading.value = true

	try {
		const data = await ConversationService.setGroupName(conversation, newGroupName.value)
		Object.assign(conversation, data)
	} catch (err) {
		console.log(err.toString())
	}

	clearNameChange()
	loading.value = false
}

async function changePhoto() {
	loading.value = true

	try {
		const data = await ConversationService.setGroupPhoto(conversation, newGroupImage.value)
		Object.assign(conversation, data)
	} catch (err) {
		console.log(err.toString())
	}

	clearImageChange()
	loading.value = false
}

function updateParticipants(addedParticipants) {
	conversation.participants = conversation.participants.concat(addedParticipants)
}

function initNameChange() {
	editName.value = true

	nextTick(() => {
		newNameField.value.focus()
	})
}

function clearNameChange() {
	editName.value = false
	newGroupName.value = conversation.name
}

function clearImageChange() {
	newGroupImage.value = null
	newImageField.value.value = null
}

function initializePage() {
	loading.value = true
	Object.assign(conversation, {})
	clearImageChange()
	clearNameChange()
}

function fileUploaded() {
	newGroupImage.value = newImageField.value.files.item(0)
}
</script>

<template>
	<div class="detail-page" v-if="!loading">
		<div class="page-header mb-4">
			<h3 class="text-primary">Group detail</h3>
		</div>

		<div class="page-body">
			<div class="mb-4 d-flex align-items-center gap-2">
				<div
					class="group-image-box flex-shrink-0"
					data-bs-target="#image-modal"
					data-bs-toggle="modal"
				>
					<img
						:src="newImagePreviewUrl"
						:width="conversation.image.width"
						:height="conversation.image.height"
						alt="Gruppo"
						class="group-image me-3"
					/>
					<img class="form-icon image-edit" src="@/assets/images/edit-pencil.svg"/>
				</div>
				<span class="group-name-box" v-if="!editName">
					<span class="group-name d-inline mb-1">{{ conversation.name }}</span>
					<img class="d-inline form-icon" src="@/assets/images/edit-pencil.svg" @click="initNameChange"/>
				</span>
				<form class="group-name-box" v-else @submit.prevent="changeName">
					<input
						v-model="newGroupName"
						ref="group-name"
						class="group-name group-name-input d-inline mb-1"
						minlength="3"
						maxlength="16"
						required
					/>
					<label for="change-name">
						<img class="d-inline form-icon" src="@/assets/images/checkmark.svg"/>
					</label>
					<input type="submit" id="change-name" class="d-none"/>
					<img class="d-inline form-icon" src="@/assets/images/cancel.svg" @click="clearNameChange"/>
				</form>
			</div>

			<div class="participants-list-box mb-4">
				<div class="d-flex justify-content-between px-2">
					<label class="form-label fw-semibold">Participants</label>
					<img data-bs-target="#participants-modal"
						 data-bs-toggle="modal"
						 alt="add participants"
						 class="form-icon"
						 src="@/assets/images/plus.svg"
					/>
				</div>
				<div class="participants-list">
					<ul class="list-group list-group-flush">
						<li
							v-for="participant in conversation.participants"
							:key="participant.uuid"
							class="list-group-item d-flex justify-content-between align-items-center border-0"
						>
							<span class="fw-light">{{ participant.username }}</span>
						</li>
					</ul>
				</div>
			</div>
		</div>

		<div class="page-footer d-flex justify-content-between flex-shrink-0">
			<button class="btn btn-danger rounded-pill px-4" @click="leaveGroup">
				Leave group
			</button>
		</div>
	<add-participant-modal :conversation="conversation" :single-mode="false" @added-participants="updateParticipants"/>
	</div>

	<div
			class="modal fade"
			ref="image-modal"
			id="image-modal"
			tabindex="-1"
			aria-labelledby="imageModalLabel"
			aria-hidden="true"
			data-bs-backdrop="static"
			data-bs-keyboard="false"
		>
			<div class="modal-dialog modal-dialog-centered">
				<div class="modal-content">
					<div class="modal-header">
						<h5 class="modal-title" id="imageModalLabel">Change group image</h5>
					</div>
					<form @submit.prevent="changePhoto">
						<div class="modal-body">
							<input type="file"
								   ref="file-upload"
								   class="form-control border-primary rounded-pill"
								   @change="fileUploaded"
							/>
						</div>

						<div class="modal-footer">
							<button type="button" class="btn btn-secondary rounded-pill" data-bs-dismiss="modal"
									@click="clearImageChange">Close
							</button>
							<button type="submit" class="btn btn-primary rounded-pill" data-bs-dismiss="modal"
									:disabled="!newGroupImage">Change image
							</button>
						</div>
					</form>
				</div>
			</div>
		</div>
</template>

<style scoped>
.form-icon {
	width: 1.5rem;
	height: 1.5rem;
	object-fit: contain;
	flex-shrink: 0;
}

.detail-page {
	padding: 1rem;
	display: flex;
	flex-flow: column nowrap;
	width: 100%;
	height: 100%;
	flex-grow: 0;
}

.page-header, .page-footer {
	flex-shrink: 0;
}
.page-body {
	display: flex;
	flex-flow: column nowrap;
	height: 100%;
	overflow: hidden;
}

.group-name-box {
	display: flex;
	flex-flow: row nowrap;
	align-items: center;
	gap: 1rem;
	overflow: hidden;
}

.group-name {
	font-size: 2rem;
	color: #000;
}

.group-name-input {
	width: 290px;
	background-color: transparent;
	border: none;
	box-shadow: none;
	white-space: nowrap;
	text-overflow: ellipsis;
	overflow: hidden;
}

.group-name-input:focus {
	box-shadow: none;
	border: none;
	outline: none;
}

.group-image-box {
	position: relative;
	width: 150px;
	height: 150px;
}

.group-image-box:hover {
	filter: brightness(.8);

	.image-edit {
		display: block;
	}
}


.image-edit {
	display: none;
	position: absolute;
	top: calc(50% - 11px);
	left: calc(50% - 7px);
}

.group-image {
	width: 100%;
	height: 100%;
	border-radius: 50%;
	object-fit: cover;
}

.participants-list-box {
	display: flex;
	flex-flow: column nowrap;
	overflow: hidden;
	height: 100%;
}

.participants-list {
	flex-grow: 0;
	overflow-y: auto;
}

.list-group {
	margin: 0;
	padding: 0;
}

.list-group-item {
	background-color: white;
}


</style>
