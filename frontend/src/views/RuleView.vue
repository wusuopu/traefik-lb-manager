<template>
  <Container :key="route.params.id as string">
    <RuleEditor v-if="workspaceStore.detail?.Category === 'custom'" />
    <ServerManage v-else />
  </Container>
</template>

<script setup lang="ts">
import _ from 'lodash';
import { useRoute } from 'vue-router';
import { useWorkspaceStore } from '@/stores/workspace';
import Container from '@/components/workspaces/Container.vue';
import ServerManage from '@/components/workspaces/ServerManage.vue';
import RuleEditor from '@/components/workspaces/RuleEditor.vue';

const workspaceStore = useWorkspaceStore()
const route = useRoute()
</script>

<style scoped>
</style>