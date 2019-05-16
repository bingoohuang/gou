package gou

import (
	"errors"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
)

/*
When using Nginx as a reverse proxy you may want to pass through the IP address of the remote user to your backend web server.
This must be done using the X-Forwarded-For header. You have a couple of options on how to set this information with Nginx.
You can either append the remote hosts IP address to any existing X-Forwarded-For values, or you can simply
set the X-Forwarded-For value, which clears out any previous IP’s that would have been on the request.

Edit the nginx configuration file, and add one of the follow lines in where appropriate.
To set the X-Forwarded-For to only contain the remote users IP:
proxy_set_header X-Forwarded-For $remote_addr;

To append the remote users IP to any existing X-Forwarded-For value:
proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
*/
func GetClientIp(req *http.Request) string {
	// "X-Forwarded-For"/ "x-forwarded-for"/"X-FORWARDED-FOR"  // capitalisation  doesn't matter
	xForwardedFor := req.Header.Get("X-FORWARDED-FOR")
	if xForwardedFor != "" {
		proxyIps := strings.Split(xForwardedFor, ",")
		return proxyIps[0]
	}

	ip, _, _ := net.SplitHostPort(req.RemoteAddr)
	return ip
}

func IsPrivateIP(ip string) (bool, error) {
	IP := net.ParseIP(ip)
	if IP == nil {
		return false, errors.New("invalid IP")
	}

	networks := []string{
		"0.0.0.0/8",
		"10.0.0.0/8",
		"100.64.0.0/10",
		"127.0.0.0/8",
		"169.254.0.0/16",
		"172.16.0.0/12",
		"192.0.0.0/24",
		"192.0.2.0/24",
		"192.88.99.0/24",
		"192.168.0.0/16",
		"198.18.0.0/15",
		"198.51.100.0/24",
		"203.0.113.0/24",
		"240.0.0.0/4",
		"255.255.255.255/32",
		"224.0.0.0/4",
	}

	for _, network := range networks {
		_, privateBitBlock, _ := net.ParseCIDR(network)
		if privateBitBlock.Contains(IP) {
			return true, nil
		}
	}

	return false, nil
}

func IsIP4Valid(ipv4 string) bool {
	re, _ := regexp.Compile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
	return re.MatchString(ipv4)
}

func GetLocalIps() []string {
	ips := make([]string, 0)
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		return ips
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}

	return ips
}

// IfaceAddr 表示一个IP地址和网卡名称的结构
type IfaceAddr struct {
	IP        string
	IfaceName string
}

// ListIP 列出本机所有IP
func ListIP() ([]IfaceAddr, error) {
	list, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	ret := make([]IfaceAddr, 0)
	for _, iface := range list {
		addrs, err := iface.Addrs()
		if err != nil {
			return ret, err
		}

		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if !ok || ipnet.IP.IsLoopback() || ipnet.IP.To4() == nil {
				continue
			}

			ret = append(ret, IfaceAddr{
				IP:        ipnet.IP.String(),
				IfaceName: iface.Name,
			})

		}
	}

	return ret, nil
}