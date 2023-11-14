package ipcalc

import (
	"errors"
	"fmt"
	"math"
	"net"
	"racent.com/pkg/logger"
)

// ipcalc命令的全称是：Calculate IP information for a host（计算主机的IP信息）

type IPCalc struct {
	NetMask       string
	NetMaskPrefix int
	NetWork       string
	Broadcast     string
	HostSum       uint
	MinIP         string
	MaxIP         string
}

// ParseCIDR 解析无类域间路由，IP地址/网络ID的位数 192.168.23.35/21
func ParseCIDR(s string) (*IPCalc, error) {
	_, ipv4net, err := net.ParseCIDR(s)
	if err != nil {
		return nil, err
	}
	info := &IPCalc{}
	info.NetMaskPrefix, _ = ipv4net.Mask.Size()
	info.NetWork = ipv4net.IP.String()
	info.setHostSum()
	info.setNetMask()
	info.setMinIP()
	info.setMaxIP()
	info.setBroadcast()
	return info, nil
}

func IPString2Long(ip string) (uint, error) {
	b := net.ParseIP(ip).To4()
	if b == nil {
		return 0, errors.New("invalid ipv4 format")
	}

	return uint(b[3]) | uint(b[2])<<8 | uint(b[1])<<16 | uint(b[0])<<24, nil
}

func Long2IPString(i uint) (string, error) {
	if i > math.MaxUint32 {
		return "", errors.New("beyond the scope of ipv4")
	}

	ip := make(net.IP, net.IPv4len)
	ip[0] = byte(i >> 24)
	ip[1] = byte(i >> 16)
	ip[2] = byte(i >> 8)
	ip[3] = byte(i)

	return ip.String(), nil
}

func (i *IPCalc) setHostSum() {
	switch i.NetMaskPrefix {
	case 32:
		i.HostSum = 1
	case 31:
		i.HostSum = 2
	default:
		i.HostSum = 1<<(32-i.NetMaskPrefix) - 2
	}
}

func (i *IPCalc) setNetMask() {
	// ^uint32(0)二进制为32个比特1，通过向左位移，得到CIDR掩码的二进制
	cidrMask := ^uint32(0) << uint(32-i.NetMaskPrefix)

	logger.DebugString("ipcalc", "setNetMask", fmt.Sprintf("%d", cidrMask))

	//计算CIDR掩码的四个片段，将想要得到的片段移动到内存最低8位后，将其强转为8位整型，从而得到
	cidrMaskSeg1 := uint8(cidrMask >> 24)
	cidrMaskSeg2 := uint8(cidrMask >> 16)
	cidrMaskSeg3 := uint8(cidrMask >> 8)
	cidrMaskSeg4 := uint8(cidrMask & uint32(255))

	i.NetMask = fmt.Sprintf("%d.%d.%d.%d", cidrMaskSeg1, cidrMaskSeg2, cidrMaskSeg3, cidrMaskSeg4)
}

func (i *IPCalc) setMinIP() {
	if i.NetMaskPrefix == 31 || i.NetMaskPrefix == 32 {
		i.MinIP = i.NetWork
	} else {
		tmp, _ := IPString2Long(i.NetWork)
		i.MinIP, _ = Long2IPString(tmp + 1)
	}
}

func (i *IPCalc) setMaxIP() {
	switch {
	case i.NetMaskPrefix == 31:
		tmp, _ := IPString2Long(i.NetWork)
		i.MaxIP, _ = Long2IPString(tmp + 1)
	case i.NetMaskPrefix == 32:
		i.MaxIP = i.NetWork
	default:
		tmp, _ := IPString2Long(i.NetWork)
		i.MaxIP, _ = Long2IPString(tmp + i.HostSum)
	}
}

func (i *IPCalc) setBroadcast() {
	if i.NetMaskPrefix < 31 {
		tmp, _ := IPString2Long(i.NetWork)
		i.Broadcast, _ = Long2IPString(tmp + i.HostSum + 1)
	}
}
