package bootstrap

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/tx7do/go-utils/id"
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

// printAppInfo 打印应用信息
func printAppInfo() {
	ai := GetAppInfo()
	ts := time.Now().Format(time.RFC3339)
	host, _ := os.Hostname()
	pid := os.Getpid()

	// JSON 输出（便于日志采集/自动化）
	if os.Getenv("APPINFO_FORMAT") == "json" {
		out := map[string]interface{}{
			"timestamp":   ts,
			"host":        host,
			"pid":         pid,
			"name":        ai.Name,
			"version":     ai.Version,
			"app_id":      ai.AppId,
			"instance_id": ai.InstanceId,
			"metadata":    ai.Metadata,
		}
		if b, err := json.Marshal(out); err == nil {
			fmt.Println(string(b))
		} else {
			fmt.Printf("Application info marshal error: %v\n", err)
		}
		return
	}

	// 人类可读输出，元数据按键排序
	fmt.Printf("[%s] %s (pid:%d@%s)\n", ts, ai.Name, pid, host)
	fmt.Printf("  Version: %s\n", ai.Version)
	fmt.Printf("  AppId: %s\n", ai.AppId)
	fmt.Printf("  InstanceId: %s\n", ai.InstanceId)
	if len(ai.Metadata) > 0 {
		fmt.Println("  Metadata:")
		keys := make([]string, 0, len(ai.Metadata))
		for k := range ai.Metadata {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Printf("    %s=%s\n", k, ai.Metadata[k])
		}
	}
}

// NewInstanceId 生成实例ID 格式：appId-version@host:port@xid
func NewInstanceId(appId, version, host, port string) string {
	return fmt.Sprintf("%s-%s@%s:%s@%s", appId, version, host, port, id.NewXID())
}
