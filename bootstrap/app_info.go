package bootstrap

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/tx7do/go-utils/id"
	"github.com/tx7do/go-utils/stringcase"
	"github.com/tx7do/go-utils/timeutil"
	"github.com/tx7do/go-utils/trans"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

const (
	defaultProject = "gowind"
	defaultAppName = "Unknown Service"
	defaultAppId   = "unknown-service"
	defaultName    = "GoWind Unknown Service"
	defaultVersion = "1.0.0"
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

	if appName == nil {
		ai.Name = defaultAppName
	} else {
		ai.Name = *appName
	}

	return ai
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

// NewInstanceId 生成实例ID 格式：project-appId-version@host@xid
func NewInstanceId(project, appId, version, host string) string {
	return fmt.Sprintf("%s-%s-%s@%s@%s", project, appId, version, host, id.NewXID())
}

// NewAppName 生成应用名称
func NewAppName(project, appId string) string {
	return stringcase.UpperCamelCase(project) + " " + stringcase.UpperCamelCase(appId)
}

// AdjustAppInfo 调整应用信息，设置默认值
func AdjustAppInfo(ai *conf.AppInfo) {
	if ai == nil {
		return
	}

	if ai.Project == "" {
		ai.Project = defaultProject
	}

	if ai.AppId == "" {
		ai.AppId = defaultAppId
	}

	if ai.Version == "" {
		ai.Version = defaultVersion
	}

	if ai.Name == "" {
		if ai.Project != "" && ai.AppId != "" {
			ai.Name = NewAppName(ai.Project, ai.AppId)
		} else {
			ai.Name = defaultName
		}
	}

	if ai.Hostname == "" {
		host, _ := os.Hostname()
		ai.Hostname = host
	}

	if ai.Metadata == nil {
		ai.Metadata = map[string]string{}
	}

	if ai.InstanceId == "" {
		host := ResolveHost()
		ai.InstanceId = NewInstanceId(ai.Project, ai.AppId, ai.Version, host)
	}

	ai.StartTime = timeutil.TimeToTimestamppb(trans.Ptr(time.Now()))
}

// ResolveHost 返回优先级选择的 host 标识：POD_NAME -> HOSTNAME env -> os.Hostname() -> 首个非 loopback IPv4 -> "unknown-host"
func ResolveHost() string {
	if v := os.Getenv("POD_NAME"); v != "" {
		return v
	}
	if v := os.Getenv("HOSTNAME"); v != "" {
		return v
	}
	if h, err := os.Hostname(); err == nil && h != "" {
		return h
	}
	// 回退到首个非 loopback IPv4
	if ip := firstNonLoopbackIP(); ip != "" {
		return ip
	}
	return "unknown-host"
}

// firstNonLoopbackIP 返回第一个符合条件的非 loopback IPv4：跳过 down、loopback、常见 docker/CNI/bridge/veth/virbr 接口，
// 并过滤链路本地与 docker 默认桥接网段 172.17.*。
func firstNonLoopbackIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, iface := range ifaces {
		// 跳过未启用或 loopback 接口
		if (iface.Flags&net.FlagUp) == 0 || (iface.Flags&net.FlagLoopback) != 0 {
			continue
		}
		// 跳过常见的容器/桥接/CNI 接口
		if isContainerLikeInterface(iface.Name) {
			continue
		}
		addrs, _ := iface.Addrs()
		for _, a := range addrs {
			var ip net.IP
			switch v := a.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				// 非 IPv4 忽略
				continue
			}
			// 过滤链路本地地址 (169.254.*) 与 IPv6 链路本地（已通过 To4() 排除）等
			if ip[0] == 169 && ip[1] == 254 {
				continue
			}
			// 过滤 docker 默认桥接网段 172.17.*（可按需调整或移除）
			if ip[0] == 172 && ip[1] == 17 {
				continue
			}
			return ip.String()
		}
	}
	return ""
}

// isContainerLikeInterface 判断接口名称是否为常见容器/桥接/CNI/veth/虚拟网卡 接口
func isContainerLikeInterface(name string) bool {
	n := strings.ToLower(name)
	// 常见容器/桥接/CNI/veth/虚拟网卡 前缀或名称
	if n == "docker0" ||
		strings.HasPrefix(n, "veth") ||
		strings.HasPrefix(n, "br-") ||
		strings.HasPrefix(n, "cni0") ||
		strings.HasPrefix(n, "flannel") ||
		strings.HasPrefix(n, "weave") ||
		strings.HasPrefix(n, "virbr") ||
		strings.HasPrefix(n, "cbr0") {
		return true
	}
	return false
}
