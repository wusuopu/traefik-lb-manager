<template>
  <el-container class="min-h-[100vh] mx-auto">
    <el-header height="60px" class="">
      <div class="w-full mx-auto text-white">
        <el-menu v-if="!_.isEmpty(workspaceStore.currentWorkspace)" mode="horizontal" :default-active="$route.name">
          <el-menu-item index="0" route="" class="!mr-auto">
            <el-select :model-value="workspaceStore.currentWorkspace?.id" placeholder="Select workspace" class="!w-[140px]">
              <template #header>All workspaces</template>
              <el-option v-for="item in state.workspaces" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </el-menu-item>

          <el-menu-item index="Home" @click="router.push('/')">Home</el-menu-item>
          <el-menu-item
            v-for="(item, index) in menus"
            :key="index"
            :index="item.name"
            @click="handleClickMenu(item)"
          >{{ item.label }}</el-menu-item>
        </el-menu>
        <h2 v-else class="header-title">
          Traefik LB Manager
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
import { reactive } from 'vue';
import { useRouter } from 'vue-router';
import { RouterLink, RouterView } from 'vue-router'
import { useWorkspaceStore } from '@/stores/workspace';

const version = __APP_VERSION__

const router = useRouter()
const workspaceStore = useWorkspaceStore()

const state = reactive({
  workspaces: [],
})

const menus = [
  {label: 'Rules', name: 'Rules'},
  {label: 'Authentications', name: 'Authentications'},
  {label: 'SSL Certificates', name: 'Certificates'},
]


const handleClickMenu = (menu: any) => {
  router.push({name: menu.name, params: {id: workspaceStore.currentWorkspace?.id}})
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