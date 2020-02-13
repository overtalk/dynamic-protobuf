# golang dynamic protobuf

`golang` 中如何不使用 `protoc` 先对 `.proto` 文件进行编译，而直接使用 `.proto` 文件对 `protobuf` 进行序列化/反序列化

## 简介
- 主要借助于 [protoreflect库](https://github.com/jhump/protoreflect)
- 首先对 `.proto` 文件进行解析，得到 `FileDescriptor`,`MessageDescriptor`,`EnumDescriptor`
- 可以根据 `MessageDescriptor` 构造出 `dynamic message`, 它完全实现了官方的 `proto.Message` 接口，可以使用官方库直接对它进行序列化/反序列化

## Tips
- 上述的 `dynamic` 软件包([protoreflect库](https://github.com/jhump/protoreflect))提供了动态消息实现。它实现proto.Message但由消息描述符和fields-> values映射（而不是生成的结构）支持。这对于一般地使用协议缓冲区消息有用，而不必为每种消息生成并链接Go代码。这对于需要在任意协议缓冲区模式上运行的通用工具特别有用。通过使工具在运行时加载描述符，可以做到这一点。