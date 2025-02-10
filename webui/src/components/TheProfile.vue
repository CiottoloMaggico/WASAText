<script setup>
import {ref, watch, watchEffect, useTemplateRef, nextTick, onBeforeMount, reactive, computed} from "vue"
import {getApiUrl} from "../services/axios";
import {useProfileStore} from "@/stores/profileStore";
import {storeToRefs} from 'pinia'

const newNameField = useTemplateRef("profile-username")
const newImageField = useTemplateRef("file-upload")

const profileStore = useProfileStore()
const { profile } = storeToRefs(profileStore)

const loading = ref(false)
const editUsername = ref(false)
const newUsername = ref(profile.value.username)
const newProfileImage = ref(null)
const editImage = ref(false)

const newImagePreviewUrl = computed(() => {
	if (!newProfileImage.value) {
		return ""
	}
	return URL.createObjectURL(newProfileImage.value)
})

async function changeUsername() {
	loading.value = true
	await profileStore.changeUsername(newUsername.value)

	clearUsernameChange()
	loading.value = false
}

function initUsernameChange() {
	editUsername.value = true
	newUsername.value = profile.value.username

	nextTick(() => {
		newNameField.value.focus()
	})
}

function clearUsernameChange() {
	editUsername.value = false
}

async function changePhoto() {
	loading.value = true
	await profileStore.changeAvatar(newProfileImage.value)

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
	<div v-if="!loading" class="sidebar-content">
		<div class="sidebar-header">
			<h3 class="header-title">Profile</h3>
		</div>

		<div class="sidebar-body">
			<div
				class="avatar-box flex-shrink-0"
				data-bs-target="#image-modal"
				data-bs-toggle="modal"
				@click="initChangePhoto"
			>
				<img
					:src="(!editImage || !newImagePreviewUrl) ? getApiUrl(profile.photo.fullUrl) : newImagePreviewUrl"
					:width="profile.photo.width"
					:height="profile.photo.height"
					class="avatar me-3"
				/>
				<img class="form-icon image-edit" src="@/assets/images/edit-pencil.svg"/>
			</div>
			<span class="profile-username-box" v-if="!editUsername">
					<span class="profile-username d-inline mb-1">{{ profile.username }}</span>
					<img @click="initUsernameChange" class="d-inline form-icon" src="@/assets/images/edit-pencil.svg"/>
				</span>
			<form @submit.prevent="changeUsername" class="profile-username-box" v-else>
				<input
					v-model="newUsername"
					ref="profile-username"
					class="profile-username profile-username-input d-inline mb-1"
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
						<h5 class="modal-title" id="imageModalLabel">Change image</h5>
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

.sidebar-body {
	display: flex;
	flex-flow: column nowrap;
	height: 100%;
	overflow: hidden;
	align-items: center;
	gap: 1rem;
	padding: 1rem;
}

.profile-username-box {
	display: flex;
	flex-flow: row nowrap;
	align-items: center;
	gap: 1rem;
	overflow: hidden;
}

.profile-username {
	font-size: 2rem;
	color: #000;
}

.profile-username-input {
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

</style>
