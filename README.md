# grpc_auth


* 客户端调用
  
  通过grpc.WithPerRPCCredentials()设置认证选项。

* 服务端调用
   
   服务端需要通过调用authintercepter中NewUnaryServerInterceptor注册拦截器。比如：
   ```
   grpcserver := grpc.NewServer(
		grpc.Creds(c),
		grpc.UnaryInterceptor(serversdk.NewUnaryServerInterceptor),
	)
   