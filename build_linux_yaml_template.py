# -*- coding: UTF-8 -*-
# yaml 配置文件的模板
# 当前脚本只能在 Windows 下运行
# 可以使用 Python 3 的语法，Windows 下大家约定用 Python 3.7 以上

TEMPLATE_HALL = """# 若要加配置、改配置名,请一定在打包脚本 build_linux_yaml_template.py 的模板上也同时加上

INCLUDE: ["config_share.yaml","config_hall_secret.yaml"] # 可以包含多个其他 *.yaml 文件

SERVER_TYPE: 2     #服务器类型
SERVER_ID: 200{group}     #服务器编号:2000开头
SERVER_NAME: "大厅服{group}"  #服务器名称
SERVER_ADDR: "{host}"  #服务器外部地址
SERVER_ADDR_INTERNAL: "{host}" #服务器内部地址
LISTEN_PORT_FOR_CLIENT: 200{group} # 监听 WebSocket 客户端的链接(U3D,H5)取值区间(2001-2299)
LISTEN_PORT_FOR_CLIENT_TCP: 230{group}  #监听 TCP 客户端的链接
LISTEN_PORT_FOR_GAME: 255{group} # 监听 其他大厅的链接
LISTEN_PORT_FOR_HALL: 265{group} # 监听 其他大厅的链接
LISTEN_PORT_FOR_BACKSTAGE: 275{group} # 监听 后台的连接
LISTEN_PORT_FOR_SHOP: 285{group} # 监听 商场的连接
LISTEN_PORT_FOR_SQUARE: 295{group} # 监听 社交广场的连接
LISTEN_ADDR_FOR_TCP_SOCKET_CLIENT: 250{group} # 监听 TcpSocket 客户端的链接(机器人)
LISTEN_ADDR_FOR_WEB_API: 260{group} # Web API 地址

PAYMENT_ADDRESS: "http://54.92.13.178:151{group}/payment" # 第三方支付回调监听端口(端口要和上面的 LISTEN_ADDR_FOR_WEB_API 一样)

TLS_CRT_FILE_PATH: "" # crt 文件路径. 比如 www.test.com.crt 或 ./test/www.test.com.crt, 值为空串时变成走 ws 协议
TLS_KEY_FILE_PATH: "" # key 文件路径. 比如 www.test.com.key 或 ./test/www.test.com.key, 值为空串时变成走 ws 协议

SHARE_QR_CODE_URL: "http://54.92.13.178:8001/#/signIn?site=%s&parentid=%d&channel=%d" # 扫码跳转注册新玩家页面
CLIENT_ARTICLE_URL: "http://appshare.lemonchat.cn/article.html"  # 文章访问地址
VERSION_NUMBER: "1.0.0" # 版本号.升级了版本号记得到 build_linux_yaml_template.py 打包模板修改一下


"""

#-----------------------------------------------------------------------------------------------------

TEMPLATE_LOGIN = """# 若要加配置、改配置名,请一定在打包脚本 build_linux_yaml_template.py 的模板上也同时加上

INCLUDE: ["config_share.yaml","config_hall_secret.yaml"] # 可以包含多个其他 *.yaml 文件
SERVER_TYPE: 1     #服务器类型
SERVER_ID: 100{group}     #服务器编号 新增服id递增取值区间(1001-1999)
SERVER_NAME: "登录服{group}"  #服务器名称
SERVER_ADDR: "{host}"  #服务器外部地址
SERVER_ADDR_INTERNAL: "{host}"  #服务器内部地址
LISTEN_PORT_FOR_CLIENT: 111{group}  #监听客户端连接端口
LISTEN_PORT_FOR_CLIENT_TCP: 131{group}  #监听 TCP 客户端的链接
LISTEN_PORT_FOR_HALL: 150{group}    #监听大厅连接端口
TLS_CRT_FILE_PATH: "" # crt 文件路径. 比如 www.test.com.crt 或 ./test/www.test.com.crt, 值为空串时变成走 ws 协议
TLS_KEY_FILE_PATH: "" # key 文件路径. 比如 www.test.com.key 或 ./test/www.test.com.key, 值为空串时变成走 ws 协议
VERSION_NUMBER: "1.0.0" # 版本号.升级了版本号记得到 build_linux_yaml_template.py 打包模板修改一下
"""

#-----------------------------------------------------------------------------------------------------

TEMPLATE_BACKSTAGE = """# 若要加配置、改配置名,请一定在打包脚本 build_linux_yaml_template.py 的模板上也同时加上

INCLUDE: ["config_share.yaml","config_hall_secret.yaml"] # 可以包含多个其他 *.yaml 文件
SERVER_TYPE: 3     #服务器类型
SERVER_ID: 400{group}     #服务器编号:4000开头
SERVER_NAME: "后台{group}"  #服务器名称
SERVER_ADDR: "{host}"  #服务器外部地址
SERVER_ADDR_INTERNAL: "{host}"  #服务器内部地址
LISTEN_PORT_FOR_CLIENT: 400{group}  #监听浏览器客户端连接端口
LISTEN_PORT_FOR_API: 450{group} # WebAPI 地址,外部请求

# 超级帐户谷歌验证器配置
IS_GOOGLE_AUTH: False # 是否开启超级帐户谷歌验证器
GOOGLE_SECRET: "MQE3RYVEOI6GYMSBA3MDI2OXQKP27JUH" # 谷歌验证器秘钥

# 客户端版本管理JSON位置
CLIENT_VERSION_DATA: "./version.json"
CLIENT_TFSERVER_ADDR: "./tfserver.json"
"""
#-----------------------------------------------------------------------------------------------------

TEMPLATE_SHARE = """# 若要加配置、改配置名,请一定在打包脚本 build_linux_yaml_template.py 的模板上也同时加上

EDITION: "common" # 发行版本.有 common,yicai 等
IS_FORMAL_SERVER: False # 是否正式服
SERVER_CENTER_ADDR: "{center}:2371"        #etcd服务器中心地址
REDIS_SERVER_ADDR: "{center}:6379"       #redis服务器地址
REDIS_SERVER_PASS: "redis2020"             #redis服务器密码
IS_TFSERVER: False    #是否走转发
# MongoDB 数据库的链接配置
# 可以为每个站点配置不同的数据库信息，否则都用 share 相同的数据库信息
MONGODB_MASTER:
  ningmeng:
    user: "root"
    password: "MongoDB.2019"
    host: "{db}"
    port: "20201"
    max_pool_size: "512" # 最大连接数。mgo driver 代码里面的默认值是 4096
MONGODB_SLAVE:
  ningmeng_log:
    user: "root"
    password: "MongoDB.2019"
    host: "{db}"
    port: "20301"
    max_pool_size: "512" # 最大连接数。mgo driver 代码里面的默认值是 4096
"""
#-----------------------------------------------------------------------------------------------------
TEMPLATE_SHOP = """
INCLUDE: ["config_share.yaml","config_hall_secret.yaml"] # 可以包含多个其他 *.yaml 文件

SERVER_TYPE: 4     #服务器类型
SERVER_ID: 500{group}     #服务器编号:5000开头
SERVER_NAME: "商场{group}"  #服务器名称
SERVER_ADDR: "{host}"  #服务器外部地址
SERVER_ADDR_INTERNAL: "{host}"  #服务器内部地址
LISTEN_PORT_FOR_CLIENT: 511{group} #监听Websocket客户端连接端口
LISTEN_PORT_FOR_CLIENT_TCP: 531{group}  #监听 TCP 客户端的链接
LISTEN_ADDR_FOR_WEB_API: 570{group}  # Web API 地址
SEND_MAIL_ADDRESS: "cash77825@sina.com"     #H5邮箱发送邮箱地址
SEND_MAIL_PASS: "be6309dca39bd8fa"          #H5邮箱发送密码或者是授权码(根据邮箱配置)
SEND_MAIL_HOST: "smtp.sina.com"             #H5邮箱发送HOST
SEND_MAIL_PORT: "465"                       #H5邮箱发送PORT
"""
#-----------------------------------------------------------------------------------------------------

TEMPLATE_HALL_SECRET = """# 若要加配置、改配置名,请一定在打包脚本 build_linux_yaml_template.py 的模板上也同时加上
INCLUDE: [] # 可以包含多个其他 *.yaml 文件

# 若要加配置、改配置名,请一定在打包脚本 build_linux_yaml_template.py 的模板上也同时加上

INCLUDE: [] # 可以包含多个其他 *.yaml 文件

##国内版短信模板
ALIYU_ACCESSKEYID: "LTAI4FdUu8fXogvv9penpDCK"                           #阿里云短信accesskeyId
ALIYU_ACCESSKEYSECRET: "mUe3dSYqxCAo3X5HvSdiInQt7690aG"         #阿里云短信accesskeySecret
ALIYU_AREA: "cn-hangzhou"                                       #阿里云地区
ALIYU_SIGNNAME: "柠檬畅聊"                                        #阿里云短信模板名字
ALIYU_TEMPLEATECODE: "SMS_180350717"                            #阿里云短信模板代号

##国际版短信模板
ALIYU_SIGNNAME_INTERNATIONAL: "柠檬畅聊国际"                                        #阿里云短信模板名字
ALIYU_TEMPLEATECODE_INTERNATIONAL: "SMS_197896950"                            #阿里云短信模板代号

##腾讯云短信模板
TC_TEMPLATE_ID_CODE: "柠檬畅聊国际"                                        #阿里云短信模板名字
TC_SMS_SDK_APPID_CODE: "SMS_197896950"                            #阿里云短信模板代号

TC_SMS_LOGOUT_FAILED_TEMPLATE_ID: "705285"                                        #腾讯云用户注销失败短信模板名字
TC_SMS_LOGOUT_SUCCESS_TEMPLATE_ID: "705282"                                        #腾讯云用户注销成功短信模板名字
TC_SMS_SDK_APPID: "1400111500"                            #腾讯云SDKAppId
TC_SMS_SIGN: "柠檬畅聊"  #腾讯云签名
TC_SMS_SECRETID: "AKIDLRs45skfuLZZ7O87etzoOk6KHyPp7E1Z"
TC_SMS_SECRETKEY: "a3F9rXM44rYWENBy94JI3Wryn5SZPffS"

##应用微信appid
# appcode=100 使用
WEIXIN_APPID_100: "wx27682f51b810f0c2"                              #微信appid
WEIXIN_APPSECRET_100: "28c78426a9e5ac167c56554259430efb"            #微信appsecret
#apkcode=101 使用
WEIXIN_APPID_101: "wxce0ea7527a63f28d"                              #微信appid
WEIXIN_APPSECRET_101: "9df014d7ed57b41c0c1f3cc6b4f85137"            #微信appsecret
WEIXIN_URL: "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"#微信验证url
WEIXIN_USERURL: "https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s"
##小程序微信appid
# appcode=100 使用
WEIXIN_SGAME_APPID_100: "wx414f06b4024739c8"
WEIXIN_SGAME_SECRET_100: "645661043e2b05624c3e1bb9bec71c3b"
# appcode=101 使用
WEIXIN_SGAME_APPID_101: "wx905904883fc074d3"
WEIXIN_SGAME_SECRET_101: "98844f416ebfdbd23a7238c592e08a30"
##极光快速登录appid
# appcode=100 使用
JG_APPKEY_100: "36d53261f51f9fffe4980a77"
JG_SECRET_100: "8e6854a3c07bba4e0363fc37"
#apkcode=101 使用
JG_APPKEY_101: "95279e7cbe6af3de678647ee"
JG_SECRET_101: "4c87a3377d9594b297be41c9"

JH_AUTHPEOPLEIDKEY: "fd28afdb60abcb05f3b221f55a5f1c18"
JH_AUTHBANKIDKEY: "3f15b88acd0a636682f8648c118238ef"

MOB_APPKEY: "2f34428f2351b"
MOB_SECRET: "bd15bf06f8dcba4aa18a14d8b088b9e6"
MOB_URL: "http://api.push.mob.com/v3/push/createPush"
MOB_URL_V2: "http://api.push.mob.com/v2/push"
#秒到支付配置信息
MD_PAY_APPID: "wx414f06b4024739c8"
MD_PAY_KEY: "aa403187f550527bf7cda6affb0660ab"
MD_PAY_MERCHANTNO: "15099973008"
MD_PAY_AGENCYCODE: "86038810"
MD_PAY_NOTIFY_URL: "http://tf1.lemonchat.cn:2601/miaodao"

#通联支付配置信息
TL_PAY_APPID: "wx414f06b4024739c8"
TL_PAY_KEY: "64336113WEIwei"
TL_PAY_MERCHANTNO: "56058104816WR4N"
TL_PAY_MERCHANTAPPID: "00184603"
TL_PAY_NOTIFY_URL: "http://tf1.lemonchat.cn:2601/tonglian"

# 汇潮配置信息
HUICHAO_URL: "https://gwapi.yemadai.com/" # url
HUICHAO_MERCHANTNO: "50592" # 商户号
HUICHAO_PAY_CALLBACK_URL: "http://47.112.142.227:2601/huichaopay" # 汇潮支付回调地址
HUICHAO_PAY_DF_CALLBACK_URL: "http://47.112.142.227:2601/huichaoDF" # 汇潮支付 代付回调地址
HUICHAO_WXCOMPANYNO: "sweep-b90d814d534a4b219ab1fe0983f248e9"
HUICHAO_ZFBCOMPANYNO: "sweep-d3a493fa16d5475f9c655eda30c385a7"
HUICHAO_WXAPPID: "wx905904883fc074d3"

#统统付支付配置
TTP_URL: "http://117.48.192.183/cgi-bin/" # url
TTP_PAY_CALLBACK_URL: "http://47.112.142.227:2601/tongtongpay" # 统统付支付回调地址
TTP_SYSCODE: "20000115" #系统号
TTP_MERCHANTNO: "9062000058" # 商户号
TTP_MD5_KEY: "b11519815gzgas115fip6gzpu2k13"

TTP_WXAPPID: "wx905904883fc074d3"

"""
#-----------------------------------------------------------------------------------------------------
TEMPLATE_STATISTICS = """
INCLUDE: ["config_share.yaml","config_hall_secret.yaml"] # 可以包含多个其他 *.yaml 文件

SERVER_TYPE: 5     #服务器类型
SERVER_ID: 600{group}     #服务器编号:6000开头
SERVER_NAME: "统计服{group}"  #服务器名称
"""

#-----------------------------------------------------------------------------------------------------
TEMPLATE_SQUARE = """
INCLUDE: ["config_share.yaml"] # 可以包含多个其他 *.yaml 文件

SERVER_TYPE: 6     #服务器类型
SERVER_ID: 700{group}     #服务器编号:7000开头

REDIS_SERVER_DATABASE: 0             #redis数据库序号

SERVER_NAME: "社交广场1"  #服务器名称
SERVER_ADDR: "{host}"  #服务器外部地址
SERVER_ADDR_INTERNAL: "{host}"  #服务器内部地址
LISTEN_PORT_FOR_CLIENT: 711{group}  #监听Websocket客户端连接端口
LISTEN_PORT_FOR_CLIENT_TCP: 731{group}  #监听 TCP 客户端的链接
"""
#-----------------------------------------------------------------------------------------------------
