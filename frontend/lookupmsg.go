package frontend

import "github.com/lw000/gocommon/network/ws/packet"

// 检查消息code
func CheckMessageCode() MsgHooks {
	return func(pk *typacket.Packet) bool {
		return true
	}
}
