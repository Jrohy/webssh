<template>
    <div>
        <div id="terminal"></div>
    </div>
</template>

<script>
import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'
import { AttachAddon } from 'xterm-addon-attach'

export default {
    name: 'WebSSH',
    data () {
        return {
            term: null,
            ws: null,
            resetClose: false
        }
    },
    methods: {
        openTerm() {
            const sshReq = this.$store.getters.sshReq
            if (sshReq === undefined || sshReq === '') {
                return
            }
            this.close()
            const prefix = process.env.NODE_ENV === 'production' ? '' : '/ws'
            const fitAddon = new FitAddon()
            this.term = new Terminal({
                rows: Math.floor(document.documentElement.clientHeight / 18) - 4
            })
            this.term.loadAddon(fitAddon)
            this.term.open(document.getElementById('terminal'))
            fitAddon.fit()
            const self = this
            // open websocket
            this.ws = new WebSocket(`${(location.protocol === 'http:' ? 'ws' : 'wss')}://${location.host}${prefix}/term?sshInfo=${sshReq}&rows=${this.term.rows}&cols=${this.term.cols}`)
            this.ws.onopen = e => {
                console.log(Date(), 'onopen')
                self.connected()
            }
            this.ws.onclose = e => {
                console.log(Date(), 'onclose')
                if (!self.resetClose) {
                    this.$store.commit('SET_PASS', '')
                    this.$message({
                        message: 'websocket连接已断开!',
                        type: 'warning',
                        duration: 0,
                        showClose: true
                    })
                    this.ws = null
                }
                self.resetClose = false
            }
            this.ws.onerror = e => {
                console.log(Date(), 'onerror')
            }
            const attachAddon = new AttachAddon(this.ws)
            this.term.loadAddon(attachAddon)
            this.term.attachCustomKeyEventHandler((e) => {
                // ctrl + v
                if (e.ctrlKey && e.key === 'v') {
                    document.execCommand('copy')
                    return false
                }
                // ctrl + c
                if (e.ctrlKey && e.key === 'c' && self.term.hasSelection()) {
                    document.execCommand('copy')
                    return false
                }
            })
        },
        connected() {
            const sshInfo = this.$store.state.sshInfo
            document.title = sshInfo.host
            let sshList = this.$store.state.sshList
            if (sshList === null) {
                sshList = `[{"host": "${sshInfo.host}", "username": "${sshInfo.username}", "port":${sshInfo.port}}]`
            } else {
                const sshListObj = JSON.parse(window.atob(sshList))
                sshListObj.forEach((v, i) => {
                    if (v.host === sshInfo.host) {
                        sshListObj.splice(i, 1)
                    }
                })
                sshListObj.push({
                    host: sshInfo.host,
                    username: sshInfo.username,
                    port: sshInfo.port
                })
                sshList = JSON.stringify(sshListObj)
            }
            this.$store.commit('SET_LIST', window.btoa(sshList))
        },
        close() {
            if (this.ws !== null) {
                this.ws.close()
                this.resetClose = true
            }
            if (this.term !== null) {
                this.term.dispose()
            }
        }
    },
    beforeDestroy () {
        this.close()
    }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
