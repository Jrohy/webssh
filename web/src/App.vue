<template>
  <div id="app">
    <el-container>
      <el-header style="height: auto">
        <vheader @ssh-select="() => {
          this.$refs.tabs.openTerm()
        }"/>
      </el-header>
      <el-container>
        <el-main style="padding: 0">
            <tabs ref="tabs"></tabs>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script>
import Header from '@/components/Header'
import Tabs from '@/components/Tabs'

export default {
    name: 'App',
    components: {
        vheader: Header,
        tabs: Tabs
    }
}
</script>

<style lang="scss">
#app {
    height: 100%;
    > div {
        height: 100%;
    }
}
</style>
