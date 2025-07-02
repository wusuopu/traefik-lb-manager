<template>
  <TopInfo :workspace="workspaceStore.currentWorkspace">
    <template #before>
      <p>
        HTTP Provide URL: <a :href="httpProvideUrl" target="_blank">{{ httpProvideUrl }}</a> <br>
        It needs to generate traefik configuration after change config.
      </p>
    </template>

    <el-button type="primary" @click="handleAdd">Add Server</el-button>
    <el-button type="success" @click="handleRefresh">Refresh</el-button>
    <el-popconfirm @confirm="handleGenerate" title="Are you sure to re-generate traefik config from this rule?">
      <template #reference>
        <el-button type="warning">Generate Traefik Config</el-button>
      </template>
    </el-popconfirm>
  </TopInfo>

  <div class="section-box-dark mb-3">
    <el-table :data="serverStore.servers" border style="width: 100%">
      <el-table-column type="expand">
        <template #default="props">
          <RuleManage :server="props.row" @edit="handleEditRule" />
        </template>
      </el-table-column>

      <el-table-column prop="ID" label="ID" width="100" />
      <el-table-column prop="Name" label="Name" min-width="150" />
      <el-table-column prop="Enable" label="Enable" width="100">
        <template #default="scope">
          <el-switch :model-value="scope.row.Enable" size="small" disabled/>
        </template>
      </el-table-column>
      <el-table-column prop="Host" label="Host" min-width="250">
        <template #default="scope">
          {{ scope.row.Host?.join(" | ") }}
        </template>
      </el-table-column>
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
          <el-button type="success" size="small" @click="handleAddRule(scope.row)">
            Add Rule
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

  <RuleForm v-if="state.loaded" ref="ruleFormRef" />
</template>

<script setup lang="ts">
import _ from 'lodash';
import { reactive, onMounted, ref, computed } from 'vue';
import VueForm from '@lljj/vue3-form-element';
import { ElMessage } from 'element-plus';
import { useWorkspaceStore } from '@/stores/workspace';
import { useMiddlewareStore } from '@/stores/middlewares';
import { useServiceStore } from '@/stores/services';
import { useServerStore } from '@/stores/servers';
import { useRuleStore } from '@/stores/rules';
import TopInfo from './TopInfo.vue';
import RuleManage from './RuleManage.vue';
import RuleForm from './RuleForm.vue';
import format from '@/lib/format';

const workspaceStore = useWorkspaceStore()
const middlewareStore = useMiddlewareStore()
const serviceStore = useServiceStore()
const serverStore = useServerStore()
const ruleStore = useRuleStore()
const formRef = ref()
const ruleFormRef = ref<typeof RuleForm>()

const state = reactive({
  loading: false,
  loaded: false,
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

const httpProvideUrl = computed(() => {
  let baseUrl = workspaceStore.detail?.ManagerBaseUrl
  if (_.endsWith(baseUrl, '/')) { baseUrl = _.trimEnd(baseUrl, '/') }
  baseUrl += location.pathname
  if (_.endsWith(baseUrl, '/')) { baseUrl = _.trimEnd(baseUrl, '/') }
  return `${baseUrl}/workspaces/${workspaceStore.detail?.ID!}/traefik.yml?name=${encodeURIComponent(workspaceStore.detail?.Name!)}`
})

onMounted(async () => {
  await handleRefresh()
  state.loaded = true
})

const handleFetchList = () => {
  return serverStore.fetchIndexAsync(workspaceStore.detail?.ID!)
}
const handleRefresh = async () => {
  await handleFetchList()
  await workspaceStore.fetchShowAsync(workspaceStore.detail?.ID!)
  await middlewareStore.fetchIndexAsync(workspaceStore.detail?.ID!)
  await serviceStore.fetchIndexAsync(workspaceStore.detail?.ID!)
  await ruleStore.fetchIndexAsync(workspaceStore.detail?.ID!)
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

const handleAddRule = (row: Server) => {
  ruleFormRef.value?.Add(row)
}
const handleEditRule = (rule: Rule, server: Server) => {
  ruleFormRef.value?.Edit(rule, server)
}

// ==================
const handleGenerate = async () => {
  state.loading = true
  try {
    await workspaceStore.generateTraefikConfigAsync(workspaceStore.detail?.ID!)
    ElMessage.success('traefik.yaml has generated')
  } catch (error: any) {
    ElMessage.error(error.message)
  } finally {
    state.loading = false
  }
}
</script>

<style scoped>
</style>