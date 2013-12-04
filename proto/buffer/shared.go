package buffer

func BothClosed(in *Inbound, out *Outbound) (closed bool) {
	in.L.Lock()
	defer in.L.Unlock()

	out.L.Lock()
	defer out.L.Unlock()

	closed = (in.err != nil && out.err != nil)
	return
}
