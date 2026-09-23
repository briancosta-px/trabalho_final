package models

type Aluno struct {
	Matricula int    `json:"matricula"`
	Nome      string `json:"nome"`
	Email     string `json:"email"`
}

var Alunos []Aluno
