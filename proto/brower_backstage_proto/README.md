# brower_server_proto

## brower_backstage.proto
系统管理后台的proto，系统后台的rpc在这里。

## brower_backstage_team.proto
群主后台的proto，群主后台的rpc在这里。

## common.proto
所有rpc公用的请求和返回数据结构，非私有的数据结构请在这个文件上补充

# backstage_api

## 后台提供的api接口

### /activity 活动api
1、集卡活动首页 ?t=1&id=1&pid=1887436002&sid=1887436003 (id:活动id 默认1,pid:用户自己的id,sid:邀请人id)
2、抽卡 ?t=2&id=1&pid=1887436002 (id:活动id 默认1,pid:用户id)
3、送卡记录 ?t=3&pid=1 (pid:用户id)
4、赠送卡 ?t=4&id=1887436002&pid=1887436003&aid=1 (aid:活动id,id:送卡用户,pid:获赠用户,cid:卡id)
5、发送验证码 ?t=5&phone=13736597264 (phone:手机号)
6、激活邀请 ?t=6&id=1&sid=1887436002&phone=13736597264&code=4352 (id:活动id,sid:邀请人,phone:手机号,code:验证码)
7、开奖 ?t=7&pid=1&id=1 (pid:用户id,id:活动id)
8、签到 ?t=8&id=1&pid=1 (id:活动id,pid:用户id)
9、分享活动 ?t=9&id=1&pid=1&sid=1887436003 (id:活动id,pid:用户id,sid:分享人id)
10、好友列表 ?t=10&pid=1 (pid:用户id)

### /DeviceCode 激活设备api 参数:code设备码,channle渠道号
激活设备api ?code=xxx-xx-xxxx&channle=1000058

### /registered 注册api 参数: t=1获取验证码,2注册 ，phone手机号,code验证码,channelno渠道号
1、获取验证码：?t=1&phone=1389898988
2、注册账号：?t=2&phone=1389898988&code=4588&channelno=1000058 

### /advidfa 广告idfa上报API 参数：t 上报平台类型
#### t=ks 快手上报接口 参数：code设备码MD5后，advid广告计划id，os系统类型，ip地址，scenesid广告场景，callback回调链接，ts时间毫秒
安卓：https://testapi.lemonchat.cn/advidfa?t=ks&code=__ANDROIDID2__&callback=__CALLBACK__&advid=__DID__&os=__OS__&ip=__IP__&scenesid=__CSITE__&ts=__TS__
苹果：https://testapi.lemonchat.cn/advidfa?t=ks&code=__IDFA2__&callback=__CALLBACK__&advid=__DID__&os=__OS__&ip=__IP__&scenesid=__CSITE__&ts=__TS__
#### t=qtt 趣头条上报接口  参数：andid安卓设备码MD5后，idfa苹果idfa原值，advid广告计划id，os系统类型，ip地址，callback回调链接，ts时间毫秒
https://testapi.lemonchat.cn/advidfa?t=qtt&andid=__ANDROIDIDMD5__&idfa=__IDFA__&callback=__CALLBACK_URL__&advid=__CID__&os=__OS__&ip=__IP__&ts=__TSMS__

### /uprecall 召回api 参数: t=1打开链接,2下载 3设备打开，now=当前时间戳 （2个参数签名）
1、打开：?t=1&now=1606103184
2、下载：?t=2&now=1606103184
3、设备打开：?t=3&now=1606103184