<script setup>
import {ref, watch, watchEffect, useTemplateRef, nextTick, onBeforeMount, reactive, computed} from "vue"
import UserService from "../services/userService";
import {getApiUrl} from "../services/axios";

const loading = ref(true)
const profile = reactive({})

const editUsername = ref(false)
const newUsername = ref("")
const newNameField = useTemplateRef("profile-username")

const newProfileImage = ref(null)
const editImage = ref(false)
const newImageField = useTemplateRef("file-upload")
const newImagePreviewUrl = computed(() => {
	if (!newProfileImage.value) {
		return ""
	}
	return URL.createObjectURL(newProfileImage.value)
})

await getUserProfile()

async function getUserProfile() {
	loading.value = true
	try {
		const response = await UserService.getProfile()
		Object.assign(profile, response.data)
		newUsername.value = profile.username
	} catch (error) {
		console.error(error.toString())
	}
	loading.value = false
}

async function changeUsername() {
	loading.value = true
	try {
		const response = await UserService.setMyUsername(newUsername.value)
		Object.assign(profile, response.data)
	} catch (error) {
		console.error(error.toString())
	}

	clearUsernameChange()
	loading.value = false
}

function initUsernameChange() {
	editUsername.value = true

	nextTick(() => {
		newNameField.value.focus()
	})
}

function clearUsernameChange() {
	editUsername.value = false
	newUsername.value = profile.username
}

async function changePhoto() {
	if (!newProfileImage.value) {
		return
	}

	loading.value = true
	try {
		const response = await UserService.setMyPhoto(newProfileImage.value)
		Object.assign(profile, response.data)
	} catch (err) {
		console.log(err.toString())
	}

	editImage.value = false
	newProfileImage.value = null
	loading.value = false
}

function initChangePhoto() {
	editImage.value = true
}

function clearImageChange() {
	editImage.value = false
	newProfileImage.value = null
	newImageField.value.value = null
}

function fileUploaded() {
	newProfileImage.value = newImageField.value.files.item(0)
}

</script>

<template>
	<div class="detail-page" v-if="!loading">
		<div class="page-header mb-4">
			<h3 class="text-primary">Profile</h3>
		</div>

		<div class="page-body">
			<div class="mb-4 d-flex align-items-center gap-2">
				<div
					class="group-image-box flex-shrink-0"
					data-bs-target="#image-modal"
					data-bs-toggle="modal"
					@click="initChangePhoto"
				>
					<img
						:src="(!editImage || !newImagePreviewUrl) ? getApiUrl(profile.photo.fullUrl) : newImagePreviewUrl"
						:width="profile.photo.width"
						:height="profile.photo.height"
						alt="Gruppo"
						class="group-image me-3"
					/>
					<img class="form-icon image-edit" src="@/assets/images/edit-pencil.svg"/>
				</div>
				<span class="group-name-box" v-if="!editUsername">
					<span class="group-name d-inline mb-1">{{ profile.username }}</span>
					<img @click="initUsernameChange" class="d-inline form-icon" src="@/assets/images/edit-pencil.svg"/>
				</span>
				<form @submit.prevent="changeUsername" class="group-name-box" v-else>
					<input
						v-model="newUsername"
						ref="profile-username"
						class="group-name group-name-input d-inline mb-1"
						minlength="3"
						maxlength="16"
						required
					/>
					<label for="change-name">
						<img class="d-inline form-icon" src="@/assets/images/checkmark.svg"/>
					</label>
					<input type="submit" id="change-name" class="d-none"/>
					<img @click="clearUsernameChange" class="d-inline form-icon" src="@/assets/images/cancel.svg"/>
				</form>
			</div>
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
									:disabled="!newProfileImage">Change image
							</button>
						</div>
					</form>
				</div>
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
	top: calc(50% - 11.5px);
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
