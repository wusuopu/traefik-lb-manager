import http from './_http';

export default {
  index (workspaceId: number) {
    return http.get(`/workspaces/${workspaceId}/rules/`);
  },
  create (workspaceId: number, data: any) {
    return http.post(`/workspaces/${workspaceId}/rules/`, data);
  },
  update (workspaceId: number, id: number, data: any) {
    return http.put(`/workspaces/${workspaceId}/rules/${id}`, data);
  },
  delete (workspaceId: number, id: number) {
    return http.delete(`/workspaces/${workspaceId}/rules/${id}`);
  },
}