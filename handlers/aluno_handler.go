package handlers

import (
	"api-gin/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CriarAluno(c *gin.Context) {
	var aluno models.Aluno

	err := c.ShouldBindJSON(&aluno)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "Dados inválidos.",
		})
		return
	}

	for _, alunoExistente := range models.Alunos {
		if alunoExistente.Matricula == aluno.Matricula {
			c.JSON(http.StatusConflict, gin.H{
				"erro": "Matrícula já cadastrada.",
			})
			return
		}
	}
	models.Alunos = append(models.Alunos, aluno)
	c.JSON(http.StatusCreated, aluno)
}

func ListarAlunos(c *gin.Context) {
	c.JSON(http.StatusOK, models.Alunos)
}

func BuscarAluno(c *gin.Context) {
	Matricula, err := strconv.Atoi(c.Param("matricula"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "Matrícula inválida.",
		})
		return
	}

	for _, aluno := range models.Alunos {
		if aluno.Matricula == Matricula {
			c.JSON(http.StatusOK, aluno)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{
		"erro": "Aluno não encontrado.",
	})
}
