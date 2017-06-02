package netinfo

import (
	"net"
	"testing"
)

func TestGetNames(t *testing.T) {
	ip := "207.171.7.1"
	expect := "gw.develooper.com."

	names, err := GetNames(ip)
	if err != nil {
		t.Errorf("GetNames(%q): %s", ip, err)
	}

	if names[0] != expect {
		t.Errorf("Got name %q, expected %q", names[0], expect)
	}
}

func TestGetASN(t *testing.T) {
	ip := "207.171.3.1"
	expect := 7012

	asn, err := GetASN(ip)
	if err != nil {
		t.Errorf("GetASN(%q): %s", ip, err)
	}

	if asn != expect {
		t.Errorf("Got ASN %d, expected %d", asn, expect)
	}
}

func TestReverseIP(t *testing.T) {

	tests := []struct {
		ip     net.IP
		expect string
	}{
		{net.IP{192, 168, 0, 2}, "2.0.168.192"},
		{net.ParseIP("2607:f238:2::ff:4"), "4.0.0.0.f.f.0.0.0.0.0.0.0.0.0.0.0.0.0.0.2.0.0.0.8.3.2.f.7.0.6.2"},
	}

	for _, test := range tests {
		r := reverseIP(test.ip)
		if r != test.expect {
			t.Errorf("Got for %q\nreverse IP %q\nexpected   %q",
				test.ip.String(), r, test.expect)
		}
	}
}
