package idhelper

import (
	"errors"
	"sync"
	"time"
)

// Snowflake 结构体
type Snowflake struct {
	mu            sync.Mutex // 互斥锁，保证并发安全
	startTime     int64      // 起始时间（毫秒）
	machineID     int64      // 机器ID
	sequence      int64      // 序列号
	lastTimestamp int64      // 上一次生成ID的时间戳
}

// 常量
const (
	machineIDBits  = 10                           // 机器ID占用的位数
	sequenceBits   = 12                           // 序列号占用的位数
	machineIDShift = sequenceBits                 // 机器ID左移位数
	timestampShift = machineIDBits + sequenceBits // 时间戳左移位数
	maxMachineID   = -1 ^ (-1 << machineIDBits)   // 最大机器ID
	maxSequence    = -1 ^ (-1 << sequenceBits)    // 最大序列号
)

// NewSnowflake 初始化雪花算法
func NewSnowflake(machineID int64) (*Snowflake, error) {
	if machineID < 0 || machineID > maxMachineID {
		return nil, errors.New("machine ID out of range")
	}
	return &Snowflake{
		startTime:     time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC).UnixNano() / 1e6,
		machineID:     machineID,
		sequence:      0,
		lastTimestamp: -1,
	}, nil
}

// NextID 生成下一个ID
func (s *Snowflake) NextID() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 获取当前时间戳
	currentTimestamp := time.Now().UnixNano() / 1e6

	// 如果当前时间小于上一次生成ID的时间，说明时钟回拨
	if currentTimestamp < s.lastTimestamp {
		panic("clock moved backwards")
	}

	// 如果是同一毫秒内生成的ID
	if currentTimestamp == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & maxSequence
		if s.sequence == 0 { // 如果序列号超出范围，等待下一毫秒
			for currentTimestamp <= s.lastTimestamp {
				currentTimestamp = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		s.sequence = 0
	}

	// 更新上一次生成ID的时间戳
	s.lastTimestamp = currentTimestamp

	// 生成ID
	id := (currentTimestamp-s.startTime)<<timestampShift |
		(s.machineID << machineIDShift) |
		s.sequence
	return id
}

//func main() {
//	// 初始化雪花算法，机器ID为1
//	snowflake, err := NewSnowflake(1)
//	if err != nil {
//		panic(err)
//	}
//
//	// 生成10个ID
//	for i := 0; i < 10; i++ {
//		id := snowflake.NextID()
//		fmt.Println(id)
//	}
//}
