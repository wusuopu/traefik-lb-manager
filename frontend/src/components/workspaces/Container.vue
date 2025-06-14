<template>
  <div v-if="workspaceStore.detail?.ID === Number(route.params.id)" class="w-full">
    <slot></slot>
  </div>
</template>
<script setup lang="ts">
import { onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { useWorkspaceStore } from '@/stores/workspace';

const route = useRoute()
const workspaceStore = useWorkspaceStore()

onMounted(async () => {
  if (workspaceStore.detail?.ID !== Number(route.params.id)) {
    await workspaceStore.fetchShowAsync(Number(route.params.id))
  }
})

</script>
<style lang="scss" scoped>
</style>
