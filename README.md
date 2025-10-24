# shelter_adder

避難所データを CSV ファイルから PostgreSQL データベースに登録するための Go アプリケーション

## 📋 必要なもの

- Go 1.21 以上
- Railway にデプロイされた PostgreSQL データベース
- CSV ファイル（`172014_evacuation_space.csv`）

## 🚀 セットアップ手順

### 1. 環境変数の設定

`.env.example` を `.env` にコピーして、Railway の PostgreSQL 接続情報を設定します：

```bash
cp .env.example .env
```

`.env` ファイルを開いて、Railway のダッシュボードから取得した情報を入力してください：

```
DB_HOST=あなたのRailwayホスト名.railway.app
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=あなたのパスワード
DB_NAME=railway
```

### 2. 依存パッケージのインストール

```bash
go mod download
```

### 3. データベーステーブルの確認

PostgreSQL に `place` テーブルが以下の構造で作成されていることを確認してください：

```sql
CREATE TABLE place (
    id VARCHAR(4) PRIMARY KEY,    -- 例: "0001"
    name VARCHAR(255) NOT NULL,   -- 避難所名
    name_kana VARCHAR(255),       -- フリガナ
    address TEXT,                 -- 住所
    lat DOUBLE PRECISION,         -- 緯度
    lon DOUBLE PRECISION,         -- 経度
    url TEXT,                     -- URL（任意）
    tel VARCHAR(50)               -- 電話番号（任意）
);
```

### 4. アプリケーションの実行

```bash
go run main.go
```

## 📊 処理内容

1. CSV ファイルから避難所データを読み込みます
2. 各行のデータを解析します
3. データベースに 1 件ずつ挿入します
4. ID は `0001` から順番に自動採番されます
5. データが空の場合（URL や電話番号など）は NULL として保存されます

## 🎯 実行結果の例

```
✓ データベースに接続しました
✓ CSVファイルを読み込みました（1839行）
✓ 2行目: 三谷公民館（指定避難場所） を登録しました (ID: 0001)
✓ 3行目: 三谷小学校（指定避難場所） を登録しました (ID: 0002)
...
========== 処理完了 ==========
成功: 1838件
失敗: 0件
合計: 1838件
```

## ⚠️ 注意事項

- `.env` ファイルには機密情報が含まれるため、Git にコミットしないでください
- すでに `.gitignore` に追加されています
- データベーステーブルが事前に作成されている必要があります
- CSV ファイルは UTF-8 エンコーディングで保存されている必要があります

## 🔧 トラブルシューティング

### データベースに接続できない場合

- Railway のダッシュボードで接続情報が正しいか確認してください
- ファイアウォールやネットワーク設定を確認してください

### CSV ファイルが見つからない場合

- ファイル名が `172014_evacuation_space.csv` であることを確認してください
- ファイルがプロジェクトのルートディレクトリにあることを確認してください

### データの挿入に失敗する場合

- テーブルのスキーマが正しいか確認してください
- 緯度・経度が数値として正しく変換できているか確認してください
