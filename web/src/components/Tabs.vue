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
    </div>
</template>

<script>
import Terminal from '@/components/Terminal'

export default {
    name: 'Tabs',
    components: {
        terminal: Terminal
    },
    data () {
        return {
            termList: [],
            currentTerm: ''
        }
    },
    methods: {
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
                label: sshInfo.host
            })
            this.currentTerm = this.termList[this.termList.length - 1].name
        },
        clickTab(tab) {
            this.$refs[`${tab.name}`][0].setSSH()
            document.title = tab.label
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
            }
            this.termList = tabs.filter(tab => tab.name !== targetName)
        }
    }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
