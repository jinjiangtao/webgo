package models

import (
	"time"

	"gorm.io/gorm"
)

type Department struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"size:100;not null"`
	ParentID  *uint          `json:"parentId" gorm:"index"`
	Sort      int            `json:"sort" gorm:"default:0"`
	Status    int            `json:"status" gorm:"default:1"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Children  []Department   `json:"children,omitempty" gorm:"-"`
	Users     []User         `json:"users,omitempty" gorm:"foreignKey:DeptID"`
}

type Role struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"size:100;not null"`
	Code      string         `json:"code" gorm:"size:50;uniqueIndex;not null"`
	Remark    string         `json:"remark" gorm:"size:255"`
	Status    int            `json:"status" gorm:"default:1"`
	Sort      int            `json:"sort" gorm:"default:0"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Users     []User         `json:"users,omitempty" gorm:"many2many:user_roles;"`
	Menus     []Menu         `json:"menus,omitempty" gorm:"many2many:role_menus;"`
	Buttons   []Button       `json:"buttons,omitempty" gorm:"many2many:role_buttons;"`
}

type Menu struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"size:100;not null"`
	Path      string         `json:"path" gorm:"size:255"`
	Component string         `json:"component" gorm:"size:255"`
	Icon      string         `json:"icon" gorm:"size:100"`
	ParentID  *uint          `json:"parentId" gorm:"index"`
	Sort      int            `json:"sort" gorm:"default:0"`
	Visible   int            `json:"visible" gorm:"default:1"`
	Status    int            `json:"status" gorm:"default:1"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Children  []Menu         `json:"children,omitempty" gorm:"-"`
	Buttons   []Button       `json:"buttons,omitempty" gorm:"foreignKey:MenuID"`
}

type Button struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"size:100;not null"`
	Code      string         `json:"code" gorm:"size:100;not null"`
	MenuID    uint           `json:"menuId" gorm:"index;not null"`
	Sort      int            `json:"sort" gorm:"default:0"`
	Status    int            `json:"status" gorm:"default:1"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type User struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	Username   string         `json:"username" gorm:"size:50;uniqueIndex;not null"`
	Password   string         `json:"-" gorm:"size:255;not null"`
	Nickname   string         `json:"nickname" gorm:"size:100"`
	Email      string         `json:"email" gorm:"size:100"`
	Phone      string         `json:"phone" gorm:"size:20"`
	Avatar     string         `json:"avatar" gorm:"size:255"`
	DeptID     *uint          `json:"deptId" gorm:"index"`
	Status     int            `json:"status" gorm:"default:1"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	Department *Department    `json:"department,omitempty" gorm:"foreignKey:DeptID"`
	Roles      []Role         `json:"roles,omitempty" gorm:"many2many:user_roles;"`
}

type UserRole struct {
	UserID uint `json:"userId" gorm:"primaryKey"`
	RoleID uint `json:"roleId" gorm:"primaryKey"`
}

type RoleMenu struct {
	RoleID uint `json:"roleId" gorm:"primaryKey"`
	MenuID uint `json:"menuId" gorm:"primaryKey"`
}

type RoleButton struct {
	RoleID   uint `json:"roleId" gorm:"primaryKey"`
	ButtonID uint `json:"buttonId" gorm:"primaryKey"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type TreeNode struct {
	ID       uint       `json:"id"`
	Label    string     `json:"label"`
	ParentID *uint      `json:"parentId,omitempty"`
	Children []TreeNode `json:"children,omitempty"`
	Type     string     `json:"type,omitempty"`
	Checked  bool       `json:"checked,omitempty"`
	Sort     int        `json:"sort,omitempty"`
}

type PermissionPreview struct {
	Roles   []Role           `json:"roles"`
	Menus   []Menu           `json:"menus"`
	Buttons map[string][]string `json:"buttons"`
}
