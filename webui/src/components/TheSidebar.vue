<script setup>
import {ref, reactive} from "vue"
import UserConversationService from "../services/userConversation"
import UserService from "../services/userService";
import {getApiUrl} from "../services/axios";
import TheConversationList from "./TheConversationList.vue";
import TheProfile from "./TheProfile.vue";
import {storeToRefs} from "pinia";
import {useProfileStore} from "@/stores/profileStore";

const profileStore = useProfileStore()
const {profile} = storeToRefs(profileStore)

const component = ref(TheConversationList)
</script>

<template>
	<div class="sidebar">
		<div class="actions-bar">
			<div class="action-box avatar-box" @click="component = TheProfile">
				<img class="avatar" :src="getApiUrl(profile.photo.fullUrl)"/>
			</div>
			<div class="action-box" @click="component = TheConversationList">
				<img class="" src="@/assets/images/chat.svg"/>
			</div>
		</div>

		<div class="sidebar-content">
			<component :is="component"></component>
		</div>
	</div>
</template>

<style scoped>
.sidebar {
	width: 40%;
	height: 100%;
	display: flex;
	flex-flow: row nowrap;
	border-right: 1px gray solid
}

.actions-bar {
	width: 5rem;
	height: 100%;
	border-right: 1px gray solid;
	flex-shrink: 0;
	display: flex;
	flex-flow: column nowrap;
	gap: 1rem;
	align-items: center;
	padding: 1rem 0;
	overflow: hidden;
}

.action-box {
	width: 2rem;
	height: 2rem;
	flex-shrink: 0;
}

.avatar-box {
	width: 4rem;
	height: 4rem;
	margin-bottom: 2rem;
}

.avatar {
	width: 100%;
	height: 100%;
	object-fit: cover;
	border-radius: 50%;
}

.sidebar-content {
	width: 100%;
	height: 100%;
}
</style>
