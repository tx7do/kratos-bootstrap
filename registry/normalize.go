package registry

import (
	"strings"
	"unicode"
)

const (
	consulMaxRunes     = 64
	k8sMaxRunes        = 63
	defaultPlaceholder = "service"
)

// NormalizeForRegistry 根据不同注册中心的命名规则，规范化服务 ID
func NormalizeForRegistry(appId string, registry string) string {
	appId = strings.TrimSpace(appId)
	if appId == "" {
		return defaultPlaceholder
	}

	regLower := strings.ToLower(strings.TrimSpace(registry))
	switch Type(regLower) {
	case Consul:
		return normalizeConsul(appId)

	case Kubernetes:
		return normalizeKubernetes(appId)

	default:
		// etcd 等支持层次化路径，返回原始 id
		return appId
	}
}

// normalizeConsul 规范化 Consul 服务 ID
// 规则：只允许小写字母、数字和破折号，其他字符视为分隔符，连续的分隔符合并为一个破折号，
// 且不以破折号开头或结尾，最长 64 个 rune。
// 如果结果为空，则返回默认占位符 "service"。
func normalizeConsul(s string) string {
	var b strings.Builder
	b.Grow(len(s))

	lastHyphen := false
	count := 0

	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(unicode.ToLower(r))
			lastHyphen = false
			count++
		} else {
			// 其他字符视为分隔符，只有在已有内容且上一个不是 '-' 时插入一个 '-'
			if count > 0 && !lastHyphen {
				b.WriteRune('-')
				lastHyphen = true
				count++
			}
		}

		if count >= consulMaxRunes {
			break
		}
	}

	res := strings.Trim(b.String(), "-")
	if res == "" {
		return defaultPlaceholder
	}

	// 再按 rune 数量严格截断，避免截断 UTF-8 字符
	rns := []rune(res)
	if len(rns) > consulMaxRunes {
		res = string(rns[:consulMaxRunes])
		res = strings.Trim(res, "-")
		if res == "" {
			return defaultPlaceholder
		}
	}

	return res
}

// normalizeKubernetes 规范化 Kubernetes 服务名称
// 规则：只允许小写字母、数字和破折号，其他字符视为分隔符，连续的分隔符合并为一个破折号，
// 且不以破折号开头或结尾，最长 63 个 rune，且首尾必须为字母或数字。
// 如果结果为空或不符合要求，则返回默认占位符 "service"。
func normalizeKubernetes(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	b.Grow(len(s))

	lastHyphen := false
	count := 0

	for _, r := range s {
		// 只允许小写 ASCII 字母、数字；其他视为分隔符；保留单个 '-'
		if (r >= 'a' && r <= 'z') || unicode.IsDigit(r) {
			b.WriteRune(r)
			lastHyphen = false
			count++
		} else {
			// 遇到 '-' 也当做分隔符处理，但避免连续 '-'
			if count > 0 && !lastHyphen {
				b.WriteRune('-')
				lastHyphen = true
				count++
			}
		}
		if count >= k8sMaxRunes {
			break
		}
	}

	res := strings.Trim(b.String(), "-")
	if res == "" {
		return defaultPlaceholder
	}

	// 按 rune 数量严格截断并再次去除首尾 '-'
	rns := []rune(res)
	if len(rns) > k8sMaxRunes {
		res = string(rns[:k8sMaxRunes])
		res = strings.Trim(res, "-")
		if res == "" {
			return defaultPlaceholder
		}
	}

	// 确保首尾为字母或数字（DNS-1123 label 要求）
	first := rune(res[0])
	last := rune(res[len(res)-1])
	if !((first >= 'a' && first <= 'z') || unicode.IsDigit(first)) ||
		!((last >= 'a' && last <= 'z') || unicode.IsDigit(last)) {
		return defaultPlaceholder
	}

	return res
}
