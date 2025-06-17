<template>
  <TopInfo :workspace="workspaceStore.currentWorkspace">
    <el-button type="primary" @click="handleAdd">Add Server</el-button>
  </TopInfo>

  <div class="section-box-dark mb-3">
    <el-table :data="serverStore.servers" style="width: 100%">
      <el-table-column type="expand">
        <template #default="props">
          <RuleManage :server="props.row" />
        </template>
      </el-table-column>

      <el-table-column prop="Name" label="Name" width="150" />
      <el-table-column prop="Enable" label="Enable" width="250" />
      <el-table-column prop="Host" label="Host" width="250" />
      <el-table-column prop="CreatedAt" label="CreatedAt" width="250" :formatter="format.tableDatetimeFormat" />
      <el-table-column prop="UpdatedAt" label="UpdatedAt" width="250" :formatter="format.tableDatetimeFormat" />
      <el-table-column fixed="right" label="Operations" min-width="150">
        <template #default="scope">
          <el-popconfirm @confirm="handleDelete(scope.row.ID)" title="Are you sure to delete this record?">
            <template #reference>
              <el-button v-loading.fullscreen.lock="state.loading" type="danger" size="small">Delete</el-button>
            </template>
          </el-popconfirm>
          <el-button type="primary" size="small" @click="handleEdit(scope.row)">
            Edit
          </el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>

  <el-drawer v-model="state.form.showDrawer" direction="rtl" class="!w-[90%] max-w-[600px]">
    <template #header>
      <h4 v-if="state.form.action === 'update'">Update server #{{ state.form.data.ID }}</h4>
      <h4 v-else>Create server</h4>
    </template>

    <template #default>
      <VueForm
        v-model="state.form.data"
        ref="formRef"
        :schema="state.schema"
        :uiSchema="state.uiSchema"
        :formFooter="{show: false}"
      />
    </template>

    <template #footer>
      <div style="flex: auto">
        <el-button @click="state.form.showDrawer = false">Cancel</el-button>
        <el-button v-loading.fullscreen.lock="state.form.loading" type="primary" @click="handleSubmit">
          {{ state.form.action == 'create' ? 'Create' : 'Update'}}
        </el-button>
      </div>
    </template>
  </el-drawer>
</template>

<script setup lang="ts">
import _ from 'lodash';
import { reactive, onMounted, ref } from 'vue';
import VueForm from '@lljj/vue3-form-element';
import { ElMessage } from 'element-plus';
import { Delete, Plus, Edit } from '@element-plus/icons-vue'
import { useWorkspaceStore } from '@/stores/workspace';
import { useServerStore } from '@/stores/servers';
import TopInfo from './TopInfo.vue';
import RuleManage from './RuleManage.vue';
import format from '@/lib/format';

const workspaceStore = useWorkspaceStore()
const serverStore = useServerStore()
const formRef = ref()

const state = reactive({
  loading: false,
  form: {
    showDrawer: false,
    loading: false,
    data: {
    } as any,
    rules: {
    },
    action: '',   // create or update
  },
  schema: {
    type: 'object',
    required: [
      "Name",
      "Host",
    ],
    properties: {
      Name: {
        title: 'Name',
        type: 'string',
        minLength: 3,
      },
      Host: {
        title: 'Host',
        type: 'array',
        uniqueItems: true,
        minItems: 1,
        items: {
          type: 'string',
        },
      },
      Enable: {
        title: 'Enable',
        type: 'boolean',
      },
    }
  },
  uiSchema: {
    Host: {
      items: {
        "ui:options": {
          placeholder: 'example.com',
        },
      }
    },
  }
})

onMounted(async () => {
  await handleFetchList()
})

const handleFetchList = () => {
  return serverStore.fetchIndexAsync(workspaceStore.detail?.ID!)
}


const handleAdd = () => {
  state.form.action = 'create'
  state.form.data = {
    Name: '',
    Host: [],
    Enable: true,
  }

  formRef.value?.$$uiFormRef?.resetFields()
  state.form.showDrawer = true
}
const handleEdit = (row: Server) => {
  state.form.action = 'update'
  state.form.data = {...row}

  formRef.value?.$$uiFormRef?.resetFields()
  state.form.showDrawer = true
}
const handleDelete = async (id: any) => {
  try {
    await serverStore.deleteAsync(workspaceStore.detail?.ID!, id)
    ElMessage.success('Server has deleted')
  } catch (error: any) {
    ElMessage.error(error.response.data.Error)
    return
  }

  await handleFetchList()
}
const handleSubmit = async () => {
  await formRef.value.$$uiFormRef.validate()

  const payload = {...state.form.data}

  try {
    state.form.loading = true

    if (state.form.action === 'create') {
      await serverStore.createAsync(workspaceStore.detail?.ID!, payload)
      ElMessage.success('Server has created')
    } else {
      await serverStore.updateAsync(workspaceStore.detail?.ID!, payload)
      ElMessage.success('Server has changed')
    }
  } catch (error: any) {
    ElMessage.error(error.response.data.Error)
    return
  } finally {
    state.form.loading = false
  }

  state.form.showDrawer = false

  await handleFetchList()
}
</script>

<style scoped>
</style>