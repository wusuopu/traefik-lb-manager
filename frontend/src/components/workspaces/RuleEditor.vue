<template>
  <div>
    Reference: <a href="https://doc.traefik.io/traefik/reference/dynamic-configuration/file/" target="_blank">https://doc.traefik.io/traefik/reference/dynamic-configuration/file/</a>
  </div>
  <div class="flex gap-1 mb-2">
    <el-button type="primary" @click="insertHTTPChallengeRule">Insert Let's Encrypt Rule</el-button>
    <el-popconfirm @confirm="handleUpdate" title="Are you sure to update this rule?">
      <template #reference>
        <el-button v-loading.fullscreen.lock="state.loading" type="danger">Update</el-button>
      </template>
    </el-popconfirm>
  </div>

  <codemirror
    v-model="code"
    placeholder="Traefik config here..."
    :style="{ height: '400px' }"
    :autofocus="true"
    :indent-with-tab="true"
    :tab-size="2"
    :extensions="extensions"
    @ready="handleReady"
  />
</template>

<script setup lang="ts">
import { onMounted, reactive, ref, shallowRef } from 'vue';
import { useWorkspaceStore } from '@/stores/workspace';
import { Codemirror } from 'vue-codemirror'
import { yaml } from '@codemirror/lang-yaml'
import { oneDark } from '@codemirror/theme-one-dark'
import { ElMessage } from 'element-plus';

const workspaceStore = useWorkspaceStore()
const code = ref('')
const editor = shallowRef()

const state = reactive({
  loading: false,
})

const extensions = [
  yaml(),
  oneDark,
]

const handleReady = (payload: any) => {
  editor.value = payload.view
}

onMounted(() => {
  code.value = (workspaceStore.detail?.TraefikConfig || `
http:
  routers:
    my-router:
      rule: "Host(\`example.com\`)"
      service: "my-service"

  services:
    my-service:
      loadBalancer:
        servers:
          - url: "http://localhost:8080"

  middlewares:
    my-basic-auth:
      basicAuth:
        users:
        - test:$apr1$H6uskkkW$IgXLP6ewTrSuBkTrqE8wj/
  `).trim()
})

const insertHTTPChallengeRule = () => {
  let content = code.value
  code.value = content.replace(/\n  routers:\s*\n/m, `
  routers:
    lets-encrypt-router:
      rule: "PathPrefix(\`/.well-known/\`)"
      service: "lets-encrypt-service"\n`).replace(/\n  services:\s*\n/m, `
  services:
    lets-encrypt-service:
      http:
        loadBalancer:
          servers:
            - url: "${workspaceStore.detail?.ManagerBaseUrl}"\n`)
}

const handleUpdate = async () => {
  state.loading = true
  try {
    await workspaceStore.updateTraefikConfigAsync(workspaceStore.detail?.ID!, code.value)
    ElMessage.success('traefik.yaml has updated')
  } catch (error: any) {
    ElMessage.error(error.message)
  } finally {
    state.loading = false
  }
}
</script>

<style lang="scss" scoped>
</style>
