package server

import (
	"github.com/travisjeffery/proglog/internal/config"
	"google.golang.org/grpc/credentials"
)

t.Helper()

l, err := net.Listen("tcp", "127.0.0.1:0")
require.NoError(t, err)

clientTLSConfig, err := config.setupTLSConfig(config.TLSConfig{
	CertFile: config.ClientCertFile,
	KeyFile:  config.ClientKeyFile,
	CAFile: config.CAFile,
})
require.NoError(t, err)

clientCreds := credentials.NewTLS(clientTLSConfig)
cc, err := grpc.Dial(
	l.Addr().String(),
	grpc.WithTransportCredentials(clientCreds),
)
require.NoError(t, err)

client = api.NewLogClient(cc)

serverTSLConfig, err := config.setupTLSConfig(config.TLSConfig{
	CertFile: config.ServerCertFile,
	KeyFile:  config.ServerKeyFile,
	CAFile:   config.CAFile,
	ServerAddress: l.Addr().String(),
	Server: true,
})
require.NoError(t, err)
serverCreds := credentials.NewTLS(serverTLSConfig)

dir, err := os.MkdirTemp("", "server-test")
require.NoError(t, err)

clog, err := log.NewLog(dir, log.Config{})
cfg = &Config{
	CommitLog: clog,
}
if fn != nil {
	fn(cfg)
}

server, err := NewGRPCServer(cfg, grpc.Creds(serverCreds))
require.NoError(t, err)

go func() {
	server.Serve(l)
}()

return client, cfg, func() {
	cc.Close()
	server.Stop()
	l.Close()
}