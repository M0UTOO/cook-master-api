package managers

type Manager struct {
	IdManager int `json:"idmanager"`
	IsItemManager bool `json:"isitemmanager"`
	IsClientManager bool `json:"isclientmanager"`
	IsContractorManager bool `json:"iscontractormanager"`
	IsSuperAdmin bool `json:"issuperadmin"`
}