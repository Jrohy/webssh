<template>
    <div>
        <div :id="id"></div>
    </div>
</template>

<script>
import { checkSSH } from '@/api/common'
import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'
import { AttachAddon } from 'xterm-addon-attach'

export default {
    name: 'Terminal',
    props: ['id'],
    data () {
        return {
            term: null,
            ws: null,
            resetClose: false,
            ssh: null,
            savePass: false,
            fontSize: 15
        }
    },
    mounted() {
        this.createTerm()
    },
    methods: {
        setSSH() {
            this.$store.commit('SET_SSH', this.ssh)
        },
        createTerm() {
            if (this.$store.state.sshInfo.password === '') {
                return
            }
            const termWeb = document.getElementById(this.id)
            termWeb.style.height = (document.body.clientHeight - 102) + 'px'
            const sshReq = this.$store.getters.sshReq
            this.close()
            const prefix = process.env.NODE_ENV === 'production' ? '' : '/ws'
            const fitAddon = new FitAddon()
            this.term = new Terminal()
            this.term.loadAddon(fitAddon)
            this.term.open(document.getElementById(this.id))
            fitAddon.fit()
            const self = this
            const heartCheck = {
                timeout: 5000, // 5s发一次心跳
                intervalObj: null,
                stop: function() {
                    clearInterval(this.intervalObj)
                },
                start: function() {
                    this.intervalObj = setInterval(function() {
                        if (self.ws !== null && self.ws.readyState === 1) {
                            self.ws.send('ping')
                        }
                    }, this.timeout)
                }
            }
            // open websocket
            this.ws = new WebSocket(`${(location.protocol === 'http:' ? 'ws' : 'wss')}://${location.host}${prefix}/term?sshInfo=${sshReq}&rows=${this.term.rows}&cols=${this.term.cols}`)
            this.ws.onopen = () => {
                console.log(Date(), 'onopen')
                self.connected()
                heartCheck.start()
            }
            this.ws.onclose = () => {
                console.log(Date(), 'onclose')
                if (!self.resetClose) {
                    if (!this.savePass) {
                        this.$store.commit('SET_PASS', '')
                        this.ssh.password = ''
                    }
                    this.$message({
                        message: 'websocket连接已断开!',
                        type: 'warning',
                        duration: 0,
                        showClose: true
                    })
                    this.ws = null
                }
                heartCheck.stop()
                self.resetClose = false
            }
            this.ws.onerror = () => {
                console.log(Date(), 'onerror')
            }
            const attachAddon = new AttachAddon(this.ws)
            this.term.loadAddon(attachAddon)
            this.term.attachCustomKeyEventHandler((e) => {
                const keyArray = ['F5', 'F11', 'F12']
                if (keyArray.indexOf(e.key) > -1) {
                    return false
                }
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
            // detect available wheel event
            // 各个厂商的高版本浏览器都支持"wheel"
            // Webkit 和 IE一定支持"mousewheel"
            // "DOMMouseScroll" 用于低版本的firefox
            const wheelSupport = 'onwheel' in document.createElement('div') ? 'wheel' : document.onmousewheel !== undefined ? 'mousewheel' : 'DOMMouseScroll'
            termWeb.addEventListener(wheelSupport, (e) => {
                if (e.ctrlKey) {
                    e.preventDefault()
                    if (e.deltaY < 0) {
                        self.term.setOption('fontSize', ++this.fontSize)
                    } else {
                        self.term.setOption('fontSize', --this.fontSize)
                    }
                    fitAddon.fit()
                    if (self.ws !== null && self.ws.readyState === 1) {
                        self.ws.send(`resize:${self.term.rows}:${self.term.cols}`)
                    }
                }
            })
            window.addEventListener('resize', () => {
                termWeb.style.height = (document.body.clientHeight - 102) + 'px'
                fitAddon.fit()
                if (self.ws !== null && self.ws.readyState === 1) {
                    self.ws.send(`resize:${self.term.rows}:${self.term.cols}`)
                }
            })
        },
        async connected() {
            const sshInfo = this.$store.state.sshInfo
            // 深度拷贝对象
            this.ssh = Object.assign({}, sshInfo)
            // 校验ssh连接信息是否正确
            const result = await checkSSH(this.$store.getters.sshReq)
            if (result.Msg !== 'success') {
                return
            } else {
                this.savePass = result.Data.savePass
            }
            document.title = sshInfo.host
            let sshList = this.$store.state.sshList
            if (sshList === null) {
                if (this.savePass) {
                    sshList = `[{"host": "${sshInfo.host}", "username": "${sshInfo.username}", "port":${sshInfo.port}, "logintype":${sshInfo.logintype}, "password":"${sshInfo.password}"}]`
                } else {
                    sshList = `[{"host": "${sshInfo.host}", "username": "${sshInfo.username}", "port":${sshInfo.port},  "logintype":${sshInfo.logintype}}]`
                }
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
                    port: sshInfo.port,
                    logintype: sshInfo.logintype
                })
                if (this.savePass) {
                    sshListObj[sshListObj.length - 1].password = sshInfo.password
                }
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
<style scoped lang="scss">
</style>
