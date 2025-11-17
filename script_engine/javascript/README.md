# JavaScript脚本引擎

基于[goja](https://github.com/dop251/goja)封装的一个JavaScript脚本引擎。

最开始是基于[otto](https://github.com/robertkrimen/otto)封装的，可惜它对很多ES5的语法都不支持，用着实在是蛋疼无比。

后来改为[goja](https://github.com/dop251/goja)，goja甚至支持一些ES6的语法。其实，goja是基于otto做的开发，可以说是青出于蓝，更甚于蓝。

## 基本交互

### 变量

变量有三种访问方式：

- 单向只写 - 将宿主的变量注入；
- 单向只读 - 读取脚本中的变量；
- 双向可读可写 - 将宿主的变量注入，脚本可操作宿主的变量并且反馈到宿主。

### 方法

- 宿主调用脚本中的方法；
- 脚本调用宿主中的方法。
