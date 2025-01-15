<script setup>
import {getApiUrl} from "../services/axios";
import {ref} from "vue";

const props = defineProps({
	message: Object,
	isAuthor: Boolean,
})

const emits = defineEmits(["reply", "delete", "forward"])

const sendAt = ref(new Date(props.message.sendAt))


function replyTo() {
	emits("reply", props.message)
}

function deleteMessage() {
	emits("delete", props.message)
}

function forward() {
	emits("forward", props.message)
}

</script>

<template>
	<div class="message-row" :class="{'author-message-row': isAuthor}">
		<div class="message-box" :class="{'author-message-box': isAuthor}">
			<div class="message">
				<div v-if="!isAuthor" class="header">
					<div class="sender-avatar-box">
						<img class="sender-avatar" :src="getApiUrl(message.author.photo.fullUrl)"
							 :width="message.author.photo.width" :height="message.author.photo.height"/>
					</div>
					<div class="sender-name-box">
						<span class="sender-name">{{ message.author.username }}</span>
					</div>
				</div>
				<div class="body">
					<div class="content">
						<div v-if="message.attachment !== null" class="message-image-box"
							 :class="{'mb-2' : message.content !== null}">
							<img class="message-image" :src="getApiUrl(message.attachment.fullUrl)"
								 :width="message.attachment.fullUrl.width" :height="message.attachment.fullUrl.height"/>
						</div>
						<span v-if="message.content !== null" class="message-text">
						{{ message.content }}
					</span>
					</div>
					<div class="info-box">
					<span class="send-at">
						{{ sendAt.toLocaleString([], {dateStyle: 'short', timeStyle: 'short'}) }}
					</span>
						<span class="checkmark-box" v-if="isAuthor">
						<img v-if="message.status === 'delivered'" class="checkmark"
							 src="@/assets/images/Checkmark.png" width="512" height="512"/>
						<img v-else-if="message.status === 'seen'" class="checkmark"
							 src="@/assets/images/seen.png" width="512" height="512"/>
					</span>
					</div>
				</div>
			</div>
		</div>
		<div class="options-box">
			<div class="option" @click="replyTo">
				<img class="option-icon" src="@/assets/images/reply.png" alt=""/>
			</div>
			<div class="option" @click="deleteMessage">
				<img class="option-icon" src="@/assets/images/bin.png" alt=""/>
			</div>
			<div class="option" @click="forward">
				<img class="option-icon" src="@/assets/images/forward.png" alt=""/>
			</div>
		</div>
	</div>
</template>

<style scoped>
.author-message-row {
	flex-direction: row-reverse !important;
}

.author-message-box {
	justify-content: end;
}

.message-row {
	display: flex;
	flex-flow: row nowrap;
	align-items: center;
	width: 100%;
	min-height: 1rem;
	margin-bottom: 10px;
	padding: 0 1rem;
}

.options-box {
	display: flex;
	flex-shrink: 0;
	flex-flow: row nowrap;
	align-items: center;
	height: 100%;
	gap: .7rem;
}

.option {
	width: 2.5rem;
	height: 2.5rem;
	outline: 1px solid black;
	border-radius: 50%;
	cursor: pointer;
}

.option-icon {
	width: 100%;
	height: 100%;
	scale: 0.8
}


.message-box {
	display: flex;
	height: 100%;
	width: 100%;
}

.message {
	display: flex;
	flex-flow: column nowrap;
	padding: 5px;
	max-width: 60%;
}

.header {
	align-items: center;
	display: flex;
	flex-flow: row nowrap;
	margin-bottom: 0.5rem;
}

.sender-avatar-box {
	flex-shrink: 0;
	height: 3rem;
	width: 3rem;
}

.sender-avatar {
	border-radius: 50%;
	height: 100%;
	width: 100%;
}

.sender-name-box {
	display: flex;
	margin-left: 10px;
	overflow: hidden;
	width: 100%;
}

.sender-name {
	font-size: 1.2rem;
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
	width: 100%;
}

.body {
	background-color: var(--TERTIARY-COLOR);
	border-radius: var(--MAIN-BORDER-RADIUS);
	min-height: 2rem;
	padding: .7rem;
	width: fit-content;
}

.content {
	margin-bottom: 5px;
}

.message-text {
	font-size: 1.2rem;
}

.message-image-box {
	max-width: 336px;
	overflow: hidden;
	background-color: #fff;
	border-radius: 4px;
}

.message-image {
	width: 100%;
}

.info-box {
	display: flex;
	flex-flow: row nowrap;
	align-items: center;
	width: 100%;
	justify-content: space-between;
}

.send-at {
	font-size: .7rem;
}

.checkmark-box {
	width: 1.1rem;
	height: 1.1rem;
}

.checkmark {
	width: 100%;
	height: 100%;
}

</style>
