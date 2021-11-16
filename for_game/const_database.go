package for_game

const MONGODB_NINGMENG = "ningmeng"
const MONGODB_NINGMENG_LOG = "ningmeng_log"

//----------------- period 数据的 key (天，周，月都用同一个 key) ----------------
const (
	CHANGE_PHONE      = "Change_Phone"
	CHECK_PAYPASSWORD = "Check_PayPassword"
	AUTOLOGIN_TOKEN   = "AutoLogin_token"
	CHECK_FRESH_NUM   = "Check_Fresh_Num" //语音名片刷新次数，每天100次
)

//--------------------------- 表名统一放这里 ---------------------------
const (
	TABLE_GOLDCHANGELOG      = "log_gold_change"        //金币变化log
	TABLE_COINCHANGELOG      = "log_coin_change"        //硬币变化log
	TABLE_ESPORTCHANGELOG    = "log_esport_coin_change" //电竞币变化log
	TABLE_SOURCETYPE         = "source_type"            //现金源类型表
	TABLE_MANAGER            = "manager"                //后台管理员表
	TABLE_MANAGER_TYPES      = "manager_types"          //后台客服类型
	TABLE_HTOPTLOG_NAME      = "log_backstage_opt"      //后台人员操作日志的表名
	TABLE_ROLEPOWER          = "role_power"             //角色权限管理
	TABLE_AUTHGROUP          = "auth_group"             //角色权限树状图
	TABLE_IPLIBRARY          = "ip_library"             //IP库
	TABLE_POS_DEVICECODE     = "pos_devicecode"         //终端设备号
	TABLE_POS_DEVICEIDFA     = "pos_deviceidfa"         //终端设备idfa号
	TABLE_PLAYER_FREEZE_LOG  = "log_player_freeze"      //用户封号记录
	TABLE_POS_ADV_DEVICEIDFA = "pos_adv_deviceidfa"     //终端设备Adv idfa号

	TABLE_GENERAL_QUOTA    = "pay_generalquota"     //通用额度配置
	TABLE_PAYTYPE          = "pay_paytype"          //支付类型
	TABLE_PAYSCENE         = "pay_payscene"         //支付场景
	TABLE_PAYMENTSETTING   = "pay_payment_setting"  //支付设定
	TABLE_PAYMENTPLATFORM  = "pay_payment_platform" //支付平台
	TABLE_PLATFORM_CHANNEL = "pay_platform_channel" //平台通道

	TABLE_ID_GENERATOR = "id_generator" //自增id表
	TABLE_ORDER        = "order"        //订单表

	TABLE_REPORTJOB                  = "report_job"                 //报表任务进度记录表
	TABLE_LOG_LOGIN_INFO             = "log_login_info"             //登陆日志
	TABLE_FREEZEIP                   = "freeze_ip"                  //冻结IP
	TABLE_FREEZEACCOUNT              = "freeze_account"             //冻结IP
	TABLE_PLAYERKEEPREPORT           = "report_player_keep"         //玩家留存报表
	TABLE_ONLINETIMELOG              = "log_onlinetime"             //玩家在线时长日志
	TABLE_PLAYER_ACTIVE_REPORT       = "report_player_active"       //玩家日活跃度报表
	TABLE_PLAYER_WEEK_ACTIVE_REPORT  = "report_player_week_active"  //玩家周活跃度报表
	TABLE_PLAYER_MONTH_ACTIVE_REPORT = "report_player_month_active" //玩家月活跃度报表
	TABLE_PLAYER_BEHAVIOR_REPORT     = "report_player_behavior"     //用户行为报表
	TABLE_INOUTCASHSUM_REPORT        = "report_inoutcashsum"        //出入款汇总报表
	TABLE_LOGIN_REGISTER_REPORT      = "report_register_login"      //埋点登录注册报表
	TABLE_VC_BURYING_POINT_REPORT    = "report_vc_burying_point"    //恋爱匹配埋点
	TABLE_LOGIN_REGISTER_LOG         = "register_login_log"         //埋点登录注册日志
	TABLE_OPERATION_CHANNEL          = "operation_channel"          //运营渠道表
	TABLE_OPERATION_CHANNEL_USE      = "operation_channel_use"      //运营渠道表
	TABLE_OPERATION_CHANNEL_REPORT   = "report_operation_channel"   //运营渠道汇总报表
	TABLE_CHANNEL_REPORT             = "report_channel"             //渠道报表
	TABLE_ARTICLE_REPORT             = "report_article"             //文章报表
	TABLE_NOTICE_REPORT              = "report_notice"              //通知报表
	TABLE_SQUARE_REPORT              = "report_square"              //社交广场报表
	TABLE_PLAYERONLINE_REPORT        = "report_playerOnline"        //玩家在线报表
	TABLE_PLAYERLOG_LOCATION_REPORT  = "report_playerlog_location"  //玩家登录分布地域报表
	TABLE_ADV_REPORT                 = "report_adv"                 //广告报表
	TABLE_NEARBY_ADV_REPORT          = "report_nearby_adv"          //附近的人引导项报表
	TABLE_RECALL_REPORT              = "report_recall"              //召回报表
	TABLE_COIN_PRODUCT_REPORT        = "report_coin_product"        //虚拟商城报表
	TABLE_BUTTON_CLICK_REPORT        = "report_button_click"        //按钮点击行为报表
	TABLE_BUTTON_CLICK_LOG           = "log_button_click"           //按钮点击行为日志表
	TABLE_PAGE_REGLOG                = "log_page_reglog"            //注册登录页面埋点日志
	TABLE_PAGE_REGLOG_REPORT         = "report_page_reglog"         //注册登录页面埋点报表

	TABLE_SERVER_INFO         = "server_info"    //服务器信息表
	TABLE_PLAYER_BASE         = "player_base"    //玩家表
	TABLE_PLAYER_ACCOUNT      = "player_account" //玩家登录表
	TABLE_PLAYER_PERIOD       = "player_period"
	TABLE_PLAYER_COMPLAINT    = "player_complaint"    //玩家投诉
	TABLE_PLAYER_FRIEND       = "player_friend"       //玩家好友表
	TABLE_PLAYER_EMOTICON     = "player_emoticon"     //玩家表情库
	TABLE_PLAYER_EQUIPMENT    = "player_equipment"    //玩家装备表
	TABLE_PLAYER_BAG_ITEM     = "player_bag_item"     //玩家背包道具
	TABLE_PLAYER_BCOIN_LOG    = "player_bcoin_log"    //玩家绑定硬币获得记录
	TABLE_PLAYER_CHAT_SESSION = "player_chat_session" //玩家会话列表
	TABLE_PLAYER_INTIMACY     = "player_intimacy"     //玩家亲密度

	TABLE_TEAM_DATA              = "team_data"         //聊天群组表
	TABLE_TEAMMEMBER_DATA        = "teammember_data"   //个人群里数据
	TABLE_TEAMMEMBER             = "team_members"      //群成员
	TABLE_TEAM_CHAT_LOG          = "team_chat_log"     //群聊天记录
	TABLE_PERSONAL_CHAT_LOG      = "personal_chat_log" //个人聊天记录
	TABLE_MONGODB_CHAT_SESSION   = "chat_session"
	TABLE_RED_PACKET             = "red_packet"             //红包表数据
	TABLE_RED_PACKET_LOG         = "red_packet_log"         //红包领取记录
	TABLE_TRANSFER_MONEY         = "transfer_money"         //转账数据
	TABLE_NEARBY_LOG             = "nearby_log"             //附近的人打招呼日志
	TABLE_NEARBY_SESSIOIN_LIST   = "nearby_session_list"    // 新版本附近的人打招呼会话列表
	TABLE_NEARBY_MESSAGE_NEW_LOG = "nearby_message_new_log" // 附近的人打招呼的消息
	TABLE_SYSTEM_NOTICE          = "system_notice"          // 系统公告
	TABLE_ASSISTANT              = "assistant"              // 畅聊助手
	TABLE_SYS_PARAMETER          = "sys_parameter"          //系统功能参数
	TABLE_INTERESTTAG            = "interest_tag"           //兴趣标签
	TABLE_INTERESTTYPE           = "interest_type"          //兴趣分类
	TABLE_INTERESTGROUP          = "interest_group"         //兴趣组合
	TABLE_REDPACKET_STATISTICS   = "redpacket_statistics"   //红包玩家统计表
	TABLE_ARTICLE                = "article"                //文章
	TABLE_ARTICLE_ZAN            = "article_zan"            //文章赞
	TABLE_ARTICLE_COMMENT        = "article_comment"        //文章评论
	TABLE_TWEETS                 = "tweets"                 //推文
	TABLE_REGISTER_PUSH          = "register_push"          //注册推送
	TABLE_PLAYER_TWEETS          = "player_tweets"          //用户推文
	TABLE_CUSTOMTAG              = "custom_tag"             //自定义标签
	TABLE_CRAWL_WORDS            = "crawl_words"            //抓取词
	TABLE_GRABTAG                = "grab_tag"               //抓取标签
	TABLE_SYSTEM_LOG             = "system_log"             //抓取词日志
	TABLE_PLAYER_CRAWL_WORDS     = "player_crawl_words"     //玩家抓取词日志
	TABLE_PLAYER_TALK_LOG        = "player_talk_log"        //玩家敏感屏蔽日志
	TABLE_CONTROL_MODERATIONS    = "control_moderations"    //屏蔽控制值
	TABLE_DIRTY_WORDS            = "dirty_words"            //屏蔽词库
	TABLE_NEAR_LEAD              = "near_lead"              //附近的人引导项
	TABLE_NEAR_FAST_TERM         = "near_fast_term"         //附近的人快捷打招呼
	TABLE_SIGNATURE              = "signature"              // 个性签名库

	TABLE_FEATURES_APPPUSHMSG = "features_apppushmsg" //推送APP消息
	TABLE_FEATURES_HELPTYPE   = "features_helptype"   //帮组类型
	TABLE_FEATURES_HELPMSG    = "features_helpmsg"    //帮组设置

	TABLE_SHOP_ITEMS          = "shop_items"           // 商品表
	TABLE_SHOP_BILLS          = "shop_bills"           // 付钱订单表
	TABLE_SHOP_CACHE_EXPRESS  = "shop_cache_express"   // 物品订单表的物流缓存
	TABLE_SHOP_ORDERS         = "shop_orders"          // 物品订单表
	TABLE_RECEIVE_ADDRESS     = "shop_receive_address" // 收货地址表
	TABLE_DELIVER_ADDRESS     = "shop_deliver_address" // 发货地址表
	TABLE_PLAYER_CART         = "shop_player_cart"     // 用户购物车表
	TABLE_ITEM_COMMENT        = "shop_item_comment"    // 商品留言
	TABLE_SHOP_MESSAGE        = "shop_messages"        // 商店消息
	TABLE_SHOP_ALI_AUDIT_FAIL = "shop_ali_audit_fail"  // 商城上架审核失败信息
	TABLE_SHOP_PLAYER         = "shop_player"          // 商城发布商品认证用户信息表
	TABLE_SHOP_LIKE           = "shop_like_record"     // 留言点赞记录
	TABLE_PLAYER_STORE        = "shop_player_store"    // 用户收藏表
	TABLE_SHOP_POINT_CARD     = "shop_point_card"      // 商品点卡表

	TABLE_WAITER_MESSAGE     = "waiter_message"     //客服消息表
	TABLE_WAITER_PERFORMANCE = "waiter_performance" //客服绩效表
	TABLE_WAITER_FAQ         = "waiter_faq"         //客服常见问题表
	TABLE_WAITER_FASTREPLY   = "waiter_fastreply"   //客服快捷语表

	TABLE_SQUARE_DYNAMIC      = "square_dynamic"      //社交广场动态信息
	TABLE_SQUARE_DYNAMIC_SNAP = "square_dynamic_snap" // 延时发布的动态临时表
	TABLE_SQUARE_COMMENT      = "square_comment"      //社交广场评论信息
	TABLE_SQUARE_COMMENT_ZAN  = "square_comment_zan"  //社交广场评论点赞
	TABLE_SQUARE_ZAN          = "square_zan"          //社区点赞信息
	TABLE_SQUARE_ATTENTION    = "square_attention"    //社区关注信息
	TABLE_READ_DYNAMIC_DEVICE = "read_dynamic_device" //社区读取动态设备库

	TABLE_CANCEL_ACCOUNT       = "player_cancel_account" //注销请求表
	TABLE_CANCEL_ACCOUNT_LIST  = "cancel_account_list"   //注销成功列表
	TABLE_SQUARE_TOP_TIMER_MGR = "square_top_timer_mgr"  // 社交广场置顶定时任务管理器

	TABLE_DATA_COUNTRY = "data_country" //国家库
	TABLE_DATA_AREA    = "data_area"    //区域库
	TABLE_DATA_REGION  = "data_region"  //省库
	TABLE_DATA_CITY    = "data_city"    //市库

	TABLE_ADV_DATA       = "adv_data"       //广告库
	TABLE_ADV_LOG        = "adv_log"        //广告埋点日志数据
	TABLE_NEARBY_ADV_LOG = "nearby_adv_log" //附近的人引导项日志数据

	TABLE_VV_DURSTION_LOG = "log_vv_duration" //视频语音使用时长日志

	//=======话题=========
	TABLE_TOPIC_TYPE                   = "topic_type"                   // 话题类型
	TABLE_TOPIC                        = "topic"                        // 具体的话题
	TABLE_TOPIC_PLAYER_ATTENTION       = "topic_player_attention"       // 玩家所关注的话题.
	TABLE_TOPIC_PLAYER_DEVOTE_DAY      = "topic_player_devote_day"      // 玩家贡献日榜.
	TABLE_TOPIC_PLAYER_DEVOTE_MONTH    = "topic_player_devote_month"    // 玩家贡献月榜.
	TABLE_TOPIC_PLAYER_DEVOTE_TOTAL    = "topic_player_devote_total"    // 玩家贡献总榜.
	TABLE_TOPIC_APPLY_LOG              = "apply_edit_topic_info"        // 话题信息修改申请表.
	TABLE_TOPIC_APPLY_TOPIC_MASTER     = "apply_topic_master"           // 申请话题主记录表.
	TABLE_TOPIC_MASTER_DEL_DYNAMIC_LOG = "topic_master_del_dynamic_log" // 话题主删除动态操作日志表.

	TABLE_RECALLPLAYER_LOG = "log_recall_player" //玩家回归日志

	//====硬币商场=========
	TABLE_COIN_RECHARGE       = "coin_recharge"        //硬币配置表
	TABLE_PROPS_ITEM          = "props_item"           //道具配置表
	TABLE_COIN_PRODUCT        = "coin_product"         //虚拟商品表
	TABLE_PLAYER_BAGITEM      = "player_bag_item"      //用户背包
	TABLE_PLAYER_GETPROPS_LOG = "log_player_get_props" //用户道具获得日志表

	//====恋爱交友=======================
	TABLE_STARSIGNS_TAG          = "starsigns_tag"          //星座标签表
	TABLE_CHARACTER_TAG          = "character_tag"          //个性标签表
	TABLE_BG_VOICE_TAG           = "bg_voice_tag"           //录音背景标签
	TABLE_PLAYER_MIX_VOICE_VIDEO = "player_mix_voice_video" //玩家录音作品表
	TABLE_BG_VOICE_VIDEO         = "bg_voice_video"         //录音背景资源表
	TABLE_PLAYER_VC_ZAN_LOG      = "player_vc_zan_log"      //玩家点赞日志
	TABLE_PLAYER_VC_SAY_HI_LOG   = "player_vc_say_hi_log"   //玩家sayHi日志
	TABLE_PLAYER_ATTENTION_LOG   = "player_attention_log"   //玩家关注日志
	TABLE_MATCH_GUIDE            = "match_guide"            //匹配引导语
	TABLE_SAY_HI                 = "say_hi"                 //打招呼
	TABLE_SYSTEM_BG_IMAGE        = "system_bg_image"        //系统背景图
	TABLE_INTIMACY_COINFIG       = "intimacy_config"        //亲密度配置
	TABLE_PLAYER_OPERATE         = "player_operate"         //运营号配置
	//埋点表
	TABLE_BURYING_POINT_LOG = "burying_point_log" //埋点日志
	//=========弹窗广告==============
	TABLE_INDEX_TIPS  = "index_tips"  //主页菜单选项
	TABLE_POP_SUSPEND = "pop_suspend" //悬浮窗广告

)
