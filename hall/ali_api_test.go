package hall

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"game_server/easygo"
	"game_server/for_game"
	"game_server/pb/share_message"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/iGoogle-ink/gopay"
	"github.com/iGoogle-ink/gopay/alipay"
)

func TestAliAPILogin(t *testing.T) {
	//PAliApiMgr.AliAPILogin(nil, "f73b4c0144fc4e2d940002fa01f4WX01", nil)
	//result := `{"error_response":{"code":"20001","msg":"Insufficient Token Permissions","sub_code":"aop.invalid-app-auth-token","sub_msg":"无效的应用授权令牌"},"sign":"WnrZKyL/pSa3Ve8fN4q0/ER6zXMhSE6bPmZklZHV618bo5JIEfm+Hdhkb2cDiSZHDrXrm/a9ObL9cO047NNjxtNcbTUufX+jfjQP2mLvIJjT5YjjNlInVCC/NYPScJWTvq2RwHrxhYEe34sPmtseeoCBqMgmmGWStHIyEC45YhmAJALJXXjwfloLIR0I9+/PJPsrDXDOmwqKbsNOnUOQbS4YXz1oIfjrvOa4stGYrpcv/Yl/foaG/mDduHthaOrB0PQO1ViNkVcHTV0L2dWKMOEi5fB16JKpVEScQskZ1vmfcx+bjRPc+cwb+b1olbbNKbVO7QdvMsOPt35d0WjS2w=="}`
	re := `{"alipay_system_oauth_token_response":{"access_token":"authusrB94739c226e9c400196af5bc4740a9X01","alipay_user_id":"20880063008624437006283650118501","expires_in":1296000,"re_expires_in":2592000,"refresh_token":"authusrBff1f0be8cae1436dace07e6e7cd7cX01","user_id":"2088902487861012"},"sign":"W/ZQgCAX2ExObxWYWire/+vfGfRqiA4IIBZC9oUzNLdHF+rCC2aZN1ERqymUomqBVjaFsFSB0S0XzqbTFJNf+JkQpwvYwJr0mCudA2mjolpnD/jIZJy+to/w20Kw/wp6vvdIcZ/pXGULi2s98Ndu1mnmC3phZNfNGgAUHnzG1QVeo2ap/yT9iul04oMsiRHddcAWHKY44SE8Bpj4k6MdqBufuBeC4fbfbPDH/4UfXjN3Dj6FJNAjy5n9+gofXnzECOP0Jvj63I12/ag273d1FJYLVppCuyw+TUGbnEOiGh7ZtNsDWnBmZRPv2uG7DvsiogEEuZYM0Fm3X1fUMLoDPw=="}`

	var r *share_message.AliLoginResult
	err := json.Unmarshal([]byte(re), &r)
	if err != nil {
		logs.Info("err: %s", err.Error())
		return

	}
	fmt.Println("----->", r.AlipaySystemOauthTokenResponse.GetUserId())
}

func TestAliPayOrder(t *testing.T) {
	s := `{"PlayerId":1887436001,"Amount":"0.01","PayId":7,"PayWay":1,"PaySence":1,"ProduceName":"测试测试","PayType":2}`
	info := &share_message.PayOrderInfo{}
	if err := json.Unmarshal([]byte(s), info); err != nil {
		logs.Info("----->%s", err.Error())
		return
	}
	b, _ := PAliApiMgr.AliAPILogin(WebServreMgr, "e23db62ad6c64b2781e9507e76edWX01", info)
	var m map[string]interface{}
	err := json.Unmarshal(b, &m)
	if err != nil {
		logs.Info("----->%s", err.Error())
		return
	}
	logs.Info("===============> %+v", m)
}

func TestUrlDesc(t *testing.T) {
	s := `{"PlayerId":1887436002,"Amount":"0.01","PayId":2,"PayWay":1,"PayScene":2,"ProduceName":"充钱","PayType":2,"PayTargetId":1887436002,"Content":"给服务器发充钱信息"}`

	data := url.QueryEscape(base64.StdEncoding.EncodeToString([]byte(s)))
	//data := "eyJQbGF5ZXJJZCI6MTg4NzQzNjAwMiwiQW1vdW50IjoiMC4wMSIsIlBheUlkIjoyLCJQYXlXYXkiOjEsIlBheVNjZW5lIjoyLCJQcm9kdWNlTmFtZSI6IuWFhemSsSIsIlBheVR5cGUiOjIsIlBheVRhcmdldElkIjoxODg3NDM2MDAyLCJDb250ZW50Ijoi57uZ5pyN5Yqh5Zmo5Y+R5YWF6ZKx5L+h5oGvIn0="
	enEscapeUrl, e := url.QueryUnescape(data)
	easygo.PanicError(e)
	if len(enEscapeUrl) == 0 {
		return
	}
	enEscapeUrl = strings.Replace(enEscapeUrl, " ", "", -1)
	decodeBytes, err := base64.StdEncoding.DecodeString(enEscapeUrl)
	if err != nil {
		logs.Info("-------> %s", err.Error())
		return
	}
	logs.Info("decodeBytes----------->%s", string(decodeBytes))
}

func TestAliPayQuery(t *testing.T) {
	appID := "2016102700770128"
	//privateKey := `MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCN8VHpXP1ZSJSUZeQtZZNDQrHqr88nn7A+NmeKrsuu5PkQ2CvOc0deBLbL42Znz1PqaVK3ioHDOCTJHrX+xBsmDfhqE0s2QUoJDWrhqGw+HRQ8anlgUQ2C6UYKsII26jDYPFywtB86ZrykOZ8Cd+su0CmGjnE9htCOdjEql5jmmmTWdz8FZePLN7XNja8LbumVTG5TBzOornOm5Cz3IqrEqlJMMS+mO7yxpasb2JOXgrt+lyE3MRAv7wS9pCPgtwC0TveLN6FYe/lvb8kHnq9ROq0phBJlMdHpZ79CT41PpQoNTAxVLN5bFmbdO6syjgn0pHRLvLIBaiLprVdqRnSfAgMBAAECggEAeebyjhSKkI9A62HGYSaHHpC88+0hX8pJNmTK79PGoeGL9edxV9CxThGGW/xkCmuIih0CKRcO8nXZQdDaRH5vQnNlENSZF3Ni/ftD+6EFtSKMKobWzt1NWUy2FqAYdMkUQeE1SZyn5SQuhmvmH9yVYpLr1t+maUzK+E6RUx729bQCb2OMkfP8tZtfcMLWiFwk6U+p9zo7WlIzxJqpRyu2KodQtlo94Qsx0X0EIcHkR72MqdTOCdC7KcWd04FMSW9NPfLEejykbpKiF9T8HR1mK2dC0em5eSTLsYcknduWUgiBD3S2gM9qWFIDcWFChcV5yTs73iNULUHcRYJuMvWjYQKBgQDQ92MuQFz4CEFLretNSx2v8W22cwhIFBkGBYSLo1yAbV2UxOuSmi9JVkEP6My3iETzDFGh9tXAQxbfNqFEcwQ8t423maNvhqybU47nJn/F84uxNuFZOH9fK14RD9QSRpWtwnTTZoOGPlR7U7+qLqFLnEICLoUVR8UP8tp9avr3dQKBgQCt5AlCP79DApWs20Vqoh9HQqjnw+LRXHMpFTfnw5F9Qa3g7N5v94osKbZE7oCdqrIuzjboEfNnrKd/x6y6TSuYIquECYQwv+4ci/RLhOjTYnY/NCErq83qoGfuUxWGEkoX1KE7yB56OeuC6MdJ9Q+Ms4LJ6NLh/CaavqXu6dPNQwKBgET1TnKF5Ogo+Ts7Moo4PpzAJD9wKIx4rWVSTtIx36W18YrVjRO889vUrfXNEjmCq5Y1O38iUJl4ykRw57kJ550NyaOL/OYh4DYF1gOrrcCqRS/+91CVF1tVmV4yBf7d8ij8IcddbgvP59sm4PoNF0c3UoUbyukh3QMNVlLLCfS9AoGBAJQQ5EFg/n8UqFYzr3wI6BFJlYEjrvMOgZCt3JigUjYRwvkPOKimYyUPr4AqhaG7Q1XPibk578SLo2SOpWlNZJ16iAk6ATFxfFMaaL4VQhsccAuJW+VPuVrbkyO/40fyMtzv1QqOcEUrJHqns2oqHT91axx5/3cluclyJOC2gf75AoGAChJrZtH2lYsB9XFLHQnr1YYe/AFWmDNNWYIySvoNdduxbamdmFMtqXQDtu9G328EcZE74p6RdhXRlJe8hkwxPBDdR6VZfMvlEjpKtSnzQWwG7tYE5Zn5NP5LJ/YPPenwQMXM099kwFo9Z3YxMjM7o6D72wu27JGiw19BhkPoFuQ=`
	client := alipay.NewClient(appID, for_game.AliPrivateKeyNoPreEnd, false)

	// 初始化 BodyMap
	bm := make(gopay.BodyMap)
	bm.Set("trade_no", "2020072022001461011420303814")
	//bm.Set("query_options", "TRADE_SETTLE_INFO")
	client.PrivateKeyType = 2
	query, err := client.TradeQuery(bm)
	if err != nil {
		logs.Info("err------->", err.Error())
		return
	}
	logs.Info("------------>%+v", query)
}

func TestName(t *testing.T) {
	PAliApiMgr.TradeQuery("", 2*time.Second)
}
