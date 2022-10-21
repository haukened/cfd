# cfd
A lightweight Cloudflare DDNS Client, written in Go.

Initially i created this because of how frequently ddclient broke and wouldn't update cloudflare. I decided by using the official golang cloudflare package, i could minimize breaking with upstream changes.

All configuration should be in a single config file, which by default is /etc/cfd.yml
The program will then update the DNS records for all zone/host combinations in the file with the discovered public IP address of the machine running the service.

# TODO

- ~~systemd unit file~~
- deb packaging
- IPv6 support