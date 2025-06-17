<template>
  <TopInfo :workspace="workspaceStore.currentWorkspace">
    <el-button type="primary" @click="handleAdd">Add service</el-button>
    <el-button type="success" @click="handleFetchList">Refresh</el-button>
  </TopInfo>

  <div class="section-box-dark mb-3">
    <el-table :data="serviceStore.services" style="width: 100%">
        <el-table-column prop="Name" label="Name" width="150" />
        <el-table-column prop="LBServers" label="LBServers" min-width="250">
          <template #default="scope">
            <p v-for="(item, index) in scope.row.LBServers" :key="index">
              url: {{ item.Url }} <br />
            </p>
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
          </template>
        </el-table-column>
      </el-table>
  </div>

  <el-drawer v-model="state.form.showDrawer" direction="rtl" class="!w-[90%] max-w-[600px]">
    <template #header>
      <h4 v-if="state.form.action === 'update'">Update service #{{ state.form.data.ID }}</h4>
      <h4 v-else>Create service</h4>
    </template>

    <template #default>
      <ServiceForm ref="formRef" :value="state.form.data" :action="state.form.action" />
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
import { onMounted, reactive, ref } from 'vue';
import { ElMessage, type FormInstance } from 'element-plus';
import { useWorkspaceStore } from '@/stores/workspace';
import { useServiceStore } from '@/stores/services';
import format from '@/lib/format';
import ServiceForm from './ServiceForm.vue';
import TopInfo from '../workspaces/TopInfo.vue';


const workspaceStore = useWorkspaceStore()
const serviceStore = useServiceStore()
const formRef = ref<typeof ServiceForm>()

const state = reactive({
  loading: false,
  form: {
    showDrawer: false,
    loading: false,
    data: {
      Name: '',
      LBServers: [],
    } as Service,
    action: '',   // create or update
  },
})

onMounted(async () => {
  await handleFetchList()
})

const handleFetchList = () => {
  return serviceStore.fetchIndexAsync(workspaceStore.detail?.ID!)
}

const handleAdd = () => {
  state.form.action = 'create'
  state.form.data = {
    Name: '',
    LBServers: [{Url: '', PreservePath: true, Weight: 1, HostName: '', Port: '', PathName: ''}],
  }

  formRef.value?.resetFields()
  state.form.showDrawer = true
}
const handleEdit = (row: Service) => {
  state.form.action = 'update'
  state.form.data = {...row}
  if (_.isEmpty(state.form.data.LBServers) || !_.isArray(state.form.data.LBServers)) {
    state.form.data.LBServers = [{Url: '', PreservePath: true, Weight: 1, HostName: '', Port: '', PathName: ''}]
  }
  if (_.includes(['rancher_v1', 'portainer_swarm'], workspaceStore.detail?.Category)) {
    // 选择内部服务
    _.each(state.form.data.LBServers, (item) => {
      if (!item.Url) { return }

      try {
        const u = new URL(item.Url)
        item.HostName = u.hostname
        item.Port = u.port
        item.PathName = u.pathname
      } catch (error) {
      }
    })
  }

  formRef.value?.resetFields()
  state.form.showDrawer = true
}

const handleSubmit = async () => {
  await formRef.value!.validate()
  const payload = {...state.form.data}

  _.each(state.form.data.LBServers, (item, index) => {
    item.Weight = state.form.data.LBServers?.length! - index
    item.PreservePath = true

    if (_.includes(['rancher_v1', 'portainer_swarm'], workspaceStore.detail?.Category)) {
      item.Url = `http://${item.HostName}:${item.Port || '80'}${item.PathName || ''}`
    }
    delete item.HostName
    delete item.Port
    delete item.PathName
  })


  try {
    state.form.loading = true

    if (state.form.action === 'create') {
      await serviceStore.createAsync(workspaceStore.detail?.ID!, payload)
      ElMessage.success('Service has created')
    } else {
      await serviceStore.updateAsync(workspaceStore.detail?.ID!, payload)
      ElMessage.success('Service has changed')
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

const handleDelete = async (id: number) => {
  try {
    await serviceStore.deleteAsync(workspaceStore.detail?.ID!, id)
    ElMessage.success('Service has deleted')
  } catch (error: any) {
    ElMessage.error(error.response.data.Error)
    return
  }

  await handleFetchList()
}

</script>

<style lang="scss" scoped>
</style>
