package service

import (
	"crypto/rand"
	"dashboard/internal/model"
	"dashboard/internal/repository"
	"database/sql"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
)

type DataSeedService struct{}

func NewDataSeedService() *DataSeedService {
	return &DataSeedService{}
}

func (s *DataSeedService) SeedAllData() error {
	db := repository.GetDB()

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM dim_time").Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	if err := s.seedDimTime(); err != nil {
		return err
	}
	if err := s.seedDimRegion(); err != nil {
		return err
	}
	if err := s.seedDimBusiness(); err != nil {
		return err
	}
	if err := s.seedRawData(100000); err != nil {
		return err
	}
	if err := s.seedAggregatedData(); err != nil {
		return err
	}

	return nil
}

func (s *DataSeedService) seedDimTime() error {
	db := repository.GetDB()
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO dim_time (date, year, quarter, month, week, month_name) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	current := startDate
	for !current.After(endDate) {
		_, week := current.ISOWeek()
		_, err := stmt.Exec(
			current.Format("2006-01-02"),
			current.Year(),
			(int(current.Month())-1)/3+1,
			int(current.Month()),
			week,
			current.Month().String(),
		)
		if err != nil {
			return err
		}
		current = current.AddDate(0, 0, 1)
	}

	return tx.Commit()
}

func (s *DataSeedService) seedDimRegion() error {
	db := repository.GetDB()

	regions := []model.DimRegion{
		{RegionCode: "CN", RegionName: "中国", Level: 1, Path: "CN"},
		{RegionCode: "CN-HD", RegionName: "华东", ParentCode: strPtr("CN"), Level: 2, Path: "CN/CN-HD"},
		{RegionCode: "CN-HB", RegionName: "华北", ParentCode: strPtr("CN"), Level: 2, Path: "CN/CN-HB"},
		{RegionCode: "CN-HN", RegionName: "华南", ParentCode: strPtr("CN"), Level: 2, Path: "CN/CN-HN"},
		{RegionCode: "CN-HZ", RegionName: "华中", ParentCode: strPtr("CN"), Level: 2, Path: "CN/CN-HZ"},
		{RegionCode: "CN-XB", RegionName: "西北", ParentCode: strPtr("CN"), Level: 2, Path: "CN/CN-XB"},
		{RegionCode: "CN-XN", RegionName: "西南", ParentCode: strPtr("CN"), Level: 2, Path: "CN/CN-XN"},
		{RegionCode: "CN-DB", RegionName: "东北", ParentCode: strPtr("CN"), Level: 2, Path: "CN/CN-DB"},
		{RegionCode: "CN-HD-SH", RegionName: "上海", ParentCode: strPtr("CN-HD"), Level: 3, Path: "CN/CN-HD/CN-HD-SH"},
		{RegionCode: "CN-HD-JS", RegionName: "江苏", ParentCode: strPtr("CN-HD"), Level: 3, Path: "CN/CN-HD/CN-HD-JS"},
		{RegionCode: "CN-HD-ZJ", RegionName: "浙江", ParentCode: strPtr("CN-HD"), Level: 3, Path: "CN/CN-HD/CN-HD-ZJ"},
		{RegionCode: "CN-HB-BJ", RegionName: "北京", ParentCode: strPtr("CN-HB"), Level: 3, Path: "CN/CN-HB/CN-HB-BJ"},
		{RegionCode: "CN-HB-TJ", RegionName: "天津", ParentCode: strPtr("CN-HB"), Level: 3, Path: "CN/CN-HB/CN-HB-TJ"},
		{RegionCode: "CN-HN-GD", RegionName: "广东", ParentCode: strPtr("CN-HN"), Level: 3, Path: "CN/CN-HN/CN-HN-GD"},
		{RegionCode: "CN-HN-SZ", RegionName: "深圳", ParentCode: strPtr("CN-HN"), Level: 3, Path: "CN/CN-HN/CN-HN-SZ"},
		{RegionCode: "CN-HD-SH-PD", RegionName: "浦东", ParentCode: strPtr("CN-HD-SH"), Level: 4, Path: "CN/CN-HD/CN-HD-SH/CN-HD-SH-PD"},
		{RegionCode: "CN-HD-SH-HP", RegionName: "黄浦", ParentCode: strPtr("CN-HD-SH"), Level: 4, Path: "CN/CN-HD/CN-HD-SH/CN-HD-SH-HP"},
		{RegionCode: "CN-HD-JS-NJ", RegionName: "南京", ParentCode: strPtr("CN-HD-JS"), Level: 4, Path: "CN/CN-HD/CN-HD-JS/CN-HD-JS-NJ"},
		{RegionCode: "CN-HD-JS-SZ", RegionName: "苏州", ParentCode: strPtr("CN-HD-JS"), Level: 4, Path: "CN/CN-HD/CN-HD-JS/CN-HD-JS-SZ"},
		{RegionCode: "CN-HB-BJ-CY", RegionName: "朝阳", ParentCode: strPtr("CN-HB-BJ"), Level: 4, Path: "CN/CN-HB/CN-HB-BJ/CN-HB-BJ-CY"},
		{RegionCode: "CN-HB-BJ-HD", RegionName: "海淀", ParentCode: strPtr("CN-HB-BJ"), Level: 4, Path: "CN/CN-HB/CN-HB-BJ/CN-HB-BJ-HD"},
		{RegionCode: "CN-HN-GD-GZ", RegionName: "广州", ParentCode: strPtr("CN-HN-GD"), Level: 4, Path: "CN/CN-HN/CN-HN-GD/CN-HN-GD-GZ"},
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO dim_region (region_code, region_name, parent_code, level, path) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, r := range regions {
		_, err := stmt.Exec(r.RegionCode, r.RegionName, r.ParentCode, r.Level, r.Path)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *DataSeedService) seedDimBusiness() error {
	db := repository.GetDB()

	businesses := []model.DimBusiness{
		{BusinessCode: "ALL", BusinessName: "全部业务", Level: 1, Category: "总览"},
		{BusinessCode: "RETAIL", BusinessName: "零售业务", ParentCode: strPtr("ALL"), Level: 2, Category: "ToC"},
		{BusinessCode: "WHOLESALE", BusinessName: "批发业务", ParentCode: strPtr("ALL"), Level: 2, Category: "ToB"},
		{BusinessCode: "ONLINE", BusinessName: "线上业务", ParentCode: strPtr("ALL"), Level: 2, Category: "电商"},
		{BusinessCode: "RETAIL-FOOD", BusinessName: "食品零售", ParentCode: strPtr("RETAIL"), Level: 3, Category: "消费品"},
		{BusinessCode: "RETAIL-ELEC", BusinessName: "电子零售", ParentCode: strPtr("RETAIL"), Level: 3, Category: "3C"},
		{BusinessCode: "RETAIL-CLOTH", BusinessName: "服装零售", ParentCode: strPtr("RETAIL"), Level: 3, Category: "服饰"},
		{BusinessCode: "WHOLESALE-RAW", BusinessName: "原料批发", ParentCode: strPtr("WHOLESALE"), Level: 3, Category: "原材料"},
		{BusinessCode: "WHOLESALE-FIN", BusinessName: "成品批发", ParentCode: strPtr("WHOLESALE"), Level: 3, Category: "成品"},
		{BusinessCode: "ONLINE-B2C", BusinessName: "B2C商城", ParentCode: strPtr("ONLINE"), Level: 3, Category: "直营"},
		{BusinessCode: "ONLINE-MKT", BusinessName: "营销活动", ParentCode: strPtr("ONLINE"), Level: 3, Category: "市场"},
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO dim_business (business_code, business_name, parent_code, level, category) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, b := range businesses {
		_, err := stmt.Exec(b.BusinessCode, b.BusinessName, b.ParentCode, b.Level, b.Category)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *DataSeedService) seedRawData(count int) error {
	db := repository.GetDB()

	regionCodes := []string{
		"CN-HD-SH-PD", "CN-HD-SH-HP", "CN-HD-JS-NJ", "CN-HD-JS-SZ",
		"CN-HB-BJ-CY", "CN-HB-BJ-HD", "CN-HN-GD-GZ", "CN-HN-SZ",
	}
	businessCodes := []string{
		"RETAIL-FOOD", "RETAIL-ELEC", "RETAIL-CLOTH",
		"WHOLESALE-RAW", "WHOLESALE-FIN",
		"ONLINE-B2C", "ONLINE-MKT",
	}
	products := []string{
		"智能手机", "笔记本电脑", "蓝牙耳机", "运动服饰",
		"休闲食品", "生鲜水果", "家用电器", "护肤美妆",
		"家居用品", "办公用品",
	}

	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
	totalDays := int(endDate.Sub(startDate).Hours() / 24)

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO raw_data (trace_no, date, region_code, business_code, order_no, amount, user_id, product_name, quantity, trade_time)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	batchSize := 1000
	for i := 0; i < count; i++ {
		daysOffset, _ := rand.Int(rand.Reader, big.NewInt(int64(totalDays)))
		tradeDate := startDate.AddDate(0, 0, int(daysOffset.Int64()))

		regionIdx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(regionCodes))))
		businessIdx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(businessCodes))))
		productIdx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(products))))

		amount, _ := rand.Int(rand.Reader, big.NewInt(10000))
		amountFloat := float64(amount.Int64()) + 50.0

		quantity, _ := rand.Int(rand.Reader, big.NewInt(10))
		userID := fmt.Sprintf("U%06d", i%10000+1)
		orderNo := fmt.Sprintf("ORD%s%08d", tradeDate.Format("20060102"), i+1)

		hour, _ := rand.Int(rand.Reader, big.NewInt(24))
		minute, _ := rand.Int(rand.Reader, big.NewInt(60))
		tradeTime := time.Date(tradeDate.Year(), tradeDate.Month(), tradeDate.Day(), int(hour.Int64()), int(minute.Int64()), 0, 0, time.UTC)

		_, err := stmt.Exec(
			uuid.New().String(),
			tradeDate.Format("2006-01-02"),
			regionCodes[regionIdx.Int64()],
			businessCodes[businessIdx.Int64()],
			orderNo,
			amountFloat,
			userID,
			products[productIdx.Int64()],
			int(quantity.Int64())+1,
			tradeTime.Format("2006-01-02 15:04:05"),
		)
		if err != nil {
			return err
		}

		if (i+1)%batchSize == 0 {
			if err := tx.Commit(); err != nil {
				return err
			}
			tx, err = db.Begin()
			if err != nil {
				return err
			}
			stmt, err = tx.Prepare(`
				INSERT INTO raw_data (trace_no, date, region_code, business_code, order_no, amount, user_id, product_name, quantity, trade_time)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			`)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (s *DataSeedService) seedAggregatedData() error {
	db := repository.GetDB()

	for aggLevel := 1; aggLevel <= 4; aggLevel++ {
		sql := fmt.Sprintf(`
			INSERT INTO agg_data (date, region_code, business_code, sales, orders, users, amount, agg_level)
			SELECT 
				r.date,
				r.region_code,
				r.business_code,
				SUM(r.amount) as sales,
				COUNT(*) as orders,
				COUNT(DISTINCT r.user_id) as users,
				SUM(r.amount) as amount,
				%d as agg_level
			FROM raw_data r
			JOIN dim_region reg ON r.region_code = reg.region_code
			JOIN dim_business biz ON r.business_code = biz.business_code
			WHERE reg.level <= %d AND biz.level <= %d
			GROUP BY r.date, r.region_code, r.business_code
		`, aggLevel, 5-aggLevel, 4-aggLevel)

		_, err := db.Exec(sql)
		if err != nil {
			return err
		}
	}

	return nil
}

func strPtr(s string) *string {
	return &s
}
