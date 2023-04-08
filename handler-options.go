package main

type HandlerType int

const (
	/* 本地文件系统 */
	LocalFileSystem HandlerType = 10000
	/* 远程Http服务 */
	RemoteHttpServer HandlerType = 20000
)

type HandlerOptions struct {
	Id int32
	/* 处理器名称 */
	Name string
	/* 处理器类型 */
	Type HandlerType
	/* 处理路径 */
	Path string
	/* 禁用缓存 */
	DisableCache bool
	/* 目标服务器,或本地文件目录 */
	Target string
	/* 子处理器列表 */
	SubHandlers []*HandlerOptions
	/* 是否启用 */
	Enabled bool
}

var handlers = []*HandlerOptions{
	{
		Id:           1,
		Name:         "app1",
		Path:         "^/app1",
		Type:         RemoteHttpServer,
		Target:       "https://192.168.1.50:12580/app1",
		Enabled:      true,
		DisableCache: true,
		SubHandlers: []*HandlerOptions{
			{
				Id:           3,
				Name:         "登录接口",
				Path:         "^/app1/login",
				Type:         RemoteHttpServer,
				Target:       "http://192.168.1.50:8500/login",
				Enabled:      true,
				DisableCache: true,
				SubHandlers:  []*HandlerOptions{},
			},
			{
				Id:           4,
				Name:         "API",
				Path:         "^/app1/api/v1/softlog/last/",
				Type:         RemoteHttpServer,
				Target:       "https://192.168.1.60:8000/api/v1/softlog/last/",
				Enabled:      true,
				DisableCache: true,
				SubHandlers:  []*HandlerOptions{},
			},
			{
				Id:           5,
				Name:         "后台API",
				Path:         "^/app1/api/",
				Type:         RemoteHttpServer,
				Target:       "http://192.168.1.50:8500/api/",
				Enabled:      true,
				DisableCache: true,
				SubHandlers:  []*HandlerOptions{},
			},
			{
				Id:           6,
				Name:         "WEBSOCKET",
				Path:         "^/app1/websocket/",
				Type:         RemoteHttpServer,
				Target:       "http://192.168.1.50:8500/websocket/",
				Enabled:      true,
				DisableCache: true,
				SubHandlers:  []*HandlerOptions{},
			},
			{
				Id:           7,
				Name:         "后台",
				Path:         "^/app1/power/",
				Type:         RemoteHttpServer,
				Target:       "http://192.168.1.50:8500/power/",
				Enabled:      true,
				DisableCache: true,
				SubHandlers:  []*HandlerOptions{},
			},
		},
	},
	{
		Id:           2,
		Name:         "本地测试",
		Path:         "^/test",
		Type:         LocalFileSystem,
		Target:       "/root/workspace/go-admin/static/www",
		Enabled:      true,
		DisableCache: true,
		SubHandlers: []*HandlerOptions{
			{
				Id:           5,
				Name:         "api",
				Path:         "^/test/api/",
				Type:         RemoteHttpServer,
				Target:       "https://192.168.1.100:8000/api/",
				Enabled:      true,
				DisableCache: true,
				SubHandlers:  []*HandlerOptions{},
			},
			{
				Id:           5,
				Name:         "静态文件",
				Path:         "^/test/static/",
				Type:         LocalFileSystem,
				Target:       "/root/workspace/go-admin/static",
				Enabled:      true,
				DisableCache: true,
				SubHandlers:  []*HandlerOptions{},
			},
			{
				Id:           5,
				Name:         "websocket相关",
				Path:         "^/test/ws",
				Type:         RemoteHttpServer,
				Target:       "https://192.168.1.100:8000/ws",
				Enabled:      true,
				DisableCache: true,
				SubHandlers:  []*HandlerOptions{},
			},
		},
	},
}
