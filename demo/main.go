package main

import (
	"github.com/yuin/gopher-lua"
	"fmt"
	"strconv"
	"time"
	"os"
	"bufio"
	"github.com/yuin/gopher-lua/parse"
	"unsafe"
	"runtime"
)

var num = 0
var codeToShare *lua.FunctionProto
func main() {
	runtime.GOMAXPROCS(1)

	//defer luaPool.Shutdown()
	var err error
	codeToShare,err = CompileLua("main.lua")
	if err!=nil{
		fmt.Println("加载main.lua文件出错了！")
	}

	go start(1)
	//go start(2)


	//goCallLua(L)

	for{
		select {

		}
	}


}

// 把指针传递过去给dll
func IntPtr(L *lua.LState) uintptr {
	return uintptr(unsafe.Pointer(L))
}
func start(timer time.Duration) {


	L := lua.NewState()


	//L := lua.NewState()
	defer L.Close()


	//L := luaPool.Get()
	//defer luaPool.Put(L)


	// 直接调用luaopen_pb
	luaopen_pb(L)

	// Lua调用go函数声明
	// 声明double函数为Lua的全局函数，绑定go函数Double
	////L.SetGlobal("double", L.NewFunction(Double))
	//L.Register("double", Double)
	//DoCompiledFile(L, codeToShare)

	//// 执行lua文件
	if err := L.DoFile("main.lua"); err != nil {
		fmt.Println("加载main.lua文件出错了！")
		fmt.Println(err.Error())
	}

	// 通过dll加载luaopen_pb
	//DllTestDef := syscall.MustLoadDLL("libpb.dll")
	//add := DllTestDef.MustFindProc("luaopen_pb")
	//ret, _, err := add.Call(IntPtr(L))
	//if err!=nil{
	//	fmt.Println("返回",ret)
	//}





	tickerCheckUpdateData := time.NewTicker(time.Second * timer)
	defer tickerCheckUpdateData.Stop()

	for{
		select {
		case <-tickerCheckUpdateData.C:
			timerFunc(L,timer)

		}
	}

}



//------------------编译lua文件------------------------------

// CompileLua reads the passed lua file from disk and compiles it.
func CompileLua(filePath string) (*lua.FunctionProto, error) {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)
	chunk, err := parse.Parse(reader, filePath)
	if err != nil {
		return nil, err
	}
	proto, err := lua.Compile(chunk, filePath)
	if err != nil {
		return nil, err
	}
	return proto, nil
}

// DoCompiledFile takes a FunctionProto, as returned by CompileLua, and runs it in the LState. It is equivalent
// to calling DoFile on the LState with the original source file.
func DoCompiledFile(L *lua.LState, proto *lua.FunctionProto) error {
	lfunc := L.NewFunctionFromProto(proto)
	L.Push(lfunc)
	return L.PCall(0, lua.MultRet, nil)
}


//// ---------------------lua 文件池-------------------------------------------
//type lStatePool struct {
//	m     sync.Mutex
//	saved []*lua.LState
//}
//
//func (pl *lStatePool) Get() *lua.LState {
//	pl.m.Lock()
//	defer pl.m.Unlock()
//	n := len(pl.saved)
//	if n == 0 {
//		return pl.New()
//	}
//	x := pl.saved[n-1]
//	pl.saved = pl.saved[0 : n-1]
//	return x
//}
//
//func (pl *lStatePool) New() *lua.LState {
//	L := lua.NewState()
//
//	// 执行lua文件
//	if err := L.DoFile("main.lua"); err != nil {
//		fmt.Println("加载main.lua文件出错了！")
//		fmt.Println(err.Error())
//	}
//	// setting the L up here.
//	// load scripts, set global variables, share channels, etc...
//	return L
//}
//
//func (pl *lStatePool) Put(L *lua.LState) {
//	pl.m.Lock()
//	defer pl.m.Unlock()
//	pl.saved = append(pl.saved, L)
//}
//
//func (pl *lStatePool) Shutdown() {
//	for _, L := range pl.saved {
//		L.Close()
//	}
//}
//
//// Global LState pool
//var luaPool = &lStatePool{
//	saved: make([]*lua.LState, 0, 4),
//}
//





//-------------计时器------------------------
func timerFunc(L *lua.LState,timer time.Duration)  {
	//fmt.Println("timer--------")
	//goCallLuaReload(L)
	//goCallLua(L,int(timer))

	num++
	//goCallLua(L)
}

// Lua重新加载，Lua的热更新按钮
func goCallLuaReload(L *lua.LState)  {
	//fmt.Println("----------lua reload--------------")
	if err := L.CallByParam(lua.P{
		Fn: L.GetGlobal("ReloadAll"), //reloadUp  ReloadAll
		NRet: 0,
		Protect: true,
	}); err != nil {
		fmt.Println("",err.Error())
	}
}

// go调用lua函数
func goCallLua(L *lua.LState, num int)  {
	fmt.Println("----------go call lua--------------")
	// 这里是go调用lua的函数
	if err := L.CallByParam(lua.P{
		Fn: L.GetGlobal("Zsw2"),
		NRet: 2,
		Protect: true,
	}, lua.LNumber(num),lua.LNumber(num)); err != nil {
		fmt.Println("---------------")
		fmt.Println("",err.Error())
		fmt.Println("----------------")
	}

	ret := L.Get(1) // returned value
	fmt.Println("lua return: ",ret)
	ret = L.Get(2) // returned value
	fmt.Println("lua return: ",ret)
	L.Pop(1)  // remove received value
	L.Pop(1)  // remove received value
}

//-----------------------------------------------------
//Type name	Go type	Type() value	Constants
//LNilType	(constants)	LTNil	LNil
//LBool	(constants)	LTBool	LTrue, LFalse
//LNumber	float64	LTNumber	-
//LString	string	LTString	-
//LFunction	struct pointer	LTFunction	-
//LUserData	struct pointer	LTUserData	-
//LState	struct pointer	LTThread	-
//LTable	struct pointer	LTTable	-
//LChannel	chan LValue	LTChannel	-
//-----------------------------------------------------


// Lua调用的go函数
func Double(L *lua.LState) int {
	lv := L.ToInt(1)             //第一个参数
	lv2 :=  L.ToInt(2)			 //第一个参数
	str := L.ToString(3)

	L.Push(lua.LString(str+"  call "+strconv.Itoa(lv * lv2))) /* push result */
	L.Push(lua.LString(str+"  call "+strconv.Itoa(lv * lv2))) /* push result */

	return 2                    /* number of results */
}

