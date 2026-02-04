package models

import (
	"strings"

	"uu/config"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UItem struct {
	HashName string `json:"commodityHashName" gorm:"type:varchar(255);uniqueIndex"`
	Name     string `json:"commodityName" gorm:"type:varchar(255)"`
	IconUrl  string `json:"iconUrl" gorm:"type:varchar(255)"`
	Id       int64  `json:"id" gorm:"type:int;primaryKey"`
	//LongLeaseUnitPrice string `json:"longLeaseUnitPrice" gorm:"type:decimal(10,2)"`
	//OnLeaseCount       int64  `json:"onLeaseCount" gorm:"type:int"`
	Count int64  `json:"onSaleCount" gorm:"type:int"`
	Price string `json:"price" gorm:"type:decimal(10,2)"`
	//RarityColor        string `json:"rarityColor" gorm:"type:varchar(20)"`
	//Rent               string `json:"rent" gorm:"type:decimal(10,2)"`
	//SortId             int64  `json:"sortId" gorm:"type:int"`
	SteamPrice string `json:"steamPrice" gorm:"type:decimal(10,2)"`
	TypeName   string `json:"typeName" gorm:"type:varchar(255)"`
}

type UItemsInfo struct {
	MarketHashName      string  `json:"marketHashName" gorm:"type:varchar(255);uniqueIndex"`
	Name                string  `json:"Name" gorm:"type:varchar(255)"`
	ImageUrl            string  `json:"imageUrl" gorm:"type:varchar(255)"`
	Id                  int64   `gorm:"type:int;primaryKey"`
	CacheExpirationDesc string  `json:"cacheExpirationDesc" gorm:"type:varchar(20)"`
	AssetMergeCount     int64   `json:"assetMergeCount" gorm:"type:int"`
	Price               float64 `json:"price" gorm:"type:decimal(10,2)"`
}

type UBaseInfo struct {
	Id          int    `json:"item_id" gorm:"primaryKey"`
	HashName    string `json:"hash_name" gorm:"type:varchar(255);uniqueIndex;not null"`
	IconUrl     string `json:"icon_url" gorm:"index"`
	RarityName  string `json:"rarity_name"`
	QualityName string `json:"quality_name"`
	TypeName    string `json:"type_name" gorm:"type:varchar(50)"`
}

// TableName 指定表名为 u_base_info
func (UBaseInfo) TableName() string {
	return "u_base_info"
}

// 武器类型映射
var weaponTypeMap = map[string]string{
	// 匕首
	"Butterfly Knife": "匕首",
	"Karambit":        "匕首",
	"M9 Bayonet":      "匕首",
	"Skeleton Knife":  "匕首",
	"Bayonet":         "匕首",
	"Flip Knife":      "匕首",
	"Stiletto Knife":  "匕首",
	"Talon Knife":     "匕首",
	"Nomad Knife":     "匕首",
	"Ursus Knife":     "匕首",
	"Classic Knife":   "匕首",
	"Huntsman Knife":  "匕首",
	"Paracord Knife":  "匕首",
	"Survival Knife":  "匕首",
	"Falchion Knife":  "匕首",
	"Shadow Daggers":  "匕首",
	"Bowie Knife":     "匕首",
	"Gut Knife":       "匕首",
	"Navaja Knife":    "匕首",
	"Kukri Knife":     "匕首",

	// 手套
	"Sport Gloves":       "手套",
	"Specialist Gloves":  "手套",
	"Moto Gloves":        "手套",
	"Driver Gloves":      "手套",
	"Hand Wraps":         "手套",
	"Broken Fang Gloves": "手套",
	"Hydra Gloves":       "手套",
	"Bloodhound Gloves":  "手套",

	// 步枪
	"AK-47":    "步枪",
	"AWP":      "步枪",
	"M4A1-S":   "步枪",
	"M4A4":     "步枪",
	"Galil AR": "步枪",
	"FAMAS":    "步枪",
	"SSG 08":   "步枪",
	"AUG":      "步枪",
	"SG 553":   "步枪",
	"SCAR-20":  "步枪",
	"G3SG1":    "步枪",

	// 手枪
	"Desert Eagle":  "手枪",
	"USP-S":         "手枪",
	"Glock-18":      "手枪",
	"Tec-9":         "手枪",
	"Five-SeveN":    "手枪",
	"P250":          "手枪",
	"Dual Berettas": "手枪",
	"CZ75-Auto":     "手枪",
	"R8 Revolver":   "手枪",
	"P2000":         "手枪",

	// 微型冲锋枪
	"MP9":      "微型冲锋枪",
	"MAC-10":   "微型冲锋枪",
	"P90":      "微型冲锋枪",
	"UMP-45":   "微型冲锋枪",
	"MP7":      "微型冲锋枪",
	"PP-Bizon": "微型冲锋枪",
	"MP5-SD":   "微型冲锋枪",

	// 霰弹枪
	"MAG-7":     "霰弹枪",
	"XM1014":    "霰弹枪",
	"Sawed-Off": "霰弹枪",
	"Nova":      "霰弹枪",

	// 机枪
	"Negev": "机枪",
	"M249":  "机枪",
}

// InferTypeFromHashName 根据 hash_name 推断饰品类型
func InferTypeFromHashName(hashName string) string {
	// 去掉前缀 ★, StatTrak™, Souvenir
	cleanName := hashName
	cleanName = strings.TrimPrefix(cleanName, "★ StatTrak™ ")
	cleanName = strings.TrimPrefix(cleanName, "★ ")
	cleanName = strings.TrimPrefix(cleanName, "StatTrak™ ")
	cleanName = strings.TrimPrefix(cleanName, "Souvenir ")

	// 检查特殊类型（通过前缀判断）
	if strings.HasPrefix(hashName, "Sticker |") || strings.HasPrefix(hashName, "Sticker Slab |") {
		return "印花"
	}
	if strings.HasPrefix(cleanName, "Sealed Graffiti |") || strings.HasPrefix(cleanName, "Graffiti |") {
		return "涂鸦"
	}
	if strings.HasPrefix(cleanName, "Music Kit |") {
		return "音乐盒"
	}
	if strings.HasPrefix(cleanName, "Patch |") {
		return "布章"
	}
	if strings.HasPrefix(cleanName, "Charm |") {
		return "挂件"
	}
	if strings.Contains(hashName, "Case") && !strings.Contains(hashName, "|") {
		return "武器箱"
	}
	if strings.Contains(hashName, "Capsule") && !strings.Contains(hashName, "|") {
		return "胶囊"
	}
	if strings.Contains(hashName, "Package") && !strings.Contains(hashName, "|") {
		return "包裹"
	}
	if strings.Contains(hashName, "Key") && !strings.Contains(hashName, "|") {
		return "钥匙"
	}

	// 提取武器名称（格式：WeaponName | SkinName (Wear)）
	parts := strings.Split(cleanName, " | ")
	if len(parts) >= 1 {
		weaponName := strings.TrimSpace(parts[0])

		// 直接匹配武器名称
		if typeName, ok := weaponTypeMap[weaponName]; ok {
			return typeName
		}

		// 部分匹配（处理如 "Butterfly Knife" 在 hash_name 中的情况）
		for weapon, typeName := range weaponTypeMap {
			if strings.Contains(weaponName, weapon) {
				return typeName
			}
		}
	}

	// 探员检查
	if strings.Contains(hashName, "Agent") || strings.Contains(hashName, "Master Agent") ||
		strings.Contains(hashName, "| FBI") || strings.Contains(hashName, "| SAS") ||
		strings.Contains(hashName, "| SEAL") || strings.Contains(hashName, "| KSK") ||
		strings.Contains(hashName, "| Phoenix") || strings.Contains(hashName, "| Elite Crew") ||
		strings.Contains(hashName, "| Professional") || strings.Contains(hashName, "| Ground Rebel") {
		return "探员"
	}

	return "其他"
}

func BatchAddUUItem(uu []*UItem) {
	err := config.DB.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(uu, 100).Error
	if err != nil {
		config.Log.Errorf("batch insert uu item fail: %s", err)
	}
}

func BatchAddUUInventory(uu []*UItemsInfo) {
	err := config.DB.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(uu, 100).Error
	if err != nil {
		config.Log.Errorf("barch insert uu inventory fail: %s", err)
	}
}

// -------------------------------------------------------v2------------------------------------------------------------
// data from steamDT

type U struct {
	Id             string  `json:"platformItemId" gorm:"primaryKey"`
	MarketHashName string  `json:"marketHashName" gorm:"type:varchar(255);uniqueIndex;not null"`
	SellPrice      float64 `json:"sellPrice" gorm:"index"`
	SellCount      int64   `json:"sellCount" gorm:"index"`
	BiddingPrice   float64 `json:"biddingPrice" gorm:"index"`
	BiddingCount   int64   `json:"biddingCount"`
	UpdateTime     int64   `json:"updateTime"`
	BeforeTime     int64   `json:"beforeTime"`
	BeforeCount    int64   `json:"beforeCount"`
	TurnOver       int64   `json:"turn_over"`
	Link           string  `json:"link"`
}

func BatchGetUUGoods(hashNames []string) map[string]*U {
	var uList []U
	err := config.DB.Where("market_hash_name in ?", hashNames).Find(&uList).Error
	if err != nil {
		config.Log.Errorf("Batch get uu goods error: %v", err)
	}
	result := make(map[string]*U)
	for i := range uList {
		result[uList[i].MarketHashName] = &uList[i]
	}
	return result
}

func BatchUpdateUUGoods(uu []*U) {
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(uu, 100).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		config.Log.Errorf("Update UU Goods fail: %v", err)
		return
	}
	//config.Log.Info("Update UU Goods Success")
}

func BatchQueryHashIcon() ([]UBaseInfo, error) {
	var Infos []UBaseInfo
	err := config.DB.Select("hash_name, icon_url").Find(&Infos).Error
	if err != nil {
		config.Log.Errorf("Get uu icon url error: %v", err)
	}
	return Infos, err
}

func BatchUpdateIcon(infos []*UBaseInfo) {
	// 在插入前设置类型
	for i := range infos {
		if infos[i].TypeName == "" {
			infos[i].TypeName = InferTypeFromHashName(infos[i].HashName)
		}
	}

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(infos, 200).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		config.Log.Errorf("Update UU Goods Icon fail: %v", err)
		return
	}
}

func QueryAllUUHashName() []string {
	var hashNames []string
	err := config.DB.Model(&U{}).Pluck("market_hash_name", &hashNames).Error
	if err != nil {
		config.Log.Errorf("Get all uu hash name error: %v", err)
	}
	return hashNames
}
