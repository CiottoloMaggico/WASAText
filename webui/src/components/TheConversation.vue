<script setup>
import {getApiUrl} from "../services/axios";
import {getAuthentication} from "../services/sessionService";
import {computed} from "vue";

const props = defineProps({
	conversation: Object,
})

const latestMessage = computed(() => props.conversation.latestMessage)
</script>

<template>
	<div class="conversation-card">
		<picture class="conversation-image-box">
			<img :src="getApiUrl(conversation.image.fullUrl)" :width="conversation.image.width"
				 :height="conversation.image.height" class="conversation-image"/>
		</picture>
		<div class="conversation-card-body">
			<div class="conversation-title-box">
				<span class="conversation-title">{{ conversation.name }}</span>
				<span v-if="!conversation.read" class="unread-dot"/>
			</div>
			<div class="conversation-latestMessage-box" v-if="latestMessage">
				<span class="checkmark-box" v-if="latestMessage.author.uuid === getAuthentication()">
					<img v-if="latestMessage.status === 'delivered'" class="checkmark"
						 src="@/assets/images/Checkmark.png" width="512" height="512"/>
					<img v-else-if="latestMessage.status === 'seen'" class="checkmark"
						 src="@/assets/images/seen.png" width="512" height="512"/>
				</span>
				<div class="latestMessage-content-box">
					<span class="latestMessage-content-author">
						{{ conversation.latestMessage.author.username }}:&nbsp;
					</span>
						<span v-if="conversation.latestMessage.content !== null"
							  class="latestMessage-content-text">{{ conversation.latestMessage.content }}</span>
						<span v-else
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
	height: 100%;
	overflow: hidden;
	width: 100%;
	padding: 0 .5rem;
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
	width: 100%;
}

.conversation-title {
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
	width: 100%;
}

.unread-dot {
	width: 1rem;
	height: 1rem;
	background-color: var(--PRIMARY-COLOR);
	border-radius: 50%;
	flex-shrink: 0;
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
