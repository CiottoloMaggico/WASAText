<script setup>
import {ref} from "vue"
import {getApiUrl} from "../services/axios";
import TheConversationList from "./TheConversationList.vue";
import TheNewConversation from "@/components/TheNewConversation.vue";
import TheProfile from "@/components/TheProfile.vue";
import {storeToRefs} from "pinia";
import {useProfileStore} from "@/stores/profileStore";

const profileStore = useProfileStore()
const {profile} = storeToRefs(profileStore)

const component = ref(TheConversationList.__name)

function switchComponent(componentName) {
	component.value = componentName
}
</script>

<template>
	<div class="sidebar">
		<div class="actions-bar">
			<div class="actions-group">
				<div class="action-box" @click="switchComponent(TheConversationList.__name)">
					<img class="" src="@/assets/images/chat.svg"/>
				</div>
				<div class="action-box" @click="switchComponent(TheNewConversation.__name)">
					<img class="" src="@/assets/images/add-chat.svg"/>
				</div>
			</div>
			<div class="action-box avatar-box" @click="switchComponent(TheProfile.__name)">
				<img class="avatar" :src="getApiUrl(profile.photo.fullUrl)"/>
			</div>
		</div>

		<component :is="component" @switch="switchComponent"/>
	</div>
</template>

<style scoped>
:deep(.sidebar-header) {
	display: flex;
	flex-flow: row nowrap;
	justify-content: space-between;
	padding: 1rem;
	gap: .5rem;
	align-items: center;
}

:deep(.header-title) {
	font-weight: bolder;
	margin-bottom: 0;
}

:deep(.sidebar-content) {
	display: flex;
	flex-flow: column nowrap;
	width: 100%;
	height: 100%;
}

:deep(.sidebar-body) {
	display: flex;
	flex-flow: column nowrap;
	height: 100%;
	width: 100%;
	overflow: hidden;
	align-items: center;
	max-height: 100%;
}

:deep(.sidebar-group) {
	display: flex;
	width: 100%;
	flex-flow: column nowrap;
	max-height: 100%;
	overflow: hidden;
}

:deep(.sidebar-group-header) {
	width: 100%;
	padding: 1rem;
	display: flex;
	flex-flow: column nowrap;
}

:deep(.sidebar-group-title) {
	font-size: 1.3rem;
	font-weight: 600;
}

:deep(.sidebar-item) {
	display: flex;
	flex-shrink: 0;
	align-items: center;
	width: 100%;
	height: 5rem;
}

:deep(.selected), :deep(.sidebar-item:hover) {
	background: var(--SECONDARY-COLOR);
	border: none;
}

:deep(.sidebar-item:hover) {
	cursor: pointer;
}

:deep(.sidebar-item-title) {
	font-size: 1.2rem;
	font-weight: 600;
}

.sidebar {
	width: 30%;
	height: 100%;
	display: flex;
	flex-flow: row nowrap;
	border-right: 1px gray solid;
	flex-shrink: 0;
}

.actions-bar, .actions-group {
	display: flex;
	flex-flow: column nowrap;
	align-items: center;
	flex-shrink: 0;
}

.actions-bar {
	width: 5rem;
	height: 100%;
	border-right: 1px gray solid;
	padding: 1rem 0;
	justify-content: space-between;
	overflow: hidden;
}

.actions-group {
	width: 100%;
	padding: 1rem 0;
	gap: 1.5rem;
}

.action-box {
	width: 2rem;
	height: 2rem;
	flex-shrink: 0;
}

.avatar-box {
	justify-self: end;
	width: 4rem;
	height: 4rem;
	margin-top: 2rem;
}

.avatar {
	width: 100%;
	height: 100%;
	object-fit: cover;
	border-radius: 50%;
}

</style>
