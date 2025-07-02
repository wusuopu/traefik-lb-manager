<template>
  <div>
    <div class="text-[14px]">Custom Plugin Middleware Config:</div>
    <codemirror
      v-model="pluginOption"
      placeholder="plugin config..."
      :style="{ minHeight: '400px', height: '100%', }"
      :autofocus="true"
      :indent-with-tab="true"
      :tab-size="2"
      :extensions="editorExtensions"
      @update:modelValue="handleOptionChange"
    />
    <div v-if="state.jsonError">
      <el-text class="mx-1" type="danger">{{ state.jsonError }}</el-text>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, defineProps, onMounted, reactive, ref, shallowRef } from 'vue';
import { Codemirror } from 'vue-codemirror'
import { oneDark } from '@codemirror/theme-one-dark'
import { json } from '@codemirror/lang-json'

const props = defineProps<{
  formData: any,
}>()

const pluginOption = ref('')
const editorExtensions = [
  json(),
  oneDark,
]
const state = reactive({
  jsonError: '',
})

onMounted(() => {
  if (!_.isPlainObject(props.formData.Options) || _.isEmpty(props.formData.Options)) {
    pluginOption.value = '{\n}'
  } else {
    pluginOption.value = JSON.stringify(props.formData.Options, undefined, 2)
  }
})

const handleOptionChange = (value: string) => {
  try {
    const data = JSON.parse(value)
    props.formData.Options = data
    state.jsonError = ''
  } catch (error: any) {
    state.jsonError = error.message
  }
}
</script>

<style lang="scss" scoped>
</style>
