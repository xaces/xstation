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