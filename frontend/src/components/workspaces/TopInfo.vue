<template>
  <div v-if="!_.isEmpty(workspace)" class="section-box mb-3">
    <div class="mb-2 break-all">
      <p class="font-bold text-lg">{{ workspace.Name }} ({{ workspace.Category }})</p>
      <p>{{ workspace.Description }}</p>
      <p v-if="workspace.ApiBaseUrl">{{ workspace.ApiBaseUrl }}</p>
    </div>
    <div class="flex gap-1 items-center">
      <slot></slot>
    </div>
  </div>
</template>

<script setup lang="ts">
import _ from 'lodash';

defineProps<{
  workspace: Workspace | null
}>()
</script>

<style lang="scss" scoped>
</style>
