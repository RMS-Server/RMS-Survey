<template>
  <a-config-provider :locale="zhCN" :theme="themeConfig">
    <router-view />
  </a-config-provider>
</template>

<script setup lang="ts">
import zhCN from 'ant-design-vue/es/locale/zh_CN'

const themeConfig = {
  token: {
    colorPrimary: '#FC7961',
    colorPrimaryHover: '#D65645',
    colorPrimaryActive: '#C44A3A',
    borderRadius: 12,
  },
}
</script>

<style>
#app {
  min-height: 100vh;
}
</style>
