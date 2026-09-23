package models

type Turma struct {
	ID               int     `json:"id"`
	Nome             string  `json:"nome"`
	Disciplina       string  `json:"disciplina"`
	Professor        string  `json:"professor"`
	QuantidadeAlunos int     `json:"quantidade_alunos"`
	Alunos           []Aluno `json:"alunos"`
	Sala             *Sala   `json:"sala"`
	DiaDaSemana      string  `json:"dia_da_semana"`
	HorarioInicio    string  `json:"horario_inicio"`
	HorarioTermino   string  `json:"horario_termino"`
	Ativa            bool    `json:"ativa"`
}

var Turmas []Turma
