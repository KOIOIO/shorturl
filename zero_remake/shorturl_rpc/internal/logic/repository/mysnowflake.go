package repository

import (
	"errors"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

// 自定义雪花算法结构体
// 41位时间戳(毫秒)+4位逻辑时钟+10位机器ID+8位序列号
// 支持时钟回拨时逻辑时钟补偿

type MySnowflake struct {
	startTime     int64      // 起始时间戳（毫秒）
	machineID     uint16     // 机器ID（0-1023）
	lastTimestamp int64      // 上次生成ID的时间戳
	logicClock    uint16     // 逻辑时钟（0-15）
	sequence      uint16     // 序列号（0-255）
	lock          sync.Mutex // 并发锁
}

const (
	timestampBits  = 41
	logicClockBits = 4
	machineIDBits  = 10
	sequenceBits   = 8

	maxLogicClock = -1 ^ (-1 << logicClockBits) // 15
	maxMachineID  = -1 ^ (-1 << machineIDBits)  // 1023
	maxSequence   = -1 ^ (-1 << sequenceBits)   // 255
)

var (
	myFlake     *MySnowflake
	myFlakeOnce sync.Once
)

func getMachineID() uint16 {
	idStr := os.Getenv("MACHINE_ID")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 0 || id > maxMachineID {
		return uint16(rand.Intn(maxMachineID + 1))
	}
	return uint16(id)
}

func NewMySnowflake() *MySnowflake {
	return &MySnowflake{
		startTime:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli(),
		machineID:  getMachineID(),
		logicClock: 0,
		sequence:   0,
	}
}

func GetMyFlake() *MySnowflake {
	myFlakeOnce.Do(func() {
		myFlake = NewMySnowflake()
	})
	return myFlake
}

func (sf *MySnowflake) NextID() (uint64, error) {
	sf.lock.Lock()
	defer sf.lock.Unlock()
	now := time.Now().UnixMilli()
	if now < sf.lastTimestamp {
		// 时钟回拨，逻辑时钟补偿
		sf.logicClock++
		if sf.logicClock > maxLogicClock {
			return 0, errors.New("时钟回拨过多，逻辑时钟溢出")
		}
		now = sf.lastTimestamp // 逻辑时钟补偿后，时间戳不变
		sf.sequence = 0
	} else {
		if now == sf.lastTimestamp {
			sf.sequence++
			if sf.sequence > maxSequence {
				// 当前毫秒内序列号用完，等待下一毫秒
				for now <= sf.lastTimestamp {
					now = time.Now().UnixMilli()
				}
				sf.sequence = 0
				sf.logicClock = 0
			}
		} else {
			sf.sequence = 0
			sf.logicClock = 0
		}
	}
	sf.lastTimestamp = now
	id := ((now - sf.startTime) & ((1 << timestampBits) - 1)) << (logicClockBits + machineIDBits + sequenceBits)
	id |= (int64(sf.logicClock) & ((1 << logicClockBits) - 1)) << (machineIDBits + sequenceBits)
	id |= (int64(sf.machineID) & ((1 << machineIDBits) - 1)) << sequenceBits
	id |= int64(sf.sequence) & ((1 << sequenceBits) - 1)
	return uint64(id), nil
}
