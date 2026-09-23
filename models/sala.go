package models

type Sala struct {
	ID         int      `json:"id"`
	Nome       string   `json:"nome"`
	Capacidade int      `json:"capacidade"`
	Recursos   Recursos `json:"recursos"`
	Ativa      bool     `json:"ativa"`
}

var Salas []Sala
