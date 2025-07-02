import { fileURLToPath, URL } from 'node:url'

// import { defineConfig, loadEnv } from 'vite'
import { type ConfigEnv, type UserConfigExport, loadEnv } from "vite"
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'
import tailwindcss from '@tailwindcss/vite'

// https://vite.dev/config/
export default ({ mode }: ConfigEnv): UserConfigExport => {
  const viteEnv = loadEnv(mode, process.cwd())

  return {
    base: "", // relative path
    define: {
      __APP_VERSION__: JSON.stringify("v0.1.0"),
    },
    plugins: [vue(), vueDevTools(), tailwindcss()],
    resolve: {
      alias: {
        "@": fileURLToPath(new URL("./src", import.meta.url)),
      },
    },
    server: {
      host: true,
      open: false,
      cors: true,
      strictPort: true,
      proxy: {
        "/workspaces": {
          target: viteEnv.VITE_BASE_API || "http://127.0.0.1:8080",
          ws: true,
          /** 是否允许跨域 */
          changeOrigin: true,
        },
        "/api": {
          target: viteEnv.VITE_BASE_API || "http://127.0.0.1:8080",
          ws: true,
          /** 是否允许跨域 */
          changeOrigin: true,
        },
      },
    },
    build: {
      assetsDir: 'statics/assets'
    }
  };
}
