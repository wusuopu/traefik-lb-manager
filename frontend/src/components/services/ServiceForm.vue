<template>
  <el-form ref="formRef" :model="value" :rules="state.rules" label-position="top">
    <el-form-item label="Name" prop="Name">
      <el-input v-model="value.Name" placeholder="Name" />
    </el-form-item>

    <template v-if="!showExternalService">
      <div v-for="(item, index) in value.LBServers" :key="index" class="border-b-[1px]">
        <h3 class="text-[14px]">LoadBalancer Server {{ index }}</h3>
        <el-form-item label="Url" :prop="`LBServers.${index}.url`" required :rules="[
          { validator: validateUrl, trigger: 'blur' },
        ]">
          <el-input v-model="item.url" placeholder="http://host:port/path" />
        </el-form-item>

        <p v-if="index !== 0" class="mt-1">
          <el-button type="danger" size="small" @click="handleRemove(index)">Remove</el-button>
        </p>
      </div>
    </template>
    <template v-else>
      <div v-for="(item, index) in value.LBServers" :key="index" class="border-b-[1px]">
        <h3 class="text-[14px]">LoadBalancer Server {{ index }}</h3>
        <div class="flex gap-1">
          <el-form-item label="Target" :prop="`LBServers.${index}.HostName`" required class="flex-1">
            <el-select v-model="item.HostName" filterable clearable class="w-full">
              <el-option-group v-for="(group, index) in externalServices" :key="index" :label="`Stack: ${group.stackName || '<none>'}`">
                <el-option v-for="(option, idx) in group.services" :key="idx" :label="option.Label" :value="option.HostName" />
              </el-option-group>
            </el-select>
          </el-form-item>

          <el-form-item label="Port" :prop="`LBServers.${index}.Port`" class="flex-1">
            <el-input v-model="item.Port" placeholder="80" />
          </el-form-item>

          <el-form-item label="Path" :prop="`LBServers.${index}.PathName`" class="flex-1">
            <el-input v-model="item.PathName" placeholder="/" />
          </el-form-item>
        </div>

        <p v-if="index !== 0" class="mt-1">
          <el-button type="danger" size="small" @click="handleRemove(index)">Remove</el-button>
        </p>
      </div>
    </template>

    <div class="flex">
      <el-button type="primary" size="small" @click="handleAdd">Add server</el-button>
      <el-button v-if="showExternalService" type="success" size="small" @click="fetchExternalSerevices">Refresh Target</el-button>
    </div>
  </el-form>
</template>

<script setup lang="ts">
import _ from 'lodash';
import { reactive, onMounted, computed, ref } from 'vue';
import { ElMessage, type FormInstance } from 'element-plus';
import { useWorkspaceStore } from '@/stores/workspace';
import { useServiceStore } from '@/stores/services';

const workspaceStore = useWorkspaceStore()
const serviceStore = useServiceStore()
const formRef = ref<FormInstance>()

const props = defineProps<{
  value: Service,
  action: string,
}>()

const state = reactive({
  loading: false,
  rules: {
    Name: [
      { required: true, message: 'Name is required' },
    ],
  },
})

const showExternalService = computed(() => {
  return _.includes(['rancher_v1', 'portainer_swarm'], workspaceStore.detail?.Category)
})
const externalServices = computed(() => {
  const group = _.groupBy(serviceStore.externals, "Stack")
  return _.map(_.orderBy(_.keys(group)), (stackName: string) => {
    return {
      stackName,
      services: _.orderBy(group[stackName], "Name"),
    }
  })
})

const validateUrl = (rule: any, value: string, callback: any) => {
  if (!value) {
    callback(new Error('url is required'))
  } else if (!/^https?:\/\/.+$/.test(value)) {
    callback(new Error('url must start with http:// or https://'))
  } else {
    callback()
  }
}

onMounted(async () => {
  if (_.includes(['rancher_v1', 'portainer_swarm'], workspaceStore.detail?.Category)) {
    await fetchExternalSerevices()
  }
})

const handleAdd = () => {
  props.value.LBServers?.push({ url: '', preservePath: true, weight: 1, HostName: '', Port: '', PathName: '' })
}
const handleRemove = (index: number) => {
  props.value.LBServers?.splice(index, 1)
}

const fetchExternalSerevices = async () => {
  await serviceStore.fetchExternalIndexAsync(workspaceStore.detail?.ID!)
}

// =============================================
const validate = async () => {
  await formRef.value!.validate()
}
const resetFields = () => {
  formRef.value?.resetFields()
}

defineExpose({
  validate,
  resetFields,
})

</script>

<style lang="scss" scoped>
</style>
