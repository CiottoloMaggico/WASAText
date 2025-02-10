<script setup>
import {computed, reactive, ref, useTemplateRef, watch} from "vue"
import UserService from "@/services/userService";
import {getAuthentication} from "@/services/session";
import ConversationService from "@/services/conversationService";
import {useProfileStore} from "@/stores/profileStore";
import TheConversationList from "@/components/TheConversationList.vue";
import router from "@/router";
import AddParticipantForm from "@/components/AddParticipantForm.vue";

const emits = defineEmits(["switch"])

const profileStore = useProfileStore()
const newImageField = useTemplateRef("file-upload")

const newGroup = reactive({
	name: null,
	image: null,
	participants: [],
})

const newImagePreviewUrl = computed(() => {
	if (!newGroup.image) {
		return "/src/assets/images/default_group_image.jpg"
	}
	return URL.createObjectURL(newGroup.image)
})

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
						class="text-input d-inline mb-1"
						minlength="3"
						maxlength="16"
						placeholder="Group name"
						required
					/>
				</div>
				<AddParticipantForm :participants="newGroup.participants" :single-mode="false"/>
			</div>
			<div class="form-footer">
				<button type="submit" class="btn btn-primary">Crea gruppo</button>
			</div>
		</form>
	</div>
</template>

<style scoped>

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

.text-input {
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

.text-input:focus {
	box-shadow: none;
	border: none;
	outline: none;
}

.image-edit {
	display: none;
	position: absolute;
	top: calc(50% - 11.5px);
	left: calc(50% - 7px);
}

.avatar-box {
	position: relative;
	width: 250px;
	height: 250px;

	&:hover {
		filter: brightness(.8);

		.image-edit {
			display: block;
		}
	}
}

.avatar {
	width: 100%;
	height: 100%;
	border-radius: 50%;
	object-fit: cover;
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
