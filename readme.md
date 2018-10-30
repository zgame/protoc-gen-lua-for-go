# protoc-gen-lua-for-go干什么用的？

	在lua代码中使用protocolbuffer，官方不支持， 需要使用第三方工具proto-gen-lua， 其中pb.c是c语言版本， 本例中pb.go为go语言版本，另外支持了int64，uint64


# 例子demo

	demo目录下面有一个例子， 可以运行查看protobuf的输入和输出
	入口点都是main,  包括main.go, main.lua
	build目录下面可以编辑和修改proto文件， 编译成lua的pb文件之后， protocol_test.lua是具体的protobuf调用文件

# 流程

	build目录下面修改proto文件，然后build.bat(需要python环境) ，protocol_test.lua调用
	在goland IDE打开目录， 运行的时候选择运行整个目录即可





