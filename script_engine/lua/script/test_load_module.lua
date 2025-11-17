local my_path = "D:\\GoProject\\go-lua\\script"
package.path = package.path .. [[;]] .. GetLuaPath() .. [[/?.lua;]]
package.path = package.path .. [[;]] .. my_path .. [[/?.lua;]]
print(package.path)
print(GetLuaPath())

require("test_module")
print(module.constant)
module.func3()

local test = require("test_module1")
test:test()
