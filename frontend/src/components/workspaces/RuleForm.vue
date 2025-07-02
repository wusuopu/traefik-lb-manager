<template>
  <el-drawer v-model="state.form.showDrawer" direction="rtl" class="!w-[90%] max-w-[600px]">
    <template #header>
      <h4 v-if="state.form.action === 'update'">Update rule #{{ state.form.data.ID }}</h4>
      <h4 v-else>Create rule for Server #{{ state.server.ID }}</h4>
    </template>

    <template #default>
      <div class="flex gap-1 items-center border-b-2 mb-2">
        Enable: <el-switch v-model="state.form.data.Enable" />
      </div>
      <VueForm
        v-model="state.form.data.Options"
        ref="formRef"
        :schema="schema"
        :uiSchema="uiSchema"
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
import { computed, reactive, ref } from 'vue';
import VueForm from '@lljj/vue3-form-element';
import { useWorkspaceStore } from '@/stores/workspace';
import { useRuleStore } from '@/stores/rules';
import { useMiddlewareStore } from '@/stores/middlewares';
import { useServiceStore } from '@/stores/services';
import { httpRouterSchema } from './rule-schema';
import { ElMessage } from 'element-plus';

const workspaceStore = useWorkspaceStore()
const middlewareStore = useMiddlewareStore()
const serviceStore = useServiceStore()
const ruleStore = useRuleStore()
const formRef = ref()
const state = reactive({
  server: {} as Server,
  form: {
    showDrawer: false,
    loading: false,
    data: {

    } as Rule,
    action: '',   // create or update
  },
})

const schema = computed(() => {
  const data = _.cloneDeep(httpRouterSchema)

  const serviceEnum: number[] = []
  const serviceEnumName: string[] = []
  _.each(serviceStore.services, (s) => {
    serviceEnum.push(s.ID!)
    serviceEnumName.push(`#${s.ID} ${s.Name}`)
  })
  data.properties.service.enum = serviceEnum
  data.properties.service.enumNames = serviceEnumName

  const middlewareEnum: number[] = []
  const middlewareEnumName: string[] = []
  _.each(middlewareStore.middlewares, (s) => {
    middlewareEnum.push(s.ID!)
    middlewareEnumName.push(`#${s.ID} ${s.Name}`)
  })
  data.properties.middlewares.items.enum = middlewareEnum
  data.properties.middlewares.items.enumNames = middlewareEnumName

  data.properties.entryPoints.items.enum = workspaceStore.detail?.Entrypoints || []

  return data
})
const uiSchema = computed(() => {
  const data = {
    service: {
      'ui:options': {
        attrs: {
          filterable: true,
        },
      },
    },
  }
  return data
})

const handleSubmit = async () => {
  await formRef.value.$$uiFormRef.validate()

  const payload = {...state.form.data}

  try {
    state.form.loading = true

    if (state.form.action === 'create') {
      payload.ServerID = state.server.ID
      await ruleStore.createAsync(workspaceStore.detail?.ID!, payload)
      ElMessage.success('Server has created')
    } else {
      await ruleStore.updateAsync(workspaceStore.detail?.ID!, payload)
      ElMessage.success('Server has changed')
    }
  } catch (error: any) {
    ElMessage.error(error.response.data.Error)
    return
  } finally {
    state.form.loading = false
  }

  state.form.showDrawer = false

  await ruleStore.fetchIndexAsync(workspaceStore.detail?.ID!)
}

const Add = (server: Server) => {
  state.form.action = 'create'
  state.server = server

  state.form.data = {
    Enable: true,
    Options: {
      priority: 1,
    },
  }

  formRef.value?.$$uiFormRef?.resetFields()
  state.form.showDrawer = true
}
const Edit = (item: Rule, server: Server) => {
  state.form.action = 'update'
  state.server = server

  state.form.data = {
    ...item,
    Options: {...item.Options}
  }

  formRef.value?.$$uiFormRef?.resetFields()
  state.form.showDrawer = true
}

defineExpose({
  Add,
  Edit,
})
</script>

<style lang="scss" scoped>
</style>
