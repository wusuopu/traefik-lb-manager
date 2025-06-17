import http from './_http';

export default {
  index (workspaceId: number) {
    return http.get(`/workspaces/${workspaceId}/middlewares/`);
  },
  create (workspaceId: number, data: any) {
    return http.post(`/workspaces/${workspaceId}/middlewares/`, data);
  },
  update (workspaceId: number, id: number, data: any) {
    return http.put(`/workspaces/${workspaceId}/middlewares/${id}`, data);
  },
  delete (workspaceId: number, id: number) {
    return http.delete(`/workspaces/${workspaceId}/middlewares/${id}`);
  },
}