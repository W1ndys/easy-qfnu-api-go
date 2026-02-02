package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite" // 使用项目已有的纯 Go SQLite 驱动
)

// OldReview 对应旧数据库 campus.db 中的 course_reviews 表
type OldReview struct {
	ID          int64
	CourseName  sql.NullString
	TeacherName sql.NullString
	Campus      sql.NullString
	Semester    sql.NullString
	Reason      sql.NullString
	Nickname    sql.NullString
	Status      int
	IsVisible   int
	CreatedAt   sql.NullFloat64 // 旧库是 REAL 类型
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// 1. 定义数据库路径
	sourceDBPath := filepath.Join(cwd, "data", "campus.db")
	targetDBPath := filepath.Join(cwd, "data", "course_recommendation.db")

	fmt.Println("=== 开始迁移数据 ===")
	fmt.Printf("源数据库 (旧): %s\n", sourceDBPath)
	fmt.Printf("目标数据库 (新): %s\n", targetDBPath)

	// 2. 连接源数据库
	if _, err := os.Stat(sourceDBPath); os.IsNotExist(err) {
		log.Fatalf("错误: 源数据库文件不存在: %s", sourceDBPath)
	}
	sourceDB, err := sql.Open("sqlite", sourceDBPath)
	if err != nil {
		log.Fatalf("无法打开源数据库: %v", err)
	}
	defer sourceDB.Close()

	// 3. 连接目标数据库
	// 确保目录存在
	os.MkdirAll(filepath.Dir(targetDBPath), 0755)
	targetDB, err := sql.Open("sqlite", targetDBPath)
	if err != nil {
		log.Fatalf("无法打开目标数据库: %v", err)
	}
	defer targetDB.Close()

	// 4. 确保目标表存在 (根据 internal/database/database.go 的定义)
	initSQL := `
	CREATE TABLE IF NOT EXISTS course_recommendations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		course_name TEXT NOT NULL,
		teacher_name TEXT NOT NULL,
		recommendation_reason TEXT NOT NULL,
		recommender_nickname TEXT NOT NULL,
		recommendation_time INTEGER NOT NULL,
		is_visible INTEGER DEFAULT 0,
		campus TEXT DEFAULT '',
		recommendation_year TEXT DEFAULT ''
	);
	CREATE INDEX IF NOT EXISTS idx_course_rec_course ON course_recommendations(course_name);
	CREATE INDEX IF NOT EXISTS idx_course_rec_teacher ON course_recommendations(teacher_name);
	CREATE INDEX IF NOT EXISTS idx_course_rec_visible ON course_recommendations(is_visible);
	`
	if _, err := targetDB.Exec(initSQL); err != nil {
		log.Fatalf("初始化目标表失败: %v", err)
	}

	// 5. 从旧库读取数据
	fmt.Println("正在读取旧数据...")
	rows, err := sourceDB.Query(`SELECT id, course_name, teacher_name, campus, semester, reason, nickname, status, is_visible, created_at FROM course_reviews`)
	if err != nil {
		log.Fatalf("读取 course_reviews 失败: %v", err)
	}
	defer rows.Close()

	// 6. 开启事务准备写入
	tx, err := targetDB.Begin()
	if err != nil {
		log.Fatal(err)
	}

	insertSQL := `INSERT INTO course_recommendations
		(course_name, teacher_name, recommendation_reason, recommender_nickname, recommendation_time, is_visible, campus, recommendation_year)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	stmt, err := tx.Prepare(insertSQL)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	count := 0
	success := 0

	for rows.Next() {
		var r OldReview
		// 使用 NullString 和 NullFloat64 处理可能的 NULL 值
		if err := rows.Scan(&r.ID, &r.CourseName, &r.TeacherName, &r.Campus, &r.Semester, &r.Reason, &r.Nickname, &r.Status, &r.IsVisible, &r.CreatedAt); err != nil {
			log.Printf("警告: 读取第 %d 行失败: %v", count+1, err)
			continue
		}

		// --- 字段映射与转换逻辑 ---

		// 1. 可见性转换: 旧库 status=1 且 is_visible=1 才算可见
		newIsVisible := 0
		if r.Status == 1 && r.IsVisible == 1 {
			newIsVisible = 1
		}

		// 2. 时间转换: float64 (Unix timestamp) -> int64
		newTime := int64(0)
		if r.CreatedAt.Valid {
			newTime = int64(r.CreatedAt.Float64)
		}
		if newTime == 0 {
			newTime = time.Now().Unix()
		}

		// 4. 年份格式标准化处理
		yearStr := r.Semester.String
		// 处理 xxxx-xxxx-1/2 格式
		if strings.Contains(yearStr, "-") {
			parts := strings.Split(yearStr, "-")
			// 如果是 xxxx-xxxx-x 格式 (例如 2023-2024-1)
			if len(parts) >= 3 {
				term := parts[len(parts)-1]
				// 第一学期取第一个年份，第二学期取第二个年份
				if term == "1" {
					yearStr = parts[0]
				} else if term == "2" {
					// 确保第二个年份存在
					if len(parts) >= 2 {
						yearStr = parts[1]
					}
				} else {
					// 其他情况默认取第一个年份
					yearStr = parts[0]
				}
			} else if len(parts) == 2 {
				// 简单的 xxxx-xxxx 格式，默认取第一个年份
				yearStr = parts[0]
			}
		}

		// 3. 执行插入 (String 字段处理 NULL -> 空字符串)
		_, err = stmt.Exec(
			r.CourseName.String,
			r.TeacherName.String,
			r.Reason.String,
			r.Nickname.String,
			newTime,
			newIsVisible,
			r.Campus.String, // 如果是 NULL，这里会是 ""
			yearStr,         // 使用处理后的年份
		)

		if err != nil {
			log.Printf("写入 ID %d 失败: %v", r.ID, err)
		} else {
			success++
		}
		count++
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("提交事务失败: %v", err)
	}

	fmt.Printf("=== 迁移完成 ===\n共扫描: %d 条\n成功导入: %d 条\n", count, success)
}
