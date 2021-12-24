package models

/**
 * @Author: Tao Jun
 * @Description: models
 * @File:  alertmanager_rules
 * @Version: 1.0.0
 * @Date: 2021/12/6 下午5:30
 */
type RuleGroup struct {
	Groups []SubGroup `json:"groups"`
}
type SubGroup struct {
	Name  string   `json:"name"`
	Rules []Metric `json:"rules"`
}

type Metric struct {
	Expr        string                 `json:"expr"`
	Alert       string                 `json:"alert"`
	For         string                 `json:"for"`
	Labels      map[string]interface{} `json:"labels"`
	Annotations Annotations            `json:"annotations"`
}
type Annotations struct {
	Summary     string `json:"summary"`
	Description string `json:"description"`
}
