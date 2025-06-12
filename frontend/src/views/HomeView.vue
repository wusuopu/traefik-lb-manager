<template>
  <div class="flex mb-3 section-box">
    <el-button type="primary" @click="handleAddWorkspace">Create workspace</el-button>
  </div>

  <div class="section-box mb-3">

  </div>

  <el-drawer v-model="state.workspaceForm.showDrawer" direction="rtl" class="!w-[90%] max-w-[600px]">
    <template #header>
      <h4 v-if="state.workspaceForm.action === 'update'">Update workspace #{{ state.workspaceForm.data.id }}</h4>
      <h4 v-else>Create workspace</h4>
    </template>

    <template #default>

    </template>

    <template #footer>
      <div style="flex: auto">
        <el-button @click="state.workspaceForm.showDrawer = false">Cancel</el-button>
        <el-button v-loading.fullscreen.lock="state.workspaceForm.loading" type="primary" @click="handleSubmitWorkspace">
          {{ state.workspaceForm.action == 'create' ? 'Create' : 'Update'}}
        </el-button>
      </div>
    </template>
  </el-drawer>
</template>

<script setup lang="ts">
import { reactive, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useWorkspaceStore } from '@/stores/workspace';

const router = useRouter()
const workspaceStore = useWorkspaceStore()

const state = reactive({
  workspaceForm: {
    showDrawer: false,
    loading: false,
    data: {
    } as any,
    rules: {
    },
    action: '',   // create or update
  },
})

onMounted(() => {
});

const checkAllowDrag = (node: any) => {
  // only allow drag for route
  return node.level === 2
}
const checkAllowDrop = (draggingNode: any, dropNode: any, type: string) => {
  if (dropNode.level === 0) {
    return false
  }
  if (dropNode.level === 1 && type !== 'inner') {
    // route cannot drop to workspace
    return false
  }
  if (dropNode.level === 2 && type === 'inner') {
    // route cannot contain children
    return false
  }
  console.log(draggingNode, draggingNode.data, dropNode.level, dropNode.data, type)
  return true
}

const handleAddWorkspace = () => {
  state.workspaceForm.action = 'create'
  state.workspaceForm.showDrawer = true
}
const handleEditWorkspace = (workspace: any) => {
  state.workspaceForm.action = 'update'
  state.workspaceForm.showDrawer = true
  state.workspaceForm.data = workspace
}
const handleDeleteWorkspace = (workspace: any) => {
}
const handleSubmitWorkspace = async () => {
}
</script>

<style scoped>
</style>