package utils

import (
	"log"
	"quanxian/models"

	"golang.org/x/crypto/bcrypt"
)

func InitSchema() {
	err := DB.AutoMigrate(
		&models.Department{},
		&models.Role{},
		&models.Menu{},
		&models.Button{},
		&models.User{},
		&models.UserRole{},
		&models.RoleMenu{},
		&models.RoleButton{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}
	log.Println("Database schema migrated successfully")
}

func SeedData() {
	var count int64
	DB.Model(&models.User{}).Count(&count)
	if count > 0 {
		log.Println("Seed data already exists, skipping")
		return
	}

	log.Println("Seeding initial data...")

	depts := []models.Department{
		{Name: "总公司", ParentID: nil, Sort: 1, Status: 1},
		{Name: "技术部", ParentID: uintPtr(1), Sort: 1, Status: 1},
		{Name: "市场部", ParentID: uintPtr(1), Sort: 2, Status: 1},
		{Name: "财务部", ParentID: uintPtr(1), Sort: 3, Status: 1},
		{Name: "前端组", ParentID: uintPtr(2), Sort: 1, Status: 1},
		{Name: "后端组", ParentID: uintPtr(2), Sort: 2, Status: 1},
	}
	for i := range depts {
		DB.Create(&depts[i])
	}

	menus := []models.Menu{
		{Name: "系统管理", Path: "/system", Component: "Layout", Icon: "setting", ParentID: nil, Sort: 1, Visible: 1, Status: 1},
		{Name: "用户管理", Path: "/system/user", Component: "system/user/index", Icon: "user", ParentID: uintPtr(1), Sort: 1, Visible: 1, Status: 1},
		{Name: "角色管理", Path: "/system/role", Component: "system/role/index", Icon: "peoples", ParentID: uintPtr(1), Sort: 2, Visible: 1, Status: 1},
		{Name: "菜单管理", Path: "/system/menu", Component: "system/menu/index", Icon: "menu", ParentID: uintPtr(1), Sort: 3, Visible: 1, Status: 1},
		{Name: "部门管理", Path: "/system/dept", Component: "system/dept/index", Icon: "tree", ParentID: uintPtr(1), Sort: 4, Visible: 1, Status: 1},
		{Name: "权限配置", Path: "/system/permission", Component: "system/permission/index", Icon: "key", ParentID: uintPtr(1), Sort: 5, Visible: 1, Status: 1},
	}
	for i := range menus {
		DB.Create(&menus[i])
	}

	buttons := []models.Button{
		{Name: "用户查询", Code: "system:user:query", MenuID: 2, Sort: 1, Status: 1},
		{Name: "用户新增", Code: "system:user:add", MenuID: 2, Sort: 2, Status: 1},
		{Name: "用户修改", Code: "system:user:edit", MenuID: 2, Sort: 3, Status: 1},
		{Name: "用户删除", Code: "system:user:delete", MenuID: 2, Sort: 4, Status: 1},
		{Name: "角色查询", Code: "system:role:query", MenuID: 3, Sort: 1, Status: 1},
		{Name: "角色新增", Code: "system:role:add", MenuID: 3, Sort: 2, Status: 1},
		{Name: "角色修改", Code: "system:role:edit", MenuID: 3, Sort: 3, Status: 1},
		{Name: "角色删除", Code: "system:role:delete", MenuID: 3, Sort: 4, Status: 1},
		{Name: "菜单查询", Code: "system:menu:query", MenuID: 4, Sort: 1, Status: 1},
		{Name: "菜单新增", Code: "system:menu:add", MenuID: 4, Sort: 2, Status: 1},
		{Name: "菜单修改", Code: "system:menu:edit", MenuID: 4, Sort: 3, Status: 1},
		{Name: "菜单删除", Code: "system:menu:delete", MenuID: 4, Sort: 4, Status: 1},
		{Name: "部门查询", Code: "system:dept:query", MenuID: 5, Sort: 1, Status: 1},
		{Name: "部门新增", Code: "system:dept:add", MenuID: 5, Sort: 2, Status: 1},
		{Name: "部门修改", Code: "system:dept:edit", MenuID: 5, Sort: 3, Status: 1},
		{Name: "部门删除", Code: "system:dept:delete", MenuID: 5, Sort: 4, Status: 1},
		{Name: "权限查询", Code: "system:permission:query", MenuID: 6, Sort: 1, Status: 1},
		{Name: "权限配置", Code: "system:permission:config", MenuID: 6, Sort: 2, Status: 1},
	}
	for i := range buttons {
		DB.Create(&buttons[i])
	}

	roles := []models.Role{
		{Name: "超级管理员", Code: "admin", Remark: "拥有所有权限", Sort: 1, Status: 1},
		{Name: "系统管理员", Code: "system", Remark: "系统管理权限", Sort: 2, Status: 1},
		{Name: "普通用户", Code: "user", Remark: "基础权限", Sort: 3, Status: 1},
	}
	for i := range roles {
		DB.Create(&roles[i])
	}

	DB.Model(&roles[0]).Association("Menus").Replace(&menus)
	DB.Model(&roles[0]).Association("Buttons").Replace(&buttons)

	DB.Model(&roles[1]).Association("Menus").Replace(&menus)
	systemButtons := buttons[:8]
	DB.Model(&roles[1]).Association("Buttons").Replace(&systemButtons)

	userMenus := []models.Menu{menus[0], menus[1]}
	DB.Model(&roles[2]).Association("Menus").Replace(&userMenus)
	userButtons := []models.Button{buttons[0]}
	DB.Model(&roles[2]).Association("Buttons").Replace(&userButtons)

	password, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	users := []models.User{
		{Username: "admin", Password: string(password), Nickname: "超级管理员", Email: "admin@example.com", Phone: "13800138000", DeptID: uintPtr(2), Status: 1},
		{Username: "system", Password: string(password), Nickname: "系统管理员", Email: "system@example.com", Phone: "13800138001", DeptID: uintPtr(2), Status: 1},
		{Username: "user", Password: string(password), Nickname: "普通用户", Email: "user@example.com", Phone: "13800138002", DeptID: uintPtr(5), Status: 1},
	}
	for i := range users {
		DB.Create(&users[i])
	}

	DB.Model(&users[0]).Association("Roles").Replace([]models.Role{roles[0]})
	DB.Model(&users[1]).Association("Roles").Replace([]models.Role{roles[1]})
	DB.Model(&users[2]).Association("Roles").Replace([]models.Role{roles[2]})

	log.Println("Seed data initialized successfully")
}

func uintPtr(v uint) *uint {
	return &v
}

func GetRoleMenuIds(roleID uint) []uint {
	var roleMenus []models.RoleMenu
	DB.Where("role_id = ?", roleID).Find(&roleMenus)
	ids := make([]uint, len(roleMenus))
	for i, rm := range roleMenus {
		ids[i] = rm.MenuID
	}
	return ids
}

func GetRoleButtonIds(roleID uint) []uint {
	var roleButtons []models.RoleButton
	DB.Where("role_id = ?", roleID).Find(&roleButtons)
	ids := make([]uint, len(roleButtons))
	for i, rb := range roleButtons {
		ids[i] = rb.ButtonID
	}
	return ids
}

func GetUserRoleIds(userID uint) []uint {
	var userRoles []models.UserRole
	DB.Where("user_id = ?", userID).Find(&userRoles)
	ids := make([]uint, len(userRoles))
	for i, ur := range userRoles {
		ids[i] = ur.RoleID
	}
	return ids
}

func GetUserAllMenuIds(userID uint) []uint {
	roleIds := GetUserRoleIds(userID)
	if len(roleIds) == 0 {
		return []uint{}
	}

	var roleMenus []models.RoleMenu
	DB.Where("role_id IN ?", roleIds).Find(&roleMenus)
	idSet := make(map[uint]bool)
	for _, rm := range roleMenus {
		idSet[rm.MenuID] = true
	}

	ids := make([]uint, 0, len(idSet))
	for id := range idSet {
		ids = append(ids, id)
	}
	return ids
}

func GetUserAllButtonCodes(userID uint) []string {
	roleIds := GetUserRoleIds(userID)
	if len(roleIds) == 0 {
		return []string{}
	}

	var buttons []models.Button
	DB.Joins("JOIN role_buttons ON role_buttons.button_id = buttons.id").
		Where("role_buttons.role_id IN ?", roleIds).
		Find(&buttons)

	codes := make([]string, len(buttons))
	for i, btn := range buttons {
		codes[i] = btn.Code
	}
	return codes
}

func IsSuperAdmin(userID uint) bool {
	roleIds := GetUserRoleIds(userID)
	for _, rid := range roleIds {
		var role models.Role
		if err := DB.First(&role, rid).Error; err == nil {
			if role.Code == "admin" {
				return true
			}
		}
	}
	return false
}

func HasPermission(userID uint, permissionCode string) bool {
	if IsSuperAdmin(userID) {
		return true
	}
	codes := GetUserAllButtonCodes(userID)
	for _, code := range codes {
		if code == permissionCode {
			return true
		}
	}
	return false
}

func GetUserMenus(userID uint) ([]models.Menu, error) {
	var allMenus []models.Menu
	if err := DB.Where("status = ? AND visible = ?", 1, 1).Order("sort ASC").Find(&allMenus).Error; err != nil {
		return nil, err
	}

	if IsSuperAdmin(userID) {
		return BuildMenuTree(allMenus), nil
	}

	menuIds := GetUserAllMenuIds(userID)
	filteredMenus := FilterMenusByRole(allMenus, menuIds)
	return filteredMenus, nil
}

func IsSuperAdminByRole(roleID uint) bool {
	var role models.Role
	if err := DB.First(&role, roleID).Error; err == nil {
		if role.Code == "admin" {
			return true
		}
	}
	return false
}
