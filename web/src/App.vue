<template>
  <div id="app">
    <el-container>
      <el-header>
        <vheader @ssh-select="openTerm"/>
      </el-header>
      <el-container>
        <el-main style="padding: 0">
            <webssh ref="webssh"></webssh>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script>
import Header from '@/components/Header'
import WebSSH from '@/components/WebSSH'

export default {
    name: 'App',
    components: {
        vheader: Header,
        webssh: WebSSH
    },
    methods: {
        openTerm() {
            this.$refs.webssh.openTerm()
        }
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
