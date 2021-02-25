export default {
    sshInfo: {
        host: '',
        username: 'root',
        port: 22,
        password: ''
    },
    sshList: Object.prototype.hasOwnProperty.call(localStorage, 'sshList') ? localStorage.getItem('sshList') : null,
    termList: [],
    currentTab: {}
}
