package main

import (
	"context"
	"log"

	"github.com/cloudflare/cloudflare-go"
)

func UpdateCloudflareIP(ctx context.Context, conf *Config) error {
	ip, err := GetPublicIPv4()
	if err != nil {
		return err
	}
	if conf.lastIP == ip {
		// if the ip address hasn't changed just return
		Debugf("IP address hasn't changed. Skipping Update.")
		return nil
	}
	Debugf("public ip changed to %s", ip)
	// store the change
	conf.lastIP = ip
	api, err := cloudflare.NewWithAPIToken(conf.Token)
	if err != nil {
		return err
	}
	for _, z := range conf.Zones {
		// first we need to get the Zone ID
		id, err := api.ZoneIDByName(z.Name)
		if err != nil {
			return err
		}
		Debugf("fetched zone info for %s", z.Name)
		// then we need to fetch DNS A records for that zone
		records, err := api.DNSRecords(ctx, id, cloudflare.DNSRecord{})
		if err != nil {
			return err
		}
		Debugf("fetched DNS records for zone %s", z.Name)
		for _, hostname := range z.Hosts {
			found := false
			for _, record := range records {
				if record.Type == "A" {
					if record.Name == hostname {
						found = true
						if record.Content == ip {
							Debugf("not updating host %s because ip is already %s", hostname, ip)
						} else {
							record.Content = ip
							err := api.UpdateDNSRecord(ctx, id, record.ID, record)
							if err != nil {
								log.Printf("unable to update ip for host %s: %v", hostname, err)
							}
							log.Printf("updated host %s to %s", hostname, ip)
						}
					}
				}
			}
			if !found {
				log.Printf("no DNS A record found for host %s", hostname)
			}
		}
	}
	return nil
}
