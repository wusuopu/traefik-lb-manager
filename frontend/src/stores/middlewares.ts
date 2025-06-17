import { defineStore } from 'pinia'
import API from '@/api/middleware';

export const useMiddlewareStore = defineStore('middlewares', {
  state: () => ({
    middlewares: [] as Middleware[],
  }),
  actions: {
    async fetchIndexAsync (workspaceId: number) {
      const resp = await API.index(workspaceId);
      this.middlewares = resp.data.Data;
    },
    async createAsync (workspaceId: number, middleware: Middleware) {
      await API.create(workspaceId, middleware);
    },
    async updateAsync (workspaceId: number, middleware: Middleware) {
      await API.update(workspaceId, middleware.ID!, middleware);
    },
    async deleteAsync (workspaceId: number, id: number) {
      await API.delete(workspaceId, id);
    },
  },
})