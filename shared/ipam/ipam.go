package ipam

import (
	"fmt"
	"net/netip"
)

var (
	ErrNoAvailableSubnet = fmt.Errorf("no available subnet")
)

func AllocateSubnet(
	parent netip.Prefix,
	childBits int,
	used []netip.Prefix,
) (netip.Prefix, error) {

	parent = parent.Masked()

	if !parent.Addr().Is4() {
		return netip.Prefix{}, fmt.Errorf("only IPv4 is supported")
	}

	if childBits <= parent.Bits() || childBits > 32 {
		return netip.Prefix{}, fmt.Errorf(
			"invalid child prefix /%d for parent %s",
			childBits,
			parent,
		)
	}

	usedSet := make(map[netip.Prefix]struct{}, len(used))

	for _, prefix := range used {
		usedSet[prefix.Masked()] = struct{}{}
	}

	addr := parent.Addr()

	addr = nextSubnet(addr, childBits)

	for parent.Contains(addr) {
		candidate := netip.PrefixFrom(addr, childBits).Masked()

		if _, exists := usedSet[candidate]; !exists {
			return candidate, nil
		}

		addr = nextSubnet(addr, childBits)
	}

	return netip.Prefix{}, ErrNoAvailableSubnet
}

func nextSubnet(
	addr netip.Addr,
	prefixBits int,
) netip.Addr {

	size := uint32(1) << (32 - prefixBits)

	raw := addr.As4()

	value :=
		uint32(raw[0])<<24 |
			uint32(raw[1])<<16 |
			uint32(raw[2])<<8 |
			uint32(raw[3])

	value += size

	return netip.AddrFrom4([4]byte{
		byte(value >> 24),
		byte(value >> 16),
		byte(value >> 8),
		byte(value),
	})
}
