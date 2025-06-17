<template>
  <div class="flex mb-3 section-box">
    <el-button type="primary" @click="handleAddWorkspace">Create workspace</el-button>
  </div>

  <div class="section-box-dark mb-3">
    <el-table :data="workspaceStore.workspaces" style="width: 100%">
        <el-table-column prop="Name" label="Name" width="150">
          <template #default="scope">
            <el-button link type="primary" @click="handleSelectWorkspace(scope.row)">{{ scope.row.Name }}</el-button>
          </template>
        </el-table-column>
        <el-table-column prop="Category" label="Category" width="150" />
        <el-table-column prop="Entrypoints" label="Entrypoints" width="150" />
        <el-table-column prop="ApiBaseUrl" label="ApiBaseUrl" width="450" />
        <el-table-column prop="CreatedAt" label="CreatedAt" width="250" :formatter="format.tableDatetimeFormat" />
        <el-table-column prop="UpdatedAt" label="UpdatedAt" width="250" :formatter="format.tableDatetimeFormat" />
        <el-table-column fixed="right" label="Operations" min-width="120">
          <template #default="scope">
            <el-button type="primary" size="small" @click="handleEditWorkspace(scope.row)">
              Edit
            </el-button>
          </template>
        </el-table-column>
      </el-table>
  </div>

  <el-drawer v-model="state.form.showDrawer" direction="rtl" class="!w-[90%] max-w-[600px]">
    <template #header>
      <h4 v-if="state.form.action === 'update'">Update workspace #{{ state.form.data.ID }}</h4>
      <h4 v-else>Create workspace</h4>
    </template>

    <template #default>
        <div>
          <el-form ref="formRef" :model="state.form.data" :rules="state.form.rules" label-position="top">
            <el-form-item label="Name" prop="Name">
              <el-input v-model="state.form.data.Name" placeholder="Name" />
            </el-form-item>

            <el-form-item label="Description" prop="Description">
              <el-input v-model="state.form.data.Description" placeholder="Description" />
            </el-form-item>

            <el-form-item label="ManagerBaseUrl" prop="ManagerBaseUrl" required>
              <el-input v-model="state.form.data.ManagerBaseUrl" placeholder="http://manager:8080" />
            </el-form-item>
            <el-alert type="info" show-icon :closable="false">
              <p>This will use to receive request from Let's Encrypt</p>
            </el-alert>

            <el-form-item label="Category" prop="Category" required>
              <el-radio-group v-model="state.form.data.Category" :disabled="state.form.action === 'update'">
                <el-radio-button label="Rancher V1" value="rancher_v1" />
                <el-radio-button label="Portainer Swarm" value="portainer_swarm" />
                <el-radio-button label="Common" value="common" />
                <el-radio-button label="Custom" value="custom" />
              </el-radio-group>
            </el-form-item>

            <el-form-item label="Entrypoints" prop="Entrypoints" required>
              <el-input-tag v-model="state.form.data.Entrypoints" clearable placeholder="Please input" />
            </el-form-item>
            <el-alert type="info" show-icon :closable="false">
              <p>Please define these entrypoints in traefik Static Configuration</p>
            </el-alert>

            <template v-if="state.form.data.Category === 'rancher_v1' || state.form.data.Category === 'portainer_swarm'">
              <el-form-item label="ApiBaseUrl" prop="ApiBaseUrl" required>
                <el-input
                  v-model="state.form.data.ApiBaseUrl"
                  :placeholder="state.form.data.Category === 'rancher_v1' ? 'http://<host>/v2-beta/projects/<env>' : 'http://host:port/api/endpoints/<envId>'"
                />
              </el-form-item>

              <el-form-item label="ApiKey" prop="ApiKey" required>
                <el-input v-model="state.form.data.ApiKey" placeholder="" />
              </el-form-item>

              <el-form-item v-if="state.form.data.Category === 'rancher_v1'" label="ApiSecret" prop="ApiSecret" required>
                <el-input v-model="state.form.data.ApiSecret" placeholder="" />
              </el-form-item>
            </template>
          </el-form>
        </div>
    </template>

    <template #footer>
      <div style="flex: auto">
        <el-button @click="state.form.showDrawer = false">Cancel</el-button>
        <el-button v-loading.fullscreen.lock="state.form.loading" type="primary" @click="handleSubmitWorkspace">
          {{ state.form.action == 'create' ? 'Create' : 'Update'}}
        </el-button>
      </div>
    </template>
  </el-drawer>
</template>

<script setup lang="ts">
import _ from 'lodash';
import { reactive, ref } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessage, type FormInstance } from 'element-plus'
import { useWorkspaceStore } from '@/stores/workspace';
import format from '@/lib/format';

const router = useRouter()
const workspaceStore = useWorkspaceStore()
const formRef = ref<FormInstance>()

const state = reactive({
  form: {
    showDrawer: false,
    loading: false,
    data: {
      Name: '',
      Description: '',
      ManagerBaseUrl: '',
      Category: '',
      ApiBaseUrl: '',
      ApiKey: '',
      ApiSecret: '',
    } as Workspace,
    rules: {
      Name: [
        { required: true, message: 'Name is required' },
      ],
    },
    action: '',   // create or update
  },
})

const handleAddWorkspace = () => {
  state.form.action = 'create'
  state.form.data = {
    Name: '',
    Description: '',
    ManagerBaseUrl: `${location.protocol}//${location.host}${location.pathname}`,
    Category: 'common',
    ApiBaseUrl: '',
    ApiKey: '',
    ApiSecret: '',
    Entrypoints: ['web', 'websecure'],
  }

  formRef.value?.resetFields()
  state.form.showDrawer = true
}

const handleSelectWorkspace = (workspace: Workspace) => {
  workspaceStore.setCurrentWorkspace(workspace)
  router.push(`/workspaces/${workspace.ID}/rules`)
}
const handleEditWorkspace = (workspace: Workspace) => {
  state.form.action = 'update'
  state.form.data = {...workspace}

  formRef.value?.resetFields()
  state.form.showDrawer = true
}
const handleDeleteWorkspace = (workspace: any) => {
}
const handleSubmitWorkspace = async () => {
  await formRef.value!.validate()

  const payload = {...state.form.data}
  if (payload.Category === 'custom' || payload.Category === 'common') {
    delete payload.ApiBaseUrl
    delete payload.ApiKey
    delete payload.ApiSecret
  } else {
    payload.ApiBaseUrl = _.trimEnd(payload.ApiBaseUrl, '/')
  }
  payload.Entrypoints = _.uniq(payload.Entrypoints)

  try {
    state.form.loading = true

    if (state.form.action === 'create') {
      await workspaceStore.createAsync(payload)
      ElMessage.success('Workspace has created')
    } else {
      await workspaceStore.updateAsync(payload)
      ElMessage.success('Workspace has changed')
    }
  } catch (error: any) {
    ElMessage.error(error.response.data.Error)
    return
  } finally {
    state.form.loading = false
  }

  state.form.showDrawer = false

  await workspaceStore.fetchIndexAsync()
}
</script>

<style scoped>
</style>