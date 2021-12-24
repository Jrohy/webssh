<template>
    <div class="MainContainer">
        <el-button type="primary" size="small" @click="getFileList(); dialogVisible = true">文件管理</el-button>
        <el-dialog :title="'文件管理(' + this.$store.state.sshInfo.host + ')'" :visible.sync="dialogVisible" top="5vh" :width="dialogWidth">
            <el-row>
                <el-col :span="18">
                    <el-input v-model="currentPath" @keyup.enter.native="getFileList()"></el-input>
                </el-col>
                <el-col :span="6">
                    <el-button-group>
                        <el-button-group>
                            <el-button type="primary" size="mini" icon="el-icon-arrow-up" @click="upDirectory()"></el-button>
                            <el-button type="primary" size="mini" icon="el-icon-refresh" @click="getFileList()"></el-button>
                            <el-button type="primary" size="mini" icon="el-icon-upload" @click="openUploadDialog()"></el-button>
                        </el-button-group>
                        <el-dialog title="文件上传" :visible.sync="uploadVisible" append-to-body :width="uploadWidth">
                            <el-upload class="upload-demo" multiple drag :action="uploadUrl" :data="uploadData" :before-upload="beforeUpload" :on-progress="uploadProgress" :on-success="uploadSuccess">
                                <i class="el-icon-upload"></i>
                                <div class="el-upload__text">将文件拖到此处，或<em>点击上传</em></div>
                                <div class="el-upload__tip" slot="tip">{{ this.uploadTip }}</div>
                            </el-upload>
                        </el-dialog>
                    </el-button-group>
                </el-col>
            </el-row>
            <el-table :data="fileList" :height="clientHeight" @row-dblclick="rowClick">
                <el-table-column
                    label="名字"
                    :width="nameWidth"
                    sortable :sort-method="nameSort">
                    <template slot-scope="scope">
                        <p v-if="scope.row.IsDir === true" style="color:#0c60b5" class="el-icon-folder"> {{ scope.row.Name }}</p>
                        <p v-else-if="scope.row.IsDir === false" class="el-icon-document"> {{ scope.row.Name }}</p>
                    </template>
                </el-table-column>
                <el-table-column label="大小" prop="Size"></el-table-column>
                <el-table-column label="修改时间" prop="ModifyTime" sortable></el-table-column>
            </el-table>
        </el-dialog>
    </div>
</template>

<script>
import { fileList } from '@/api/file'
import { mapState } from 'vuex'

export default {
    name: 'FileList',
    data () {
        return {
            uploadVisible: false,
            dialogVisible: false,
            fileList: [],
            downloadFilePath: '',
            currentPath: '',
            clientHeight: 0,
            uploadTip: '',
            dialogWidth: '50%',
            uploadWidth: '32%',
            nameWidth: 260,
            progressPercent: 0
        }
    },
    created() {
        this.clientHeight = document.body.clientHeight - 200
    },
    mounted() {
        this.setDialogWidth()
        window.onresize = () => {
            this.setDialogWidth()
            this.clientHeight = document.body.clientHeight - 200
        }
    },
    computed: {
        ...mapState(['currentTab']),
        uploadUrl: () => {
            return `${process.env.NODE_ENV === 'production' ? `${location.origin}` : 'api'}/file/upload`
        },
        uploadData: function() {
            return {
                sshInfo: this.$store.getters.sshReq,
                path: this.currentPath
            }
        }
    },
    watch: {
        currentTab: function() {
            this.fileList = []
            this.currentPath = this.currentTab.path
        }
    },
    methods: {
        setDialogWidth() {
            const clientWith = document.body.clientWidth
            if (clientWith < 600) {
                this.dialogWidth = '98%'
                this.uploadWidth = '100%'
                this.nameWidth = 120
            } else if (clientWith >= 600 && clientWith < 1000) {
                this.dialogWidth = '80%'
                this.uploadWidth = '58%'
                this.nameWidth = 220
            } else {
                this.dialogWidth = '50%'
                this.uploadWidth = '32%'
                this.nameWidth = 260
            }
        },
        openUploadDialog() {
            this.uploadTip = `当前上传目录: ${this.currentPath}`
            this.uploadVisible = true
        },
        beforeUpload(file) {
            this.uploadTip = `正在上传${file.name} 到 ${this.currentPath}, 请勿关闭窗口..`
            this.uploadData.id = file.uid
            return true
        },
        uploadSuccess(r, file) {
            this.uploadTip = `${file.name}上传完成!`
        },
        uploadProgress(e, f) {
            e.percent = e.percent / 2
            f.percentage = f.percentage / 2
            if (e.percent === 50) {
                const ws = new WebSocket(`${(location.protocol === 'http:' ? 'ws' : 'wss')}://${location.host}${process.env.NODE_ENV === 'production' ? '' : '/ws'}/file/progress?id=${f.uid}`)
                ws.onmessage = e1 => {
                    f.percentage = (f.size + Number(e1.data)) / (f.size * 2) * 100
                }
                ws.onclose = () => {
                    console.log(Date(), 'onclose')
                }
                ws.onerror = () => {
                    console.log(Date(), 'onerror')
                }
            }
        },
        nameSort(a, b) {
            return a.Name > b.Name
        },
        rowClick(row) {
            if (row.IsDir) {
                // 文件夹处理
                this.currentPath = this.currentPath.charAt(this.currentPath.length - 1) === '/' ? this.currentPath + row.Name : this.currentPath + '/' + row.Name
                this.getFileList()
            } else {
                // 文件处理
                this.downloadFilePath = this.currentPath.charAt(this.currentPath.length - 1) === '/' ? this.currentPath + row.Name : this.currentPath + '/' + row.Name
                this.downloadFile()
            }
        },
        async getFileList() {
            this.currentPath = this.currentPath.replace(/\/+/g, '/')
            if (this.currentPath === '') {
                this.currentPath = '/'
            }
            const result = await fileList(this.currentPath, this.$store.getters.sshReq)
            if (result.Msg === 'success') {
                if (result.Data.list === null) {
                    this.fileList = []
                } else {
                    this.fileList = result.Data.list
                }
                this.updatePath(this.currentPath)
            } else {
                this.fileList = []
                this.$message.error(result.Msg)
                this.updatePath('/')
            }
        },
        upDirectory() {
            if (this.currentPath === '/') {
                return
            }
            let pathList = this.currentPath.split('/')
            if (pathList[pathList.length - 1] === '') {
                pathList = pathList.slice(0, pathList.length - 2)
            } else {
                pathList = pathList.slice(0, pathList.length - 1)
            }
            this.currentPath = pathList.length === 1 ? '/' : pathList.join('/')
            this.getFileList()
        },
        downloadFile() {
            const prefix = process.env.NODE_ENV === 'production' ? `${location.origin}` : 'api'
            const downloadUrl = `${prefix}/file/download?path=${this.downloadFilePath}&sshInfo=${this.$store.getters.sshReq}`
            window.open(downloadUrl)
        },
        updatePath(path) {
            const termList = this.$store.state.termList
            for (let i = 0; i < termList.length; ++i) {
                if (termList[i].name === this.currentTab.name) {
                    termList[i].path = path
                    break
                }
            }
            this.$store.commit('SET_TERMLIST', termList)
        }
    }
}
</script>

<style lang="scss">
.MainContainer {
    .el-dialog__wrapper {
        overflow: hidden;
    }
    .el-input__inner {
        border: 0 none;
        border-bottom: 1px solid #ccc;
        border-radius: 0px;
        width: 80%;
    }
    .el-table--border tr,td{
        border: none!important;
    }
    .el-table::before{
        height:0;
    }
   .el-table td, .el-table th {
        padding: 2px 0;
    }
}
</style>
