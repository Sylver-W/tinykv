## Tinykv Project1 Standalone Kv

### 官方文档重点

#### 列族（**CF Column family**）

Column Family，也叫 CF，这个概念从 HBase 中来，就是将多个列合并为一个CF进行管理。这样读取一行数据时，你可以按照 CF 加载列，不需要加载所有列（通常同一个CF的列会保存在同一个文件中，所以这样有很高的效率）。此外因为同一列的数据格式相同，你可以针对某种格式采用高效的压缩算法。

简单来说，每一个列族对应一个类型的信息，包含多个键值存储对，需要获取键值信息时，只需要加载对应列族即可

由于Tinykv的底层badger并不支持列族，这里使用“CF_Key”作为key存储来实现列族，engine_util中有函数`KeyWithCF`实现，如下

```go
func KeyWithCF(cf string, key []byte) []byte {
	return append([]byte(cf+"_"), key...)
}
```

#### gRPC服务

Project1构建的就是基于列族的独立键值存储gRPC服务

**gRPC服务**，即google Remote Procedure Call，是由 Google 开发的一种**高性能、开源和通用的远程过程调用（RPC）框架**，用于在分布式系统中进行服务间通信。它允许客户端像调用本地方法一样调用远程服务器上的方法。

tinykv中的proto部分就是gRPC相关的代码，tinykv是一个分布式系统，需要通信方式，就需要使用gRPC，这里面定义了gRPC接口

其中重点关注message和service，message定义传输数据结构，例如在kvrpcpb.proto中，有

```protobuf
message RawGetRequest {
    Context context = 1;
    bytes key = 2;
    string cf = 3;
}
```

这里的 =1是每个字段的编号，在传输时更高效编码

在tinykvpb.proto文件中，有类似于

```protobuf
service KV {
  rpc Get(GetRequest) returns (GetResponse);
}
```

定义了一个函数，当接收到GetRequst需要返回GetResponse

proto可以生成相应语言的代码，剩下需要做的就是调用或实现相应的函数或结构就行，例如在project1中，我们就实现了RawGet等方法，这在gRPC中已经封装完成，proto定义了该接口的方法，我们做的是实现

通常情况下，你不需要修改 proto 文件，因为所有必要的字段都已经为你定义好了。但如果仍然需要修改，你可以修改 proto 文件并运行`make proto`以更新`proto/pkg/xxx/xxx.pb.go`中生成的相关 go 代码

