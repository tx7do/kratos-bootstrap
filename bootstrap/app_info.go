package bootstrap

import (
	"os"
	"sync"

	"github.com/tx7do/go-utils/trans"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

const (
	defaultAppName = "Unknown Service"
	defaultAppId   = "gowind-unknown-service"
	defaultVersion = "1.0.0"
)

var (
	// appInfo 应用信息
	appInfo = NewAppInfo(
		trans.Ptr(defaultAppId),
		trans.Ptr(defaultVersion),
		trans.Ptr(defaultAppName),
	)
	appInfoMu sync.RWMutex
)

// NewAppInfo 创建应用信息
func NewAppInfo(appId, version, appName *string) *conf.AppInfo {
	ai := &conf.AppInfo{
		Metadata: map[string]string{},
	}

	if appId == nil {
		ai.AppId = defaultAppId
	} else {
		ai.AppId = *appId
	}

	if version == nil {
		ai.Version = defaultVersion
	} else {
		ai.Version = *version
	}

	if appId != nil && *appId != "" &&
		version != nil && *version != "" {
		SetInstanceId(ai, *appId, *version)
	} else {
		hostName, _ := os.Hostname()
		appId = &hostName
	}

	if appName == nil {
		ai.Name = defaultAppName
	} else {
		ai.Name = *appName
	}

	return ai
}

// GetAppInfo 返回当前 appInfo 的只读副本（避免外部直接修改内部状态）
func GetAppInfo() *conf.AppInfo {
	appInfoMu.RLock()
	defer appInfoMu.RUnlock()

	return &conf.AppInfo{
		Name:       appInfo.Name,
		Version:    appInfo.Version,
		AppId:      appInfo.AppId,
		InstanceId: appInfo.InstanceId,
		Metadata:   cloneMetadata(appInfo.Metadata),
	}
}

// SetAppInfo 用受控方式替换整个 appInfo（可选）
func SetAppInfo(src *conf.AppInfo) {
	if src == nil {
		return
	}

	appInfoMu.Lock()
	defer appInfoMu.Unlock()

	appInfo = &conf.AppInfo{
		Name:       src.Name,
		Version:    src.Version,
		AppId:      src.AppId,
		InstanceId: src.InstanceId,
		Metadata:   cloneMetadata(src.Metadata),
	}
}

// SetInstanceId 设置实例ID
func SetInstanceId(appInfo *conf.AppInfo, appId, version string) string {
	if appInfo == nil {
		return ""
	}

	if appId == "" {
		return ""
	}

	appInfo.InstanceId = appId + "-" + version
	return appInfo.InstanceId
}

// mergeFrom 更新内部 appInfo（仅合并非空字段）
func mergeFrom(src *conf.AppInfo) {
	if src == nil {
		return
	}
	appInfoMu.Lock()
	defer appInfoMu.Unlock()

	if src.Name != "" {
		appInfo.Name = src.Name
	}
	if src.Version != "" {
		appInfo.Version = src.Version
	}
	if src.InstanceId != "" {
		appInfo.InstanceId = src.InstanceId
	}
	if src.Metadata != nil {
		appInfo.Metadata = src.Metadata
	}
}

func cloneMetadata(m map[string]string) map[string]string {
	if m == nil {
		return nil
	}
	clone := make(map[string]string, len(m))
	for k, v := range m {
		clone[k] = v
	}
	return clone
}
