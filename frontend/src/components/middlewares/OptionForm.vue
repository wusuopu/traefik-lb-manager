<template>
  <div>
    <VueForm
      v-model="formData.Options"
      :schema="schema"
      :uiSchema="uiSchema"
      :formFooter="{show: false}"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, defineProps } from 'vue';
import VueForm from '@lljj/vue3-form-element';
import ObjectField from './ObjectField.vue';
import allSchema from './schema';

const props = defineProps<{
  category: string,
  formData: any,
}>()

const schema = computed(() => {
  return allSchema[props.category + 'Middleware']
})
const uiSchema = computed(() => {
  const ui: any = {}
  _.each(schema.value?.properties, (p, k) => {
    if (p.type === 'object' && p.additionalProperties) {
      ui[k] = {
        'ui:field': ObjectField,
      }
    }
  })
  return ui
})
</script>

<style lang="scss" scoped>
</style>
