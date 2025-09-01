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
	// 强制使用IPv4优先，避免IPv6回环地址问题
	if customDNSServer != "" {
		// 使用自定义DNS服务器，明确指定IPv4
		return lookupWithDNSServer(host, customDNSServer)
	}
	
	// 系统默认DNS，但避免使用IPv6回环
	// 使用备用DNS服务器列表
	for _, dnsServer := range defaultServers {
		ips, err := lookupWithDNSServer(host, dnsServer)
		if err == nil && len(ips) > 0 {
			return ips, nil
		}
	}
	
	// 最后回退到系统默认，但使用IPv4
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

// lookupWithSystemDNS 使用系统DNS但避免IPv6回环
func lookupWithSystemDNS(host string) ([]net.IP, error) {
	// 创建自定义解析器，避免IPv6回环
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			// 明确使用IPv4地址，避免[::1]:53
			if network == "udp" || network == "tcp" {
				network = "udp4"
			}
			
			// 如果地址是本地DNS回环，使用公共DNS
			if address == "[::1]:53" || address == "127.0.0.1:53" {
				// 使用公共DNS服务器
				dialer := net.Dialer{Timeout: 5 * time.Second}
				return dialer.DialContext(ctx, "udp4", "8.8.8.8:53")
			}
			
			dialer := net.Dialer{Timeout: 5 * time.Second}
			return dialer.DialContext(ctx, network, address)
		},
	}
	
	// 优先使用IPv4解析
	ips, err := resolver.LookupIP(context.Background(), "ip4", host)
	if err == nil && len(ips) > 0 {
		return ips, nil
	}
	
	// 如果IPv4失败，尝试获取所有IP
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