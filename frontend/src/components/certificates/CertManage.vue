<template>
  <div class="flex mb-3 section-box">
    <el-button type="primary" @click="handleAddCert">Add Cert</el-button>
    <el-button type="success" @click="handleFetchList">Refresh</el-button>
  </div>

  <div class="section-box-dark mb-3">
    <el-table :data="certificateStore.certificates" style="width: 100%">
        <el-table-column prop="Name" label="Name" width="150" />
        <el-table-column prop="Domain" label="Domain" width="250" />
        <el-table-column prop="Status" label="Status" width="150" />
        <el-table-column prop="Enable" label="Enable" width="80" />
        <el-table-column prop="ExpiredAt" label="ExpiredAt" width="250" :formatter="format.tableDatetimeFormat" />
        <el-table-column prop="CreatedAt" label="CreatedAt" width="250" :formatter="format.tableDatetimeFormat" />
        <el-table-column prop="UpdatedAt" label="UpdatedAt" width="250" :formatter="format.tableDatetimeFormat" />
        <el-table-column fixed="right" label="Operations" min-width="120">
          <template #default="scope">
            <el-popconfirm @confirm="handleDelete(scope.row.ID)" title="Are you sure to delete this record?">
              <template #reference>
                <el-button v-loading.fullscreen.lock="state.loading" type="danger" size="small">Delete</el-button>
              </template>
            </el-popconfirm>
            <el-button type="primary" size="small" @click="handleEdit(scope.row)">
              Edit
            </el-button>
            <el-button type="success" size="small" @click="handleRenew(scope.row.ID)">
              Renew
            </el-button>
          </template>
        </el-table-column>
      </el-table>
  </div>

  <el-drawer v-model="state.form.showDrawer" direction="rtl" class="!w-[90%] max-w-[600px]">
    <template #header>
      <h4 v-if="state.form.action === 'update'">Update certificate #{{ state.form.data.ID }}</h4>
      <h4 v-else>Create certificate</h4>
    </template>

    <template #default>
        <div>
          <el-form ref="formRef" :model="state.form.data" :rules="state.form.rules" label-position="top">
            <el-form-item label="Name" prop="Name">
              <el-input v-model="state.form.data.Name" placeholder="Name" />
            </el-form-item>

            <el-form-item label="Domain" prop="Domain" required>
              <el-input v-model="state.form.data.Domain" placeholder="www.example.com" :disabled="state.form.action === 'update'" />
            </el-form-item>

            <el-form-item label="Enable" prop="Enable" required>
              <el-switch v-model="state.form.data.Enable" />
            </el-form-item>
          </el-form>
        </div>
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
import { onMounted, reactive, ref, shallowRef } from 'vue';
import { ElMessage, type FormInstance } from 'element-plus';
import { useWorkspaceStore } from '@/stores/workspace';
import { useCertificateStore } from '@/stores/certificate';
import format from '@/lib/format';

const workspaceStore = useWorkspaceStore()
const certificateStore = useCertificateStore()
const formRef = ref<FormInstance>()

const state = reactive({
  loading: false,
  form: {
    showDrawer: false,
    loading: false,
    data: {
      Name: '',
      Domain: '',
      Enable: true,
    } as Certificate,
    rules: {
      Name: [
        { required: true, message: 'Name is required' },
      ],
    },
    action: '',   // create or update
  },
})

onMounted(async () => {
  await handleFetchList()
})

const handleFetchList = () => {
  return certificateStore.fetchIndexAsync(workspaceStore.detail?.ID!)
}
const handleAddCert = () => {
  state.form.action = 'create'
  state.form.data = {
    Name: '',
    Domain: '',
    Enable: true,
  }

  formRef.value?.resetFields()
  state.form.showDrawer = true
}
const handleEdit = (row: Certificate) => {
  state.form.action = 'update'
  state.form.data = {...row}

  formRef.value?.resetFields()
  state.form.showDrawer = true
}

const handleSubmit = async () => {
  await formRef.value!.validate()
  const payload = {...state.form.data}
  try {
    const u = new URL(payload.Domain)
    payload.Domain = u.hostname
  } catch (error) {
  }

  try {
    state.form.loading = true

    if (state.form.action === 'create') {
      await certificateStore.createAsync(workspaceStore.detail?.ID!, payload)
      ElMessage.success('Certificate has created')
    } else {
      await certificateStore.updateAsync(workspaceStore.detail?.ID!, payload)
      ElMessage.success('Certificate has changed')
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
    await certificateStore.deleteAsync(workspaceStore.detail?.ID!, id)
    ElMessage.success('Certificate has deleted')
  } catch (error: any) {
    ElMessage.error(error.response.data.Error)
    return
  }

  await handleFetchList()
}

const handleRenew = async (id: number) => {
  try {
    await certificateStore.renewAsync(workspaceStore.detail?.ID!, id)
    ElMessage.success('Certificate will try to renew')
  } catch (error: any) {
    ElMessage.error(error.response.data.Error)
    return
  }
}
</script>

<style lang="scss" scoped>
</style>
