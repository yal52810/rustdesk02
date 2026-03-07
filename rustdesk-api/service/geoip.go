package service

import (
	"net"
	"strings"
)

// GeoIPService Geo-IP 地理位置识别服务
type GeoIPService struct{}

// RegionType 地区类型
type RegionType string

const (
	RegionCN       RegionType = "CN"       // 中国大陆
	RegionHK       RegionType = "HK"       // 香港
	RegionTW       RegionType = "TW"       // 台湾
	RegionMO       RegionType = "MO"       // 澳门
	RegionAsia     RegionType = "ASIA"     // 亚洲其他
	RegionUS       RegionType = "US"       // 美国
	RegionEU       RegionType = "EU"       // 欧洲
	RegionOther    RegionType = "OTHER"    // 其他
	RegionInternal RegionType = "INTERNAL" // 内网
)

// GetRegionByIP 根据 IP 地址识别地区
func (s *GeoIPService) GetRegionByIP(ip string) RegionType {
	// 检查是否为内网 IP
	if s.isPrivateIP(ip) {
		return RegionInternal
	}

	// 简化版本：基于 IP 段判断
	// 生产环境建议使用 GeoIP2 数据库或第三方 API
	ipAddr := net.ParseIP(ip)
	if ipAddr == nil {
		return RegionOther
	}

	// 中国大陆 IP 段（示例，实际需要完整的 IP 库）
	cnRanges := []string{
		"1.0.0.0/8", "14.0.0.0/8", "27.0.0.0/8", "36.0.0.0/8",
		"42.0.0.0/8", "49.0.0.0/8", "58.0.0.0/8", "59.0.0.0/8",
		"60.0.0.0/8", "61.0.0.0/8", "101.0.0.0/8", "106.0.0.0/8",
		"110.0.0.0/8", "111.0.0.0/8", "112.0.0.0/8", "113.0.0.0/8",
		"114.0.0.0/8", "115.0.0.0/8", "116.0.0.0/8", "117.0.0.0/8",
		"118.0.0.0/8", "119.0.0.0/8", "120.0.0.0/8", "121.0.0.0/8",
		"122.0.0.0/8", "123.0.0.0/8", "124.0.0.0/8", "125.0.0.0/8",
		"175.0.0.0/8", "180.0.0.0/8", "182.0.0.0/8", "183.0.0.0/8",
		"202.0.0.0/8", "203.0.0.0/8", "210.0.0.0/8", "211.0.0.0/8",
		"218.0.0.0/8", "219.0.0.0/8", "220.0.0.0/8", "221.0.0.0/8",
		"222.0.0.0/8", "223.0.0.0/8",
	}

	for _, cidr := range cnRanges {
		if s.ipInRange(ip, cidr) {
			return RegionCN
		}
	}

	// 其他地区可以继续添加
	return RegionOther
}

// isPrivateIP 检查是否为内网 IP
func (s *GeoIPService) isPrivateIP(ip string) bool {
	privateRanges := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"127.0.0.0/8",
		"169.254.0.0/16",
		"::1/128",
		"fc00::/7",
		"fe80::/10",
	}

	for _, cidr := range privateRanges {
		if s.ipInRange(ip, cidr) {
			return true
		}
	}
	return false
}

// ipInRange 检查 IP 是否在 CIDR 范围内
func (s *GeoIPService) ipInRange(ip, cidr string) bool {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return false
	}
	ipAddr := net.ParseIP(ip)
	if ipAddr == nil {
		return false
	}
	return ipNet.Contains(ipAddr)
}

// IsChinaRegion 判断是否为中国地区（包括大陆、香港、台湾、澳门）
func (s *GeoIPService) IsChinaRegion(region RegionType) bool {
	return region == RegionCN || region == RegionHK || region == RegionTW || region == RegionMO
}

// GetClientIP 从请求中获取真实客户端 IP
func (s *GeoIPService) GetClientIP(remoteAddr, xForwardedFor, xRealIP string) string {
	// 优先使用 X-Real-IP
	if xRealIP != "" && xRealIP != "unknown" {
		return strings.TrimSpace(strings.Split(xRealIP, ",")[0])
	}

	// 其次使用 X-Forwarded-For
	if xForwardedFor != "" && xForwardedFor != "unknown" {
		return strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	}

	// 最后使用 RemoteAddr
	if remoteAddr != "" {
		// 移除端口号
		if idx := strings.LastIndex(remoteAddr, ":"); idx != -1 {
			return remoteAddr[:idx]
		}
		return remoteAddr
	}

	return ""
}
