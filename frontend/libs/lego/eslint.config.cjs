const js = require("@eslint/js");
const { FlatCompat } = require("@eslint/eslintrc");
const baseConfig = require('../../../eslint.config.cjs')

const compat = new FlatCompat({
    baseDirectory: __dirname,
    recommendedConfig: js.configs.recommended,
    allConfig: js.configs.all
});

module.exports = [
...baseConfig,
...compat.extends("plugin:@nx/react"),
{
    ignores: ["!**/*"],
},
{
    files: ["**/*.ts", "**/*.tsx", "**/*.js", "**/*.jsx"],
    rules: {},
}, {
    files: ["**/*.ts", "**/*.tsx"],
    rules: {},
}, {
    files: ["**/*.js", "**/*.jsx"],
    rules: {},
}];