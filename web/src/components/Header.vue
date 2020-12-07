<template>
    <el-row>
        <el-col>
            <el-form :inline="true" style="padding-top: 15px" :model="sshInfo" :rules="checkRules">
                <el-form-item label="Host" size="small" prop="host">
                    <el-input v-model="sshInfo.host" placeholder="请输入远程host地址" @keyup.enter.native="$emit('ssh-select')"></el-input>
                </el-form-item>
                <el-form-item label="Port" size="small" prop="port">
                    <el-input v-model="sshInfo.port" placeholder="请输入端口" @keyup.enter.native="$emit('ssh-select')" style="width: 100px"></el-input>
                </el-form-item>
                <el-form-item label="Username" size="small" prop="username">
                    <el-input v-model="sshInfo.username" placeholder="请输入用户名" @keyup.enter.native="$emit('ssh-select')" style="width: 110px"></el-input>
                </el-form-item>
                <el-form-item label="Password" size="small" prop="password">
                    <el-input v-model="sshInfo.password" @keyup.enter.native="$emit('ssh-select')" placeholder="请输入密码" show-password></el-input>
                </el-form-item>
                <el-form-item  size="small">
                    <el-button type="primary" @click="$emit('ssh-select')" plain>连接</el-button>
                </el-form-item>
                <el-form-item  size="small">
                    <file-list></file-list>
                </el-form-item>
                <el-form-item size="small">
                    <el-dropdown @command="handleCommand">
                        <el-button type="primary">
                            历史记录
                        </el-button>
                        <el-dropdown-menu slot="dropdown">
                            <el-dropdown-item
                                v-for="item in sshList"
                                :key="item.host" :command="item" style="padding:0px 5px 0px 15px">
                                {{item.host}}
                                <i @click="cleanHistory(item)" class="el-icon-close"></i>
                            </el-dropdown-item>
                        </el-dropdown-menu>
                    </el-dropdown>
                </el-form-item>
            </el-form>
        </el-col>
    </el-row>
</template>

<script>
import FileList from '@/components/FileList'
import { mapState } from 'vuex'

export default {
    components: {
        'file-list': FileList
    },
    data() {
        return {
            checkRules: {
                host: [
                    { required: true, trigger: 'blur' }
                ],
                port: [
                    { required: true, trigger: 'blur', type: 'number', transform(value) { return Number(value) } }
                ],
                username: [
                    { required: true, trigger: 'blur' }
                ],
                password: [
                    { required: true, trigger: 'blur' }
                ]
            }
        }
    },
    methods: {
        handleCommand(command) {
            this.$store.commit('SET_SSH', command)
            if (command.password === undefined) {
                this.$store.commit('SET_PASS', '')
            }
        },
        cleanHistory(command) {
            const sshListObj = this.sshList
            sshListObj.forEach((v, i) => {
                if (v.host === command.host) {
                    sshListObj.splice(i, 1)
                }
            })
            this.$store.commit('SET_LIST', window.btoa(JSON.stringify(sshListObj)))
        }
    },
    mounted() {
        if (this.sshList.length > 0) {
            const latestSSH = this.sshList[this.sshList.length - 1]
            this.$store.commit('SET_SSH', latestSSH)
            if (latestSSH.password === undefined) {
                this.$store.commit('SET_PASS', '')
            }
        }
    },
    computed: {
        ...mapState(['sshInfo']),
        sshList() {
            const sshList = this.$store.state.sshList
            if (sshList === null) {
                return []
            } else {
                return JSON.parse(window.atob(sshList))
            }
        }
    }
}
</script>

<style scoped>
</style>
