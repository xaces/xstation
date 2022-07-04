# [接入站点](https://github.com/xaces/xstation)

## 编译

 * `go.mod`删除自定义引入package

```go
    replace	github.com/xaces/xproto => ../../xproto
    replace	github.com/xaces/xutils => ../xutils
```

 * 由于ttx/ho协议暂未开放，在`app/access/xproto.go`中删除

```go
	_ "github.com/xaces/xproto/ttx"
	_ "github.com/xaces/xproto/ho"
```

## 子服务
 - <a target="_blank" href="https://github.com/xaces/xstream">媒体服务</a>
 - <a target="_blank" href="https://github.com/xaces/xstorage">存储服务</a>
 - <a target="_blank" href="https://github.com/xaces/xdownload">下载服务</a>