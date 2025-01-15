<script setup>
import {getApiUrl} from "../services/axios";
import {getAuthentication} from "../services/session";

const props = defineProps({
	conversation: Object,
	selected: Boolean,
})
</script>

<template>
	<div class="conversation-card" :class="{active: selected}">
		<picture class="conversation-image-box">
			<img :src="getApiUrl(conversation.image.fullUrl)" :width="conversation.image.width"
				 :height="conversation.image.height" class="conversation-image"/>
		</picture>
		<div class="conversation-card-body">
			<div class="conversation-title-box">
				<span class="conversation-title">{{ conversation.name }}</span>
				<img v-if="!conversation.read" class="unread-dot" src="@/assets/images/icons8-green-dot-48.png"
					 width="48" height="48"/>
			</div>
			<div class="conversation-latestMessage-box" v-if="conversation.latestMessage !== null">
				<span class="checkmark-box" v-if="conversation.latestMessage.author.uuid === getAuthentication()">
					<img v-if="conversation.latestMessage.status === 'delivered'" class="checkmark"
						 src="@/assets/images/Checkmark.png" width="512" height="512"/>
					<img v-else-if="conversation.latestMessage.status === 'seen'" class="checkmark"
						 src="@/assets/images/seen.png" width="512" height="512"/>
				</span>
				<div class="latestMessage-content-box">
					<span class="latestMessage-content-author">
						{{ conversation.latestMessage.author.username }}:&nbsp
					</span>
						<span v-if="conversation.latestMessage.content !== null"
							  class="latestMessage-content-text">{{ conversation.latestMessage.content }}</span>
						<span v-if="conversation.latestMessage.attachment !== null"
							  class="latestMessage-content-text">image</span>
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped>
.conversation-card {
	align-items: center;
	display: flex;
	flex-flow: row nowrap;
	height: 5rem;
	overflow: hidden;
	width: 100%;
	padding-left: 10px;
}

.conversation-card:hover{
	cursor: pointer;
	background: var(--SECONDARY-COLOR);
}

.active {
	background-color: var(--SECONDARY-COLOR);
}

.conversation-image-box {
	align-items: center;
	display: flex;
	flex-shrink: 0;
	height: 4rem;
	justify-content: center;
	width: 4rem;
}

.conversation-image {
	border-radius: 50%;
	height: 100%;
	object-fit: cover;
	width: 100%;
}

.conversation-card-body {
	align-items: center;
	border-bottom: 1px solid rgba(128, 128, 128, 0.26);
	display: flex;
	flex-flow: column nowrap;
	height: 100%;
	justify-content: center;
	margin-left: 10px;
	max-width: 100%;
	overflow: hidden;
	width: 100%;
}

.conversation-title-box {
	display: flex;
	flex-flow: row nowrap;
	flex-shrink: 0;
	width: 100%;
}

.conversation-title {
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
	width: 100%;
}

.unread-dot {
	width: 1.5rem;
	height: 1.5rem;
}

.conversation-latestMessage-box {
	display: flex;
	align-items: center;
	width: 100%;
	flex-grow: 0;
}

.checkmark-box {
	flex-shrink: 0;
	width: 1rem;
	height: 1rem;
	margin-right: 5px;
}

.checkmark {
	width: 100%;
	height: 100%;
}

.latestMessage-content-box {
	display: flex;
	width: 100%;
	overflow: hidden;
}

.latestMessage-content-author {
	flex-shrink: 0;
}

.latestMessage-content-author, .latestMessage-content-text {
	overflow: hidden;
	white-space: nowrap;
	text-overflow: ellipsis;
}

</style>
