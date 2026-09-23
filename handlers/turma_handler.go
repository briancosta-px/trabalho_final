package handlers

import (
	"api-gin/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func CriarTurma(c *gin.Context) {
	var turma models.Turma

	err := c.ShouldBindJSON(&turma)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "Dados inválidos.",
		})
		return
	}

	for _, turmaExistente := range models.Turmas {
		if turmaExistente.ID == turma.ID {
			c.JSON(http.StatusConflict, gin.H{
				"erro": "ID da turma já cadastrado.",
			})
			return
		}
	}

	models.Turmas = append(models.Turmas, turma)
	c.JSON(http.StatusCreated, turma)
}

func ListarTurmas(c *gin.Context) {

	var resposta []gin.H

	for _, turma := range models.Turmas {

		statusAlocacao := "Não alocada"

		if turma.Sala != nil {
			statusAlocacao = "Alocada"
		}

		resposta = append(resposta, gin.H{
			"id":                turma.ID,
			"nome":              turma.Nome,
			"disciplina":        turma.Disciplina,
			"professor":         turma.Professor,
			"quantidade_alunos": turma.QuantidadeAlunos,
			"status_alocacao":   statusAlocacao,
		})
	}

	c.JSON(http.StatusOK, resposta)
}

func BuscarTurma(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "ID inválido.",
		})
		return
	}

	for _, turma := range models.Turmas {
		if turma.ID == id {
			c.JSON(http.StatusOK, turma)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"erro": "Turma não encontrada.",
	})
}

func AdicionarAluno(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "ID da turma inválido.",
		})
		return
	}

	var dados struct {
		AlunoID int `json:"aluno_id"`
	}

	err = c.ShouldBindJSON(&dados)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "Dados inválidos.",
		})
		return
	}

	for i := range models.Turmas {
		if models.Turmas[i].ID == id {
			for _, aluno := range models.Alunos {
				if aluno.Matricula == dados.AlunoID {
					for _, alunoDaTurma := range models.Turmas[i].Alunos {
						if alunoDaTurma.Matricula == aluno.Matricula {
							c.JSON(http.StatusConflict, gin.H{
								"erro": "Aluno já está matriculado nessa turma.",
							})
							return
						}
					}

					for _, outraTurma := range models.Turmas {

						if outraTurma.ID == models.Turmas[i].ID {
							continue
						}

						if outraTurma.DiaDaSemana != models.Turmas[i].DiaDaSemana {
							continue
						}

						for _, alunoDaOutraTurma := range outraTurma.Alunos {
							if alunoDaOutraTurma.Matricula == aluno.Matricula {

								inicioNovo, err := time.Parse("15:04", models.Turmas[i].HorarioInicio)

								if err != nil {
									c.JSON(http.StatusBadRequest, gin.H{
										"erro": "Horário de início inválido.",
									})
									return
								}

								fimNovo, err := time.Parse("15:04", models.Turmas[i].HorarioTermino)

								if err != nil {
									c.JSON(http.StatusBadRequest, gin.H{
										"erro": "Horário de término inválido.",
									})
									return
								}

								inicioExistente, err := time.Parse("15:04", outraTurma.HorarioInicio)

								if err != nil {
									continue
								}

								fimExistente, err := time.Parse("15:04", outraTurma.HorarioTermino)

								if err != nil {
									continue
								}

								if inicioNovo.Before(fimExistente) && fimNovo.After(inicioExistente) {
									c.JSON(http.StatusConflict, gin.H{
										"erro": "Aluno possui conflito de horário com outra turma.",
									})
									return
								}
							}
						}
					}

					if models.Turmas[i].Sala != nil {
						if models.Turmas[i].QuantidadeAlunos >= models.Turmas[i].Sala.Capacidade {
							c.JSON(http.StatusUnprocessableEntity, gin.H{
								"erro": "Capacidade da sala atingida.",
							})
							return
						}
					}

					models.Turmas[i].Alunos = append(models.Turmas[i].Alunos, aluno)
					models.Turmas[i].QuantidadeAlunos++
					c.JSON(http.StatusOK, models.Turmas[i])
					return
				}
			}
			c.JSON(http.StatusNotFound, gin.H{
				"erro": "Aluno não encontrado.",
			})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{
		"erro": "Turma não encontrada.",
	})
}

func AlocarSala(c *gin.Context) {
	idTurma, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "ID da turma inválido.",
		})
		return
	}

	var dados struct {
		SalaID         int    `json:"sala_id"`
		DiaDaSemana    string `json:"dia_da_semana"`
		HorarioInicio  string `json:"horario_inicio"`
		HorarioTermino string `json:"horario_termino"`
	}

	err = c.ShouldBindJSON(&dados)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "Dados inválidos.",
		})
		return
	}

	for i := range models.Turmas {
		if models.Turmas[i].ID == idTurma {

			if !models.Turmas[i].Ativa {
				c.JSON(http.StatusConflict, gin.H{
					"erro": "Turma está inativa.",
				})
				return
			}

			for j := range models.Salas {
				if models.Salas[j].ID == dados.SalaID {

					sala := &models.Salas[j]

					if !sala.Ativa {
						c.JSON(http.StatusConflict, gin.H{
							"erro": "Sala está inativa.",
						})
						return
					}

					if models.Turmas[i].QuantidadeAlunos > sala.Capacidade {
						c.JSON(http.StatusUnprocessableEntity, gin.H{
							"erro": "Quantidade de alunos excede a capacidade da sala.",
						})
						return
					}

					inicioNovo, err := time.Parse("15:04", dados.HorarioInicio)

					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"erro": "Horário de início inválido.",
						})
						return
					}

					fimNovo, err := time.Parse("15:04", dados.HorarioTermino)

					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"erro": "Horário de término inválido.",
						})
						return
					}

					if !fimNovo.After(inicioNovo) {
						c.JSON(http.StatusBadRequest, gin.H{
							"erro": "Horário de término deve ser depois do horário de início.",
						})
						return
					}

					for _, outraTurma := range models.Turmas {

						if outraTurma.ID == models.Turmas[i].ID {
							continue
						}

						if outraTurma.Sala == nil {
							continue
						}

						if outraTurma.Sala.ID != sala.ID {
							continue
						}

						if outraTurma.DiaDaSemana != dados.DiaDaSemana {
							continue
						}

						inicioExistente, err := time.Parse("15:04", outraTurma.HorarioInicio)

						if err != nil {
							continue
						}

						fimExistente, err := time.Parse("15:04", outraTurma.HorarioTermino)

						if err != nil {
							continue
						}

						if inicioNovo.Before(fimExistente) && fimNovo.After(inicioExistente) {
							c.JSON(http.StatusConflict, gin.H{
								"erro": "Conflito de horário da sala.",
							})
							return
						}
					}

					models.Turmas[i].DiaDaSemana = dados.DiaDaSemana
					models.Turmas[i].HorarioInicio = dados.HorarioInicio
					models.Turmas[i].HorarioTermino = dados.HorarioTermino
					models.Turmas[i].Sala = sala

					c.JSON(http.StatusOK, models.Turmas[i])
					return
				}
			}

			c.JSON(http.StatusNotFound, gin.H{
				"erro": "Sala não encontrada.",
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"erro": "Turma não encontrada.",
	})
}

func ListarAlunosDaTurma(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "ID da turma inválido.",
		})
		return
	}

	for _, turma := range models.Turmas {
		if turma.ID == id {
			c.JSON(http.StatusOK, turma.Alunos)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"erro": "Turma não encontrada.",
	})
}
