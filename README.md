# proxy
golang  web serve and proxy
由于业务需要实现对多个web应用做同域二级目录代理，用NGINX的又感觉太重了，而且不好做配置页面，用golang来实现代理功能

- 支持正则表达式匹配机制
- 支持多应用多级目录代理。
- 支持应用子路由代理
- 支持webapi代理
- 支持websocket代理
- 支持禁用缓存设置
- 支持本地文件服务
- 支持http、https混合使用
- 支持/dir/app 重定向为 /dir/app/
- 支持简单的路由热度升级