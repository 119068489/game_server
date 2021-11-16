package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"game_server/easygo"
	"game_server/easygo/base"
	"game_server/for_game"
	"game_server/hall"
	"game_server/pb/h5_wish"
	"game_server/pb/server_server"
	"game_server/wish"
	"io/ioutil"
	"net/http"
	url2 "net/url"
	"runtime"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
)

const (
	ApiUrl = "http://127.0.0.1:15511/api"
)

var addr = flag.String("addr", "127.0.0.1:15511/api", "http service address")

type Socket struct {
	C   *websocket.Conn
	Err error
}

var socket *Socket

func init() {
	socket = &Socket{}
	initializer := hall.NewInitializer()
	defer func() { // 若是异常了,确保异步日志有成功写盘
		logger := initializer.GetBeeLogger()
		if logger != nil {
			logger.Flush()
		}
	}()
	dict := easygo.KWAT{
		"logName":  "hall",
		"yamlPath": "config_hall.yaml",
	}
	initializer.Execute(dict)
	hall.Initialize()
	wish.Initialize()
	////启动etcd
	//hall.PClient3KVMgr.StartClintTV3()
	//defer hall.PClient3KVMgr.Close() //关闭etcd
	////把已启动的服务增加到内存管理
	//for_game.InitExistServer(hall.PClient3KVMgr, hall.PServerInfoMgr, hall.PServerInfo)
	//hall.PWebApiForServer = hall.NewWebHttpForServer(hall.PServerInfo.GetServerApiPort())

}
func TestWishLogin() {
	reqMsg := &h5_wish.LoginReq{
		//Account:  easygo.NewString("lm70804e6a"),
		Channel:  easygo.NewInt32(1001),
		NickName: easygo.NewString("祖宁"),
		HeadUrl:  easygo.NewString(""),
		PlayerId: easygo.NewInt64(1887436059),
		Token:    easygo.NewString("harMJaBxDzzrxnhI1887436059"),
	}
	common := &base.Common{
		Version: easygo.NewString("1.0.1"),
		UserId:  easygo.NewInt64(1887436059),
		Token:   easygo.NewString("harMJaBxDzzrxnhI1887436059"),
	}
	resp, err := testSendToServer(ApiUrl, "RpcLogin", common, reqMsg)
	if err != nil {
		err = easygo.NewFailMsg(byteByRuneString(err.GetReason()))
		logs.Error("err", err)
		return
	}

	data := resp.(*h5_wish.LoginResp)
	marshal, _ := json.Marshal(data)
	fmt.Println("--------->", string(marshal))
}

func TestRpcWish() {
	reqMsg := &h5_wish.WishReq{
		BoxId:     easygo.NewInt64(190),
		ProductId: easygo.NewInt64(1473),
		OpType:    easygo.NewInt32(1),
	}
	p := for_game.GetRedisWishPlayer(19240012)
	pb := for_game.GetRedisPlayerBase(p.GetPlayerId())

	common := &base.Common{
		Version: easygo.NewString("1.0.1"),
		UserId:  easygo.NewInt64(p.GetId()),      //1887440774  10085;   1887440279   10086
		Token:   easygo.NewString(pb.GetToken()), // HJhC3e4RChZAjsMM10086  10086;  Pc5zsiTV3vVWVXpE10085   10085
	}
	resp, err := testSendToServer(ApiUrl, "RpcWish", common, reqMsg)
	if err != nil {
		err = easygo.NewFailMsg(byteByRuneString(err.GetReason()))
		logs.Error("err", err)
		return
	}

	data := resp.(*h5_wish.WishResp)
	logs.Info("GetResult:", data.GetResult())
}

func TestRpcRecycleGoods() {
	/*	reqMsgA := &h5_wish.MyWishReq{
			Type:     easygo.NewInt32(0),
			Page:     easygo.NewInt32(1),
			PageSize: easygo.NewInt32(100),
		}
		items, _ := wish.MyWishService(18805014, reqMsgA)
		ids := make([]int64, 0)
		for _, v := range items {
			ids = append(ids, v.GetPlayerWishItemId())
		}
		logs.Info("ids: ", ids)
	*/
	st := for_game.GetMillSecond()

	reqMsg := &h5_wish.WishBoxReq{
		IdList: []int64{738},
		//BankCardId: easygo.NewString("6228270081238239374"),
	}
	common := &base.Common{
		Version: easygo.NewString("1.0.1"),
		UserId:  easygo.NewInt64(18885065),            //1887440774  10085;   1887440279   10086
		Token:   easygo.NewString("RLfPVHs2eEtalXd1"), // HJhC3e4RChZAjsMM10086  10086;  Pc5zsiTV3vVWVXpE10085   10085
	}
	resp, err := testSendToServer(ApiUrl, "RpcRecycleGoods", common, reqMsg)
	if err != nil {
		err = easygo.NewFailMsg(byteByRuneString(err.GetReason()))
		logs.Error("err", err)
		return
	}
	ed := for_game.GetMillSecond()
	fmt.Println("耗时----------->", ed-st)
	data := resp.(*h5_wish.DefaultResp)
	logs.Info("GetResult:", data.GetResult())
}

// 发起挑战
func TestRpcDoDare() easygo.IMessage {
	st := for_game.GetMillSecond()
	reqMsg := &h5_wish.DoDareReq{
		DareType:  easygo.NewInt32(2),
		WishBoxId: easygo.NewInt64(190),
	}
	p := for_game.GetRedisWishPlayer(19240012)
	//p := for_game.GetRedisWishPlayer(19193011)
	pb := for_game.GetRedisPlayerBase(p.GetPlayerId())
	common := &base.Common{
		Version: easygo.NewString("1.0.1"),
		UserId:  easygo.NewInt64(p.GetId()),      //1887440774  10085;   1887440279   10086
		Token:   easygo.NewString(pb.GetToken()), // HJhC3e4RChZAjsMM10086  10086;  Pc5zsiTV3vVWVXpE10085   10085
	}
	resp, err := testSendToServer(ApiUrl, "RpcDoDare", common, reqMsg)
	if err != nil {
		err = easygo.NewFailMsg(byteByRuneString(err.GetReason()))
		logs.Error("err", err.GetReason())
		return nil
	}

	ed := for_game.GetMillSecond()
	fmt.Println("耗时----------->", ed-st)
	return resp

}

//奖池列表
func TesRpcActPoolList() easygo.IMessage {

	reqMsg := &base.Empty{}
	common := &base.Common{
		Version: easygo.NewString("1.0.1"),
		UserId:  easygo.NewInt64(18801002),                    //1887440774  10085;   1887440279   10086
		Token:   easygo.NewString("ZqTDrB3JxUNo8b8P10085000"), // HJhC3e4RChZAjsMM10086  10086;  Pc5zsiTV3vVWVXpE10085   10085
	}
	resp, err := testSendToServer(ApiUrl, "RpcActPoolList", common, reqMsg)
	if err != nil {
		err = easygo.NewFailMsg(byteByRuneString(err.GetReason()))
		logs.Error("err", err.GetReason())
		return nil
	}

	data := resp.(*h5_wish.ActPoolResp)
	marshal, _ := json.Marshal(data)
	fmt.Println("------>", string(marshal))
	return resp

}

// 奖池规则查询返回
func TesRpcActPoolRule() easygo.IMessage {
	reqMsg := &h5_wish.ActPoolRuleReq{
		ActPoolId: easygo.NewInt64(1),
		Type:      easygo.NewInt32(1),
	}
	p := for_game.GetRedisWishPlayer(18801002)
	pb := for_game.GetRedisPlayerBase(p.GetPlayerId())
	common := &base.Common{
		Version: easygo.NewString("1.0.1"),
		UserId:  easygo.NewInt64(p.GetId()),      //1887440774  10085;   1887440279   10086
		Token:   easygo.NewString(pb.GetToken()), // HJhC3e4RChZAjsMM10086  10086;  Pc5zsiTV3vVWVXpE10085   10085
	}
	resp, err := testSendToServer(ApiUrl, "RpcActPoolRule", common, reqMsg)
	if err != nil {
		err = easygo.NewFailMsg(byteByRuneString(err.GetReason()))
		logs.Error("err", err.GetReason())
		return nil
	}

	data := resp.(*h5_wish.ActPoolRuleResp)
	marshal, _ := json.Marshal(data)
	fmt.Println("------>", string(marshal))
	return resp

}

// 奖池规则查询返回
func TesRpcSumDay() easygo.IMessage {
	reqMsg := &h5_wish.SumReq{
		ActPoolId: easygo.NewInt64(1),
		Type:      easygo.NewInt32(2),
	}
	b1 := for_game.GetRedisWishPlayer(19267015)
	pb := for_game.GetRedisPlayerBase(b1.GetPlayerId())
	common := &base.Common{
		Version: easygo.NewString("1.0.1"),
		UserId:  easygo.NewInt64(19267015),       //1887440774  10085;   1887440279   10086
		Token:   easygo.NewString(pb.GetToken()), // HJhC3e4RChZAjsMM10086  10086;  Pc5zsiTV3vVWVXpE10085   10085
	}
	resp, err := testSendToServer(ApiUrl, "RpcSumNum", common, reqMsg)
	if err != nil {
		err = easygo.NewFailMsg(byteByRuneString(err.GetReason()))
		logs.Error("err", err.GetReason())
		return nil
	}

	data := resp.(*h5_wish.SumNumResp)
	marshal, _ := json.Marshal(data)
	fmt.Println("------>", string(marshal))
	return resp

}

// 奖池规则查询返回
func TesGetPlayerDiamond() easygo.IMessage {
	reqMsg := &server_server.PlayerSI{
		PlayerId: easygo.NewInt64(1887440774),
	}
	common := &base.Common{
		Version: easygo.NewString("1.0.1"),
		UserId:  easygo.NewInt64(18801002),                      //1887440774  10085;   1887440279   10086
		Token:   easygo.NewString("ZqTDrB3JxUNo8b8P1008500000"), // HJhC3e4RChZAjsMM10086  10086;  Pc5zsiTV3vVWVXpE10085   10085
	}
	resp, err := testSendToServer(ApiUrl, "GetPlayerDiamond", common, reqMsg)
	if err != nil {
		err = easygo.NewFailMsg(byteByRuneString(err.GetReason()))
		logs.Error("err", err.GetReason())
		return nil
	}

	data := resp.(*h5_wish.SumNumResp)
	marshal, _ := json.Marshal(data)
	fmt.Println("------>", string(marshal))
	return resp

}
func TesRpcSumMoney() easygo.IMessage {
	reqMsg := &h5_wish.SumMoneyReq{
		DataType: easygo.NewInt64(2),
		Page:     easygo.NewInt32(1),
		PageSize: easygo.NewInt32(10),
	}
	common := &base.Common{
		Version: easygo.NewString("1.0.1"),
		UserId:  easygo.NewInt64(19250014),                       //1887440774  10085;   1887440279   10086
		Token:   easygo.NewString("PeHtbHu12uelT79C15623232323"), // HJhC3e4RChZAjsMM10086  10086;  Pc5zsiTV3vVWVXpE10085   10085
	}
	resp, err := testSendToServer(ApiUrl, "RpcSumMoney", common, reqMsg)
	if err != nil {
		err = easygo.NewFailMsg(byteByRuneString(err.GetReason()))
		logs.Error("err", err.GetReason())
		return nil
	}

	data := resp.(*h5_wish.SumMoneyResp)
	marshal, _ := json.Marshal(data)
	fmt.Println("------>", string(marshal))
	return resp

}
func TesRpcGive() easygo.IMessage {
	reqMsg := &h5_wish.GiveReq{
		PrizeLogId: easygo.NewInt64(129),
	}
	wishPlsyer := for_game.GetRedisWishPlayer(19250014)
	base1 := for_game.GetRedisPlayerBase(wishPlsyer.GetPlayerId())
	common := &base.Common{
		Version: easygo.NewString("1.0.1"),
		UserId:  easygo.NewInt64(wishPlsyer.GetId()), //1887440774  10085;   1887440279   10086
		Token:   easygo.NewString(base1.GetToken()),  // HJhC3e4RChZAjsMM10086  10086;  Pc5zsiTV3vVWVXpE10085   10085
		//UserId: easygo.NewInt64(18877004),                       //1887440774  10085;   1887440279   10086
		//Token:  easygo.NewString("KSHEMm3Kzurhe8yO18934895027"), // HJhC3e4RChZAjsMM10086  10086;  Pc5zsiTV3vVWVXpE10085   10085
	}
	resp, err := testSendToServer(ApiUrl, "RpcGive", common, reqMsg)
	if err != nil {
		err = easygo.NewFailMsg(byteByRuneString(err.GetReason()))
		logs.Error("err", err.GetReason())
		return nil
	}

	data := resp.(*h5_wish.GiveResp)
	marshal, _ := json.Marshal(data)
	fmt.Println("------>", string(marshal))
	return resp

}

func TestRpcGetUserIdBankCards() {
	st := for_game.GetMillSecond()
	reqMsg := &base.Empty{}
	common := &base.Common{
		Version: easygo.NewString("1.0.1"),
		UserId:  easygo.NewInt64(18805011), //1887440774  10085;   1887440279   10086
		Token:   easygo.NewString("rSL46RYmbmY4jKkC"),
	}
	resp, err := testSendToServer(ApiUrl, "RpcGetUserIdBankCards", common, reqMsg)
	if err != nil {
		err = easygo.NewFailMsg(byteByRuneString(err.GetReason()))
		logs.Error("err", err)
		return
	}
	ed := for_game.GetMillSecond()
	fmt.Println("耗时----------->", ed-st)
	data := resp.(*h5_wish.BankCardResp)
	marshal, _ := json.Marshal(data)
	logs.Info("------------>", string(marshal))
}

// 测试数据
func TestDareData() {
	// 发起挑战
	reqMsg := &h5_wish.DoDareReq{
		DareType:  easygo.NewInt32(1),
		WishBoxId: easygo.NewInt64(44),
	}
	common := &base.Common{
		Version: easygo.NewString("1.0.1"),
		UserId:  easygo.NewInt64(1887570292),               //1887440774  10085;   1887440279   10086
		Token:   easygo.NewString("Pc5zsiTV3vVWVXpE10085"), // HJhC3e4RChZAjsMM10086  10086;  Pc5zsiTV3vVWVXpE10085   10085
	}
	_, err := testSendToServer(ApiUrl, "RpcDoDare", common, reqMsg)
	if err != nil {
		logs.Error("err", err.GetReason()) // 先进行许愿
		return
	}
}

func TestAddGold() {
	reqMsg := &h5_wish.AddGoldReq{
		UserId:     easygo.NewInt64(1887440774),
		Coin:       easygo.NewInt64(10),
		SourceType: easygo.NewInt32(for_game.COIN_TYPE_WISH_ADD),
	}

	common := &base.Common{
		Version: easygo.NewString("1.0.1"),
		UserId:  easygo.NewInt64(1887440774),               //1887440774  10085;   1887440279   10086
		Token:   easygo.NewString("Pc5zsiTV3vVWVXpE10085"), // HJhC3e4RChZAjsMM10086  10086;  Pc5zsiTV3vVWVXpE10085   10085
	}
	_, err := testSendToServer(ApiUrl, "RpcAddGold", common, reqMsg)
	if err != nil {
		logs.Error("err", err.GetReason()) // 先进行许愿
		return
	}
}
func TestRpcMyDare() {
	reqMsg := &h5_wish.MyDareReq{
		Page:     nil,
		PageSize: nil,
	}

	common := &base.Common{
		Version: easygo.NewString("1.0.1"),
		UserId:  easygo.NewInt64(1887440774),               //1887440774  10085;   1887440279   10086
		Token:   easygo.NewString("Pc5zsiTV3vVWVXpE10085"), // HJhC3e4RChZAjsMM10086  10086;  Pc5zsiTV3vVWVXpE10085   10085
	}
	_, err := testSendToServer(ApiUrl, "RpcMyDare", common, reqMsg)
	if err != nil {
		logs.Error("err", err.GetReason()) // 先进行许愿
		return
	}
}
func TestRpcDareList() {
	reqMsg := &h5_wish.DareReq{
		BoxId: easygo.NewInt64(18),
	}

	common := &base.Common{
		Version: easygo.NewString("1.0.1"),
		UserId:  easygo.NewInt64(1887440774),               //1887440774  10085;   1887440279   10086
		Token:   easygo.NewString("Pc5zsiTV3vVWVXpE10085"), // HJhC3e4RChZAjsMM10086  10086;  Pc5zsiTV3vVWVXpE10085   10085
	}
	_, err := testSendToServer(ApiUrl, "RpcDareList", common, reqMsg)
	if err != nil {
		logs.Error("err", err.GetReason()) // 先进行许愿
		return
	}
}
func TestRpcGetRandProduct() {
	reqMsg := &h5_wish.DareReq{}

	common := &base.Common{
		Version: easygo.NewString("1.0.1"),
		UserId:  easygo.NewInt64(1887440774),               //1887440774  10085;   1887440279   10086
		Token:   easygo.NewString("Pc5zsiTV3vVWVXpE10085"), // HJhC3e4RChZAjsMM10086  10086;  Pc5zsiTV3vVWVXpE10085   10085
	}
	resp, err := testSendToServer(ApiUrl, "RpcGetRandProduct", common, reqMsg)
	if err != nil {
		logs.Error("err", err.GetReason()) // 先进行许愿
		return
	}
	logs.Info(resp)
}
func TestRpcCoinToDiamond() {
	reqMsg := &h5_wish.CoinToDiamondReq{
		Coin: nil,
		Id:   easygo.NewInt64(2),
	}

	common := &base.Common{
		Version: easygo.NewString("1.0.1"),
		UserId:  easygo.NewInt64(18801002),                      //1887440774  10085;   1887440279   10086
		Token:   easygo.NewString("Tw3zikSD2kUP6mjIlm70800f86"), // HJhC3e4RChZAjsMM10086  10086;  Pc5zsiTV3vVWVXpE10085   10085
	}
	resp, err := testSendToServer(ApiUrl, "RpcCoinToDiamond", common, reqMsg)
	if err != nil {
		logs.Error("err", err.GetReason()) // 先进行许愿
		return
	}
	logs.Info(resp)
}
func TestRpcDiamondRechargeList() {
	reqMsg := &base.Empty{}

	common := &base.Common{
		Version: easygo.NewString("1.0.1"),
		UserId:  easygo.NewInt64(18801002),                      //1887440774  10085;   1887440279   10086
		Token:   easygo.NewString("Tw3zikSD2kUP6mjIlm70800f86"), // HJhC3e4RChZAjsMM10086  10086;  Pc5zsiTV3vVWVXpE10085   10085
	}
	resp, err := testSendToServer(ApiUrl, "RpcDiamondRechargeList", common, reqMsg)
	if err != nil {
		logs.Error("err", err.GetReason()) // 先进行许愿
		return
	}
	data := resp.(*h5_wish.DiamondRechargeResp)
	marshal, _ := json.Marshal(data)

	logs.Info("------------>", string(marshal))
}
func TestRpcDiamondChangeLogList() {
	reqMsg := &h5_wish.DiamondChangeLogReq{
		Page:     easygo.NewInt32(1),
		PageSize: easygo.NewInt32(10),
		Type:     easygo.NewInt32(2),
	}
	common := &base.Common{
		Version: easygo.NewString("1.0.1"),
		UserId:  easygo.NewInt64(18801002),                      //1887440774  10085;   1887440279   10086
		Token:   easygo.NewString("Tw3zikSD2kUP6mjIlm70800f86"), // HJhC3e4RChZAjsMM10086  10086;  Pc5zsiTV3vVWVXpE10085   10085
	}
	resp, err := testSendToServer(ApiUrl, "RpcDiamondChangeLogList", common, reqMsg)
	if err != nil {
		logs.Error("err", err.GetReason()) // 先进行许愿
		return
	}
	data := resp.(*h5_wish.DiamondChangeLogResp)
	marshal, _ := json.Marshal(data)

	logs.Info("------------>", string(marshal))
}

func TestRpcGetConfig() {
	reqMsg := &base.Empty{}
	common := &base.Common{
		Version: easygo.NewString("1.0.1"),
		UserId:  easygo.NewInt64(18805014),                       //1887440774  10085;   1887440279   10086
		Token:   easygo.NewString("ILPrK2fc8aTxXysy15099973008"), // HJhC3e4RChZAjsMM10086  10086;  Pc5zsiTV3vVWVXpE10085   10085
	}
	resp, err := testSendToServer(ApiUrl, "RpcGetConfig", common, reqMsg)
	if err != nil {
		logs.Error("err", err.GetReason()) // 先进行许愿
		return
	}
	data := resp.(*h5_wish.ConfigResp)
	marshal, _ := json.Marshal(data)

	logs.Info("------------>", string(marshal))
}

// 配合策划做大数据测试
func TestDoDareByCH() {
	count := 0
	dc := 0
	//rateCount := 0
	for {
		if count == 100 {
			break
		}
		resp := TestRpcDoDare()
		data := resp.(*h5_wish.DoDareResp)
		logs.Info("GetProductId:", data.GetProductId())
		logs.Info("GetImage:", data.GetImage())
		logs.Info("GetIsLucky:", data.GetIsLucky())
		logs.Info("GetIsOnce:", data.GetIsOnce())
		logs.Info("GetProductName:", data.GetProductName())
		logs.Info("GetProductType:", data.GetProductType())
		if data.GetProductId() >= 1053 && data.GetProductId() <= 1060 {
			dc++
		}
		count++
		//id := data.GetProductId()
		//if id == 68 || id == 69 || id == 70 || id == 71 || id == 72 {
		//	rateCount++
		//}
		time.Sleep(100 * time.Millisecond)
		// 判断水池状态,如果不是普通,则停止
		//status := wish.GetPoolStatus(4)
		//if status != 3 {
		//	fmt.Println("抽奖次数为----->", count)
		//	fmt.Println("rateCount----->", rateCount)
		//	return
		//}
	}
	//TestRpcDiamondChangeLogList()
	logs.Info("大奖的个数为-------->", dc)
	//fmt.Println("rateCount----->", rateCount)
}

func TestBoxInfo() {
	reqMsg := &h5_wish.BoxReq{
		BoxId: easygo.NewInt64(190),
		Type:  easygo.NewInt32(2),
	}
	p := for_game.GetRedisWishPlayer(18940006)
	pb := for_game.GetRedisPlayerBase(p.GetPlayerId())
	common := &base.Common{
		Version: easygo.NewString("1.0.1"),
		UserId:  easygo.NewInt64(p.GetId()),      //1887440774  10085;   1887440279   10086
		Token:   easygo.NewString(pb.GetToken()), // HJhC3e4RChZAjsMM10086  10086;  Pc5zsiTV3vVWVXpE10085   10085
	}

	resp, err := testSendToServer(ApiUrl, "RpcBoxInfo", common, reqMsg)
	if err != nil {
		logs.Error("err", err.GetReason()) // 先进行许愿
		return
	}

	data := resp.(*h5_wish.BoxResp)
	//logs.Info("boxId:", data.GetBoxId())
	//logs.Info("Protector:", data.GetProtector())
	//logs.Info("ProtectorHeadUrl:", data.GetProtectorHeadUrl())
	//logs.Info("ProtectorId:", data.GetProtectorId())
	//logs.Info("ProtectorTime:", data.GetProtectorTime())
	//logs.Info("ProductList:", data.GetProductList())
	marshal, _ := json.Marshal(data)
	fmt.Println("---------->", string(marshal))
}
func TestRpcDareRecord() {
	reqMsg := &h5_wish.DareRecordReq{
		BoxId:    easygo.NewInt64(19),
		Type:     easygo.NewInt32(2),
		Page:     easygo.NewInt32(1),
		PageSize: easygo.NewInt32(15),
	}
	common := &base.Common{
		Version: easygo.NewString("1.0.1"),
		UserId:  easygo.NewInt64(18801002),                      //1887440774  10085;   1887440279   10086
		Token:   easygo.NewString("ZqTDrB3JxUNo8b8P1008500000"), // HJhC3e4RChZAjsMM10086  10086;  Pc5zsiTV3vVWVXpE10085   10085
	}

	resp, err := testSendToServer(ApiUrl, "RpcDareRecord", common, reqMsg)
	if err != nil {
		logs.Error("err", err.GetReason()) // 先进行许愿
		return
	}

	data := resp.(*h5_wish.DareRecordResp)
	//logs.Info("boxId:", data.GetBoxId())
	//logs.Info("Protector:", data.GetProtector())
	//logs.Info("ProtectorHeadUrl:", data.GetProtectorHeadUrl())
	//logs.Info("ProtectorId:", data.GetProtectorId())
	//logs.Info("ProtectorTime:", data.GetProtectorTime())
	//logs.Info("ProductList:", data.GetProductList())
	marshal, _ := json.Marshal(data)
	fmt.Println("---------->", string(marshal))
}

// 批量抽奖
func TestRpcBatchDare() {
	count := 10
	st := for_game.GetMillSecond()
	reqMsg := &h5_wish.BatchDareReq{
		BoxId: easygo.NewInt64(83),
		Uid:   easygo.NewInt64(18801002),
		Count: easygo.NewInt32(count),
	}
	common := &base.Common{
		Version: easygo.NewString("1.0.1"),
		UserId:  easygo.NewInt64(18801002),                      //1887440774  10085;   1887440279   10086
		Token:   easygo.NewString("ZqTDrB3JxUNo8b8P1008500000"), // HJhC3e4RChZAjsMM10086  10086;  Pc5zsiTV3vVWVXpE10085   10085
	}

	_, err := testSendToServer(ApiUrl, "RpcBatchDare", common, reqMsg)
	if err != nil {
		logs.Error("err", err.GetReason())
		return
	}
	ed := for_game.GetMillSecond()
	fmt.Printf("%d 次总耗时为: %d 毫秒", count, ed-st)
}
func TestRpcProductShow() {
	reqMsg := &base.Empty{}
	common := &base.Common{
		Version: easygo.NewString("1.0.1"),
		UserId:  easygo.NewInt64(18805020),                       //1887440774  10085;   1887440279   10086
		Token:   easygo.NewString("LYE6YULqFNeNMlum16620173354"), // HJhC3e4RChZAjsMM10086  10086;  Pc5zsiTV3vVWVXpE10085   10085
	}

	resp, err := testSendToServer(ApiUrl, "RpcProductShow", common, reqMsg)
	if err != nil {
		logs.Error("err", err.GetReason())
		return
	}
	data := resp.(*h5_wish.ProductShowResp)
	marshal, _ := json.Marshal(data)
	fmt.Println("----------->", string(marshal))
}
func TestRpcSearchBox() {
	reqMsg := &h5_wish.SearchBoxReq{
		Complex:       nil,
		Condition:     nil,
		ProductStatus: nil,
		MinPrice:      nil,
		MaxPrice:      nil,
		//WishBrandId:    []int64{3, 4},
		//WishItemTypeId: []int64{3},
		Label:    easygo.NewInt32(1),
		Page:     easygo.NewInt32(1),
		PageSize: easygo.NewInt32(20),
	}
	p := for_game.GetRedisWishPlayer(18940006)
	pb := for_game.GetRedisPlayerBase(p.GetPlayerId())
	common := &base.Common{
		Version: easygo.NewString("1.0.1"),
		UserId:  easygo.NewInt64(p.GetId()),      //1887440774  10085;   1887440279   10086
		Token:   easygo.NewString(pb.GetToken()), // HJhC3e4RChZAjsMM10086  10086;  Pc5zsiTV3vVWVXpE10085   10085
	}

	resp, err := testSendToServer(ApiUrl, "RpcSearchBox", common, reqMsg)
	if err != nil {
		logs.Error("err", err.GetReason())
		return
	}
	data := resp.(*h5_wish.SearchBoxResp)
	marshal, _ := json.Marshal(data)
	fmt.Println("----------->", string(marshal))
}

func TestRpcGetRechargeAct() {
	reqMsg := &h5_wish.TypeReq{
		Type: easygo.NewInt32(0),
	}
	common := &base.Common{
		Version: easygo.NewString("1.0.1"),
		UserId:  easygo.NewInt64(1887571827),                     //1887440774  10085;   1887440279   10086
		Token:   easygo.NewString("KSHEMm3Kzurhe8yO18934895027"), // HJhC3e4RChZAjsMM10086  10086;  Pc5zsiTV3vVWVXpE10085   10085
	}

	resp, err := testSendToServer(ApiUrl, "RpcGetRechargeAct", common, reqMsg)
	if err != nil {
		logs.Error("err", err.GetReason())
		return
	}
	data := resp.(*h5_wish.RechargeActResp)
	marshal, _ := json.Marshal(data)
	fmt.Println("----------->", string(marshal))
}
func TestRpcActOpenStatus() {
	reqMsg := &base.Empty{}
	common := &base.Common{
		Version: easygo.NewString("1.0.1"),
		UserId:  easygo.NewInt64(19250014),                       //1887440774  10085;   1887440279   10086
		Token:   easygo.NewString("PeHtbHu12uelT79C15623232323"), // HJhC3e4RChZAjsMM10086  10086;  Pc5zsiTV3vVWVXpE10085   10085
	}

	resp, err := testSendToServer(ApiUrl, "RpcActOpenStatus", common, reqMsg)
	if err != nil {
		logs.Error("err", err.GetReason())
		return
	}
	data := resp.(*h5_wish.ActOpenStatusResp)
	marshal, _ := json.Marshal(data)
	fmt.Println("----------->", string(marshal))
}
func TestRpcRechargeActStatus() {
	reqMsg := &base.Empty{}
	common := &base.Common{
		Version: easygo.NewString("1.0.1"),
		UserId:  easygo.NewInt64(19250014),                       //1887440774  10085;   1887440279   10086
		Token:   easygo.NewString("PeHtbHu12uelT79C15623232323"), // HJhC3e4RChZAjsMM10086  10086;  Pc5zsiTV3vVWVXpE10085   10085
	}

	resp, err := testSendToServer(ApiUrl, "RpcRechargeActStatus", common, reqMsg)
	if err != nil {
		logs.Error("err", err.GetReason())
		return
	}
	data := resp.(*h5_wish.ActOpenStatusResp)
	marshal, _ := json.Marshal(data)
	fmt.Println("----------->", string(marshal))
}

func TestRpcBackstageSetGuardian() {
	reqMsg := &h5_wish.BackstageSetGuardianReq{
		Account:  easygo.NewString(10010000024),
		Channel:  easygo.NewInt32(1001),
		NickName: nil,
		HeadUrl:  nil,
		PlayerId: nil,
		Token:    nil,
		BoxId:    easygo.NewInt64(18),
		OpType:   easygo.NewInt32(1),
	}

	common := &base.Common{
		Version: easygo.NewString("1.0.1"),
		UserId:  easygo.NewInt64(19250014),                       //1887440774  10085;   1887440279   10086
		Token:   easygo.NewString("PeHtbHu12uelT79C15623232323"), // HJhC3e4RChZAjsMM10086  10086;  Pc5zsiTV3vVWVXpE10085   10085
	}

	resp, err := testSendToServer(ApiUrl, "RpcBackstageSetGuardian", common, reqMsg)
	if err != nil {
		logs.Error("err", err.GetReason())
		return
	}
	data := resp.(*h5_wish.BackstageSetGuardianResp)
	marshal, _ := json.Marshal(data)
	fmt.Println("----------->", string(marshal))
}
func main() {
	/*	reqMsg := &h5_wish.BoxReq{
			BoxId: easygo.NewInt64(3),
			Type:  easygo.NewInt32(1),
		}
		resp, err := testSendToServer(ApiUrl, "RpcBoxInfo", reqMsg)
		if err != nil {
			err = easygo.NewFailMsg(byteByRuneString(err.GetReason()))
			logs.Error("err", err)
			return
		}

		data := resp.(*h5_wish.BoxResp)
		logs.Info("boxId:", data.GetBoxId())
		logs.Info("Protector:", data.GetProtector())
		logs.Info("ProtectorHeadUrl:", data.GetProtectorHeadUrl())
		logs.Info("ProtectorId:", data.GetProtectorId())
		logs.Info("ProtectorTime:", data.GetProtectorTime())
		logs.Info("ProductList:", data.GetProductList())*/
	//TestWishLogin()
	//time.Sleep(300 * time.Millisecond)
	//TestRpcGetUseInfo()
	//TestRpcDiamondRechargeList()
	//TestBoxInfo()
	//TestRpcBoxProduct()
	//TestDoDareByCH()
	//TestRpcCoinToDiamond()
	//TestRpcMyAllWish()
	//TestRpcBatchDare()
	//TestRpcProductShow()
	TestRpcDoDare()
	//TestRpcWish()
	//TestRpcSearchBox()
	//TestRpcDoDare()
	//TestRpcDiamondChangeLogList()
	//easygo.Spawn(TestRpcGetRechargeAct)
	//TestRpcDareRecord()
	//time.Sleep(time.Second)
	//TestRpcGetConfig()
	//TesRpcActPoolList()
	//TesRpcActPoolRule()
	//TesRpcSumDay()
	//TesRpcSumMoney()
	//TesGetPlayerDiamond()
	//TesRpcGive()
	//TestRpcRecycleGoods()
	//TestRpcRechargeActStatus()
	//TestRpcBackstageSetGuardian()
	time.Sleep(time.Hour * 10)
}

func TestRpcAreaPostage() {
	reqMsg := &base.Empty{}
	common := &base.Common{
		Version: easygo.NewString("1.0.1"),
		UserId:  easygo.NewInt64(1887436008),               //1887440774  10085;   1887440279   10086
		Token:   easygo.NewString("Pc5zsiTV3vVWVXpE10085"), // HJhC3e4RChZAjsMM10086  10086;  Pc5zsiTV3vVWVXpE10085   10085
	}
	resp, err := testSendToServer(ApiUrl, "RpcAreaPostage", common, reqMsg)
	if err != nil {
		err = easygo.NewFailMsg(byteByRuneString(err.GetReason()))
		logs.Error("err", err)
		return
	}
	logs.Info(resp)
}

func TestRpcRecycleRatio() {
	reqMsg := &base.Empty{}
	common := &base.Common{
		Version: easygo.NewString("1.0.1"),
		UserId:  easygo.NewInt64(1887436008),               //1887440774  10085;   1887440279   10086
		Token:   easygo.NewString("Pc5zsiTV3vVWVXpE10085"), // HJhC3e4RChZAjsMM10086  10086;  Pc5zsiTV3vVWVXpE10085   10085
	}
	resp, err := testSendToServer(ApiUrl, "RpcRecycleRatio", common, reqMsg)
	if err != nil {
		err = easygo.NewFailMsg(byteByRuneString(err.GetReason()))
		logs.Error("err", err)
		return
	}
	logs.Info(resp)
}

func TestRpcBoxProduct() {
	reqMsg := &h5_wish.DareReq{BoxId: easygo.NewInt64(18)}
	common := &base.Common{
		Version: easygo.NewString("1.0.1"),
		UserId:  easygo.NewInt64(4004),                     //1887440774  10085;   1887440279   10086
		Token:   easygo.NewString("Pc5zsiTV3vVWVXpE10085"), // HJhC3e4RChZAjsMM10086  10086;  Pc5zsiTV3vVWVXpE10085   10085
	}
	resp, err := testSendToServer(ApiUrl, "RpcBoxProduct", common, reqMsg)
	if err != nil {
		err = easygo.NewFailMsg(byteByRuneString(err.GetReason()))
		logs.Error("err", err)
		return
	}
	data := resp.(*h5_wish.BoxProductResp)
	marshal, _ := json.Marshal(data)

	fmt.Println("TestRpcBoxProduct--------------->", string(marshal))
}
func TestRpcGetUseInfo() {
	reqMsg := &h5_wish.UserInfoReq{
		UserId: easygo.NewInt64(1887440774),
	}
	common := &base.Common{
		Version: easygo.NewString("1.0.1"),
		UserId:  easygo.NewInt64(4004),                     //1887440774  10085;   1887440279   10086
		Token:   easygo.NewString("Pc5zsiTV3vVWVXpE10085"), // HJhC3e4RChZAjsMM10086  10086;  Pc5zsiTV3vVWVXpE10085   10085
	}
	resp, err := testSendToServer(ApiUrl, "RpcGetUseInfo", common, reqMsg)
	if err != nil {
		err = easygo.NewFailMsg(byteByRuneString(err.GetReason()))
		logs.Error("err", err)
		return
	}
	data := resp.(*h5_wish.UserInfoResp)
	marshal, _ := json.Marshal(data)

	logs.Info("---------->", string(marshal))
}

func TestRpcExchangeBox() {
	reqMsg := &h5_wish.WishBoxReq{
		IdList:    []int64{13433, 13434},
		AddressId: easygo.NewInt64(1617175737),
	}
	common := &base.Common{
		Version: easygo.NewString("1.0.1"),
		UserId:  easygo.NewInt64(18805014), //1887440774  10085;   1887440279   10086
		Token:   easygo.NewString("ILPrK2fc8aTxXysy15099973008"),
	}
	resp, err := testSendToServer(ApiUrl, "RpcExchangeBox", common, reqMsg)
	if err != nil {
		err = easygo.NewFailMsg(byteByRuneString(err.GetReason()))
		logs.Error("err", err)
		return
	}
	data := resp.(*h5_wish.DefaultResp)
	if data.GetResult() == 1 {
		logs.Info("errMsg:", data.GetMsg())
	}
}

func TestRpcAddAddress() {
	reqMsg := &h5_wish.WishAddress{
		Name:      easygo.NewString("sdfs"),
		Phone:     easygo.NewString("21312431"),
		Detail:    easygo.NewString("saefasefasf"),
		IfDefault: easygo.NewBool(true),
		AddressId: easygo.NewInt64(123432532),
		Province:  easygo.NewString("guangdong"),
		City:      easygo.NewString("guangdong"),
		Area:      easygo.NewString("guangdong"),
	}

	common := &base.Common{
		Version: easygo.NewString("1.0.1"),
		UserId:  easygo.NewInt64(3011),                           //1887440774  10085;   1887440279   10086
		Token:   easygo.NewString("xpMkxXOQDiSOy4Lq18565036899"), // HJhC3e4RChZAjsMM10086  10086;  Pc5zsiTV3vVWVXpE10085   10085
	}
	resp, err := testSendToServer(ApiUrl, "RpcAddAddress", common, reqMsg)
	if err != nil {
		logs.Error("err", err.GetReason()) // 先进行许愿
		return
	}
	data := resp.(*h5_wish.DefaultResp)
	marshal, _ := json.Marshal(data)

	logs.Info("------------>", string(marshal))
}

func TestRpcMyAllWish() {
	reqMsg := &h5_wish.MyWishReq{
		Type:     easygo.NewInt32(3),
		Page:     easygo.NewInt32(1),
		PageSize: easygo.NewInt32(10),
	}

	common := &base.Common{
		Version: easygo.NewString("1.0.1"),
		UserId:  easygo.NewInt64(18805014),                       //1887440774  10085;   1887440279   10086
		Token:   easygo.NewString("ILPrK2fc8aTxXysy15099973008"), // HJhC3e4RChZAjsMM10086  10086;  Pc5zsiTV3vVWVXpE10085   10085
	}
	resp, err := testSendToServer(ApiUrl, "RpcMyAllWish", common, reqMsg)
	if err != nil {
		logs.Error("err", err.GetReason()) // 先进行许愿
		return
	}
	data := resp.(*h5_wish.ProductResp)
	marshal, _ := json.Marshal(data)

	logs.Info("------------>", string(marshal))
}

func TestRpcTryOne() {
	reqMsg := &base.Empty{}

	common := &base.Common{
		Version: easygo.NewString("1.0.1"),
		UserId:  easygo.NewInt64(18814029),                       //1887440774  10085;   1887440279   10086
		Token:   easygo.NewString("xzDfvIrS5z8vJFwa13333333333"), // HJhC3e4RChZAjsMM10086  10086;  Pc5zsiTV3vVWVXpE10085   10085
	}
	_, err := testSendToServer(ApiUrl, "RpcTryOnce", common, reqMsg)
	if err != nil {
		logs.Error("err", err.GetReason())
		return
	}

}

func TestRpcEditAddress() {
	reqMsg := &h5_wish.WishAddress{
		Name:      easygo.NewString("sdfs"),
		Phone:     easygo.NewString("21312431"),
		Detail:    easygo.NewString("saefasefasf"),
		IfDefault: easygo.NewBool(true),
		AddressId: easygo.NewInt64(5555555),
		Province:  easygo.NewString("guangdong"),
		City:      easygo.NewString("guangdong"),
		Area:      easygo.NewString("guangdong"),
	}

	common := &base.Common{
		Version: easygo.NewString("1.0.1"),
		UserId:  easygo.NewInt64(3011),                           //1887440774  10085;   1887440279   10086
		Token:   easygo.NewString("xpMkxXOQDiSOy4Lq18565036899"), // HJhC3e4RChZAjsMM10086  10086;  Pc5zsiTV3vVWVXpE10085   10085
	}
	resp, err := testSendToServer(ApiUrl, "RpcEditAddress", common, reqMsg)
	if err != nil {
		logs.Error("err", err.GetReason()) // 先进行许愿
		return
	}
	data := resp.(*h5_wish.DefaultResp)
	marshal, _ := json.Marshal(data)

	logs.Info("------------>", string(marshal))
}

//func main() {
//	m := &h5_wish.DoDareReq{
//		DareType:easygo.NewInt32(1),
//		WishBoxId:easygo.NewInt64(1),
//	}
//	resp, err := testSendToServer(ApiUrl, "RpcDoDare", m)
//	if err != nil {
//		err = easygo.NewFailMsg(byteByRuneString(err.GetReason()))
//		logs.Error("err",err)
//		return
//	}
//	data := resp.(*h5_wish.DoDareResp)
//	logs.Info("resp:", data)
//	time.Sleep(time.Hour * 10)
//}

//序列化socket请求数据
func MarshalSocketData(methodName string, reqData easygo.IMessage) ([]byte, error) {
	msg, err := reqData.Marshal()
	easygo.PanicError(err)
	common := &base.Common{
		Version: easygo.NewString("1.0.1"),
		UserId:  easygo.NewInt64(888888),
		Token:   easygo.NewString(token),
	}
	request := base.Request{
		RequestId:  easygo.NewUint64(1),
		MethodName: easygo.NewString(methodName),
		Serialized: msg,
		Timestamp:  easygo.NewInt64(time.Now().Unix()),
		Timeout:    easygo.NewUint32(0),
		Common:     common,
	}
	d1, err := request.Marshal()
	t := base.PacketType_TYPE_REQUEST
	packet := base.Packet{
		Type:       &t,
		Serialized: d1,
	}
	data, err := packet.Marshal()
	return data, err
}

//反序列化socket返回的数据
func UnMarshalSocketData(bs []byte) (easygo.IMessage, *base.Fail) {
	b := &base.Packet{}
	err := b.Unmarshal(bs)
	if err != nil {
		logs.Error("err:", err)
		return nil, easygo.NewFailMsg(err.Error())
	}
	//logs.Info("packet:", b)
	resp := &base.Response{}
	err = resp.Unmarshal(b.GetSerialized())
	if err != nil {
		logs.Error("err:", err)
		return nil, easygo.NewFailMsg("Response反序列化失败：")
	}
	//logs.Info("resp:", resp)
	msgName := resp.GetMsgName()
	rspMsg := easygo.NewMessage(msgName)
	err = rspMsg.Unmarshal(resp.GetSerialized())
	//logs.Info("rspMsg:", rspMsg)
	if err != nil {
		return nil, easygo.NewFailMsg(err.Error())
	}
	if resp.GetSubType() == base.ResponseType_TYPE_SUCCESS {
		return rspMsg, nil
	} else {
		return nil, rspMsg.(*base.Fail)
	}
}

var token = "token_1880081_opjJZGNpnXCygvrQ"

func testSendToServer(url string, methodName string, common *base.Common, reqData easygo.IMessage) (easygo.IMessage, *base.Fail) {
	msg, err := reqData.Marshal()
	easygo.PanicError(err)
	request := base.Request{
		RequestId:  easygo.NewUint64(1),
		MethodName: easygo.NewString(methodName),
		Serialized: msg,
		Timestamp:  easygo.NewInt64(time.Now().Unix()),
		Timeout:    easygo.NewUint32(0),
	}
	d1, err := request.Marshal()
	t := base.PacketType_TYPE_REQUEST
	packet := base.Packet{
		Type:       &t,
		Serialized: d1,
	}
	data, err := packet.Marshal()

	//发起api请求
	bs, err := doBytesPost(url, data, common)
	if err != nil {
		logs.Error("err:", err)
		return nil, easygo.NewFailMsg(err.Error())
	}

	str64, _ := url2.QueryUnescape(string(bs))
	respBs, err := base64.StdEncoding.DecodeString(str64)
	if err != nil {
		logs.Error("err:", err)
		return nil, easygo.NewFailMsg(err.Error())
	}

	b := &base.Packet{}
	err = b.Unmarshal([]byte(respBs))
	if err != nil {
		logs.Error("err:", err)
		return nil, easygo.NewFailMsg(err.Error())
	}
	//logs.Info("packet:", b)
	resp := &base.Response{}
	err = resp.Unmarshal(b.GetSerialized())
	if err != nil {
		logs.Error("err:", err)
		return nil, easygo.NewFailMsg("Response反序列化失败：")
	}
	//logs.Info("resp:", resp)
	msgName := resp.GetMsgName()
	rspMsg := easygo.NewMessage(msgName)
	err = rspMsg.Unmarshal(resp.GetSerialized())
	//logs.Info("rspMsg:", rspMsg)
	if err != nil {
		return nil, easygo.NewFailMsg(err.Error())
	}
	if resp.GetSubType() == base.ResponseType_TYPE_SUCCESS {
		return rspMsg, nil
	} else {
		return nil, rspMsg.(*base.Fail)
	}

}

//body提交二进制数据
func doBytesPost(url string, data []byte, comm ...*base.Common) ([]byte, error) {
	body := bytes.NewReader(data)
	request, err := http.NewRequest("POST", url, body)
	easygo.PanicError(err)
	request.Header.Set("Connection", "Keep-Alive")
	common := &base.Common{
		Version: easygo.NewString("1.0.1"),
		UserId:  easygo.NewInt64(888888),
		Token:   easygo.NewString(token),
	}
	common = append(comm, common)[0]
	com, err := common.Marshal()
	easygo.PanicError(err)

	request.Header.Set("Common", base64.StdEncoding.EncodeToString(com))
	var resp *http.Response
	//logs.Info("request:",request)
	resp, err = http.DefaultClient.Do(request)
	easygo.PanicError(err)
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	easygo.PanicError(err)
	return b, err
}

// 获取正在运行的函数名
func runFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}

//把byte类型的数据格式化为rune类型的字符串
func byteByRuneString(str string) string {
	var text []rune
	for _, v := range str {
		text = append(text, v)
	}
	return string(text)
}
