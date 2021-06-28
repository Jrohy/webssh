<template>
    <div>
        <el-tabs v-model="currentTerm" type="card" closable @tab-remove="removeTab" @tab-click="clickTab">
            <el-tab-pane
                v-for="(item, index) in termList"
                :key="item.name"
                :label="item.label"
                :name="item.name"
            >
                <terminal :id="'Terminal' + index" :ref="item.name"></terminal>
            </el-tab-pane>
        </el-tabs>
        <div v-show="contextMenuVisible">
            <ul :style="{left:left+'px',top:top+'px'}" class="contextmenu">
                <li @click="copyTab()"><el-button type="text" size="mini">复制</el-button></li>
                <li @click="setScreenfull()"><el-button type="text" size="mini">全屏</el-button></li>
                <li @click="removeTab(menuTab)"><el-button type="text" size="mini">关闭</el-button></li>
                <el-divider></el-divider>
                <li @click="closeTabs('left')"><el-button type="text" size="mini">关闭左边</el-button></li>
                <li @click="closeTabs('right')"><el-button type="text" size="mini">关闭右边</el-button></li>
                <li @click="closeTabs('other')"><el-button type="text"  size="mini">关闭其他</el-button></li>
                <li @click="closeTabs('all')"><el-button type="text" size="mini">关闭所有</el-button></li>
            </ul>
        </div>
    </div>
</template>

<script>
import Sortable from 'sortablejs'
import screenfull from 'screenfull'
import Terminal from '@/components/Terminal'

export default {
    name: 'Tabs',
    components: {
        terminal: Terminal
    },
    data () {
        return {
            currentTerm: '',
            currentTermIndex: 0,
            menuTab: '',
            contextMenuVisible: false,
            left: '',
            top: ''
        }
    },
    computed: {
        termList: {
            get() {
                return this.$store.state.termList
            },
            set(v) {
                this.$store.commit('SET_TERMLIST', v)
            }
        }
    },
    watch: {
        contextMenuVisible() {
            if (this.contextMenuVisible) {
                document.body.addEventListener('click', this.closeContextMenu)
            } else {
                document.body.removeEventListener('click', this.closeContextMenu)
            }
        }
    },
    mounted() {
        // 使用原生js 为单个dom绑定鼠标右击事件
        const tabTop = document.body.getElementsByClassName('el-tabs__nav-scroll')
        for (let i = 0; i < tabTop.length; ++i) {
            tabTop[i].oncontextmenu = this.openContextMenu
        }
        // 实现el-tabs可拖动
        const self = this
        const el = document.querySelector('.el-tabs__nav')
        Sortable.create(el, {
            animation: 200,
            onEnd({ newIndex, oldIndex }) {
                const currRow = self.termList.splice(oldIndex, 1)[0]
                self.termList.splice(newIndex, 0, currRow)
            }
        })
    },
    methods: {
        copyTab() {
            this.$refs[`${this.menuTab}`][0].setSSH()
            this.openTerm()
        },
        setScreenfull() {
            if (!screenfull.isEnabled) {
                this.$message({
                    message: '暂不不支持全屏',
                    type: 'warning'
                })
                return false
            }
            screenfull.toggle()
        },
        closeTabs(par) {
            if (par === 'all') {
                this.termList = []
                return
            }
            let currMenuIndex = 0
            for (;currMenuIndex < this.termList.length; ++currMenuIndex) {
                if (this.termList[currMenuIndex].name === this.menuTab) {
                    break
                }
            }
            const setCurrentTerm = () => {
                this.currentTermIndex = currMenuIndex
                const tab = this.termList[currMenuIndex]
                this.currentTerm = tab.name
                document.title = tab.label
                this.$store.commit('SET_TAB', this.termList[this.currentTermIndex])
                this.$refs[`${tab.name}`][0].setSSH()
            }
            switch (par) {
            case 'left':
                // 删除左侧tab标签
                if (this.currentTermIndex < currMenuIndex) {
                    setCurrentTerm()
                }
                this.termList.splice(0, currMenuIndex)
                break
            case 'right':
                // 删除右侧tab标签
                if (this.currentTermIndex > currMenuIndex) {
                    setCurrentTerm()
                }
                this.termList.splice(currMenuIndex + 1, this.termList.length)
                break
            case 'other':
                // 删除其他所有tab标签
                if (this.currentTermIndex !== currMenuIndex) {
                    setCurrentTerm()
                }
                this.termList = this.termList.filter(tab => tab.name === this.menuTab)
                break
            }
            this.closeContextMenu()
        },
        closeContextMenu() {
            this.contextMenuVisible = false
        },
        openContextMenu(e) {
            e.preventDefault() // 防止默认菜单弹出
            const obj = e.srcElement ? e.srcElement : e.target
            if (obj.id) {
                this.menuTab = obj.id.substr(4)
                this.contextMenuVisible = true
                this.left = e.clientX
                this.top = 20
            }
        },
        genID(length) {
            return Number(Math.random().toString().substr(3, length) + Date.now()).toString(36)
        },
        openTerm() {
            const sshInfo = this.$store.state.sshInfo
            if (sshInfo.password === '') {
                return
            }
            this.termList.push({
                name: `${sshInfo.host}-${this.genID(5)}`,
                label: sshInfo.host,
                path: '/'
            })
            const tab = this.termList[this.termList.length - 1]
            this.currentTerm = tab.name
            this.currentTermIndex = this.termList.length - 1
            this.$store.commit('SET_TAB', this.termList[this.currentTermIndex])
        },
        findTerm() {
            for (let i = 0; i < this.termList.length; ++i) {
                if (this.termList[i].name === this.currentTerm) {
                    this.currentTermIndex = i
                    break
                }
            }
            this.$store.commit('SET_TAB', this.termList[this.currentTermIndex])
        },
        clickTab(tab) {
            this.$refs[`${tab.name}`][0].setSSH()
            document.title = tab.label
            this.findTerm()
        },
        removeTab(targetName) {
            const tabs = this.termList
            let activeName = this.currentTerm
            if (activeName === targetName) {
                tabs.forEach((tab, index) => {
                    if (targetName === tab.name) {
                        const nextTab = tabs[index + 1] || tabs[index - 1]
                        if (nextTab) {
                            activeName = nextTab.name
                        }
                    }
                })
                this.currentTerm = activeName
                this.$refs[`${this.currentTerm}`][0].setSSH()
            }
            this.termList = tabs.filter(tab => tab.name !== targetName)
            this.findTerm()
        }
    }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
 .el-divider--horizontal{
     margin: 3px 0;
 }
.contextmenu {
    width: 100px;
    margin: 0;
    border: 1px solid #ccc;
    background: #fff;
    z-index: 3000;
    position: absolute;
    list-style-type: none;
    padding: 5px 0;
    border-radius: 4px;
    font-size: 14px;
    color: #333;
    box-shadow: 2px 2px 3px 0 rgba(0, 0, 0, 0.2);
    li {
        margin: 0;
        padding: 0px 22px;
    }
    li:hover {
        background: #f2f2f2;
        cursor: pointer;
    }
    li button{
        color: #2c3e50;
    }
}
</style>
