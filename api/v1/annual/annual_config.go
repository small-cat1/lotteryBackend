package annual

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"lotteryBackend/global"
	"lotteryBackend/model/annual"
	annualReq "lotteryBackend/model/annual/request"
	"lotteryBackend/model/common/response"
	"lotteryBackend/service"
	"lotteryBackend/service/annual/dto"
)

type AnnualConfigApi struct{}

var annualConfigService = service.ServiceGroupApp.AnnualServiceGroup.AnnualConfigService

// GetConfigList 获取配置列表
func (a *AnnualConfigApi) GetConfigList(c *gin.Context) {
	list, err := annualConfigService.GetConfigList()
	if err != nil {
		global.GVA_LOG.Error("获取配置列表失败!", zap.Error(err))
		response.FailWithMessage("获取配置列表失败: "+err.Error(), c)
		return
	}
	response.OkWithData(list, c)
}

// GetAllConfig 获取所有配置（按分组）
func (a *AnnualConfigApi) GetAllConfig(c *gin.Context) {
	config, err := annualConfigService.GetAllConfig()
	if err != nil {
		global.GVA_LOG.Error("获取配置失败!", zap.Error(err))
		response.FailWithMessage("获取配置失败: "+err.Error(), c)
		return
	}
	response.OkWithData(config, c)
}

// GetConfigByKey 根据key获取配置
func (a *AnnualConfigApi) GetConfigByKey(c *gin.Context) {
	key := c.Param("key")
	config, err := annualConfigService.GetConfigByKey(key)
	if err != nil {
		global.GVA_LOG.Error("获取配置失败!", zap.Error(err))
		response.FailWithMessage("获取配置失败: "+err.Error(), c)
		return
	}
	response.OkWithData(config, c)
}

// SetConfig 设置配置
func (a *AnnualConfigApi) SetConfig(c *gin.Context) {
	var config annual.AnnualConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := annualConfigService.SetConfig(config.ConfigKey, config.ConfigValue, config.Description); err != nil {
		global.GVA_LOG.Error("设置配置失败!", zap.Error(err))
		response.FailWithMessage("设置配置失败: "+err.Error(), c)
		return
	}
	// 清除缓存
	annualConfigService.ClearConfigCache()
	response.OkWithMessage("设置成功", c)
}

// BatchSetConfig 批量设置配置
func (a *AnnualConfigApi) BatchSetConfig(c *gin.Context) {
	var configs []annual.AnnualConfig
	if err := c.ShouldBindJSON(&configs); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	for _, config := range configs {
		if err := annualConfigService.SetConfig(config.ConfigKey, config.ConfigValue, config.Description); err != nil {
			global.GVA_LOG.Error("设置配置失败!", zap.Error(err))
			response.FailWithMessage("设置配置失败: "+err.Error(), c)
			return
		}
	}
	// 清除缓存
	annualConfigService.ClearConfigCache()
	response.OkWithMessage("设置成功", c)
}

// SetWechatConfig 设置微信配置
func (a *AnnualConfigApi) SetWechatConfig(c *gin.Context) {
	var req annualReq.WechatConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	config := dto.WechatConfig{
		AppID:  req.AppID,
		Secret: req.Secret,
		Token:  req.Token,
	}
	if err := annualConfigService.SetWechatConfig(config); err != nil {
		global.GVA_LOG.Error("设置微信配置失败!", zap.Error(err))
		response.FailWithMessage("设置微信配置失败: "+err.Error(), c)
		return
	}
	annualConfigService.ClearConfigCache()
	response.OkWithMessage("设置成功", c)
}

// GetWechatConfig 获取微信配置
func (a *AnnualConfigApi) GetWechatConfig(c *gin.Context) {
	config, err := annualConfigService.GetWechatConfig()
	if err != nil {
		global.GVA_LOG.Error("获取微信配置失败!", zap.Error(err))
		response.FailWithMessage("获取微信配置失败: "+err.Error(), c)
		return
	}
	// 隐藏敏感信息
	if len(config.Secret) > 8 {
		config.Secret = config.Secret[:4] + "****" + config.Secret[len(config.Secret)-4:]
	}
	response.OkWithData(config, c)
}

// SetSiteConfig 设置站点配置
func (a *AnnualConfigApi) SetSiteConfig(c *gin.Context) {
	var req annualReq.SiteConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	config := dto.SiteConfig{
		Name:      req.Name,
		Logo:      req.Logo,
		Copyright: req.Copyright,
	}
	if err := annualConfigService.SetSiteConfig(config); err != nil {
		global.GVA_LOG.Error("设置站点配置失败!", zap.Error(err))
		response.FailWithMessage("设置站点配置失败: "+err.Error(), c)
		return
	}
	annualConfigService.ClearConfigCache()
	response.OkWithMessage("设置成功", c)
}

// GetSiteConfig 获取站点配置
func (a *AnnualConfigApi) GetSiteConfig(c *gin.Context) {
	config, err := annualConfigService.GetSiteConfig()
	if err != nil {
		global.GVA_LOG.Error("获取站点配置失败!", zap.Error(err))
		response.FailWithMessage("获取站点配置失败: "+err.Error(), c)
		return
	}
	response.OkWithData(config, c)
}

// DeleteConfig 删除配置
func (a *AnnualConfigApi) DeleteConfig(c *gin.Context) {
	key := c.Param("key")
	if err := annualConfigService.DeleteConfig(key); err != nil {
		global.GVA_LOG.Error("删除配置失败!", zap.Error(err))
		response.FailWithMessage("删除配置失败: "+err.Error(), c)
		return
	}
	annualConfigService.ClearConfigCache()
	response.OkWithMessage("删除成功", c)
}
