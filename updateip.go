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
	Debugf("got public ip %s", ip)
	api, err := cloudflare.NewWithAPIToken(conf.Token)
	for _, z := range conf.Zones {
		// first we need to get the Zone ID
		id, err := api.ZoneIDByName(z.Name)
		if err != nil {
			return err
		}
		// then we need to fetch DNS A records for that zone
		records, err := api.DNSRecords(ctx, id, cloudflare.DNSRecord{})
		for _, record := range records {
			if record.Type == "A" {
				for _, hostname := range z.Hosts {
					if record.Name == hostname {
						if record.Content == ip {
							Debugf("not updating %s because ip is already %s", hostname, ip)
						} else {
							log.Printf("updating %s to %s", hostname, ip)
							record.Content = ip
							err := api.UpdateDNSRecord(ctx, id, record.ID, record)
							if err != nil {
								return err
							}
						}
					}
				}
			}
		}
	}
	return nil
}
