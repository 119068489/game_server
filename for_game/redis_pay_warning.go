/**
充值提现警告
*/
package for_game

import (
	"game_server/easygo"
	"strconv"
)

// 支付警告的 key
const (
	WARNING_RECHARGE_COUNT  = "warning_recharge_count"  // hset warning_recharge_count pid : count
	WARNING_RECHARGE_AMOUNT = "warning_recharge_amount" // hset warning_recharge_amount pid : count

	WARNING_WITHDRAW_COUNT  = "warning_withdraw_count"  // hset warning_withdraw_count pid : count
	WARNING_WITHDRAW_AMOUNT = "warning_withdraw_amount" // hset warning_withdraw_amount pid : count

)

// 获取充值次数
func GetRechargeCountFromRedis(pid PLAYER_ID) int {
	b, err := easygo.RedisMgr.GetC().HExists(WARNING_RECHARGE_COUNT, easygo.AnytoA(pid))
	easygo.PanicError(err)
	if !b {
		return 0
	}
	value, err1 := easygo.RedisMgr.GetC().HGet(WARNING_RECHARGE_COUNT, easygo.AnytoA(pid))
	easygo.PanicError(err1)
	count, _ := strconv.Atoi(string(value))
	return count
}

// 设置充值次数
func SetRechargeCountToRedis(count, expire int64) int64 {
	// 先判断是否存在,如果不存,就要设置过期时间,如果存在,不用设置过期时间
	b, err := easygo.RedisMgr.GetC().Exist(WARNING_RECHARGE_COUNT)
	easygo.PanicError(err)
	value := easygo.RedisMgr.GetC().IncrBy(WARNING_RECHARGE_COUNT, count)
	if !b { // 不存在的时候,设置过期时间.
		err = easygo.RedisMgr.GetC().Expire(WARNING_RECHARGE_COUNT, expire) // 设置过期时间.
	}
	easygo.PanicError(err)
	return value
}

// 获取充值总金额 单位:分
func GetRechargeAmountFromRedis(pid PLAYER_ID) int {
	b, err := easygo.RedisMgr.GetC().HExists(WARNING_RECHARGE_AMOUNT, easygo.AnytoA(pid))
	easygo.PanicError(err)
	if !b {
		return 0
	}
	value, err1 := easygo.RedisMgr.GetC().HGet(WARNING_RECHARGE_AMOUNT, easygo.AnytoA(pid))
	easygo.PanicError(err1)
	count, _ := strconv.Atoi(string(value))
	return count
}

// 设置充值总金额 单位:分
func SetRechargeAmountToRedis(count, expire int64) int64 {
	// 先判断是否存在,如果不存,就要设置过期时间,如果存在,不用设置过期时间
	b, err := easygo.RedisMgr.GetC().Exist(WARNING_RECHARGE_AMOUNT)
	easygo.PanicError(err)
	value := easygo.RedisMgr.GetC().IncrBy(WARNING_RECHARGE_AMOUNT, count)
	if !b { // 不存在的时候,设置过期时间.
		err = easygo.RedisMgr.GetC().Expire(WARNING_RECHARGE_AMOUNT, expire) // 设置过期时间.
	}
	easygo.PanicError(err)
	return value
}

// 获取提现次数
func GetWithdrawCountFromRedis(pid PLAYER_ID) int {
	b, err := easygo.RedisMgr.GetC().HExists(WARNING_WITHDRAW_COUNT, easygo.AnytoA(pid))
	easygo.PanicError(err)
	if !b {
		return 0
	}
	value, err1 := easygo.RedisMgr.GetC().HGet(WARNING_WITHDRAW_COUNT, easygo.AnytoA(pid))
	easygo.PanicError(err1)
	count, _ := strconv.Atoi(string(value))
	return count
}

// 设置提现次数
func SetWithdrawCountToRedis(count, expire int64) int64 {
	// 先判断是否存在,如果不存,就要设置过期时间,如果存在,不用设置过期时间
	b, err := easygo.RedisMgr.GetC().Exist(WARNING_WITHDRAW_COUNT)
	easygo.PanicError(err)
	value := easygo.RedisMgr.GetC().IncrBy(WARNING_WITHDRAW_COUNT, count)
	if !b { // 不存在的时候,设置过期时间.
		err = easygo.RedisMgr.GetC().Expire(WARNING_WITHDRAW_COUNT, expire) // 设置过期时间.
	}
	easygo.PanicError(err)
	return value
}

// 获取提现总金额 单位:分
func GetWithdrawAmountFromRedis(pid PLAYER_ID) int {
	b, err := easygo.RedisMgr.GetC().HExists(WARNING_WITHDRAW_AMOUNT, easygo.AnytoA(pid))
	easygo.PanicError(err)
	if !b {
		return 0
	}
	value, err1 := easygo.RedisMgr.GetC().HGet(WARNING_WITHDRAW_AMOUNT, easygo.AnytoA(pid))
	easygo.PanicError(err1)
	count, _ := strconv.Atoi(string(value))
	return count
}

// 设置提现总金额 单位:分
func SetWithdrawAmountToRedis(amount, expire int64) int64 {
	// 先判断是否存在,如果不存,就要设置过期时间,如果存在,不用设置过期时间
	b, err := easygo.RedisMgr.GetC().Exist(WARNING_WITHDRAW_AMOUNT)
	easygo.PanicError(err)
	value := easygo.RedisMgr.GetC().IncrBy(WARNING_WITHDRAW_AMOUNT, amount)
	if !b { // 不存在的时候,设置过期时间.
		err = easygo.RedisMgr.GetC().Expire(WARNING_WITHDRAW_AMOUNT, expire) // 设置过期时间.
	}
	easygo.PanicError(err)
	return value
}
