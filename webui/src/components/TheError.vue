<script setup>

const props = defineProps(["error"])

</script>

<template>
	<div class="error-bar" v-if="error">
		<div class="error-description-container">
			<span class="error-title">An error occurred</span>
			<span class="error-description">{{ error.message }}</span>
		</div>
	</div>
</template>

<style scoped>
.error-bar {
	border-radius: var(--MAIN-BORDER-RADIUS) var(--MAIN-BORDER-RADIUS) 0 0;
	position: absolute;
	width: 100%;
	height: 5rem;
	background: var(--bs-danger);
	z-index: 1;
}

.error-description-container {
	width: 100%;
	height: 100%;
	display: flex;
	flex-flow: column nowrap;
	justify-content: space-around;
	padding : .5rem 1rem;
}

.error-title {
	color: white;
	font-size: 1.2rem;
	font-weight: 600;
}

.error-description {
	font-size: 1rem;
	color: white;
	font-weight: 600;
}

</style>
