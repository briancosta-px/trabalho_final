package handlers

import (
	"api-gin/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CriarSala(c *gin.Context) {
	var sala models.Sala

	err := c.ShouldBindJSON(&sala)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "Dados inválidos.",
		})
		return
	}

	if sala.Capacidade <= 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"erro": "A capacidade da sala deve ser maior que zero",
		})
		return
	}

	for _, salaExistente := range models.Salas {
		if salaExistente.ID == sala.ID {
			c.JSON(http.StatusConflict, gin.H{
				"erro": "ID de sala já em uso.",
			})
			return
		}
	}

	models.Salas = append(models.Salas, sala)
	c.JSON(http.StatusCreated, sala)
}

func ListarSalas(c *gin.Context) {
	c.JSON(http.StatusOK, models.Salas)
}

func BuscarSala(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "ID inválido",
		})
		return
	}

	for _, sala := range models.Salas {
		if sala.ID == id {
			c.JSON(http.StatusOK, sala)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"erro": "Sala não encontrada.",
	})
}

func AtualizarSala(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "ID inválido.",
		})
		return
	}

	var salaAtualizada models.Sala

	err = c.ShouldBindJSON(&salaAtualizada)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "Dados inválidos.",
		})
		return
	}

	if salaAtualizada.Capacidade <= 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"erro": "A capacidade da sala deve ser maior que zero",
		})
		return
	}

	for i, sala := range models.Salas {
		if sala.ID == id {
			salaAtualizada.ID = id
			models.Salas[i] = salaAtualizada

			c.JSON(http.StatusOK, salaAtualizada)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"erro": "Sala não encontrada.",
	})
}
