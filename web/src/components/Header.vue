<template>
    <el-row>
        <el-col>
            <el-form :inline="true" style="padding-top: 10px;" :model="sshInfo" :rules="checkRules">
                <el-form-item size="small" prop="host">
                    <template slot="label">
                        <el-tooltip effect="dark" placement="left">
                            <div slot="content">
                                <p>Switch Language</p>
                            </div>
                            <span @click="handleSetLanguage()">Host</span>
                        </el-tooltip>
                    </template>
                    <el-input v-model="sshInfo.host" :placeholder="$t('hostTip')" @keyup.enter.native="$emit('ssh-select')"></el-input>
                </el-form-item>
                <el-form-item label="Port" size="small" prop="port">
                    <el-input v-model="sshInfo.port" :placeholder="$t('portTip')" @keyup.enter.native="$emit('ssh-select')" style="width: 100px"></el-input>
                </el-form-item>
                <el-form-item label="Username" size="small" prop="username">
                    <el-input v-model="sshInfo.username" :placeholder="$t('nameTip')" @keyup.enter.native="$emit('ssh-select')" style="width: 110px"></el-input>
                </el-form-item>
                <el-form-item size="small" prop="password">
                    <template slot="label">
                        <el-tooltip effect="dark" placement="left">
                            <div slot="content">
                                <p>{{ `Switch to ${this.privateKey ? 'Password' : 'PrivateKey'} login` }}</p>
                            </div>
                            <span @click="sshInfo.logintype === 0 ? sshInfo.logintype=1: sshInfo.logintype=0">{{ privateKey?'PrivateKey':'Password' }}</span>
                        </el-tooltip>
                    </template>
                    <el-input v-model="sshInfo.password" @click.native="textareaVisible=privateKey" @keyup.enter.native="$emit('ssh-select')" :placeholder="$t('inputTip') + `${this.privateKey ? $t('privateKey') : $t('password')}`" show-password></el-input>
                </el-form-item>
                <el-dialog :title="$t('privateKey')" :visible.sync="textareaVisible" :close-on-click-modal="false">
                    <el-input :rows="8" v-model="sshInfo.password" type="textarea" :placeholder="$t('keyTip')"></el-input>
                    <div slot="footer" class="dialog-footer">
                        <!-- 选择密钥文件 -->
                        <input ref="pkFile" @change="handleChangePKFile" type="file" style="position: absolute;clip: rect(0 0 0 0)"/>
                        <el-button type="primary" plain @click="$refs.pkFile.click()">{{ $t('SelectFile') }}</el-button>
                        <el-button @click="sshInfo.password=''">{{ $t('Clear') }}</el-button>
                        <el-button type="primary" @click="textareaVisible = false; $emit('ssh-select')">{{ $t('Connect') }}</el-button>
                    </div>
                </el-dialog>
                <el-form-item  size="small">
                    <el-button type="primary" @click="$emit('ssh-select')" plain>{{ $t('Connect') }}</el-button>
                </el-form-item>
                <el-form-item  size="small">
                    <file-list></file-list>
                </el-form-item>
                <el-form-item size="small">
                    <el-dropdown @command="handleCommand">
                        <el-button type="primary">
                            {{ $t('History') }}
                        </el-button>
                        <el-dropdown-menu slot="dropdown">
                            <el-dropdown-item
                                v-for="item in sshList"
                                :key="item.host" :command="item" style="padding:0 5px 0 10px">
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
import { getLanguage } from '@/lang/index'
import FileList from '@/components/FileList'
import { mapState } from 'vuex'

export default {
    components: {
        'file-list': FileList
    },
    data() {
        return {
            textareaVisible: false,
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
                    { required: true, trigger: 'blur', message: 'value is required' }
                ]
            }
        }
    },
    methods: {
        handleSetLanguage() {
            const oldLang = getLanguage()
            const lang = oldLang === 'zh' ? 'en' : 'zh'
            this.$i18n.locale = lang
            this.$store.dispatch('setLanguage', lang)
        },
        handleCommand(command) {
            this.$store.commit('SET_SSH', command)
            if (command.password === undefined) {
                this.$store.commit('SET_PASS', '')
            }
            // 新开窗口
            this.$emit('ssh-select')
        },
        cleanHistory(command) {
            const sshListObj = this.sshList
            sshListObj.forEach((v, i) => {
                if (v.host === command.host) {
                    sshListObj.splice(i, 1)
                }
            })
            this.$store.commit('SET_LIST', window.btoa(JSON.stringify(sshListObj)))
        },
        handleChangePKFile(event) {
            const file = event.target.files[0]
            if (file) {
                const sshInfo = this.sshInfo
                const reader = new FileReader()
                reader.onload = e => {
                    sshInfo.password = e.target.result
                }
                reader.readAsText(file)
            }
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
        privateKey() {
            return this.sshInfo.logintype === 1
        },
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
