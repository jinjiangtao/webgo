package repository

import (
	"dashboard/internal/model"
)

type DimensionRepository struct{}

func NewDimensionRepository() *DimensionRepository {
	return &DimensionRepository{}
}

func (r *DimensionRepository) GetTimeOptions() ([]map[string]interface{}, error) {
	var times []model.DimTime
	err := DB.Distinct("year").Order("year desc").Find(&times).Error
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0)
	yearMap := make(map[int]bool)
	for _, t := range times {
		if !yearMap[t.Year] {
			yearMap[t.Year] = true
			result = append(result, map[string]interface{}{
				"value": t.Year,
				"label": t.Year,
			})
		}
	}
	return result, nil
}

func (r *DimensionRepository) GetRegionTree() ([]map[string]interface{}, error) {
	var regions []model.DimRegion
	err := DB.Where("level = ?", 1).Find(&regions).Error
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0)
	for _, region := range regions {
		children, _ := r.getRegionChildren(region.RegionCode)
		result = append(result, map[string]interface{}{
			"value":    region.RegionCode,
			"label":    region.RegionName,
			"children": children,
		})
	}
	return result, nil
}

func (r *DimensionRepository) getRegionChildren(parentCode string) ([]map[string]interface{}, error) {
	var regions []model.DimRegion
	err := DB.Where("parent_code = ?", parentCode).Find(&regions).Error
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0)
	for _, region := range regions {
		children, _ := r.getRegionChildren(region.RegionCode)
		item := map[string]interface{}{
			"value": region.RegionCode,
			"label": region.RegionName,
		}
		if len(children) > 0 {
			item["children"] = children
		}
		result = append(result, item)
	}
	return result, nil
}

func (r *DimensionRepository) GetBusinessTree() ([]map[string]interface{}, error) {
	var businesses []model.DimBusiness
	err := DB.Where("level = ?", 1).Find(&businesses).Error
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0)
	for _, biz := range businesses {
		children, _ := r.getBusinessChildren(biz.BusinessCode)
		result = append(result, map[string]interface{}{
			"value":    biz.BusinessCode,
			"label":    biz.BusinessName,
			"category": biz.Category,
			"children": children,
		})
	}
	return result, nil
}

func (r *DimensionRepository) getBusinessChildren(parentCode string) ([]map[string]interface{}, error) {
	var businesses []model.DimBusiness
	err := DB.Where("parent_code = ?", parentCode).Find(&businesses).Error
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0)
	for _, biz := range businesses {
		children, _ := r.getBusinessChildren(biz.BusinessCode)
		item := map[string]interface{}{
			"value": biz.BusinessCode,
			"label": biz.BusinessName,
		}
		if len(children) > 0 {
			item["children"] = children
		}
		result = append(result, item)
	}
	return result, nil
}

func (r *DimensionRepository) GetRegionByCode(code string) (*model.DimRegion, error) {
	var region model.DimRegion
	err := DB.Where("region_code = ?", code).First(&region).Error
	if err != nil {
		return nil, err
	}
	return &region, nil
}

func (r *DimensionRepository) GetBusinessByCode(code string) (*model.DimBusiness, error) {
	var biz model.DimBusiness
	err := DB.Where("business_code = ?", code).First(&biz).Error
	if err != nil {
		return nil, err
	}
	return &biz, nil
}

func (r *DimensionRepository) GetChildRegions(parentCode string) ([]model.DimRegion, error) {
	var regions []model.DimRegion
	err := DB.Where("parent_code = ?", parentCode).Find(&regions).Error
	return regions, err
}

func (r *DimensionRepository) GetChildBusinesses(parentCode string) ([]model.DimBusiness, error) {
	var businesses []model.DimBusiness
	err := DB.Where("parent_code = ?", parentCode).Find(&businesses).Error
	return businesses, err
}
