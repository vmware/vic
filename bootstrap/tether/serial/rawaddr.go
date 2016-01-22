package serial

type RawAddr struct {
	Net  string
	Addr string
}

func (addr RawAddr) Network() string {
	return addr.Net
}

func (addr RawAddr) String() string {
	return addr.Network() + "://" + addr.Addr
}

func NewRawAddr(net string, addr string) *RawAddr {
	return &RawAddr{
		Net:  net,
		Addr: addr,
	}
}
