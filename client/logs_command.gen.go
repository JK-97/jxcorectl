package client

// switchConfFile , 用于 解析 edit
func switchConfFile(operation string, filename string) string {
	switch filename {
	case "dnsmasqHost":
		operation += " /etc/dnsmasq.host"
	case "dnsmasqConf":
		operation += " /etc/dnsmasq.conf"
	case "dnsmasqResolv":
		operation += " /etc/dnsmasq.d/resolv.conf"
	case "initFile":
		operation += " /edge/init"
	default:
		operation += " " + filename
	}
	return operation
}
