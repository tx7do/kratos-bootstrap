local http = require("http")
local client = http.client()
local request = http.request("GET", "https://www.baidu.com")
local result, err = client:do_request(request)
if err then error(err) end
print(result.body)
