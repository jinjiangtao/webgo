package database

import (
	"log"
	"jiansuo/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	var err error
	DB, err = gorm.Open(sqlite.Open("jiansuo.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = DB.AutoMigrate(
		&models.Category{},
		&models.Keyword{},
		&models.SearchLog{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	seedData()
}

func seedData() {
	var categoryCount int64
	DB.Model(&models.Category{}).Count(&categoryCount)
	if categoryCount > 0 {
		return
	}

	categories := []models.Category{
		{Name: "技术文档", Description: "编程语言、框架、工具等技术类文档", Status: 1},
		{Name: "产品介绍", Description: "产品功能、特性、使用说明", Status: 1},
		{Name: "常见问题", Description: "FAQ、故障排除、解决方案", Status: 1},
		{Name: "教程指南", Description: "入门教程、操作指南、最佳实践", Status: 1},
		{Name: "开发规范", Description: "编码规范、流程规范、标准文档", Status: 1},
	}
	DB.Create(&categories)

	keywords := []models.Keyword{
		{Title: "Go语言入门教程", Content: "Go语言是Google开发的开源编程语言，具有简洁高效、并发支持等特点。适合构建高性能、可扩展的后端服务和分布式系统。", CategoryID: 1, Tags: "go,golang,入门,教程", Status: 1, Sort: 100},
		{Title: "Vue3组合式API详解", Content: "Vue3引入了组合式API（Composition API），包括setup、ref、reactive、computed、watch等核心概念，提供了更好的代码组织和逻辑复用能力。", CategoryID: 1, Tags: "vue,vue3,前端,组合式api", Status: 1, Sort: 99},
		{Title: "Gin框架快速上手", Content: "Gin是一个用Go语言编写的Web框架，具有高性能、易用性强的特点。支持中间件、路由分组、参数绑定、JSON验证等功能。", CategoryID: 1, Tags: "gin,go,web框架,后端", Status: 1, Sort: 98},
		{Title: "产品核心功能介绍", Content: "本平台提供智能关键词检索、全文搜索、分类管理、数据统计等核心功能。支持多条件筛选、结果高亮、联想提示等高级特性。", CategoryID: 2, Tags: "产品,功能,介绍,核心", Status: 1, Sort: 95},
		{Title: "如何使用搜索联想功能", Content: "在搜索框中输入关键词，系统会自动显示相关联想词。使用方向键或鼠标选择联想词，可快速完成搜索。支持模糊匹配和拼音检索。", CategoryID: 2, Tags: "搜索,联想,使用,帮助", Status: 1, Sort: 90},
		{Title: "搜索无结果怎么办", Content: "如果搜索无结果，建议：1.检查关键词拼写；2.尝试使用同义词或相关词；3.减少筛选条件；4.使用更宽泛的关键词。如仍有问题请联系管理员。", CategoryID: 3, Tags: "搜索,无结果,常见问题,FAQ", Status: 1, Sort: 85},
		{Title: "账号密码忘记了如何找回", Content: "点击登录页的忘记密码链接，输入注册邮箱，系统将发送重置密码邮件。如未收到邮件，请检查垃圾邮件文件夹或联系客服。", CategoryID: 3, Tags: "账号,密码,找回,登录", Status: 1, Sort: 80},
		{Title: "快速开始指南", Content: "第一步：注册账号并登录系统；第二步：在搜索框输入关键词进行检索；第三步：使用筛选条件精准定位内容；第四步：点击结果查看详情。", CategoryID: 4, Tags: "快速开始,入门,指南,教程", Status: 1, Sort: 75},
		{Title: "批量导入关键词操作步骤", Content: "进入关键词管理页面，点击批量导入按钮，下载模板文件，按照模板格式填写数据，上传文件即可完成批量导入。支持CSV和Excel格式。", CategoryID: 4, Tags: "批量导入,关键词,操作,步骤", Status: 1, Sort: 70},
		{Title: "代码编写规范", Content: "遵循以下编码规范：1.使用清晰的命名；2.添加必要的注释；3.函数职责单一；4.进行错误处理；5.编写单元测试；6.提交前代码审查。", CategoryID: 5, Tags: "编码,规范,代码,开发", Status: 1, Sort: 65},
		{Title: "Git提交规范", Content: "Git提交信息格式：type(scope): subject。type包括feat新功能、fix修复、docs文档、style格式、refactor重构、test测试、chore构建。", CategoryID: 5, Tags: "git,提交,规范,版本控制", Status: 1, Sort: 60},
		{Title: "SQLite数据库配置", Content: "SQLite是轻量级嵌入式数据库，无需独立服务。本系统使用GORM操作SQLite，支持自动迁移、事务、关联查询等ORM功能。", CategoryID: 1, Tags: "sqlite,数据库,gorm,配置", Status: 1, Sort: 55},
		{Title: "前端性能优化技巧", Content: "前端性能优化包括：懒加载、代码分割、缓存策略、图片压缩、CDN加速、减少HTTP请求、使用虚拟列表等技术手段。", CategoryID: 1, Tags: "前端,性能,优化,vue", Status: 1, Sort: 50},
		{Title: "RESTful API设计原则", Content: "RESTful API设计遵循：使用HTTP方法表示操作、资源使用名词复数、状态码正确使用、版本控制、分页支持、错误信息标准化。", CategoryID: 5, Tags: "api,restful,设计,规范", Status: 1, Sort: 45},
		{Title: "系统部署上线流程", Content: "部署流程：1.代码测试通过；2.编译构建生产版本；3.备份数据库；4.停止旧服务；5.部署新版本；6.启动服务并验证。", CategoryID: 4, Tags: "部署,上线,流程,运维", Status: 1, Sort: 40},
	}
	DB.Create(&keywords)

	log.Println("Database seeded successfully")
}
