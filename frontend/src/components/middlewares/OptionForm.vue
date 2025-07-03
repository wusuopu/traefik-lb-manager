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
import _ from 'lodash';
import { computed, defineProps } from 'vue';
import VueForm from '@lljj/vue3-form-element';
import { useMiddlewareStore } from '@/stores/middlewares';
import ObjectField from './ObjectField.vue';
import allSchema from './schema';

const middlewareStore = useMiddlewareStore()

const props = defineProps<{
  category: string,
  formData: any,
}>()

const schema = computed(() => {
  const data = _.cloneDeep(allSchema[props.category + 'Middleware'])

  if (props.category === 'chain') {
    const middlewareEnum: number[] = []
    const middlewareEnumName: string[] = []
    _.each(middlewareStore.middlewares, (s) => {
      if (s.ID === props.formData.ID) { return }

      middlewareEnum.push(s.ID!)
      middlewareEnumName.push(`#${s.ID} ${s.Name}`)
    })
    data.properties.middlewares.items.enum = middlewareEnum
    data.properties.middlewares.items.enumNames = middlewareEnumName
  }

  return data
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
