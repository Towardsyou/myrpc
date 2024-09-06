# my rpc

## 主要内容

- 网络层: 建立TCP链接，收发网络包
- 协议层
  - 序列化
  - 压缩
  - 加密
- 应用层：通过反射修改方法，实现本地调用转发远程服务

## 实现细节

### 协议层

1. 读取前 8 个字节，获得包长度，读取整个包
2. 读取包头，并根据包头解码处理包体
3. 其中处理顺序为：解密 -> 解压缩 -> 解序列化
4. 根据包头信息转发给对应服务实现

### 应用层

- 客户端需要调用 `InitService` 方法，修改映射方法到远程服务实现
- 服务端需要调用 `RegisterService` 方法，注册服务名，方法名到方法的映射

## 例子

参见 examples/helloworld

## Todo

-[ ] 池化常用资源：Request Response
-[ ] 支持多种序列化方式：protobuf
-[ ] 支持多种压缩方式：gzip
-[ ] 支持多种压缩方式：zstd
-[ ] 区分 oneway 和 normal 调用
-[ ] 通过 meta 传递超时
-[ ] 通过 mock 工具，增加测试
  

## Reference

- [极客时间：-Go实战训练营](https://gitee.com/geektime-geekbang/geektime-go)