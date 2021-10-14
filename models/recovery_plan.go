package models

type RecoveryPlanInput struct {
	Name  string                 `json:"name"`
	Steps map[string]interface{} `json:"steps"`
}
