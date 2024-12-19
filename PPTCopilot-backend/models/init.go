package models

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

func jsonObject2string(obj JsonObject) string {
	ret, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(ret)
}

// 初始化数据表
func init() {
	// 获取环境变量
	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlPort := os.Getenv("MYSQL_PORT")
	fmt.Println("MYSQL_HOST: ", mysqlHost)
	fmt.Println("MYSQL_PORT: ", mysqlPort)
	connectionString := fmt.Sprintf("root:admin@tcp(%s:%s)/now_db?charset=utf8&loc=Local", mysqlHost, mysqlPort)
	fmt.Println("Connection string: ", connectionString)

	// Retry connecting to the database
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		err := orm.RegisterDataBase("default", "mysql", connectionString)
		if err == nil {
			fmt.Println("Database connected successfully!")
			break
		}
		fmt.Printf("Failed to connect to database. Retrying in 5 seconds... (attempt %d/%d)\n", i+1, maxRetries)
		time.Sleep(5 * time.Second)
	}

	// 注册定义的model
	orm.RegisterModel(new(Outline))
	orm.RegisterModel(new(User))
	orm.RegisterModel(new(Project))
	orm.RegisterModel(new(File))
	orm.RegisterModel(new(Favorite))
	orm.RegisterModel(new(Template))

	// 如果表不存在则创建表
	orm.RunSyncdb("default", false, true)

	// Process template
	processTemplateFile("./scripts/jianyuezi.json", "简约紫")
	processTemplateFile("./scripts/taikongren.json", "太空人")
	processTemplateFile("./scripts/dangjian.json", "党建")

	// Add Admin user
	CreateUser("admin", "123456", "admin@163.com")
	CreateUser("jack", "123456", "jack@qq.com")
	CreateUser("john_doe", "password123", "john_doe@qq.com")
}

// Helper function to load and process templates
func processTemplateFile(filePath string, templateName string) {
    file, err := os.Open(filePath)
    if err != nil {
        fmt.Printf("Error opening file %s: %v\n", filePath, err)
        return
    }
    defer file.Close()

    var templateData []JsonObject
    if err := json.NewDecoder(file).Decode(&templateData); err != nil {
        fmt.Printf("Error decoding JSON from file %s: %v\n", filePath, err)
        return
    }

    if len(templateData) < 10 {
        fmt.Printf("Error: Insufficient data in template %s\n", templateName)
        return
    }

    template := Template{
        Name:       templateName,
        Cover:      jsonObject2string(templateData[0]),
        Transition: jsonObject2string(templateData[1]),
        Catalog_3:  jsonObject2string(templateData[2]),
        Catalog_4:  jsonObject2string(templateData[3]),
        Catalog_5:  jsonObject2string(templateData[4]),
        Content_1:  jsonObject2string(templateData[5]),
        Content_2:  jsonObject2string(templateData[6]),
        Content_3:  jsonObject2string(templateData[7]),
        Content_4:  jsonObject2string(templateData[8]),
        Thank:      jsonObject2string(templateData[9]),
    }

    CreateTemplate(template)
}
