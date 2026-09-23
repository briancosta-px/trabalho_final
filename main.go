package main

import (
	"api-gin/handlers"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.New()

	// Uso dos Middlewares globais nativos e personalizados
	r.Use(gin.Recovery())

	// 4. Mapeamento de Rotas sob Grupo Versionado
	v1 := r.Group("/api/v1")
	{
		// Monitoramento da API
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":    "healthy",
				"timestamp": time.Now(),
				"version":   "1.0.0",
			})
		})

		//Alunos
		v1.POST("/alunos", handlers.CriarAluno)
		v1.GET("/alunos", handlers.ListarAlunos)
		v1.GET("/alunos/:matricula", handlers.BuscarAluno)

		//Salas
		v1.POST("/salas", handlers.CriarSala)
		v1.GET("/salas", handlers.ListarSalas)
		v1.GET("/salas/:id", handlers.BuscarSala)
		v1.PUT("/salas/:id", handlers.AtualizarSala)

		// Domínio de Turmas (Classes)
		v1.POST("/turmas", handlers.CriarTurma)
		v1.GET("/turmas", handlers.ListarTurmas)
		v1.GET("/turmas/:id", handlers.BuscarTurma)
		v1.POST("/turmas/:id/alunos", handlers.AdicionarAluno)
		v1.GET("/turmas/:id/alunos", handlers.ListarAlunosDaTurma)
		v1.POST("/turmas/:id/alocar", handlers.AlocarSala)
	}

	r.Run(":8080")
}
