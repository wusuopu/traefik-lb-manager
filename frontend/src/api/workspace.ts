import http from './_http';

export default {
  index () {
    return http.get('/workspaces/');
  },
  show (id: number) {
    return http.get(`/workspaces/${id}`);
  },
  create (data: any) {
    return http.post('/workspaces/', data);
  },
  update (id: number, data: any) {
    return http.put(`/workspaces/${id}`, data);
  },
  delete (id: number) {
    return http.delete(`/workspaces/${id}`);
  },
  generateTraefikConfigAsync (id: number) {
    return http.post(`/workspaces/${id}/traefik.yml`);
  },
  updateTraefikConfigAsync (id: number, config: string) {
    return http.put(`/workspaces/${id}/traefik.yml`, {traefik_config: config});
  },
}