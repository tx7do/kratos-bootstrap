package bootstrap

import (
	"os"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

// NewAppInfo 创建应用信息
func NewAppInfo(name, version, id string) *conf.AppInfo {
	if id == "" {
		id, _ = os.Hostname()
	}
	return &conf.AppInfo{
		Name:     name,
		Version:  version,
		AppId:    id,
		Metadata: map[string]string{},
	}
}

// SetInstanceId 设置实例ID
func SetInstanceId(appInfo *conf.AppInfo, appId, name string) string {
	if appInfo == nil {
		return ""
	}
	if appId != "" {
		appInfo.AppId = appId
	}
	if name != "" {
		appInfo.Name = name
	}

	appInfo.InstanceId = appInfo.AppId + "." + appInfo.Name
	return appInfo.InstanceId
}
