// Copyright (c) 2020 tickstep.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/tickstep/cloudpan189-go/cmder/cmdtable"
	"github.com/tickstep/library-go/converter"
	"github.com/tickstep/library-go/requester"
)

// SetProxy 设置代理
func (c *PanConfig) SetProxy(proxy string) {
	c.Proxy = proxy
	requester.SetGlobalProxy(proxy)
}

// SetLocalAddrs 设置localAddrs
func (c *PanConfig) SetLocalAddrs(localAddrs string) {
	c.LocalAddrs = localAddrs
	requester.SetLocalTCPAddrList(strings.Split(localAddrs, ",")...)
}

func (c *PanConfig) SetPreferIPType(ipType string) {
	c.PreferIPType = ipType
	t := requester.IPAny
	if strings.ToLower(ipType) == "ipv4" {
		t = requester.IPv4
	} else if strings.ToLower(ipType) == "ipv6" {
		t = requester.IPv6
	}
	requester.SetPreferIPType(t)
}

// SetCacheSizeByStr 设置cache_size
func (c *PanConfig) SetCacheSizeByStr(sizeStr string) error {
	size, err := converter.ParseFileSizeStr(sizeStr)
	if err != nil {
		return err
	}
	c.CacheSize = int(size)
	return nil
}

// SetMaxDownloadRateByStr 设置 max_download_rate
func (c *PanConfig) SetMaxDownloadRateByStr(sizeStr string) error {
	size, err := converter.ParseFileSizeStr(stripPerSecond(sizeStr))
	if err != nil {
		return err
	}
	c.MaxDownloadRate = size
	return nil
}

// SetMaxUploadRateByStr 设置 max_upload_rate
func (c *PanConfig) SetMaxUploadRateByStr(sizeStr string) error {
	size, err := converter.ParseFileSizeStr(stripPerSecond(sizeStr))
	if err != nil {
		return err
	}
	c.MaxUploadRate = size
	return nil
}

// PrintTable 输出表格
func (c *PanConfig) PrintTable() {
	tb := cmdtable.NewTable(os.Stdout)
	tb.SetHeader([]string{"名称", "值", "建议值", "描述"})
	tb.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	tb.SetColumnAlignment([]int{tablewriter.ALIGN_DEFAULT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT})
	tb.AppendBulk([][]string{
		[]string{"cache_size", converter.ConvertFileSize(int64(c.CacheSize), 2), "1KB ~ 256KB", "下载缓存, 如果硬盘占用高或下载速度慢, 请尝试调大此值"},
		[]string{"max_download_parallel", strconv.Itoa(c.MaxDownloadParallel), "1 ~ 20", "最大下载并发量，即同时下载文件最大数量"},
		[]string{"max_upload_parallel", strconv.Itoa(c.MaxUploadParallel), "1 ~ 20", "最大上传并发量，即同时上传文件最大数量"},
		[]string{"max_download_rate", showMaxRate(c.MaxDownloadRate), "", "限制最大下载速度, 0代表不限制"},
		[]string{"max_upload_rate", showMaxRate(c.MaxUploadRate), "", "限制最大上传速度, 0代表不限制"},
		[]string{"savedir", c.SaveDir, "", "下载文件的储存目录"},
		[]string{"proxy", c.Proxy, "", "设置代理, 支持 http/socks5 代理，例如：http://127.0.0.1:8888"},
		[]string{"local_addrs", c.LocalAddrs, "", "设置本地网卡地址, 多个地址用逗号隔开"},
		[]string{"ip_type", c.PreferIPType, "ipv4, ipv6", "设置IP类型，优先IPv4或IPv6"},
		[]string{"dns_server", c.DNSServer, "114.114.114.114, 8.8.8.8", "设置DNS服务器地址，为空则使用系统DNS"},
	})
	tb.Render()
}
