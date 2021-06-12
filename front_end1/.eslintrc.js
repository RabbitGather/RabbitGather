module.exports = {
  root: true,
  env: {
    node: true,
  },
  extends: [
    // "eslint:recommended",

    "plugin:vue/vue3-essential",
    "@vue/typescript/recommended",
    "@vue/prettier",
    "@vue/prettier/@typescript-eslint",
    "plugin:@typescript-eslint/recommended",
    //
    "@vue/typescript",
    // "prettier",
    // "plugin:@typescript-eslint/recommended",
    // "prettier/@typescript-eslint",
    // "plugin:prettier/recommended",
    // "plugin:vue/recommended",
    // "@vue/prettier",
  ],
  parser: "vue-eslint-parser",

  parserOptions: {
    // ecmaVersion: 2020,
    parser: "@typescript-eslint/parser",
  },
  rules: {
    "no-console": process.env.NODE_ENV === "production" ? "warn" : "off",
    "no-debugger": process.env.NODE_ENV === "production" ? "warn" : "off",
  },
  overrides: [
    {
      files: [
        "**/__tests__/*.{j,t}s?(x)",
        "**/tests/unit/**/*.spec.{j,t}s?(x)",
      ],
      rules: {
        "no-undef": "off",
      },
      env: {
        jest: true,
        browser: true,
        // "node": true,
        // "jasmine": true
        // browser: true, // 浏览器环境中的全局变量
        // node: true, // Node.js全局变量和Node.js作用域。
        // commonjs: true, //CommonJS全局变量和CommonJS作用域（使用Browserify/webpack的浏览器代码）
        // "shared-node-browser": true, // Node.js 和 Browser 通用全局变量
        // es6: true, // 启用除modules以外的所有ES6特性（该选项会自动设置 ecmaVersion 解析器选项为 6）
        // worker: true, // Web Workers 全局变量
        // amd: true, // 将require() 和 define() 定义为像 amd 一样的全局变量
        // mocha: true, // 添加所有的 Mocha 测试全局变量
        // jasmine: true, // 添加所有的 Jasmine 版本 1.3 和 2.0 的测试全局变量
        // jest: true, // jest 全局变量
        // phantomjs: true, // PhantomJS 全局变量
        // jquery: true, // jQuery 全局变量
        // mongo: true, // MongoDB 全局变量
      },
    },
  ],
};
