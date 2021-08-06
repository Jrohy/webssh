export default {
    sshReq: state => window.btoa(
        `{
            "username":"${state.sshInfo.username}", 
            "ipaddress":"${state.sshInfo.host}", 
            "port":${state.sshInfo.port}, 
            "password":"${state.sshInfo.password.replace(/[\n]/g, '\\n')}",
            "logintype":${state.sshInfo.logintype === undefined ? 0 : state.sshInfo.logintype}
        }`
    )
}
