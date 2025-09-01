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
	var resolver *net.Resolver
	
	if customDNSServer != "" {
		// 使用自定义DNS服务器
		resolver = &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				dialer := net.Dialer{Timeout: 5 * time.Second}
				return dialer.DialContext(ctx, "udp", fmt.Sprintf("%s:53", customDNSServer))
			},
		}
	} else {
		// 使用系统默认DNS
		resolver = net.DefaultResolver
	}
	
	return resolver.LookupIP(context.Background(), "ip", host)
}

// GetDNSServer 获取当前DNS服务器
func GetDNSServer() string {
	if customDNSServer != "" {
		return customDNSServer
	}
	return "系统默认"
}