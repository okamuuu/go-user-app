## コマンド

json yml は生成されている

```
% swag init  -g cmd/main.go -o cmd/docs/
2025/06/04 23:51:20 Generate swagger docs....
2025/06/04 23:51:20 Generate general API Info, search dir:./
2025/06/04 23:51:20 warning: failed to get package name in dir: ./, error: execute go list command, exit status 1, stdout:, stderr:no Go files in /Users/okamuuu/Prj/go-user-app
2025/06/04 23:51:20 create docs.go at cmd/docs/docs.go
2025/06/04 23:51:20 create swagger.json at cmd/docs/swagger.json
2025/06/04 23:51:20 create swagger.yaml at cmd/docs/swagger.yaml
```

ただし、謎の warning が出ている

```
2025/06/04 23:51:20 warning: failed to get package name in dir: ./, error: execute go list command, exit status 1, stdout:, stderr:no Go files in /Users/okamuuu/Prj/go-user-app
```
