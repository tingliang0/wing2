wsserver demo
1. 服务器和客户端使用json进行数据交换，服务器需要通信的结构实现Marshal/unMarshal（struct <-> json）
2. 服务器使用boltDB作为持久化，需要持久化的结构实现Serialize/Deserialize（struct <-> []byte）
