<script setup>
import {getApiUrl} from "../services/axios";

const props = defineProps({
	"conversation" : Object,
})
</script>

<template>
	<div class="chat-info-box">
		<div class="image-box">
			<img :src="getApiUrl(props.conversation.fullUrl)" :width="props.conversation.image.width" :height="props.conversation.image.height"/>
		</div>
	</div>
</template>

<style scoped>
.chat-info-box{
	display: flex;
	flex-flow: row nowrap;
	width: 70%;
	height: 70%;
	background-color: var(--FOURTH-COLOR);
	border-radius: var(--MAIN-BORDER-RADIUS);
}
</style>
