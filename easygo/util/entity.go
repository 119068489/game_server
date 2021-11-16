package util

import (
	"fmt"
	"strconv"
)

// BuildEntityID 构建一个entityId
func BuildEntityID(eType, model, subID int) int {
	v, _ := strconv.Atoi(fmt.Sprintf("%d%02d%05d", eType, model, subID))
	return v
}

// ParseEntityID 解析entityID
func ParseEntityID(entityID int) (eType, model, subID int) {
	s := strconv.Itoa(entityID)
	if len(s) != 8 {
		return
	}
	buf := []byte(s)
	eType, _ = strconv.Atoi(string(buf[:1]))
	model, _ = strconv.Atoi(string(buf[1:3]))
	subID, _ = strconv.Atoi(string(buf[3:]))
	return
}

// RevEntityType 反转entityID的type
func RevEntityType(entityID int) (revEntityID int) {
	var eType, mode, subID = ParseEntityID(entityID)
	if eType == 0 {
		return
	}
	if eType == 1 {
		eType = 2
	} else {
		eType = 1
	}
	revEntityID = BuildEntityID(eType, mode, subID)
	return
}
