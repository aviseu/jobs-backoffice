package base

type ChannelStatus int

const (
	ChannelStatusInactive ChannelStatus = iota
	ChannelStatusActive
)

func (s ChannelStatus) String() string {
	return [...]string{"inactive", "active"}[s]
}
