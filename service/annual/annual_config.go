package annual

import (
	"encoding/json"
	"lotteryBackend/global"
	"lotteryBackend/model/annual"
	"lotteryBackend/service/annual/dto"
)

type AnnualConfigService struct{}

// GetConfigList 获取所有配置
func (s *AnnualConfigService) GetConfigList() (list []annual.AnnualConfig, err error) {
	err = global.GVA_DB.Order("id ASC").Find(&list).Error
	return
}

// GetConfigByKey 根据key获取配置
func (s *AnnualConfigService) GetConfigByKey(key string) (config annual.AnnualConfig, err error) {
	err = global.GVA_DB.Where("config_key = ?", key).First(&config).Error
	return
}

// GetConfigValueByKey 根据key获取配置值
func (s *AnnualConfigService) GetConfigValueByKey(key string) (value string, err error) {
	var config annual.AnnualConfig
	err = global.GVA_DB.Where("config_key = ?", key).First(&config).Error
	if err != nil {
		return "", err
	}
	return config.ConfigValue, nil
}

// SetConfig 设置配置（存在则更新，不存在则创建）
func (s *AnnualConfigService) SetConfig(key, value, description string) (err error) {
	var config annual.AnnualConfig
	result := global.GVA_DB.Where("config_key = ?", key).First(&config)

	if result.RowsAffected == 0 {
		// 不存在，创建
		config = annual.AnnualConfig{
			ConfigKey:   key,
			ConfigValue: value,
			Description: description,
		}
		return global.GVA_DB.Create(&config).Error
	}

	// 存在，更新
	return global.GVA_DB.Model(&config).Updates(map[string]interface{}{
		"config_value": value,
		"description":  description,
	}).Error
}

// BatchSetConfig 批量设置配置
func (s *AnnualConfigService) BatchSetConfig(configs map[string]string) (err error) {
	for key, value := range configs {
		if err = s.SetConfig(key, value, ""); err != nil {
			return err
		}
	}
	return nil
}

// DeleteConfig 删除配置
func (s *AnnualConfigService) DeleteConfig(key string) (err error) {
	return global.GVA_DB.Where("config_key = ?", key).Delete(&annual.AnnualConfig{}).Error
}

// GetWechatConfig 获取微信配置
func (s *AnnualConfigService) GetWechatConfig() (config dto.WechatConfig, err error) {
	appid, _ := s.GetConfigValueByKey("wechat_appid")
	secret, _ := s.GetConfigValueByKey("wechat_secret")
	token, _ := s.GetConfigValueByKey("wechat_token")

	config = dto.WechatConfig{
		AppID:  appid,
		Secret: secret,
		Token:  token,
	}
	return
}

// SetWechatConfig 设置微信配置
func (s *AnnualConfigService) SetWechatConfig(config dto.WechatConfig) (err error) {
	configs := map[string]string{
		"wechat_appid":  config.AppID,
		"wechat_secret": config.Secret,
		"wechat_token":  config.Token,
	}
	return s.BatchSetConfig(configs)
}

// GetSiteConfig 获取站点配置
func (s *AnnualConfigService) GetSiteConfig() (config dto.SiteConfig, err error) {
	name, _ := s.GetConfigValueByKey("site_name")
	logo, _ := s.GetConfigValueByKey("site_logo")
	copyright, _ := s.GetConfigValueByKey("site_copyright")

	config = dto.SiteConfig{
		Name:      name,
		Logo:      logo,
		Copyright: copyright,
	}
	return
}

// SetSiteConfig 设置站点配置
func (s *AnnualConfigService) SetSiteConfig(config dto.SiteConfig) (err error) {
	configs := map[string]string{
		"site_name":      config.Name,
		"site_logo":      config.Logo,
		"site_copyright": config.Copyright,
	}
	return s.BatchSetConfig(configs)
}

// GetAllConfig 获取所有配置（按分组）
func (s *AnnualConfigService) GetAllConfig() (result dto.AllConfig, err error) {
	result.Wechat, _ = s.GetWechatConfig()
	result.Site, _ = s.GetSiteConfig()
	return
}

// 缓存配置到内存（可选，提高性能）
var configCache = make(map[string]string)

// GetConfigWithCache 带缓存获取配置
func (s *AnnualConfigService) GetConfigWithCache(key string) (value string, err error) {
	if v, ok := configCache[key]; ok {
		return v, nil
	}
	value, err = s.GetConfigValueByKey(key)
	if err == nil {
		configCache[key] = value
	}
	return
}

// ClearConfigCache 清除配置缓存
func (s *AnnualConfigService) ClearConfigCache() {
	configCache = make(map[string]string)
}

// RefreshConfigCache 刷新配置缓存
func (s *AnnualConfigService) RefreshConfigCache() error {
	list, err := s.GetConfigList()
	if err != nil {
		return err
	}
	configCache = make(map[string]string)
	for _, config := range list {
		configCache[config.ConfigKey] = config.ConfigValue
	}
	return nil
}

// GetConfigAsJSON 获取配置并解析为JSON对象
func (s *AnnualConfigService) GetConfigAsJSON(key string, v interface{}) error {
	value, err := s.GetConfigValueByKey(key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(value), v)
}

// SetConfigAsJSON 将对象序列化为JSON并保存
func (s *AnnualConfigService) SetConfigAsJSON(key string, v interface{}, description string) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return s.SetConfig(key, string(data), description)
}
