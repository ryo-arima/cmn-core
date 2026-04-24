# Test Data Directory

このディレクトリには、単体テスト用のテストデータが集約されています。

## 構造

```
testdata/
├── config/          # 設定ファイルのサンプル
│   ├── app.yaml
│   ├── app_minimal.yaml
│   └── app_invalid.yaml
├── entity/          # エンティティのサンプルデータ
│   ├── user.json
│   ├── group.json
│   ├── member.json
│   └── role.json
├── request/         # リクエストのサンプル
│   ├── user_request.json
│   ├── group_request.json
│   └── login_request.json
├── response/        # レスポンスのサンプル
│   ├── user_response.json
│   ├── group_response.json
│   └── error_response.json
├── casbin/          # Casbinポリシーファイル
│   ├── model.conf
│   └── policy.csv
├── jwt/             # JWT トークンのサンプル
│   └── test_tokens.json
└── passwords/       # パスワード強度テスト用データ
    ├── valid_passwords.txt
    └── invalid_passwords.txt
```

## 使用方法

テストファイルから以下のようにアクセスします：

```go
import (
    "github.com/ryo/mysrc/github/ryo-arima/locky/test/unit/internal/testutil"
)

func TestExample(t *testing.T) {
    // JSONファイルを読み込む
    var data SomeStruct
    err := testutil.LoadJSONFile("entity/user.json", &data)
    if err != nil {
        t.Fatal(err)
    }

    // YAMLファイルを読み込む
    var config Config
    err = testutil.LoadYAMLFile("config/app.yaml", &config)
    if err != nil {
        t.Fatal(err)
    }
}
```

## 原則

1. **すべてのテストデータはこのディレクトリに配置**
2. **実際の機密情報は含めない**（ダミーデータのみ）
3. **構造化されたディレクトリで管理**
4. **testutil パッケージのヘルパー関数を使用**

## ファイル命名規則

- 設定ファイル: `app.yaml`, `app_<variant>.yaml`
- エンティティ: `<entity_name>.json`
- リクエスト: `<entity_name>_request.json`
- レスポンス: `<entity_name>_response.json`
