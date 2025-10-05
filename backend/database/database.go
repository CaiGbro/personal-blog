// /backend/database/database.go
package database

import (
    "context"
    "fmt"
    "github.com/jackc/pgx/v4/pgxpool"
    "os"
)

var Pool *pgxpool.Pool

func InitDB() error {
    // 数据库连接字符串，请根据您的实际情况修改
    // 格式: postgresql://[user]:[password]@[host]:[port]/[database]
    // 密码是 caigbro123, 端口是 5432
    connStr := os.Getenv("DATABASE_URL") // <-- 使用环境变量
    if connStr == "" {
        return fmt.Errorf("DATABASE_URL environment variable not set")
    }
    
    var err error
    Pool, err = pgxpool.Connect(context.Background(), connStr)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
        return err
    }

    // 检查连接是否成功
    err = Pool.Ping(context.Background())
    if err != nil {
         fmt.Fprintf(os.Stderr, "Database ping failed: %v\n", err)
         return err
    }

    fmt.Println("Successfully connected to the database!")
    return nil
}