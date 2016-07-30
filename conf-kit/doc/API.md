# restful api

## 增加配置

```
POST /conf/dev
{
	"Host": "127.0.0.1",
	"Port": "7200"
}
```

## 删除配置

```
DELETE /conf/dev
```

## 修改配置

```
PUT /conf/dev
{
	"Host": "127.0.0.1",
	"Port": "7200"
}
```

## 获取配置

```
GET /conf/dev
```

或

```
GET /conf/dev?hash=xxxxx
```

## 列出配置项

```
GET /list/dev/
```

或

```
GET /list/dev/?recursive=1
```

