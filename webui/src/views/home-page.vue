<script setup>
import TheSidebar from "../components/TheSidebar.vue";
</script>

<template>
	<main>
		<the-sidebar/>

		<router-view class="overflow-hidden"/>
	</main>
</template>

<style scoped>
main {
	display: flex;
	flex-flow: row nowrap;
}
</style>
