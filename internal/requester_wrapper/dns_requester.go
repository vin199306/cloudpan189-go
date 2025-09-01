package requester_wrapper

import "github.com/tickstep/cloudpan189-go/internal/dns_resolver"

// SetDNSServer 设置DNS服务器
func SetDNSServer(dnsServer string) {
	dns_resolver.SetDNSServer(dnsServer)
}