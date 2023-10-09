package core

type Location struct {
	Id        string `json:"id"`
	Address   string `json:"address"`
	City      string `json:"city"`
	St        string `json:"st"`
	Zip       string `json:"zip"`
	CompanyId int    `json:"companyId"`
}
