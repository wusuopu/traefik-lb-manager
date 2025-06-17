<template>
  <div :class="className">
    <h3 class="fieldGroupWrap_title">{{ schema.title }}</h3>
    <p class="fieldGroupWrap_des">{{ schema.description }}</p>

    <div class="section-box">
      <el-form-item
        v-for="(value, key) in rootFormData[curNodePath]"
        :index="key"
        :label="undefined"
        :prop="`${curNodePath}.${key}`"
      >
        <div class="px-4 w-full">
          <p><el-button type="danger" :icon="Delete" size="small" class="mr-2" @click="handleRemoveItem(key)" />key: {{ key }}</p>
          <el-input v-model="rootFormData[curNodePath][key]" />
        </div>
      </el-form-item>

      <div class="flex gap-2 items-center">
        <p>Key:</p>
        <el-input v-model="state.newKey" class="max-w-[200px]" />
        <el-button type="primary" :icon="Plus" :disabled="!state.newKey.trim()" @click="handleAddItem" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import _ from 'lodash';
import { computed, defineProps, reactive } from 'vue';
import { Delete, Plus } from '@element-plus/icons-vue'

defineOptions({
  name: 'ObjectInputField',
})

const props = defineProps<{
  globalOptions: any,
  formProps: any,
  schema: {
    title: string,
    description: string,
  },
  uiSchema: any,
  errorSchema: any,
  customFormats: any,
  rootSchema: any,
  rootFormData: any,
  curNodePath: string,
  required: boolean,
  needValidFieldGroup: boolean,
  fieldProps: any,
  class: any,
}>()

const className = computed(() => {
  return `${props.class || ''}`
})

const state = reactive({
  newKey: '',
})

const handleAddItem = () => {
  const key = state.newKey.trim()
  if (!key) { return }
  if (_.has(props.rootFormData[props.curNodePath], key)) { return }
  props.rootFormData[props.curNodePath][key] = ''
  state.newKey = ''
}
const handleRemoveItem = (key: any) => {
  delete props.rootFormData[props.curNodePath][key]
}
</script>

<style lang="scss" scoped>
.objectInputField {
  margin-bottom: 22px;
}
</style>
