package bootstrap

import (
	"sync"

	"github.com/google/uuid"
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
		InstanceId: uuid.NewString(),
		Metadata:   map[string]string{},
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

// cloneMetadata 克隆元数据映射，避免外部修改内部状态
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
