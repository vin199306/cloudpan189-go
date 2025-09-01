package config

import (
	"fmt"
	"net"
	"strings"
)

// SetDNSServer 设置DNS服务器
func (c *PanConfig) SetDNSServer(dnsServer string) error {
	// 验证DNS服务器地址格式
	if dnsServer != "" {
		if !strings.Contains(dnsServer, ":") {
			dnsServer = dnsServer + ":53"
		}
		_, err := net.ResolveUDPAddr("udp", dnsServer)
		if err != nil {
			return fmt.Errorf("invalid DNS server address: %s", err)
		}
	}

	c.DNSServer = dnsServer
	return nil
}

// GetCurrentDNS 获取当前使用的DNS服务器
func (c *PanConfig) GetCurrentDNS() string {
	if c.DNSServer != "" {
		return c.DNSServer
	}
	if len(c.DefaultDNS) > 0 {
		return c.DefaultDNS[0]
	}
	return ""
}

// SwitchToNextDNS 切换到下一个默认DNS服务器
func (c *PanConfig) SwitchToNextDNS() string {
	if c.DNSServer != "" {
		// 如果设置了自定义DNS，优先使用
		return c.DNSServer
	}

	currentDNS := c.GetCurrentDNS()
	if currentDNS == "" || len(c.DefaultDNS) <= 1 {
		return c.DefaultDNS[0]
	}

	// 找到当前DNS在列表中的位置
	for i, dns := range c.DefaultDNS {
		if dns == currentDNS {
			// 切换到下一个DNS，如果是最后一个则回到第一个
			if i < len(c.DefaultDNS)-1 {
				return c.DefaultDNS[i+1]
			}
			return c.DefaultDNS[0]
		}
	}

	return c.DefaultDNS[0]
}
