<template>
    <div>
        <el-tabs v-model="currentTerm" type="card" @tab-remove="removeTab" @tab-click="clickTab">
            <el-tab-pane
                v-for="(item, index) in termList"
                :key="item.name"
                :label="item.label"
                :name="item.name"
                :closable="item.closable"
            >
                <terminal :id="'Terminal' + index" :ref="item.name"></terminal>
            </el-tab-pane>
        </el-tabs>
        <div v-show="contextMenuVisible">
            <ul :style="{left:left+'px',top:top+'px'}" class="contextmenu">
                <li @click="copyTab()"><el-button type="text" size="mini">{{$t('Copy')}}</el-button></li>
                <li @click="lockTab()"><el-button type="text" size="mini">{{ lockButtonShow(menuTab) }}</el-button></li>
                <li @click="setScreenfull()"><el-button type="text" size="mini">{{ $t('FullScreen') }}</el-button></li>
                <li @click="removeTab(menuTab)"><el-button type="text" size="mini">{{$t('Close')}}</el-button></li>
                <el-divider></el-divider>
                <li @click="renameTab()"><el-button type="text" size="mini">{{$t('Rename')}}</el-button></li>
                <el-divider></el-divider>
                <li @click="closeTabs('left')"><el-button type="text" size="mini">{{$t('CloseLeft')}}</el-button></li>
                <li @click="closeTabs('right')"><el-button type="text" size="mini">{{$t('CloseRight')}}</el-button></li>
                <li @click="closeTabs('other')"><el-button type="text"  size="mini">{{$t('CloseOther')}}</el-button></li>
                <li @click="closeTabs('all')"><el-button type="text" size="mini">{{$t('CloseAll')}}</el-button></li>
            </ul>
        </div>
    </div>
</template>

<script>
import Sortable from 'sortablejs'
import screenfull from 'screenfull'
import Terminal from '@/components/Terminal'
import {MessageBox} from 'element-ui'

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
        lockButtonShow(targetName) {
            if (this.termList.length > 0 && targetName !== '') {
                const tab = this.termList.filter(tab => tab.name === targetName)[0]
                if (tab === undefined) {
                    return this.$t('Lock')
                } else {
                    return tab.closable? this.$t('Lock'):this.$t('Unlock')
                }
            } else {
                return this.$t('Lock')
            }
        },
        copyTab() {
            this.$refs[`${this.menuTab}`][0].setSSH()
            this.openTerm()
        },
        setScreenfull() {
            if (!screenfull.isEnabled) {
                this.$message({
                    message: 'not support fullscreen',
                    type: 'warning'
                })
                return false
            }
            screenfull.toggle()
        },
        getCurrMenuIndex() {
            let index = 0
            this.termList.forEach((tab, i) => { 
                if (tab.name === this.menuTab) {
                    index = i
                    return
                }
            })
            return index
        },
        closeTabs(par) {
            const setCurrentTerm = () => {
                this.currentTermIndex = currMenuIndex
                const tab = this.termList[currMenuIndex]
                this.currentTerm = tab.name
                document.title = tab.label
                this.$store.commit('SET_TAB', this.termList[this.currentTermIndex])
                this.$refs[`${tab.name}`][0].setSSH()
            }
            const filterTerm = (firstIndex, lastIndex) => {
                let tempList = []
                this.termList.forEach((tab, index) => {
                    if ((index >= firstIndex && index < lastIndex) || !tab.closable) {
                        tempList.push(tab)
                    }
                })
                this.termList = tempList
            }
            let currMenuIndex = this.getCurrMenuIndex()
            switch (par) {
            case 'left':
                // 删除左侧tab标签
                if (this.currentTermIndex < currMenuIndex) {
                    setCurrentTerm()
                }
                filterTerm(currMenuIndex, this.termList.length)
                break
            case 'right':
                // 删除右侧tab标签
                if (this.currentTermIndex > currMenuIndex) {
                    setCurrentTerm()
                }
                filterTerm(0, currMenuIndex + 1)
                break
            case 'other':
                // 删除其他所有tab标签
                if (this.currentTermIndex !== currMenuIndex) {
                    setCurrentTerm()
                }
                filterTerm(currMenuIndex, currMenuIndex + 1)
                break
            case 'all':
                filterTerm(-1, -1)
                currMenuIndex = this.getCurrMenuIndex()
                setCurrentTerm()
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
                path: '/',
                closable: true
            })
            const tab = this.termList[this.termList.length - 1]
            this.currentTerm = tab.name
            this.currentTermIndex = this.termList.length - 1
            this.$store.commit('SET_TAB', this.termList[this.currentTermIndex])
        },
        findTerm() {
            this.currentTermIndex = this.getCurrMenuIndex()
            this.$store.commit('SET_TAB', this.termList[this.currentTermIndex])
        },
        clickTab(tab) {
            this.$refs[`${tab.name}`][0].setSSH()
            document.title = tab.label
            this.findTerm()
        },
        removeTab(targetName) {
            let activeName = this.currentTerm
            for (let i = 0; i < this.termList.length; ++i) {
                if (targetName == this.termList[i].name) {
                    if (!this.termList[i].closable) {
                        this.$message({
                            message: this.$t('unlockClose'),
                            type: 'warning'
                        })
                        return
                    }
                    const nextTab = this.termList[i + 1] || this.termList[i - 1]
                    if (nextTab) {
                        activeName = nextTab.name
                    }
                }
            }
            this.currentTerm = activeName
            this.$refs[`${this.currentTerm}`][0].setSSH()
            this.termList = this.termList.filter(tab => tab.name !== targetName)
            this.findTerm()
        },
        async renameTab() {
            for (const tab of this.termList) {
                if (tab.name === this.menuTab) {
                    let {value} = await MessageBox.prompt('', this.$t('Rename'), {
                        inputValue: tab.label,
                        inputErrorMessage: 'please input value',
                        inputValidator: function (label) {
                            return label !== null && label.length > 0
                        }
                    }).catch(null)
                    tab.label = value
                    if (this.currentTerm === this.menuTab) {
                        document.title = tab.label
                    }
                    break
                }
            }
        },
        lockTab() {
            for (let tab of this.termList) {
                if (tab.name === this.menuTab) {
                    tab.closable = !tab.closable
                }
            }
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
