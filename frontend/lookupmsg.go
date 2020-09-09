package frontend

import (
	"github.com/lw000/gocommon/network/ws/packet"
)

// 检查消息code
func CheckMessageCode() MsgHooksFunc {
	return func(pk *typacket.Packet) bool {
		if pk.CheckCode() != 123456 {
			return false
		}
		return true
	}
}
