package network

import (
	"crypto/rand"
	"net"
	"net/netip"
	"strings"
)

func GenerateMac() string {
	buf := make([]byte, 6)
	var mac net.HardwareAddr
	_, err := rand.Read(buf)
	if err != nil {
	}
	buf[0] = (buf[0] | 0x02) & 0xfe
	mac = append(mac, buf[0], buf[1], buf[2], buf[3], buf[4], buf[5])
	return strings.ToUpper(mac.String())
}

func CidrToIps(cidr string) ([]string, error) {
	prefix, err := netip.ParsePrefix(cidr)
	if err != nil {
		panic(err)
	}
	var ips []string
	for addr := prefix.Addr(); prefix.Contains(addr); addr = addr.Next() {
		ips = append(ips, addr.String())
	}
	if len(ips) < 2 {
		return ips, nil
	}
	return ips[1 : len(ips)-1], nil
}

func GetIpFromDevice(device string) (ip string) {
	devices, err := net.Interfaces()
	if err != nil {
		return ip
	}
	for _, d := range devices {
		if d.Name == device {
			addrs, err := d.Addrs()
			if err != nil {
				return ip
			}
			for _, addr := range addrs {
				ip = addr.String()
				ip = strings.TrimRight(ip, "/24")
				return ip
			}
		}
	}
	return ip
}
