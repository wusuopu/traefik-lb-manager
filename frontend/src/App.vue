<template>
  <el-container class="min-h-[100vh] mx-auto">
    <el-header height="60px" class="">
      <div class="w-full mx-auto text-white">
        <el-menu v-if="!_.isEmpty(workspaceStore.currentWorkspace)" mode="horizontal" :default-active="$route.name as string">
          <el-menu-item index="0" route="" class="!mr-auto">
            <el-select :model-value="workspaceStore.currentWorkspace?.ID" placeholder="Select workspace" class="!w-[140px]" @change="handleWorkspaceChange">
              <template #header>All workspaces</template>
              <el-option v-for="item in workspaceStore.workspaces" :key="item.ID" :label="item.Name" :value="item.ID!" />
            </el-select>
          </el-menu-item>

          <el-menu-item
            v-for="(item, index) in menus"
            :key="index"
            :index="item.name"
            @click="handleClickMenu(item)"
          >{{ item.label }}</el-menu-item>
        </el-menu>
        <h2 v-else class="header-title">
          <RouterLink to="/">Traefik LB Manager</RouterLink>
        </h2>
      </div>
    </el-header>


    <el-main>
      <div class="w-full mx-auto">
        <RouterView></RouterView>
      </div>
    </el-main>

    <el-footer height="40px">
      <div class="flex justify-between mx-auto leading-10 text-white">
        <span>{{ version }}</span>
        <span><a href="https://github.com/wusuopu/traefik-lb-manager" target="_blank">Fork me on Github</a></span>
      </div>
    </el-footer>
  </el-container>
</template>

<script setup lang="ts">
import _ from 'lodash';
import { computed, onMounted, reactive } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { RouterLink, RouterView } from 'vue-router'
import { useWorkspaceStore } from '@/stores/workspace';

const version = __APP_VERSION__

const router = useRouter()
const route = useRoute()
const workspaceStore = useWorkspaceStore()

const menus = computed(() => {
  const items = [
    {label: 'Home', name: 'Home'},
    {label: 'Rules', name: 'Rules'},
    {label: 'Services', name: 'Services'},
  ]
  if (workspaceStore.detail && workspaceStore.detail?.Category !== 'custom') {
    items.push({label: 'Authentications', name: 'Authentications'})
  }
  items.push({label: 'SSL Certificates', name: 'Certificates'})
  return items
})

onMounted(async () => {
  await workspaceStore.fetchIndexAsync()
  if (_.startsWith(route.path, '/workspaces/') && route.params.id) {
    let obj = _.find(workspaceStore.workspaces, {ID: Number(route.params.id)})
    if (obj) { workspaceStore.setCurrentWorkspace(obj) }
  }
})

const handleWorkspaceChange = (id: number) => {
    let obj = _.find(workspaceStore.workspaces, {ID: id})
    if (obj) {
      workspaceStore.setCurrentWorkspace(obj)
      router.push(`/workspaces/${id}/rules`)
    }
}
const handleClickMenu = (menu: any) => {
  if (menu.name === 'Home') {
    workspaceStore.setCurrentWorkspace(null)
  }
  router.push({name: menu.name, params: {id: workspaceStore.currentWorkspace?.ID}})
}

</script>

<style scoped>
.header-title {
  border-bottom: 1px solid var(--el-menu-border-color);
  font-weight: bold;
  font-size: var(--text-2xl);
  line-height: 60px;
}
.el-footer {
  background-color: var(--el-fill-color-lighter);
}
</style>