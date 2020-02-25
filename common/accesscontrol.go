package common

import (
	"strconv"
	"sync"
	"time"
)

// 分布式控制器
type AccessControl struct {
	//用来存放用户想要存放的信息,这边记录用户上次点击秒杀的时间
	sourcesArray map[int]time.Time
	sync.RWMutex
}

//服务器间隔时间，单位秒
var interval = 20

//创建全局变量
var accessControl = &AccessControl{sourcesArray: make(map[int]time.Time)}

//根据用户名获取上次访问的时间
func (m *AccessControl) GetNewRecord(uid int) time.Time {
	m.RWMutex.RLock()
	defer m.RWMutex.RUnlock()
	return m.sourcesArray[uid]
}

//设置访问时间
func (m *AccessControl) SetNewRecord(uid int) {
	m.RWMutex.Lock()
	m.sourcesArray[uid] = time.Now()
	m.RWMutex.Unlock()
}



//黑名单
type BlackList struct {
	listArray map[int]bool
	sync.RWMutex
}

var blackList = &BlackList{listArray: make(map[int]bool)}

//获取黑名单
func (m *BlackList) GetBlackListByID(uid int) bool {
	m.RLock()
	defer m.RUnlock()
	return m.listArray[uid]
}

//添加黑名单
func (m *BlackList) SetBlackListByID(uid int) bool {
	m.Lock()
	defer m.Unlock()
	m.listArray[uid] = true
	return true
}

//获取本机map，并且处理业务逻辑，返回的结果类型为bool类型
func (m *AccessControl) GetDataFromMap(uid string) (isOk bool) {
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return false
	}
	//添加黑名单
	if blackList.GetBlackListByID(uidInt) {
		//判断是否被添加到黑名单中
		return false
	}
	//获取记录
	dataRecord := m.GetNewRecord(uidInt)
	if !dataRecord.IsZero() {
		//业务判断，是否在指定间隔之后
		if dataRecord.Add(time.Duration(interval) * time.Second).After(time.Now()) {
			return false
		}
	}
	m.SetNewRecord(uidInt)
	return true
}

func GetAccessControl() *AccessControl {
	return accessControl
}
