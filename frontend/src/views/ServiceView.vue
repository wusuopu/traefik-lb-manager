<template>
  <Container :key="route.params.id as string">
    <ServiceManage />
  </Container>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router';
import Container from '@/components/workspaces/Container.vue';
import ServiceManage from '@/components/services/ServiceManage.vue';

const route = useRoute()
</script>

<style lang="scss" scoped>
</style>
