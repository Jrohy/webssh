export default {
    SET_PASS(state, pass) {
        state.sshInfo.password = pass
    },
    SET_LIST(state, list) {
        state.sshList = list
        localStorage.setItem('sshList', list)
    },
    SET_TERMLIST(state, list) {
        state.termList = list
    },
    SET_SSH(state, ssh) {
        state.sshInfo.host = ssh.host
        state.sshInfo.username = ssh.username
        state.sshInfo.port = ssh.port
        state.sshInfo.logintype = ssh.logintype
        if (ssh.password !== undefined) {
            state.sshInfo.password = ssh.password
        }
    },
    SET_TAB(state, tab) {
        state.currentTab = tab
    }
}
