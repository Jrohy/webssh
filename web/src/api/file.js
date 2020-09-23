import request from '@/utils/request'
export function fileList(path, sshInfo) {
    return request.get(`/file/list?path=${path}&sshInfo=${sshInfo}`)
}
