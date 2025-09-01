package dns_resolver

import (
	"context"
	"fmt"
	"net"
	"time"
)

// SimpleDNSResolver 简单的DNS解析器
var (
	customDNSServer string
	defaultServers = []string{"114.114.114.114", "8.8.8.8"}
	
	// IPv6 DNS服务器配置
	ipv6DNSServers = []string{
		// 腾讯DNS IPv6地址
		"2402:4e00::",
		"2402:4e00:1::",
		// 阿里DNS IPv6地址
		"2400:3200::1",
		"2400:3200:baba::1",
	}
	
	// 是否忽略系统默认DNS（解决Linux系统/etc/resolv.conf不存在的问题）
	ignoreSystemDNS = true
)

// SetDNSServer 设置自定义DNS服务器
func SetDNSServer(server string) {
	customDNSServer = server
	if server != "" {
		fmt.Printf("DNS服务器已设置为: %s\n", server)
	} else {
		fmt.Println("DNS服务器已重置为系统默认")
	}
}

// LookupIP 解析域名到IP地址
func LookupIP(host string) ([]net.IP, error) {
	// 如果有自定义DNS服务器，优先使用
	if customDNSServer != "" {
		return lookupWithDNSServer(host, customDNSServer)
	}
	
	// 忽略系统默认DNS，直接使用公共DNS服务器
	// 解决Linux系统/etc/resolv.conf不存在的问题
	if ignoreSystemDNS {
		// 使用公共DNS服务器列表，避免依赖系统配置
		allServers := append([]string{}, defaultServers...)
		allServers = append(allServers, ipv6DNSServers...)
		
		for _, dnsServer := range allServers {
			ips, err := lookupWithDNSServer(host, dnsServer)
			if err == nil && len(ips) > 0 {
				return ips, nil
			}
		}
		
		// 如果所有公共DNS都失败，使用硬编码的可靠DNS
		fallbackServers := []string{"8.8.8.8", "1.1.1.1", "114.114.114.114"}
		for _, dnsServer := range fallbackServers {
			ips, err := lookupWithDNSServer(host, dnsServer)
			if err == nil && len(ips) > 0 {
				return ips, nil
			}
		}
		
		return nil, fmt.Errorf("所有DNS服务器都无法解析 %s", host)
	}
	
	// 回退到系统默认（不推荐，仅作兼容）
	return lookupWithSystemDNS(host)
}

// lookupWithDNSServer 使用指定DNS服务器解析域名
func lookupWithDNSServer(host, dnsServer string) ([]net.IP, error) {
	// 创建自定义解析器
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			dialer := net.Dialer{Timeout: 5 * time.Second}
			
			// 根据DNS服务器地址类型选择合适的网络协议
			var networkType string
			if isIPv6Address(dnsServer) {
				networkType = "udp6"
			} else {
				networkType = "udp4"
			}
			
			// 确保DNS服务器地址格式正确
			serverAddr := formatDNSServerAddress(dnsServer)
			return dialer.DialContext(ctx, networkType, serverAddr)
		},
	}
	
	// 根据DNS服务器类型决定解析策略
	if isIPv6Address(dnsServer) {
		// IPv6 DNS服务器，支持IPv6和IPv4解析
		return resolver.LookupIP(context.Background(), "ip", host)
	} else {
		// IPv4 DNS服务器，优先IPv4
		ips, err := resolver.LookupIP(context.Background(), "ip4", host)
		if err == nil && len(ips) > 0 {
			return ips, nil
		}
		// IPv4失败时尝试IPv6
		return resolver.LookupIP(context.Background(), "ip6", host)
	}
}

// lookupWithSystemDNS 使用系统DNS但避免localhost问题
// 注意：此函数仅作向后兼容，实际使用中会优先使用公共DNS
func lookupWithSystemDNS(host string) ([]net.IP, error) {
	// 创建一个完全独立的解析器，避免系统配置
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			// 强制使用公共DNS服务器，完全忽略系统配置
			// 避免使用localhost (127.0.0.1:53 或 [::1]:53)
			
			// 使用硬编码的公共DNS服务器列表
			publicDNS := []string{
				"8.8.8.8:53",    // Google
				"8.8.4.4:53",    // Google
				"1.1.1.1:53",    // Cloudflare
				"1.0.0.1:53",    // Cloudflare
				"114.114.114.114:53", // 114DNS
				"223.5.5.5:53",  // AliDNS
				"2402:4e00::53", // Tencent IPv6
				"2400:3200::1:53", // Ali IPv6
			}
			
			dialer := net.Dialer{Timeout: 5 * time.Second}
			
			// 根据网络类型选择合适的DNS服务器
			var server string
			if network == "udp6" || network == "tcp6" {
				// IPv6网络，使用IPv6 DNS
				server = publicDNS[6] // 2402:4e00::53
			} else {
				// IPv4网络，使用IPv4 DNS
				server = publicDNS[0] // 8.8.8.8:53
			}
			
			return dialer.DialContext(ctx, network, server)
		},
	}
	
	// 尝试IPv4解析
	ips, err := resolver.LookupIP(context.Background(), "ip4", host)
	if err == nil && len(ips) > 0 {
		return ips, nil
	}
	
	// 尝试IPv6解析
	ips, err = resolver.LookupIP(context.Background(), "ip6", host)
	if err == nil && len(ips) > 0 {
		return ips, nil
	}
	
	// 尝试双栈解析
	return resolver.LookupIP(context.Background(), "ip", host)
}

// ForceIPv4 强制使用IPv4解析，完全避免IPv6
func ForceIPv4(host string) ([]net.IP, error) {
	// 使用公共DNS服务器，完全避免系统DNS
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			dialer := net.Dialer{Timeout: 5 * time.Second}
			// 使用可靠的公共DNS
			servers := []string{"8.8.8.8:53", "114.114.114.114:53", "1.1.1.1:53"}
			for _, server := range servers {
				conn, err := dialer.DialContext(ctx, "udp4", server)
				if err == nil {
					return conn, nil
				}
			}
			return nil, fmt.Errorf("all DNS servers failed")
		},
	}
	
	return resolver.LookupIP(context.Background(), "ip4", host)
}

// isIPv6Address 判断是否为IPv6地址
func isIPv6Address(addr string) bool {
	ip := net.ParseIP(addr)
	return ip != nil && ip.To4() == nil
}

// formatDNSServerAddress 格式化DNS服务器地址
func formatDNSServerAddress(addr string) string {
	// 检查是否已经包含端口
	if _, _, err := net.SplitHostPort(addr); err == nil {
		return addr
	}
	
	// 添加默认DNS端口53
	ip := net.ParseIP(addr)
	if ip == nil {
		return addr + ":53"
	}
	
	// IPv6地址需要加括号
	if isIPv6Address(addr) {
		return "[" + addr + "]:53"
	}
	
	return addr + ":53"
}

// GetIPv6DNSServers 获取IPv6 DNS服务器列表
func GetIPv6DNSServers() []string {
	return ipv6DNSServers
}

// GetDNSServer 获取当前DNS服务器
func GetDNSServer() string {
	if customDNSServer != "" {
		return customDNSServer
	}
	return "系统默认"
}