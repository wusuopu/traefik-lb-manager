import http from './_http';

export default {
  index (workspaceId: number) {
    return http.get(`/workspaces/${workspaceId}/servers/`);
  },
  create (workspaceId: number, data: any) {
    return http.post(`/workspaces/${workspaceId}/servers/`, data);
  },
  update (workspaceId: number, id: number, data: any) {
    return http.put(`/workspaces/${workspaceId}/servers/${id}`, data);
  },
  delete (workspaceId: number, id: number) {
    return http.delete(`/workspaces/${workspaceId}/servers/${id}`);
  },
}
