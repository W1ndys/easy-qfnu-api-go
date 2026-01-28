package questions

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/W1ndys/easy-qfnu-api-go/model"
	_ "modernc.org/sqlite"
)

var (
	db   *sql.DB
	dbMu sync.Mutex
)

// getDBInstance 获取数据库单例连接
func getDBInstance() (*sql.DB, error) {
	dbMu.Lock()
	defer dbMu.Unlock()

	// 如果数据库已经连接且正常，直接返回
	if db != nil {
		if err := db.Ping(); err == nil {
			return db, nil
		}
		// 如果连接失效，关闭旧连接（忽略错误），准备重新连接
		_ = db.Close()
		db = nil
	}

	// 数据库文件路径
	dsn := "./data/freshman_questions.db"
	// 使用 modernc.org/sqlite 驱动，driverName 为 "sqlite"
	newDB, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 尝试ping一下确保连接正常
	if err = newDB.Ping(); err != nil {
		_ = newDB.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// 设置最大打开连接数
	newDB.SetMaxOpenConns(10)

	db = newDB
	return db, nil
}

// SearchQuestions 根据关键词搜索题目
func SearchQuestions(keyword string) ([]model.FreshmanQuestion, error) {
	database, err := getDBInstance()
	if err != nil {
		return nil, err
	}

	// 使用 LIKE 进行模糊查询
	query := `
		SELECT id, type, question, optionA, optionB, optionC, optionD, optionAnswer
		FROM questions
		WHERE question LIKE ?
		ORDER BY id ASC
	`

	rows, err := database.Query(query, "%"+keyword+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []model.FreshmanQuestion
	for rows.Next() {
		var q model.FreshmanQuestion
		// 注意：Scan 的顺序必须与 SQL 查询字段顺序一致
		// 数据库字段可能为 NULL，但这里 Struct 是 string，如果数据库有 NULL 会报错
		// 假设题目数据完整，如果有 NULL 需要用 sql.NullString
		err := rows.Scan(
			&q.ID,
			&q.Type,
			&q.QuestionText,
			&q.OptionA,
			&q.OptionB,
			&q.OptionC,
			&q.OptionD,
			&q.OptionAnswer,
		)
		if err != nil {
			return nil, err
		}
		questions = append(questions, q)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return questions, nil
}
