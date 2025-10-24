package main

import (
	"database/sql"    // データベース操作用
	"encoding/csv"    // CSVファイルを読み込むため
	"fmt"             // 画面に出力するため
	"log"             // エラーログを出力するため
	"os"              // ファイル操作用
	"strconv"         // 文字列を数値に変換するため

	"github.com/joho/godotenv" // .envファイルを読み込むため
	_ "github.com/lib/pq"      // PostgreSQLドライバ
)

// 避難所のデータを保持する構造体（データの入れ物）
type Shelter struct {
	Name     string  // 避難所名
	NameKana string  // 避難所名（カナ）
	Address  string  // 住所
	Lat      float64 // 緯度
	Lon      float64 // 経度
	Tel      string  // 電話番号
	URL      string  // URL
}

const CSV_FILE_PATH = "172014_evacuation_space.csv"

func main() {
	// 1. .envファイルから環境変数を読み込む
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 2. データベースに接続するための情報を取得
	dbHost := os.Getenv("DB_HOST")         // データベースのホスト名
	dbPort := os.Getenv("DB_PORT")         // ポート番号（通常5432）
	dbUser := os.Getenv("DB_USER")         // ユーザー名
	dbPassword := os.Getenv("DB_PASSWORD") // パスワード
	dbName := os.Getenv("DB_NAME")         // データベース名

	// 3. PostgreSQLに接続するための文字列を作成
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// 4. データベースに接続
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("データベースへの接続に失敗しました:", err)
	}
	defer db.Close() // 関数終了時に接続を閉じる

	// 5. 接続が成功したか確認
	err = db.Ping()
	if err != nil {
		log.Fatal("データベースとの通信に失敗しました:", err)
	}
	fmt.Println("✓ データベースに接続しました")

	// 6. CSVファイルを開く
	file, err := os.Open(CSV_FILE_PATH)
	if err != nil {
		log.Fatal("CSVファイルを開けませんでした:", err)
	}
	defer file.Close() // 関数終了時にファイルを閉じる

	// 7. CSVファイルを読み込む
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal("CSVファイルの読み込みに失敗しました:", err)
	}

	fmt.Printf("✓ CSVファイルを読み込みました（%d行）\n", len(records))

	// 8. ヘッダー行をスキップ（1行目は列名なので）
	if len(records) < 2 {
		log.Fatal("CSVファイルにデータがありません")
	}

	// 9. データを1件ずつ処理
	successCount := 0 // 成功した件数
	errorCount := 0   // エラーが発生した件数
	idCounter := 1    // ID用のカウンター（0001から始める）

	for i, record := range records {
		// ヘッダー行をスキップ
		if i == 0 {
			continue
		}

		// CSVの列番号に対応
		// 列3: 名称、列4: 名称_カナ、列8: 所在地_連結表記
		// 列14: 緯度、列15: 経度、列17: 電話番号、列33: URL

		// データを抽出
		name := record[3]      // 名称
		nameKana := record[4]  // 名称_カナ
		address := record[8]   // 所在地_連結表記
		latStr := record[14]   // 緯度（文字列）
		lonStr := record[15]   // 経度（文字列）
		tel := record[17]      // 電話番号
		url := record[33]      // URL

		// 緯度・経度を数値に変換
		lat, err := strconv.ParseFloat(latStr, 64)
		if err != nil {
			// 緯度の変換に失敗した場合はスキップ
			fmt.Printf("⚠ %d行目: 緯度の変換に失敗しました: %s\n", i+1, name)
			errorCount++
			continue
		}

		lon, err := strconv.ParseFloat(lonStr, 64)
		if err != nil {
			// 経度の変換に失敗した場合はスキップ
			fmt.Printf("⚠ %d行目: 経度の変換に失敗しました: %s\n", i+1, name)
			errorCount++
			continue
		}

		// IDを4桁のゼロ埋め文字列に変換（例: 1 → "0001"）
		id := fmt.Sprintf("%04d", idCounter)

		// 10. データベースに挿入するSQLを作成
		// データがない場合はNULLを設定
		query := `
			INSERT INTO place (id, name, name_kana, address, lat, lon, url, tel)
			VALUES ($1, $2, $3, $4, $5, $6, NULLIF($7, ''), NULLIF($8, ''))
		`

		// 11. SQLを実行してデータを挿入
		_, err = db.Exec(query, id, name, nameKana, address, lat, lon, url, tel)
		if err != nil {
			fmt.Printf("✗ %d行目: データの挿入に失敗しました: %s - エラー: %v\n", i+1, name, err)
			errorCount++
			continue
		}

		// 成功した場合
		successCount++
		idCounter++ // IDをインクリメント
		fmt.Printf("✓ %d行目: %s を登録しました (ID: %s)\n", i+1, name, id)
	}

	// 12. 結果を表示
	fmt.Println("\n========== 処理完了 ==========")
	fmt.Printf("成功: %d件\n", successCount)
	fmt.Printf("失敗: %d件\n", errorCount)
	fmt.Printf("合計: %d件\n", successCount+errorCount)
}
