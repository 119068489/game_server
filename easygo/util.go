//Package easygo ...
package easygo

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"reflect"
	"regexp"
	"time"

	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func init() {
	rand.Seed(int64(time.Now().Nanosecond()))
}

//GetCurrentDirectory 获取当前程序运行的目录
// 暂时没有任何地方用得到。留着吧
func GetCurrentDirectory() (string, string) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalf("error GetCurrentDirectory: %v\n", err)
	}
	dir = strings.Replace(dir, "\\", "/", -1)
	idx := strings.LastIndex(dir, "/")
	parentPath, filename := dir[:idx], dir[idx+1:]
	return parentPath, filename
}

var instanceKey = "!@#$%~^^&&*&^%#!!~"

//MakeTokenForInstance 生成一个token
func MakeTokenForInstance(instanceid int32, playerid int64) string {
	strInstanceID := string(instanceid)
	strPlayerID := string(playerid)

	tokenMd5 := md5.New()
	tokenMd5.Write([]byte(strInstanceID + strPlayerID + instanceKey))

	token := hex.EncodeToString(tokenMd5.Sum(nil)) //
	return token
}

const LetterBytes = "123456789abcdefghijklmnopqrstuvwxyzABCDEFJHIJKLMNOPQRSTUVWXYZ"

//生成随机字符串
func RandStringRunes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = LetterBytes[rand.Intn(len(LetterBytes))]
	}
	return string(b)
}

//float64类型取两位有效小数
func Decimal(value float64, num ...int) float64 {
	n := append(num, 4)[0]
	format := "%." + strconv.Itoa(n) + "f"
	value, _ = strconv.ParseFloat(fmt.Sprintf(format, value), 64)
	return value
}

// "fmt"
// "reflect"
// "time"
// "game_server/easygo"
// "unsafe"

// header1 := (*reflect.SliceHeader)(unsafe.Pointer(&s1))
// header2 := (*reflect.SliceHeader)(unsafe.Pointer(&s2))

// fmt.Println("header1,header2=",header1.Data==header2.Data )
//Equals
func Equals(s1 []interface{}, s2 []interface{}) {
	// header1 := (*reflect.SliceHeader)(unsafe.Pointer(&s1))
	// header2 := (*reflect.SliceHeader)(unsafe.Pointer(&s2))
	//header1.Data == header2.Data
}

type Foo struct {
}

//-----------------------
/*
func main() {
	o1 := &Foo{}
	x := unsafe.Pointer(o1)
	y := uintptr(x)

	o2 := o1
	j := unsafe.Pointer(o2)
	k := uintptr(j)
	fmt.Println(reflect.TypeOf(x))
	fmt.Println(reflect.TypeOf(y))
	fmt.Println("o1,o2=",o1,o2 )
	fmt.Println("x,y=",x,y )
	fmt.Println("j,k=",j,k )


	fmt.Println("o1==o2",o1==o2 )
	fmt.Println("x==j",x==j )
	fmt.Println("y==k",y==k )

	oo3 := &o1
	fmt.Println(reflect.TypeOf( oo3))
	test(oo3)

	// oo4 := uintptr(o1)

	s1 := []int{4}
	s2 := []int{4}





	//fmt.Println("s1,s2=",&s1,&s2 )
	//fmt.Println("unsafe(s1),unsafe(s2)=",unsafe.Pointer(&s1),unsafe.Pointer(&s2))
	//fmt.Println("s1==s2",&s1==&s2 )

}
func test(aa **Foo) {

}
*/
// --------------------------------

// 信号量
type Semaphore chan struct{}

func (self Semaphore) Acquire() {
	self <- struct{}{}
}
func (self Semaphore) Release() {
	<-self
}
func NewSemaphore(size int) Semaphore {
	// 原则上，value值不能为负数
	return make(Semaphore, size)
}

// --------------------------------

func IsInterfaceNil(i interface{}) bool {
	vi := reflect.ValueOf(i)
	if vi.Kind() == reflect.Ptr {
		return vi.IsNil()
	}
	return i == nil
}

//去空格
func CompressStr(str string) string {
	if str == "" {
		return ""
	}
	//匹配一个或多个空白符的正则表达式
	reg := regexp.MustCompile("\\s+")
	return reg.ReplaceAllString(str, "")
}

func IsEmptyStr(str interface{}) bool {
	if IsInterfaceNil(str) {
		return true
	}

	v := reflect.ValueOf(str)
	if v.Kind() == reflect.Ptr {
		str = *v.Interface().(*string)
	} else if v.Kind() == reflect.String {
		return CompressStr(v.String()) == ""
	}

	return CompressStr(fmt.Sprint(str)) == ""
}

//读取路线json配置表
func ReadJsonConfig(filename string, Data interface{}) {
	bytes, err := ioutil.ReadFile(filename)
	PanicError(err)

	err = json.Unmarshal(bytes, Data)
	PanicError(err)
}

func If(b bool, trueBack interface{}, falseBack interface{}) interface{} {
	if b {
		return trueBack
	} else {
		return falseBack
	}
}

func StringToInt(str string) int {
	i, err := strconv.Atoi(str)
	PanicError(err)
	return i
}

func StringToIntnoErr(str string) int {
	i, _ := strconv.Atoi(str)
	return i
}

func StringToInt64noErr(str string) int64 {
	i, _ := strconv.ParseInt(str, 10, 64)
	return i
}

func IntToString(i int) string {
	return strconv.FormatInt(int64(i), 10)
}

// 根据权重随机返回 key,例如  dict = {1:25, 2:20, 3:10, 4:30, 5:15}
// 未命中返回 nil,传入空字典返回 nil
// 类似 NBA 抽签选秀的方式，dict 的单元中的 value 值越高，概率就越高
func RandomByWeitht(dict map[int32]int32, totals ...int32) int32 {
	var total int32 = append(totals, 0)[0]
	if total == 0 {
		for _, v := range dict {
			total += v
		}
		if total <= 0 {
			s := fmt.Sprintf("传的权重表是什么鬼，%v", dict)
			panic(s)
		}
	}
	i := rand.Int31n(total)
	var t int32 = 0
	for k, v := range dict {
		t += v
		if i < t {
			return k
		}
	}
	return math.MaxInt32
}

//-----------------------------------------------------------------

// protected call.会把错误+调用栈写进磁盘
func PCall(function interface{}, args ...interface{}) bool {
	r, tb := PCallImp(function, args...)
	if r != nil {
		LogException(r, tb)
		return false
	}
	return true
}

// protected call
func PCallImp(function interface{}, args ...interface{}) (r interface{}, tb string) {
	defer func() {
		if r = recover(); r != nil {
			tb = CallStack()
		}
	}()

	if function == nil {
		panic("function 参数是个 nil")
	}
	f := reflect.ValueOf(function)
	if f.Kind() != reflect.Func {
		panic("function 参数必须是个 function")
	}

	in := make([]reflect.Value, len(args))
	i := 0
	for _, arg := range args {
		in[i] = reflect.ValueOf(arg)
		i++
	}

	f.Call(in)
	return r, tb
}

//-----------------------------------------------------------------

// 用于解析列表类型的命令行启动参数
type SliceValue []string

func NewSliceValue(p *[]string) *SliceValue {
	return (*SliceValue)(p)
}

func (self *SliceValue) Set(val string) error {
	*self = SliceValue(strings.Split(val, ","))
	return nil
}

func (self *SliceValue) String() string {
	return "" // 暂时没有搞懂这个函数什么时候会调用到
}

func IsPhoneStr(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}

/*营收计算器
amount=6000000	总额（需要计算的总额）
med=2000000		阶梯值（目前默认200万）
ratio=0.08 		起始比率（目前默认0.08）
times=0 		已计算次数（默认0开始）
result=0 		结果值（默认0开始）
*/
func SumRevenue(amount float64, med float64, ratio float64, times float64, result float64) float64 {
	if amount < med || times > 2 {
		result += amount * (ratio - times*0.01)
	} else {
		result += med * (ratio - times*0.01)
		times += 1
		amount = amount - med
		if amount > 0 {
			result = SumRevenue(amount, med, ratio, times, result)
		}
	}

	return result
}

func GetMapKeysValues(mapI interface{}) (keys interface{}, values interface{}) {
	typeOf := reflect.TypeOf(mapI)
	k := typeOf.Kind()
	if k != reflect.Map {
		panic("仅支持map")
	}
	valueOf := reflect.ValueOf(mapI)
	var kks []interface{}
	var vvs []interface{}
	for _, kkss := range valueOf.MapKeys() {
		kks = append(kks, kkss.Interface())
		vvs = append(vvs, valueOf.MapIndex(kkss).Interface())
	}
	return kks, vvs
}

//字符串数组拼接字符串
func StringArrayToString(str []string) string {
	return strings.Replace(strings.Trim(fmt.Sprint(str), "[]"), " ", ",", -1)
}

//int64数组拼接字符串
func Int64ArrayToString(str []int64) string {
	return strings.Replace(strings.Trim(fmt.Sprint(str), "[]"), " ", ",", -1)
}

//任意数组拼接成字符串
func ArrayToString(str []interface{}) string {
	return strings.Replace(strings.Trim(fmt.Sprint(str), "[]"), " ", ",", -1)
}

//=============================二进制计算常用方法======================只对8位数据有效，所以处理时要先转uint8
//把第pos位值设置为1
func SetBiuOne(val uint8, pos uint8) uint8 {
	result := val | (1 << (pos - 1))
	return result
}

//把第pos位值设置为0
func SetBiuZero(val uint8, pos uint8) uint8 {
	result := val &^ (1 << (pos - 1))
	return result
}

//取反第pos位值:即0变1或1变0
func OppositeBiuValue(val uint8, pos uint8) uint8 {
	result := val ^ (1 << (pos - 1))
	return result
}

//获取某一位的值
func GetBiuValue(val uint8, pos uint8) uint8 {
	result := (val << (8 - pos)) >> 7
	return result
}

//====================二进制计算常用方法====================
