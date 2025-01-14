package dns

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/li4n0/revsuit/internal/database"
	"github.com/li4n0/revsuit/internal/newdns"
	"github.com/li4n0/revsuit/internal/rule"
	"gorm.io/gorm/clause"
	log "unknwon.dev/clog/v2"
)

type Rule struct {
	rule.BaseRule `yaml:",inline"`
	Type          newdns.Type   `gorm:"default:1" form:"type" json:"type"`
	Value         string        `form:"value" json:"value"`
	TTL           time.Duration `gorm:"ttl;default:10" form:"ttl" json:"ttl"`
}

func (Rule) TableName() string {
	return "dns_rules"
}

// NewRule new dns rule struct
func NewRule(name, flagFormat, value string, pushToClient, notice bool, _type newdns.Type, ttl time.Duration) *Rule {
	return &Rule{
		BaseRule: rule.BaseRule{
			Name:         name,
			FlagFormat:   flagFormat,
			PushToClient: pushToClient,
			Notice:       notice,
		},
		Type:  _type,
		Value: value,
		TTL:   ttl,
	}
}

// CreateOrUpdate creates or updates the dns rule in database and ruleSet
func (r *Rule) CreateOrUpdate() (err error) {
	db := database.DB.Model(r)
	err = db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns(
			[]string{
				"name",
				"flag_format",
				"base_rank",
				"type",
				"value",
				"ttl",
				"push_to_client",
				"notice",
			}),
	}).Create(r).Error
	if err != nil {
		return
	}
	err = GetServer().UpdateRules()
	return err
}

// Delete deletes the dns rule in database and ruleSet
func (r *Rule) Delete() (err error) {
	db := database.DB.Model(r)
	err = db.Delete(r).Error
	if err != nil {
		return
	}
	err = GetServer().UpdateRules()
	return err
}

// ListRules lists all dns rules those satisfy the filter
func ListRules(c *gin.Context) {
	var (
		dnsRule  Rule
		res      []Rule
		count    int64
		order    = c.Query("order")
		pageSize = 10
	)

	if c.Query("pageSize") != "" {
		if n, err := strconv.Atoi(c.Query("pageSize")); err == nil {
			if n > 0 && n < 100 {
				pageSize = n
			}
		}
	}

	if err := c.ShouldBind(&dnsRule); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"result": nil,
		})
		return
	}

	db := database.DB.Model(&dnsRule)
	if dnsRule.Name != "" {
		db.Where("name = ?", dnsRule.Name)
	}
	db.Count(&count)

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"result": nil,
		})
		return
	}

	if order != "asc" {
		order = "desc"
	}

	if err := db.Order("base_rank desc").Order("id" + " " + order).Count(&count).Offset((page - 1) * pageSize).Limit(pageSize).Find(&res).Error; err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"data":   nil,
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "succeed",
		"error":  nil,
		"result": gin.H{"count": count, "data": res},
	})
}

// UpsertRules creates or updates dns rule from user submit
func UpsertRules(c *gin.Context) {
	var (
		dnsRule Rule
		update  bool
	)

	if err := c.ShouldBind(&dnsRule); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"data":   nil,
		})
		return
	}

	if dnsRule.ID != 0 {
		update = true
	}

	if err := dnsRule.CreateOrUpdate(); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"data":   nil,
		})
		return
	}

	if update {
		log.Trace("DNS rule[id:%d] has been updated", dnsRule.ID)
	} else {
		log.Trace("DNS rule[id:%d] has been created", dnsRule.ID)
	}

	c.JSON(200, gin.H{
		"status": "succeed",
		"error":  nil,
		"result": nil,
	})
}

// DeleteRules Delete dns rule from user submit
func DeleteRules(c *gin.Context) {
	var dnsRule Rule

	if err := c.ShouldBind(&dnsRule); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"data":   nil,
		})
		return
	}

	if err := dnsRule.Delete(); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"data":   nil,
		})
		return
	}

	log.Trace("DNS rule[id:%d] has been deleted", dnsRule.ID)

	c.JSON(200, gin.H{
		"status": "succeed",
		"error":  nil,
		"data":   nil,
	})
}
