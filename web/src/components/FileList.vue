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
                        <el-dialog title="文件上传" :visible.sync="uploadVisible" append-to-body width="32%">
                            <el-upload class="upload-demo" drag :action="uploadUrl" multiple :data="uploadData" :before-upload="beforeUpload" :on-success="uploadSuccess">
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
                    sortable :sort-method="nameSort">
                    <template slot-scope="scope">
                        <p v-if="scope.row.FType === 1" style="color:#0c60b5" class="el-icon-folder"> {{ scope.row.Name }}</p>
                        <p v-else-if="scope.row.FType === 0" class="el-icon-document"> {{ scope.row.Name }}</p>
                    </template>
                </el-table-column>
                <el-table-column label="大小" prop="Size" width="100"></el-table-column>
                <el-table-column label="修改时间" prop="ModifyTime"></el-table-column>
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
            dialogWidth: '50%'
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
            this.currentPath = '/'
        }
    },
    methods: {
        setDialogWidth() {
            const clientWith = document.body.clientWidth
            if (clientWith < 600) {
                this.dialogWidth = '98%'
            } else {
                this.dialogWidth = '50%'
            }
        },
        openUploadDialog() {
            this.uploadTip = `当前上传目录: ${this.currentPath}`
            this.uploadVisible = true
        },
        beforeUpload(file) {
            this.uploadTip = `正在上传${file.name} 到 ${this.currentPath}, 请勿关闭窗口..`
            return true
        },
        uploadSuccess(response, file, fileList) {
            this.uploadTip = `${file.name}上传完成!`
        },
        nameSort(a, b) {
            return a.Name > b.Name
        },
        rowClick(row) {
            if (row.FType === 0) {
                // 文件处理
                this.downloadFilePath = this.currentPath.charAt(this.currentPath.length - 1) === '/' ? this.currentPath + row.Name : this.currentPath + '/' + row.Name
                this.downloadFile()
            } else if (row.FType === 1) {
                // 文件夹处理
                this.currentPath = this.currentPath.charAt(this.currentPath.length - 1) === '/' ? this.currentPath + row.Name : this.currentPath + '/' + row.Name
                this.getFileList()
            }
        },
        async getFileList() {
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
            } else {
                this.fileList = []
                this.$message.error(result.Msg)
            }
        },
        upDirectory() {
            if (this.currentPath === '/') {
                return
            }
            let pathList = this.currentPath.split('/')
            pathList = pathList.slice(0, pathList.length - 1)
            this.currentPath = pathList.length === 1 ? '/' : pathList.join('/')
            this.getFileList()
        },
        downloadFile() {
            const prefix = process.env.NODE_ENV === 'production' ? `${location.origin}` : 'api'
            const downloadUrl = `${prefix}/file/download?path=${this.downloadFilePath}&sshInfo=${this.$store.getters.sshReq}`
            window.open(downloadUrl)
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
