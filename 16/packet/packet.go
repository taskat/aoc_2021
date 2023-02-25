package packet

import (
	"strings"
)

type Bits []bool

func (b Bits) String() string {
	s := make([]string, len(b))
	for i, value := range b {
		if value {
			s[i] = "1"
		} else {
			s[i] = "0"
		}
	}
	return strings.Join(s, "")
}

func binToInt(b Bits) int {
	num := 0
	for _, value := range b {
		num *= 2
		if value {
			num++
		}
	}
	return num
}

type TypeId int

const (
	SumType TypeId = iota
	ProductType
	MinimumType
	MaximumType
	ValueType
	GreaterThanType
	LessThanType
	EqualToType
)

type packet struct {
	version int
	typeId  TypeId
}

func (p *packet) create(bits Bits) {
	p.version = binToInt(bits[0:3])
	p.typeId = TypeId(binToInt(bits[3:6]))
}

type Packet interface {
	Create(Bits) int
	GetVersion() int
	GetTypeId() TypeId
	SumVersion() int
	Eval() int
}

type Value struct {
	packet
	value int
}

func (v *Value) Create(bits Bits) int {
	v.packet.create(bits[:6])
	numBits := make(Bits, 0)
	usedBits := 6
	for i := 6; i < len(bits); i += 5 {
		last := !bits[i]
		numBits = append(numBits, bits[i+1:i+5]...)
		usedBits += 5
		if last {
			break
		}
	}
	v.value = binToInt(numBits)
	return usedBits
}

func (v *Value) GetVersion() int {
	return v.packet.version
}

func (v *Value) GetTypeId() TypeId {
	return v.packet.typeId
}

func (v *Value) SumVersion() int {
	return v.version
}

func (v *Value) Eval() int {
	return v.value
}

type operator struct {
	packet
	subs       []Packet
	lengthType int
}

func (o *operator) Create(bits Bits) int {
	o.packet.create(bits[:6])
	o.lengthType = binToInt(bits[6:7])
	switch o.lengthType {
	case 0:
		length := binToInt(bits[7:22])
		start := 22
		end := start + length
		for start < end {
			packet, usedBits := Create(bits[start:end])
			o.subs = append(o.subs, packet)
			length  -= usedBits
			start += usedBits
		}
		return end
	case 1:
		numberOfPackets := binToInt(bits[7:18])
		o.subs = make([]Packet, numberOfPackets)
		start := 18
		usedBits := 0
		for i := 0; i < numberOfPackets; i++ {
			o.subs[i], usedBits = Create(bits[start:])
			start += usedBits
		}
		return start
	default:
		return 0
	}
}

func (o *operator) GetVersion() int {
	return o.packet.version
}

func (o *operator) GetTypeId() TypeId {
	return o.packet.typeId
}

func (o *operator) SumVersion() int {
	sum := o.version
	for _, sub := range o.subs {
		sum += sub.SumVersion()
	}
	return sum
}

func (o *operator) Eval() int {
	panic("Not implemented!")
}

type Sum struct {operator}

func (s Sum) Eval() int {
	sum := 0
	for _, sub := range s.subs {
		sum += sub.Eval()
	}
	return sum
}

type Product struct {operator}

func (p Product) Eval() int {
	product := 1
	for _, sub := range p.subs {
		product *= sub.Eval()
	}
	return product
}

type Minimum struct {operator}

func (m Minimum) Eval() int {
	min := m.subs[0].Eval()
	for i, sub := range m.subs {
		if i == 0 {
			continue
		}
		value := sub.Eval()
		if value < min {
			min = value
		}
	}
	return min
}

type Maximum struct {operator}

func (m Maximum) Eval() int {
	max := 0
	for _, sub := range m.subs {
		value := sub.Eval()
		if value > max {
			max = value
		}
	}
	return max
}

type GreaterThan struct {operator}

func (gt GreaterThan) Eval() int {
	first := gt.subs[0].Eval()
	second := gt.subs[1].Eval()
	if first > second {
		return 1
	}
	return 0
}

type LessThan struct {operator}

func (lt LessThan) Eval() int {
	first := lt.subs[0].Eval()
	second := lt.subs[1].Eval()
	if first < second {
		return 1
	}
	return 0
}

type EqualTo struct {operator}

func (et EqualTo) Eval() int {
	first := et.subs[0].Eval()
	second := et.subs[1].Eval()
	if first == second {
		return 1
	}
	return 0
}

func Create(bits Bits) (Packet, int) {
	var packet Packet
	switch TypeId(binToInt(bits[3:6])) {
	case SumType:
		packet = &Sum{}
	case ProductType:
		packet = &Product{}
	case MinimumType:
		packet = &Minimum{}
	case MaximumType:
		packet =&Maximum{}
	case ValueType:
		packet = &Value{}
	case GreaterThanType:
		packet = &GreaterThan{}
	case LessThanType:
		packet = &LessThan{}
	case EqualToType:
		packet = &EqualTo{}
	default:
		packet =&operator{}
	}
	usedbits := packet.Create(bits)
	return packet, usedbits
}
