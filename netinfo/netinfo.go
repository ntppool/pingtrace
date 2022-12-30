package netinfo

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	// "github.com/cloudflare/golibs/lrucache"
)

// type asnData struct {
// 	asnum  int
// 	asname string
// }

// func getASNData()

func GetASName(asnum int) (string, error) {
	txts, err := net.LookupTXT(fmt.Sprintf("as%d.asn.cymru.com", asnum))

	if err != nil {
		return "", fmt.Errorf("could not lookup asn: %s", err)
	}

	if len(txts) > 0 {
		fmt.Printf("txts: %#v\n", txts)
		txt := strings.Join(txts, "|")
		fields := strings.Split(txt, "|")
		if len(fields) >= 5 {
			return strings.TrimSpace(fields[4]), nil
		}
	}

	// 7012 | US | arin |  | PHYBER - Phyber Communications, LLC.,US

	return "", nil
}
func GetASN(ipStr string) (int, error) {
	rev, err := reverseIPStr(ipStr)
	if err != nil {
		return 0, err
	}
	// fmt.Printf("reverse ip: %s\n", rev)

	txts, err := net.LookupTXT(rev + ".origin.asn.cymru.com")
	if err != nil {
		return 0, fmt.Errorf("could not lookup asn: %s", err)
	}

	if len(txts) > 0 {
		// fmt.Printf("txts: %#v\n", txts)
		idx := strings.Index(txts[0], " ")
		asnStr := txts[0][:idx]
		// fmt.Printf("asnStr: %s\n", asnStr)
		return strconv.Atoi(asnStr)
	}

	// 15169 | 216.239.32.0/19 | US | arin | 2000-11-22

	return 0, nil
}

func GetNames(ipStr string) ([]string, error) {
	addrs, err := net.LookupAddr(ipStr)
	if err != nil {
		return nil, err
	}
	for i := range addrs {
		addrs[i] = strings.TrimRight(addrs[i], ".")
	}
	return addrs, nil
}

func reverseIPStr(ipStr string) (string, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "", errors.New("could not parse IP")
	}
	return reverseIP(ip), nil
}

const hexDigit = "0123456789abcdef"

func reverseIP(ip net.IP) string {
	if v4 := ip.To4(); v4 != nil {
		// fmt.Printf("v4: %#v\n", v4)
		return strconv.Itoa(int(v4[3])) + "." + strconv.Itoa(int(v4[2])) + "." +
			strconv.Itoa(int(v4[1])) + "." + strconv.Itoa(int(v4[0]))
	} else {
		// assume v6
		// ensure zeros are present in string
		buf := make([]byte, 0, len(ip)*4)
		for i := len(ip) - 1; i >= 0; i-- {
			v := ip[i]
			buf = append(buf, hexDigit[v&0xF], byte('.'))
			buf = append(buf, hexDigit[v>>4])
			if i > 0 {
				buf = append(buf, byte('.'))
			}
		}
		return string(buf)
	}
}
