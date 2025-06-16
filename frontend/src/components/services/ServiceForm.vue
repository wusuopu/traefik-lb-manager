<template>
  <el-form ref="formRef" :model="value" :rules="state.rules" label-position="top">
    <el-form-item label="Name" prop="Name">
      <el-input v-model="value.Name" placeholder="Name" />
    </el-form-item>

    <div
      v-for="(item, index) in value.LBServers"
      :key="index"
      class="border-b-[1px]"
    >
      <h3 class="text-[14px]">Server {{ index }}</h3>
      <el-form-item label="Url" :prop="`LBServers.${index}.Url`" required :rules="[
        { validator: validateUrl, trigger: 'blur' },
      ]">
        <el-input v-model="item.Url" placeholder="Url" />
      </el-form-item>

      <div class="flex gap-1 justify-between">
        <el-form-item label="PreservePath" :prop="`LBServers.${index}.PreservePath`">
          <el-switch v-model="item.PreservePath" />
        </el-form-item>
        <el-form-item label="Weight" :prop="`LBServers.${index}.Weight`" required>
          <el-input-number v-model="item.Weight" placeholder="Weight" />
        </el-form-item>
      </div>

      <p v-if="index !== 0" class="mt-1">
        <el-button type="danger" size="small" @click="handleRemove(index)">Remove</el-button>
      </p>
    </div>

    <div>
      <el-button type="primary" size="small" @click="handleAdd">Add server</el-button>
    </div>
  </el-form>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue';
import { ElMessage, type FormInstance } from 'element-plus';
import { useWorkspaceStore } from '@/stores/workspace';

const workspaceStore = useWorkspaceStore()
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

const validateUrl = (rule: any, value: string, callback: any) => {
  if (!value) {
    callback(new Error('Url is required'))
  } else if (!/^https?:\/\/.+$/.test(value)) {
    callback(new Error('Url must start with http:// or https://'))
  } else {
    callback()
  }
}

const handleAdd = () => {
  props.value.LBServers?.push({ Url: '', PreservePath: true, Weight: 1 })
}
const handleRemove = (index: number) => {
  props.value.LBServers?.splice(index, 1)
}

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
