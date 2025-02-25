package filter

import (
	"context"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/gchux/pcap-cli/pkg/pcap"
	"github.com/wissance/stringFormatter"
)

type (
	L4ProtoFilterProvider struct {
		*pcap.PcapFilter
	}
)

const (
	l4_PROTO_DEFAULT_FILTER string = "tcp or udp"
	l4_PROTO_TCP_FILTER     string = "tcp"
	l4_PROTO_UDP_FILTER     string = "udp"
)

func (p *L4ProtoFilterProvider) Get(ctx context.Context) (*string, bool) {
	if *p.Raw == "" ||
		*p.Raw == "23" || // tcp(6) + udp(17)
		*p.Raw == "0x17" || // tcp(0x06) + udp(0x11)
		strings.EqualFold(*p.Raw, "ALL") ||
		strings.EqualFold(*p.Raw, "ANY") ||
		strings.EqualFold(*p.Raw, l4_PROTO_DEFAULT_FILTER) {
		filter := string(l4_PROTO_DEFAULT_FILTER)
		return &filter, true
	}

	protos := strings.Split(*p.Raw, ",")
	if len(protos) == 0 || (len(protos) == 1 && protos[0] == "") {
		filter := string(l4_PROTO_DEFAULT_FILTER)
		return &filter, true
	}

	l4Protos := mapset.NewThreadUnsafeSet[string]()

	for _, proto := range protos {
		switch proto {
		case "tcp", "6", "0x06":
			l4Protos.Add(string(l4_PROTO_TCP_FILTER))
		case "udp", "17", "0x11":
			l4Protos.Add(string(l4_PROTO_UDP_FILTER))
		}
	}

	if l4Protos.IsEmpty() {
		filter := string(l4_PROTO_DEFAULT_FILTER)
		return &filter, true
	}

	filter := stringFormatter.Format("{0}", strings.Join(l4Protos.ToSlice(), " or "))
	return &filter, true
}

func (p *L4ProtoFilterProvider) String() string {
	if filter, ok := p.Get(context.Background()); ok {
		return stringFormatter.Format("L4Proto[{0}] => ({1})", *p.Raw, *filter)
	}
	return "L4Proto[nil]"
}

func (p *L4ProtoFilterProvider) Apply(
	ctx context.Context,
	srcFilter *string,
	mode pcap.PcapFilterMode,
) *string {
	return applyFilter(ctx, srcFilter, p, mode)
}

func newL4ProtoFilterProvider(filter *pcap.PcapFilter) pcap.PcapFilterProvider {
	provider := &L4ProtoFilterProvider{
		PcapFilter: filter,
	}
	return provider
}
