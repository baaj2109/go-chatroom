
project structure
|
| - cmd 
|    | - main.go
|    | - chatroom
| - go.mod
| - go.sum
| - logic
|    | - brodacast.go
|    | - message.go
|    | - user.go
| - server
|    | - handle.go
|    | - home.go
|    | - websocket.go
--- template
     | - home.html

cmd: 用於儲存main.main
logic: 用於儲存專案核心業務邏輯城市 和 server作用類似
server: 純存server相關程式 類似於controller
template： 儲存靜態範本