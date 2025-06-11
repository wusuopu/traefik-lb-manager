<template>
  <el-container class="min-h-[100vh] mx-auto">
    <el-header height="100px">
      <div class="w-full max-w-[1200px] mx-auto text-[#6e7687] font-bold text-2xl">
        <el-menu mode="horizontal" :router="true" :default-active="$route.name">
          <el-menu-item index="Home" route="/">Rules</el-menu-item>
          <el-menu-item index="Authentications" route="/authentications">Authentications</el-menu-item>
          <el-menu-item index="Certificates" route="/certificates">SSL Certificates</el-menu-item>
          <el-menu-item index="Caddyfile" route="/caddyfile">Caddyfile</el-menu-item>
        </el-menu>
      </div>
    </el-header>


    <el-main>
      <div class="w-full mx-auto text-[#6e7687]">
        <RouterView></RouterView>
      </div>
    </el-main>

    <el-footer height="30px">
      <div class="flex justify-between max-w-[1200px] mx-auto text-[#6e7687]">
        <span>{{ version }}</span>
        <span><a href="https://github.com/wusuopu/traefik-lb-manager" target="_blank">Fork me on Github</a></span>
      </div>
    </el-footer>
  </el-container>
</template>

<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'

const version = __APP_VERSION__
</script>


<style scoped>
</style>
