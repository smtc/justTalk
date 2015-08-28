package main

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/smtc/glog"
)

var (
	// 未分类
	UnCategory = Taxonomy{
		Id:          "0",
		Name:        "UnCategory",
		Slug:        "/uncategory",
		Description: "uncategory",
		Parent:      "",
	}
)

// create post提交的的json数据字段
type TaxInfo struct {
	Taxonomy string `json:"taxonmy"`
	TaxName  string `json:"tax_name"`
	TaxId    string `json:"tax_id"`
}

// 分类，类似于wordpress的分类，但是把wp_term_taxonomy和wp_terms合并为一个
type Taxonomy struct {
	Id          string `sql:"size:60" gorm:"primay_key" json:"id"`
	SiteId      int64  `json:"site_id"`
	Name        string `sql:"size:200" json:"name"`
	Slug        string `sql:"size:200" json:"slug"`
	TermGroup   int    `json:"term_group"`
	Taxonomy    string `sql:"size:60" json:"taxonomy"`
	Description string `sql:"size:100000" json:"description"`
	Parent      string `json:"parent`
	Count       int64  `json:"count"`
}

type TermRelation struct {
	ObjectId  string    `json:"object_id"` // post id, reply id, product id, etc...
	SiteId    int64     `json:"site_id"`
	TermId    string    `json:"term_id"`
	CreatedAt time.Time `json:"created_at"`
	TermOrder int       `json:"term_order"`
}

// 根据id查找taxonomy
func getTaxById(id string) *Taxonomy {
	var term Taxonomy

	err := db.Where("id=?", id).Find(&term).Error
	if err != nil {
		return nil
	}
	return &term
}

// 根据名称查找taxonomy
func getTaxByName(name, tax string) (*Taxonomy, error) {
	var term Taxonomy

	err := db.Where("name=? AND taxonomy=?", name, tax).Find(&term).Error
	if err == gorm.RecordNotFound {
		err = nil
	}
	return &term, nil
}

// 获取文章分类的大类
func getAllCategory() ([]*Taxonomy, error) {
	var terms []*Taxonomy

	err := db.Where("term_group=category").Find(&terms).Error
	if err == gorm.RecordNotFound {
		err = nil
	}

	return terms, err
}

// 获取一级类别下的所有子类
func getAllTaxs(tax string) ([]*Taxonomy, error) {
	var terms []*Taxonomy

	err := db.Where("taxonomy=?", tax).Find(&terms).Error
	if err == gorm.RecordNotFound {
		err = nil
	}
	return terms, err
}

// 获取某个类别下的所有post id
func getObjectsByTerm(term *Taxonomy, start, count int) ([]*TermRelation, error) {
	var rel []*TermRelation

	err := db.Where("term_id=?", term.Id).Offset(start).Order("create_at desc").Limit(count).Find(&rel).Error
	return rel, err
}

// 设置post的分类
func setPostTaxonmoy(post *Post, taxes []*Taxonomy) error {
	var (
		tr  TermRelation
		now = time.Now()
		err error
	)

	for _, tax := range taxes {
		tr = TermRelation{
			ObjectId:  post.Id,
			TermId:    tax.Id,
			CreatedAt: now,
		}
		err = db.Create(&tr).Error
		if err != nil {
			glog.Error("set post %s to taxonomy %s (%s %s) failed: %s\n",
				post.Id, tax.Id, tax.Name)
		}
	}
	return err
}

// 从[]TaxInfo中查找对应的Taxonomy, 并去重
func getTaxFromInfos(infos []TaxInfo) []*Taxonomy {
	var (
		err      error
		tax      *Taxonomy
		taxesMap = make(map[*Taxonomy]struct{})
		taxes    = []*Taxonomy{}
	)

	for _, taxInfo := range infos {
		//首先根据tax_id来查找
		if taxInfo.TaxId != "" {
			if taxInfo.TaxId == "0" {
				tax = &UnCategory
			} else {
				tax = getTaxById(taxInfo.TaxId)
			}
			taxesMap[tax] = struct{}{}
		} else {
			// 根据taxonomy和name来查找
			tax, err = getTaxByName(taxInfo.TaxName, taxInfo.Taxonomy)
			if err != nil {
				glog.Error("Not found category by name %s taxonomy %s\n",
					taxInfo.TaxName, taxInfo.Taxonomy)
			} else {
				taxesMap[tax] = struct{}{}
			}
		}
	}
	// 使用map来过滤可能出现的重复taxonomy问题
	for key, _ := range taxesMap {
		taxes = append(taxes, key)
	}

	// 如果没有找到任何分类，设置为未分类
	if len(taxes) == 0 {
		taxes = []*Taxonomy{&UnCategory}
	}
	return taxes
}

// 创建分类
func createTaxonomy(tax *Taxonomy, user *User) (err error) {
	if user == nil {
		return fmt.Errorf("create taxonomy NEED authencation")
	}
	user.parseUserCap()
	if user.capability["create_taxonomy"] == false {
		return fmt.Errorf("no authorization to create taxonomy")
	}

	return db.Create(tax).Error
}
