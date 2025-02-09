<script setup>
</script>

<template>
    <router-view v-slot="{ Component }">
        <template v-if="Component">
                <suspense timeout="0">
                    <template #default>
                        <component :is="Component"/>
                    </template>
<!--                    <template #fallback>-->
<!--                        <base-loading-animation/>-->
<!--                    </template>-->
                </suspense>
        </template>
    </router-view>
</template>
