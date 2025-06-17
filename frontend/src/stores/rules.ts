import { defineStore } from 'pinia'
import API from '@/api/rule';

export const useRuleStore = defineStore('rules', {
  state: () => ({
    rules: [] as Rule[],
  }),
  actions: {
    async fetchIndexAsync (workspaceId: number) {
      const resp = await API.index(workspaceId);
      this.rules = resp.data.Data;
    },
    async createAsync (workspaceId: number, rule: Rule) {
      await API.create(workspaceId, rule);
    },
    async updateAsync (workspaceId: number, rule: Rule) {
      await API.update(workspaceId, rule.ID!, rule);
    },
    async deleteAsync (workspaceId: number, id: number) {
      await API.delete(workspaceId, id);
    },
  },
})