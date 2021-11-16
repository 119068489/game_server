package for_game

const (
	//=================电竞===================
	TABLE_ESPORTS_CRAWL_JOB    = "esports_job_crawl"    //爬虫进度表
	TABLE_ESPORTS_NEWS_SOURCE  = "esports_news_source"  //新闻资讯库 爬虫数据
	TABLE_ESPORTS_VIDEO_SOURCE = "esports_video_source" //视频库 爬虫数据
	TABLE_ESPORTS_NEWS         = "esports_news"         //新闻资讯 发布数据
	TABLE_ESPORTS_VIDEO        = "esports_video"        //视频 发布数据
	TABLE_ESPORTS_HISTORY      = "esports_history"      //历史战绩

	TABLE_ESPORTS_SYS_MSG                   = "esports_sys_msg"              //系统消息
	TABLE_ESPORTS_GAME_ORDER_SYS_MSG        = "esports_game_order_sys_msg"   //竞猜消息待开奖
	TABLE_ESPORTS_GAME_ORDER_SYS_MSG_E      = "esports_game_order_sys_msg_e" //竞猜消息已处理
	TABLE_ESPORTS_LABEL                     = "esports_label"                //系统标签
	TABLE_ESPORTS_CAROUSEL                  = "esports_carousel"             //轮播图
	TABLE_ESPORTS_PLAYER                    = "esports_player"
	TABLE_ESPORTS_HOMEMENUCONFIG            = "esports_home_menu_config"
	TABLE_ESPORTS_COMMENT_NEWS              = "esports_comment_new"               // 资讯评论
	TABLE_ESPORTS_COMMENT_VIDEO             = "esports_comment_video"             //视频评论
	TABLE_ESPORTS_COMMENT_NEWS_REPLY        = "esports_comment_new_reply"         //资讯的评论回复
	TABLE_ESPORTS_COMMENT_VIDEO_REPLY       = "esports_comment_video_reply"       //资讯的视频的回复
	TABLE_ESPORTS_GAME                      = "esports_game"                      //比赛表
	TABLE_ESPORTS_FLOW_VIDEO_HISTORY        = "esports_flow_video_history"        //视频观看历史表
	TABLE_ESPORTS_FLOW_LIVE_HISTORY         = "esports_flow_live_history"         //放映厅观看历史表
	TABLE_ESPORTS_FLOW_LIVE_FOLLOW_HISTORY  = "esports_flow_live_follow_history"  //放映厅关注表
	TABLE_ESPORTS_GAME_DETAIL               = "esports_game_detail"               //比赛表详情
	TABLE_ESPORTS_TEAM_HIS_INFO             = "esports_team_history"              //两队历史交锋、两队胜败统计、两队天敌克制统计
	TABLE_ESPORTS_GAME_GUESS                = "esports_game_guess"                //比赛动态信息表(早盘、滚盘)
	TABLE_ESPORTS_USE_ROLL_GUESS            = "esports_game_use_roll"             //使用滚盘
	TABLE_ESPORTS_GUESS_BET_RECORD          = "esports_guess_bet_record"          //投注记录表
	TABLE_ESPORTS_BET_SLIP_REPORT           = "esports_report_bet_slip"           //注单统计报表
	TABLE_ESPORTS_BET_RISK_ONE_DAY          = "esports_bet_risk_one_day"          //风控用户当日投注额度记录表
	TABLE_ESPORTS_BET_RISK_PLATFORM_DAY_SUM = "esports_bet_risk_platform_day_sum" //风控平台当日投注额度记录表
	TABLE_ESPORTS_ROOM_CHAT_MSG_LOG         = "esports_room_chat_msg_log"         //放映厅发言记录
	TABLE_ESPORTS_BPS_CLICK_LOG             = "esports_bps_click_log"             //埋点点击记录
	TABLE_ESPORTS_BPS_DURATION_LOG          = "esports_bps_duration_log"          //埋点停留時長记录
	TABLE_ESPORTS_LOL_REAL_TIME_DATA        = "esports_lol_real_time_data"        //LOL游戏实时数据
	TABLE_ESPORTS_WZRY_REAL_TIME_DATA       = "esports_wzry_real_time_data"       //WZRY游戏实时数据
	TABLE_ESPORTS_GIVE_WHITELIST            = "esports_give_whitelist"            //充值赠送白名单
	TABLE_ESPORTS_EXCHANGE_CFG              = "esports_exchange_cfg"              //电竞币兑换配置表
	TABLE_ESPORTS_EXCHANGE_FIRST            = "esports_exchange_first"            //电竞币兑换首冲记录表

	// 报表
	TABLE_ESPORTS_BASIS_POINTS_REPORT_DAY          = "esports_report_basis_points_day"          //基础埋点日报表
	TABLE_ESPORTS_BASIS_POINTS_REPORT_WEEK         = "esports_report_basis_points_week"         //基础埋点周报表
	TABLE_ESPORTS_BASIS_POINTS_REPORT_MONTH        = "esports_report_basis_points_month"        //基础埋点月报表
	TABLE_ESPORTS_MENU_POINTS_REPORT_DAY           = "esports_report_menu_points_day"           //Tab菜单埋点日报表
	TABLE_ESPORTS_MENU_POINTS_REPORT_WEEK          = "esports_report_menu_points_week"          //Tab菜单埋点周报表
	TABLE_ESPORTS_MENU_POINTS_REPORT_MONTH         = "esports_report_menu_points_month"         //Tab菜单埋点月报表
	TABLE_ESPORTS_LABEL_POINTS_REPORT_DAY          = "esports_report_label_points_day"          //标签埋点日报表
	TABLE_ESPORTS_LABEL_POINTS_REPORT_WEEK         = "esports_report_label_points_week"         //标签埋点周报表
	TABLE_ESPORTS_LABEL_POINTS_REPORT_MONTH        = "esports_report_label_points_month"        //标签埋点月报表
	TABLE_ESPORTS_NEWS_AMUSE_POINTS_REPORT_DAY     = "esports_report_news_amuse_points_day"     //资讯娱乐埋点日报表
	TABLE_ESPORTS_NEWS_AMUSE_POINTS_REPORT_WEEK    = "esports_report_news_amuse_points_week"    //资讯娱乐埋点周报表
	TABLE_ESPORTS_NEWS_AMUSE_POINTS_REPORT_MONTH   = "esports_report_news_amuse_points_month"   //资讯娱乐埋点月报表
	TABLE_ESPORTS_VDOHALL_POINTS_REPORT_DAY        = "esports_report_vdohall_points_day"        //放映厅埋点日报表
	TABLE_ESPORTS_VDOHALL_POINTS_REPORT_WEEK       = "esports_report_vdohall_points_week"       //放映厅埋点周报表
	TABLE_ESPORTS_VDOHALL_POINTS_REPORT_MONTH      = "esports_report_vdohall_points_month"      //放映厅埋点月报表
	TABLE_ESPORTS_APPLYVDOHALL_POINTS_REPORT_DAY   = "esports_report_applyvdohall_points_day"   //申请放映厅埋点日报表
	TABLE_ESPORTS_APPLYVDOHALL_POINTS_REPORT_WEEK  = "esports_report_applyvdohall_points_week"  //申请放映厅埋点周报表
	TABLE_ESPORTS_APPLYVDOHALL_POINTS_REPORT_MONTH = "esports_report_applyvdohall_points_month" //申请放映厅埋点月报表
	TABLE_ESPORTS_MATCHLS_POINTS_REPORT_DAY        = "esports_report_matchls_points_day"        //赛事列表埋点日报表
	TABLE_ESPORTS_MATCHLS_POINTS_REPORT_WEEK       = "esports_report_matchls_points_week"       //赛事列表埋点周报表
	TABLE_ESPORTS_MATCHLS_POINTS_REPORT_MONTH      = "esports_report_matchls_points_month"      //赛事列表埋点月报表
	TABLE_ESPORTS_MATCHDIL_POINTS_REPORT_DAY       = "esports_report_matchdil_points_day"       //赛事详情埋点日报表
	TABLE_ESPORTS_MATCHDIL_POINTS_REPORT_WEEK      = "esports_report_matchdil_points_week"      //赛事详情埋点周报表
	TABLE_ESPORTS_MATCHDIL_POINTS_REPORT_MONTH     = "esports_report_matchdil_points_month"     //赛事详情埋点月报表
	TABLE_ESPORTS_GUESS_POINTS_REPORT_DAY          = "esports_report_guess_points_day"          //竞猜页埋点日报表
	TABLE_ESPORTS_GUESS_POINTS_REPORT_WEEK         = "esports_report_guess_points_week"         //竞猜页埋点周报表
	TABLE_ESPORTS_GUESS_POINTS_REPORT_MONTH        = "esports_report_guess_points_month"        //竞猜页埋点月报表
	TABLE_ESPORTS_MSG_POINTS_REPORT_DAY            = "esports_report_msg_points_day"            //消息页埋点日报表
	TABLE_ESPORTS_MSG_POINTS_REPORT_WEEK           = "esports_report_msg_points_week"           //消息页埋点周报表
	TABLE_ESPORTS_MSG_POINTS_REPORT_MONTH          = "esports_report_msg_points_month"          //消息页埋点月报表
	TABLE_ESPORTS_COIN_POINTS_REPORT_DAY           = "esports_report_coin_points_day"           //电竞币页埋点日报表
	TABLE_ESPORTS_COIN_POINTS_REPORT_WEEK          = "esports_report_coin_points_week"          //电竞币埋点周报表
	TABLE_ESPORTS_COIN_POINTS_REPORT_MONTH         = "esports_report_coin_points_month"         //电竞币埋点月报表
)
