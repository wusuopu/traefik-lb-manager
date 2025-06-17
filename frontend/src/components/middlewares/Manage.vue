<template>
  <TopInfo :workspace="workspaceStore.currentWorkspace">
    <el-dropdown :max-height="400" @command="handleAdd">
      <el-button type="primary">
        Add middleware<el-icon class="el-icon--right"><ArrowDown /></el-icon>
      </el-button>
      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item v-for="(item, index) in allMiddlewares" :key="item" :command="item">{{ item }}</el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
    <el-button type="success" @click="handleFetchList">Refresh</el-button>
  </TopInfo>

  <div class="section-box-dark mb-3">
    <el-table :data="middlewareStore.middlewares" style="width: 100%">
        <el-table-column prop="Name" label="Name" width="150" />
        <el-table-column prop="Category" label="Category" width="250" />
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
            <el-button size="small" @click="handleViewOption(scope.row)">
              View Option
            </el-button>
          </template>
        </el-table-column>
      </el-table>
  </div>

  <el-drawer v-model="state.form.showDrawer" direction="rtl" class="!w-[90%] max-w-[600px]">
    <template #header>
      <h4 v-if="state.form.action === 'update'">Update middleware #{{ state.form.data.ID }}</h4>
      <h4 v-else>Create middleware</h4>
    </template>

    <template #default>
      <div>
        <p class="mb-2">
          <a href="https://doc.traefik.io/traefik/middlewares/overview/" target="_blank">https://doc.traefik.io/traefik/middlewares/overview/</a>
        </p>

        <el-form ref="formRef" :model="state.form.data" :rules="state.form.rules" label-position="top">
          <el-form-item label="Category" prop="Category" required>
            <el-input v-model="state.form.data.Category" disabled />
          </el-form-item>
          <el-form-item label="Name" prop="Name" required>
            <el-input v-model="state.form.data.Name" placeholder="Name" />
          </el-form-item>
        </el-form>

        <div v-if="_.includes(['basicAuth', 'digestAuth'], state.form.data.Category)" class="flex justify-end">
          <el-button type="primary" size="small" @click="handleHashPassword">Hash User's Password</el-button>
        </div>

        <OptionForm
          v-if="state.form.showDrawer"
          :key="state.form.data.Category"
          :category="state.form.data.Category"
          :formData="state.form.data"
        />
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
import _ from 'lodash';
import { onMounted, reactive, ref } from 'vue';
import { ElMessage, ElMessageBox, type FormInstance } from 'element-plus';
import { Md5 } from 'ts-md5'
import bcrypt from "bcryptjs";
import { ArrowDown } from '@element-plus/icons-vue'
import { useWorkspaceStore } from '@/stores/workspace';
import { useMiddlewareStore } from '@/stores/middlewares';
import format from '@/lib/format';
import TopInfo from '../workspaces/TopInfo.vue';
import OptionForm from './OptionForm.vue';
import { allMiddlewares } from './schema';

const workspaceStore = useWorkspaceStore()
const middlewareStore = useMiddlewareStore()
const formRef = ref<FormInstance>()

const state = reactive({
  loading: false,
  form: {
    showDrawer: false,
    loading: false,
    data: {
      Name: '',
      Category: '',
    } as Middleware,
    rules: {},
    action: '',   // create or update
  },
})

onMounted(async () => {
  await handleFetchList()
})

const handleFetchList = () => {
  return middlewareStore.fetchIndexAsync(workspaceStore.detail?.ID!)
}
const handleAdd = (category: string) => {
  state.form.action = 'create'
  state.form.data = {
    Name: '',
    Category: category,
    Options: {},
  }
  formRef.value?.resetFields()
  state.form.showDrawer = true
}

const handleEdit = (row: Middleware) => {
  state.form.action = 'update'
  state.form.data = {
    ...row,
    Options: {...row.Options}
  }

  formRef.value?.resetFields()
  state.form.showDrawer = true
}

const handleViewOption = (row: Middleware) => {
  ElMessageBox({
    title: `Options for #${row.ID}`,
    message: `<pre style="max-height: 80vh; overflow: auto; width: 100%;">${JSON.stringify(row.Options, null, 2)}</pre>`,
    dangerouslyUseHTMLString: true,
  })
}

const handleHashPassword = () => {
  console.log(state.form.data.Options?.users)
  _.each(state.form.data.Options?.users, (user: string, index: number) => {
    if (!user) { return }

    if (state.form.data.Category === 'basicAuth') {
      let [username, password] = user.split(':')
      if (!username || !password) { return }
      if (_.startsWith(password, '$apr1$') || password.match(/^\$2\w\$/)) {
        // 已经加密过了
        return
      }
      const hashedPassword = bcrypt.hashSync(password)
      state.form.data.Options!.users[index] = `${username}:${hashedPassword}`
    }
    if (state.form.data.Category === 'digestAuth') {
      let [username, realm, password] = user.split(':')
      if (!username || !realm || !password) { return }
      const hashedPassword = Md5.hashStr(user)
      state.form.data.Options!.users[index] = `${username}:${realm}:${hashedPassword}`
    }
  })
}

const handleSubmit = async () => {
  await formRef.value!.validate()
  const payload = {...state.form.data}

  try {
    state.form.loading = true

    if (state.form.action === 'create') {
      await middlewareStore.createAsync(workspaceStore.detail?.ID!, payload)
      ElMessage.success('Middleware has created')
    } else {
      await middlewareStore.updateAsync(workspaceStore.detail?.ID!, payload)
      ElMessage.success('Middleware has changed')
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
    await middlewareStore.deleteAsync(workspaceStore.detail?.ID!, id)
    ElMessage.success('Middleware has deleted')
  } catch (error: any) {
    ElMessage.error(error.response.data.Error)
    return
  }

  await handleFetchList()
}
</script>

<style lang="scss" scoped>
</style>
