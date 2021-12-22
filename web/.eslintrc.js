module.exports = {
    root: true,
    env: {
        node: true
    },
    extends: [
        'plugin:vue/essential',
        'eslint:recommended'
    ],
    rules: {
        'no-debugger': process.env.NODE_ENV === 'production' ? 'error' : 'off',
        'space-before-function-paren': 0,
        'quotes': [2, 'single', { 'avoidEscape': true }],
        indent: [2, 4]
    },
    parserOptions: {
        parser: 'babel-eslint'
    }
}
