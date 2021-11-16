package shop

const (
	DATABASE_ERROR = "操作失败"

	LOGIN_PLAYER_NOT_FOUND = "玩家对象没有创建成功"
	LOGIN_TOKEN_WRONG      = "token验证码错误"

	CANCEL_ORDER_SUCCESS          = "取消订单成功"
	CANCEL_ORDER_EXPIRE_TIME_OVER = "该订单已过期"
	CANCEL_ORDER_STATE_CHANGE     = "订单状态变化,请刷新"

	UPLOAD_ITEM_NAME_IS_NULL          = "缺少好物标题信息"
	UPLOAD_ITEM_NAME_COUNT            = "好物标题只能为1~15个字符"
	UPLOAD_ITEM_IMAGE_COUNT           = "照片只能发布2~9张"
	UPLOAD_ITEM_TITLE_IS_NULL         = "缺少好物详细信息"
	UPLOAD_ITEM_TITLE_COUNT           = "好物正文只能为5~100个字符"
	UPLOAD_ITEM_ADDRESS_IS_NULL       = "缺少好物位置信息"
	UPLOAD_ITEM_COUNT_OUT_OF_RANGE    = "好物库存数量只能为1~999之间"
	UPLOAD_ITEM_SUCCESS               = "商品上架成功"
	UPLOAD_ITEM_PRICE_VAR             = "好物价格只能为1~5000之间"
	UPLOAD_ITEM_PEOPLE_VALIDATE_COUNT = "已超过单人上架最多商品数"

	EDIT_ITEM_ITEM_NOT_EXIST     = "商品不存在"
	EDIT_ITEM_NAME_IS_NULL       = "缺少好物标题信息"
	EDIT_ITEM_NAME_COUNT         = "好物标题只能为1~15个字符"
	EDIT_ITEM_TITLE_IS_NULL      = "缺少好物详细信息"
	EDIT_ITEM_TITLE_COUNT        = "好物正文只能为5~100个字符"
	EDIT_ITEM_IMAGE_NOT_ENOUGH   = "照片只能发布2~9张"
	EDIT_ITEM_SUCCESS            = "商品编辑成功"
	EDIT_ITEM_ADDRESS_IS_NULL    = "缺少好物位置信息"
	EDIT_ITEM_COUNT_OUT_OF_RANGE = "好物库存数量只能为1~999之间"
	EDIT_ITEM_PRICE_VAR          = "好物价格只能为1~5000之间"

	DELETE_ITEM_ITEM_NOT_EXIST   = "商品不存在"
	DELETE_ITEM_SUCCESS          = "商品删除成功"
	SOLD_OUT_ITEM_ITEM_NOT_EXIST = "商品不存在"
	SOLD_OUT_ITEM_SUCCESS        = "商品删除成功"

	SETTLEMENT_BUY_NULL          = "没有买任何商品"
	SETTLEMENT_ITEM_PARA_NO_SALE = "商品(%s)：已下架"
	SETTLEMENT_ITEM_NO_STOCK     = "商品(%s)：库存不足,还剩%v库存"
	SETTLEMENT_ITEM_BLACK        = "商品(%s)的商家和您存在黑名单关系"

	CREATE_ORDER_BUY_NULL      = "没有买任何商品"
	CREATE_ORDER_ITEM_NO_SALE  = "商品(%s)：商品已下架"
	CREATE_ORDER_SUCCESS       = "生成订单成功"
	CREATE_ORDER_ITEM_NO_STOCK = "商品(%s)：库存不足,还剩%v库存"
	CREATE_ORDER_ITEM_BLACK    = "商品(%s)的商家和您存在黑名单关系"

	PAY_ORDER_EXPIRE           = "订单过期"
	PAY_ORDER_CANCLE           = "订单已取消"
	PAY_ORDER_REPEATED         = "重复支付"
	PAY_ORDER_MONEY_NOT_ENOUGH = "余额不足"
	PAY_ORDER_USER_NOT_FOUND   = "找不到该用户"

	ADD_CART_COUNT_OVER             = "加购数量超过库存数量"
	ADD_CART_SUCCESS                = "加购成功"
	ADD_CART_NOT_SALE               = "商品已下架"
	ADD_CART_NOT_EXIST              = "商品不存在"
	SUB_CART_NOT_SALE               = "商品已下架"
	SUB_CART_NOT_EXIST              = "商品不存在"
	SUB_CART_SUCCESS                = "减购成功"
	REMOVE_ITEM_FROM_CART_NOT_EXIST = "没有选择任何商品"
	REMOVE_ITEM_FROM_CART_SUCCESS   = "购物车商品删除成功"
	CART_BLACK_VAR                  = "你和该商家已存在黑名单关系,无法交易"
	SELECT_CART_NOT_EXIST           = "商品不存在"
	SELECT_CART_NOT_SALE            = "商品已下架"
	QUERY_CART_NOT_SALE_WARN        = "已下架"
	QUERY_CART_NOT_STOCK_WARN       = "库存不足"
	QUERY_CART_NOT_BLACK_WARN       = "黑名单"

	ADD_STORE_SUCCESS              = "添加收藏成功"
	REMOVE_ITEM_FROM_STORE_SUCCESS = "取消收藏成功"
	BATCH_REMOVE_STORE_NOT_SELECT  = "没有选择任何商品"
	BATCH_REMOVE_STORE_SUCCESS     = "添加收藏成功"
	BATCH_REMOVE_STORE_NOT_EXIST   = "存在下架的商品"

	DETAIL_SHOP_ITEM_NOT_SALE  = "商品已下架"
	DETAIL_SHOP_ITEM_NOT_EXIST = "商品不存在"
	DETAIL_SHOP_ITEM_BLACK_VAR = "你和该商家已存在黑名单关系,无法查看此商品"

	RECEIVE_ADDRESS_ADD_SUCCESS = "添加地址成功"

	RECEIVE_ADDRESS_EDIT_SUCCESS   = "编辑地址成功"
	RECEIVE_ADDRESS_DELETE_SUCCESS = "删除地址成功"

	DELIVER_ADDRESS_ADD_SUCCESS = "添加地址成功"

	DELIVER_ADDRESS_EDIT_SUCCESS   = "编辑地址成功"
	DELIVER_ADDRESS_DELETE_SUCCESS = "删除地址成功"

	UPLOAD_ITEM_COMMENT_AUDIT_FAIL = "请勿发表敏感言论"
	UPLOAD_ITEM_COMMENT_SUCCESS    = "留言成功"
	UPLOAD_ITEM_COMMENT_NOT_BLANK  = "留言内容不能为空"

	UPLOAD_ITEM_EVALUTE_SUCCESS    = "评价成功"
	UPLOAD_ITEM_EVALUTE_AUDIT_FAIL = "请勿发表敏感言论"
	UPLOAD_ITEM_REPEATED           = "商品评价重复"
	UPLOAD_ITEM_EVALUTE_NOT_BLANK  = "评价内容不能为空"

	ORDER_DELETE_ORDER_NOT_FOUND   = "相关订单已被删除"
	ORDER_DELETE_ORDER_STATE_WRONG = "当前状态下，订单不能删除"
	ORDER_DELETE_ORDER_NOT_OWNER   = "你不是该订单的拥有者"
	ORDER_DELETE_SUCCESS           = "订单删除成功"

	DETAIL_ORDER_BLACK_VAR = "你和该商家已存在黑名单关系,无法查看此订单"

	ORDER_EXPRESS_UPLOAD_ORDER_NOT_FOUND = "订单不存在"
	ORDER_EXPRESS_UPLOAD_REPEATED        = "订单可能已经发货，请刷新"
	ORDER_EXPRESS_UPLOAD_SUCCESS         = "快递单号提交"

	ORDER_EDIT_ADDRESSS_ORDER_NOT_FOUND = "订单不存在"
	ORDER_EDIT_ADDRESSS_STATE_ERROR     = "订单已发货，收货地址不能修改"
	ORDER_EDIT_ADDRESSS_SUCCESS         = "收货地址修改成功"
	ORDER_EDIT_ADDRESSS_COUNT           = "修改收货地址次数已达到上限,修改失败"

	ORDER_EDIT_DELIVER_ADDRESSS_ORDER_NOT_FOUND = "订单不存在"
	ORDER_EDIT_DELIVER_ADDRESSS_STATE_ERROR     = "订单已发货，发货地址不能修改"
	ORDER_EDIT_DELIVER_ADDRESSS_SUCCESS         = "发货地址修改成功"

	DELAY_RECEIVE_ORDER_NOT_FOUND   = "订单不存在"
	DELAY_RECEIVE_ORDER_STATE_ERROR = "订单可能已自动收货,请刷新"
	DELAY_RECEIVE_REPEATE           = "每笔订单只能延长一次哦"
	DELAY_RECEIVE_TIME_NOT_ARRIVAL  = "亲，距离结束时间3天才可以申请哦"
	DELAY_RECEIVE_SUCCESS           = "亲，已为你延长1周的收货时间"

	CONFIRM_RECEIVE_ORDER_NOT_FOUND = "订单不存在"
	CONFIRM_RECEIVE_STATE_ERROR     = "订单可能已自动收货,请刷新"
	CONFIRM_RECEIVE_SUCCESS         = "收货地址修改成功"

	//买家消息
	MESSAGE_TO_BUYER_SEND_PUSH = "您的订单已发货"
	MESSAGE_TO_BUYER_SEND      = "您的订单已发货"

	MESSAGE_TO_BUYER_SIGN_PUSH = "您的订单已签收"
	MESSAGE_TO_BUYER_SIGN      = "您的订单已签收"

	//卖家消息
	MESSAGE_TO_SELLER_SIGN_PUSH   = "您的好物已被签收"
	MESSAGE_TO_SELLER_SIGN        = "您的好物已被签收"
	MESSAGE_TO_SELLER_NEW_PUSH    = "您收到一笔新的订单"
	MESSAGE_TO_SELLER_NEW         = "您收到一笔新的订单"
	MESSAGE_TO_SELLER_REMIND_PUSH = "您收到一条发货提醒"
	MESSAGE_TO_SELLER_REMIND      = "您收到一条发货提醒"
	MESSAGE_TO_SELLER_PAY_PUSH    = "您收到一条发货提醒"
	MESSAGE_TO_SELLER_PAY         = "您收到一条发货提醒"

	MESSAGE_TO_SELLER_ADDCHANGE_PUSH = "您的订单收货地址改变"
	MESSAGE_TO_SELLER_ADDCHANGE      = "您的订单收货地址改变"
	MESSAGE_TO_SELLER_EVALUATE_PUSH  = "您的订单已评价"
	MESSAGE_TO_SELLER_EVALUATE       = "您的订单已评价"
	MESSAGE_TYPE_AFTERSALES_BUYER_1  = "卖家回复了您的评价~"
	MESSAGE_TYPE_AFTERSALES_BUYER_2  = "卖家同意了您的退款申请~"
	MESSAGE_TYPE_AFTERSALES_BUYER_3  = "您有一笔退款到账啦~"
	MESSAGE_TYPE_AFTERSALES_SELLER_2 = "买家发起了退款申请~"

	NOTIFY_USER_FAIL              = "通知用户失败"
	NOTIFY_USER__SHIPPING_SUCCESS = "提醒发货成功"

	//ALI_AUDIT_ORIGIN_1 = "1" //发布商品
	//ALI_AUDIT_ORIGIN_2 = "2" //留言
	//ALI_AUDIT_ORIGIN_3 = "3" //评价
	//ALI_AUDIT_TYPE_1   = "1" //文本
	//ALI_AUDIT_TYPE_2   = "2" //图片
	//ALI_AUDIT_TYPE_3   = "3" //视频
	NINGMENG_SHOP = "柠檬商城"

	EXPRESS_QUERY_ERROR_CODE_999 = "204999" //请求快递接口出错
	EXPRESS_QUERY_ERROR_CODE_998 = "204998" //查询缓存数据库出错
	//底下错误返回注意：接口返回的不是很准确，基本返回的是EXPRESS_ERROR_CODE_3
	//所以在提交物流信息的时候无法正确判断是公司错误还是运单号错误(这个判断不做)
	EXPRESS_QUERY_ERROR_CODE_1 = "204301" //快递公司错误
	EXPRESS_QUERY_ERROR_CODE_2 = "204302" //运单号错误
	EXPRESS_QUERY_ERROR_CODE_3 = "204303" //查询失败
	EXPRESS_QUERY_ERROR_CODE_4 = "204304" //查不到物流信息
	EXPRESS_QUERY_ERROR_CODE_5 = "204305" //寄件人或收件人手机尾号错误

	EXPRESS_QUERY_ERROR_MSG_999 = "提交物流信息失败" //请求快递接口出错
	EXPRESS_QUERY_ERROR_MSG     = "快递单号或快递公司有误，请核实"

	ORDER_EXPRESS_QUERY_NOT_FOUND = "订单不存在"
)
