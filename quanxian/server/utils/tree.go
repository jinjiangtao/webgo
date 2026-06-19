package utils

import (
	"quanxian/models"
	"sort"
)

func BuildDepartmentTree(depts []models.Department) []models.Department {
	deptMap := make(map[uint]*models.Department)
	for i := range depts {
		deptMap[depts[i].ID] = &depts[i]
	}

	var roots []models.Department
	for i := range depts {
		if depts[i].ParentID == nil {
			roots = append(roots, depts[i])
		} else {
			if parent, exists := deptMap[*depts[i].ParentID]; exists {
				parent.Children = append(parent.Children, depts[i])
			}
		}
	}

	sortDepartments(roots)
	return roots
}

func sortDepartments(depts []models.Department) {
	sort.Slice(depts, func(i, j int) bool {
		return depts[i].Sort < depts[j].Sort
	})
	for i := range depts {
		sortDepartments(depts[i].Children)
	}
}

func BuildMenuTree(menus []models.Menu) []models.Menu {
	menuMap := make(map[uint]*models.Menu)
	for i := range menus {
		menuMap[menus[i].ID] = &menus[i]
	}

	var roots []models.Menu
	for i := range menus {
		if menus[i].ParentID == nil {
			roots = append(roots, menus[i])
		} else {
			if parent, exists := menuMap[*menus[i].ParentID]; exists {
				parent.Children = append(parent.Children, menus[i])
			}
		}
	}

	sortMenus(roots)
	return roots
}

func sortMenus(menus []models.Menu) {
	sort.Slice(menus, func(i, j int) bool {
		return menus[i].Sort < menus[j].Sort
	})
	for i := range menus {
		sortMenus(menus[i].Children)
	}
}

func BuildGenericTree(nodes []models.TreeNode) []models.TreeNode {
	nodeMap := make(map[uint]*models.TreeNode)
	for i := range nodes {
		nodeMap[nodes[i].ID] = &nodes[i]
	}

	var roots []models.TreeNode
	for i := range nodes {
		if nodes[i].ParentID == nil {
			roots = append(roots, nodes[i])
		} else {
			if parent, exists := nodeMap[*nodes[i].ParentID]; exists {
				parent.Children = append(parent.Children, nodes[i])
			}
		}
	}

	sortTreeNodes(roots)
	return roots
}

func sortTreeNodes(nodes []models.TreeNode) {
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Sort < nodes[j].Sort
	})
	for i := range nodes {
		sortTreeNodes(nodes[i].Children)
	}
}

func GetAllDepartmentIds(dept models.Department) []uint {
	ids := []uint{dept.ID}
	for _, child := range dept.Children {
		ids = append(ids, GetAllDepartmentIds(child)...)
	}
	return ids
}

func GetAllMenuIds(menu models.Menu) []uint {
	ids := []uint{menu.ID}
	for _, child := range menu.Children {
		ids = append(ids, GetAllMenuIds(child)...)
	}
	return ids
}

func FilterMenusByRole(menus []models.Menu, roleMenuIds []uint) []models.Menu {
	menuIdSet := make(map[uint]bool)
	for _, id := range roleMenuIds {
		menuIdSet[id] = true
	}

	var filtered []models.Menu
	for _, menu := range menus {
		if hasMenuPermission(menu, menuIdSet) {
			filteredMenu := filterMenuChildren(menu, menuIdSet)
			filtered = append(filtered, filteredMenu)
		}
	}
	return filtered
}

func hasMenuPermission(menu models.Menu, menuIdSet map[uint]bool) bool {
	if menuIdSet[menu.ID] {
		return true
	}
	for _, child := range menu.Children {
		if hasMenuPermission(child, menuIdSet) {
			return true
		}
	}
	return false
}

func filterMenuChildren(menu models.Menu, menuIdSet map[uint]bool) models.Menu {
	result := menu
	result.Children = nil
	for _, child := range menu.Children {
		if hasMenuPermission(child, menuIdSet) {
			filteredChild := filterMenuChildren(child, menuIdSet)
			result.Children = append(result.Children, filteredChild)
		}
	}
	return result
}
