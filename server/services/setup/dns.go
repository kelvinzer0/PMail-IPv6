package setup

import (
	"fmt"
	"github.com/Jinnrry/pmail/config"
	"net"
	"strings"

	"github.com/Jinnrry/pmail/i18n"
	"github.com/Jinnrry/pmail/services/auth"
	"github.com/Jinnrry/pmail/utils/context"
	"github.com/Jinnrry/pmail/utils/errors"
	"github.com/Jinnrry/pmail/utils/ip"
)

type DNSItem struct {
	Type  string `json:"type"`
	Host  string `json:"host"`
	Value string `json:"value"`
	TTL   int    `json:"ttl"`
	Tips  string `json:"tips"`
}

func GetDNSSettings(ctx *context.Context) (map[string][]*DNSItem, error) {
	configData, err := config.ReadConfig()
	if err != nil {
		return nil, errors.Wrap(err)
	}

	ret := make(map[string][]*DNSItem)

	for _, domain := range configData.Domains {
		var dnsItems []*DNSItem
		currentIP := ip.GetIp()

		if net.ParseIP(currentIP).To4() != nil { // It's an IPv4 address
			dnsItems = append(dnsItems, []*DNSItem{
				{Type: "A", Host: strings.ReplaceAll(configData.WebDomain, "."+configData.Domain, ""), Value: currentIP, TTL: 3600, Tips: i18n.GetText(ctx.Lang, "ip_taps")},
				{Type: "A", Host: "smtp", Value: currentIP, TTL: 3600, Tips: i18n.GetText(ctx.Lang, "ip_taps")},
				{Type: "A", Host: "imap", Value: currentIP, TTL: 3600, Tips: i18n.GetText(ctx.Lang, "ip_taps")},
				{Type: "A", Host: "pop", Value: currentIP, TTL: 3600, Tips: i18n.GetText(ctx.Lang, "ip_taps")},
				{Type: "A", Host: "@", Value: currentIP, TTL: 3600, Tips: i18n.GetText(ctx.Lang, "ip_taps")},
			}...)
		} else if net.ParseIP(currentIP).To16() != nil { // It's an IPv6 address
			dnsItems = append(dnsItems, []*DNSItem{
				{Type: "AAAA", Host: strings.ReplaceAll(configData.WebDomain, "."+configData.Domain, ""), Value: currentIP, TTL: 3600, Tips: i18n.GetText(ctx.Lang, "ip_taps")},
				{Type: "AAAA", Host: "smtp", Value: currentIP, TTL: 3600, Tips: i18n.GetText(ctx.Lang, "ip_taps")},
				{Type: "AAAA", Host: "imap", Value: currentIP, TTL: 3600, Tips: i18n.GetText(ctx.Lang, "ip_taps")},
				{Type: "AAAA", Host: "pop", Value: currentIP, TTL: 3600, Tips: i18n.GetText(ctx.Lang, "ip_taps")},
				{Type: "AAAA", Host: "@", Value: currentIP, TTL: 3600, Tips: i18n.GetText(ctx.Lang, "ip_taps")},
			}...)
		}

		dnsItems = append(dnsItems, []*DNSItem{
			{Type: "MX", Host: "@", Value: fmt.Sprintf("smtp.%s", domain), TTL: 3600},
			{Type: "TXT", Host: "@", Value: "v=spf1 a mx ~all", TTL: 3600},
			{Type: "TXT", Host: "default._domainkey", Value: auth.DkimGen(), TTL: 3600},
		}...)
		ret[domain] = dnsItems
	}

	return ret, nil
}
