<template>
  <Container :key="route.params.id as string">
    <CertManage />
  </Container>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router';
import Container from '@/components/workspaces/Container.vue';
import CertManage from '@/components/certificates/CertManage.vue';

const route = useRoute()
</script>

<style lang="scss" scoped>
</style>
