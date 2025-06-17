<template>
  <Container :key="route.params.id as string">
    <Manage />
  </Container>
</template>

<script setup lang="ts">
import _ from 'lodash';
import { useRoute } from 'vue-router';
import Container from '@/components/workspaces/Container.vue';
import Manage from '@/components/middlewares/Manage.vue';

const route = useRoute()
</script>

<style scoped>
</style>