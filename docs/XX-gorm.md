## ドメインモデルの再設計

学習用なので UserDomain と UserModel を分ける

ファクトリ関数 NewUser(name, email, password string) で必須項目を検証
ID はドメインでは扱わず、DB 側（SQLite）で管理する設計に変更

## リポジトリ設計の変更

メモリ実装から SQLite DB 実装に変更
DB 用モデル UserModel はリポジトリ層に別途定義し、ドメインモデルと変換して使う
GORM を使い、SQLite で自動的に ID（主キー）を生成・管理
Save() では UserModel に変換して DB に保存し、DB で ID が自動生成される
FindByEmail() などの検索メソッドも実装
