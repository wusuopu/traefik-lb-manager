<template>
  <div class="section-box">
    <div>
      <el-text type="primary">Rule List for {{ server.Name }}</el-text>
    </div>
    <el-table :data="ruleList" border style="width: 100%">
      <el-table-column label="Operations" width="150">
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

      <el-table-column prop="ID" label="ID" width="80" />
      <el-table-column prop="Enable" label="Enable" width="70">
        <template #default="scope">
          <el-switch :model-value="scope.row.Enable" size="small" disabled/>
        </template>
      </el-table-column>

      <el-table-column prop="Options.rule" label="Rule" min-width="200">
        <template #default="scope">
          <el-text v-if="scope.row.Options?.advanceMode" type="warning">{{ scope.row.Options.rule }}</el-text>
          <el-text v-else>{{ scope.row.Options.rule }}</el-text>
        </template>
      </el-table-column>
      <el-table-column prop="Options.priority" label="Priority" width="80" />
      <el-table-column prop="Options" label="Options" min-width="120">
        <template #default="scope">
          <p v-for="(line, index) in optionsFormat(scope.row)" :key="index">
            <span class="font-bold mr-2">{{ line[0] }}:</span>
            <span>{{ line[1] }}</span>
          </p>
        </template>
      </el-table-column>

      <el-table-column prop="CreatedAt" label="CreatedAt" width="120" :formatter="format.tableDatetimeFormat" />
      <el-table-column prop="UpdatedAt" label="UpdatedAt" width="120" :formatter="format.tableDatetimeFormat" />
    </el-table>
  </div>
</template>

<script setup lang="ts">
import _ from 'lodash';
import { reactive, onMounted, ref, computed } from 'vue';
import { useWorkspaceStore } from '@/stores/workspace';
import { useRuleStore } from '@/stores/rules';
import { useMiddlewareStore } from '@/stores/middlewares';
import { useServiceStore } from '@/stores/services';
import format from '@/lib/format';
import { ElMessage } from 'element-plus';

const workspaceStore = useWorkspaceStore()
const ruleStore = useRuleStore()
const middlewareStore = useMiddlewareStore()
const serviceStore = useServiceStore()

const props = defineProps<{
  server: Server
}>()
const emit = defineEmits<{
  edit: [Rule, Server],
}>()


const state = reactive({
  loading: false,
})

const ruleList = computed(() => {
  return _.orderBy(
    _.filter(ruleStore.rules, (item: Rule) => item.ServerID === props.server.ID),
    ['Options.priority', 'ID'], ["desc", "asc"]
  )
})

const optionsFormat = (row: Rule) => {
  const middlewares = _.reduce(row.Options?.middlewares, (ret, id: number) => {
    let obj = _.find(middlewareStore.middlewares, (item: Middleware) => item.ID === id)
    if (obj) {
      ret.push(obj)
    }
    return ret
  }, [] as Middleware[])
  const service = _.find(serviceStore.services, (item: Service) => item.ID === row.Options?.service)

  const rules = [
    '(' + _.map(props.server.Host, (host: string) => {
      return `Host(\`${host}\`)`
    }).join(' || ') + ')',
  ]
  if (row.Options?.rule) {
    if (row.Options?.advanceMode) {
      rules.push(`(${row.Options?.rule})`)
    } else {
      rules.push(`PathPrefix(\`${row.Options?.rule}\`)`)
    }
  }

  const content = [
    // ['rule', rules.join(' && ')],
    ['service', `#${service?.ID} ${service?.Name || ''}` || ''],
  ]
  if (middlewares.length > 0) {
    content.push(['middlewares', _.map(middlewares, (item: Middleware) => `#${item.ID} ${item.Name}`).join(', ')])
  }
  if (!_.isEmpty(row.Options?.entryPoints)) {
    content.push(['entryPoints', _.map(row.Options?.entryPoints, (item: string) => item).join(', ')])
  }
  return content
}
const handleEdit = (route: any) => {
  emit('edit', route, props.server)
}
const handleDelete = async (id: any) => {
  try {
    await ruleStore.deleteAsync(workspaceStore.detail?.ID!, id)
    ElMessage.success('Rule has deleted')
  } catch (error: any) {
    ElMessage.error(error.response.data.Error)
    return
  }

  await ruleStore.fetchIndexAsync(workspaceStore.detail?.ID!)
}
</script>

<style scoped>
</style>