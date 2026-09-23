package models

type Recursos struct {
	Projetor       bool `json:"projetor"`
	Tv             bool `json:"tv"`
	ArCondicionado bool `json:"ar_condicionado"`
	CaixaDeSom     bool `json:"caixa_de_som"`
}
