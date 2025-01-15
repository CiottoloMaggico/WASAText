<script setup>
</script>

<template>
<!--    <suspense>-->
<!--        <keep-alive>-->
<!--            <TheNavbar/>-->
<!--        </keep-alive>-->
<!--    </suspense>-->
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
<!--    <suspense>-->
<!--        <keep-alive>-->
<!--            <TheFooter/>-->
<!--        </keep-alive>-->
<!--    </suspense>-->
</template>
