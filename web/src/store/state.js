export default {
    sshInfo: {
        host: '',
        username: 'root',
        port: 22,
        password: '',
        logintype: 0
    },
    sshList: Object.prototype.hasOwnProperty.call(localStorage, 'sshList') ? localStorage.getItem('sshList') : null,
    termList: [],
    currentTab: {}
}
