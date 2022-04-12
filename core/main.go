package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	lambda "github.com/hedlx/doless/core/lambda"
	model "github.com/hedlx/doless/core/model"
)

func main() {
	lSvc := lambda.CreateLambdaService()

	r := gin.Default()

	r.POST("/upload", func(c *gin.Context) {
		fileHeader, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		id, err := lambda.UploadTmp(c, file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"id": id})
	})

	r.POST("/lambda", func(c *gin.Context) {
		lambda := &model.CreateLambdaM{}
		err := c.ShouldBind(lambda)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = model.ValidateCreateLambdaM(lambda)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		lambda, err = lSvc.BootstrapLambda(c, lambda)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, lambda)
	})

	r.POST("/lambda/:id/start", func(c *gin.Context) {
		if err := lSvc.Start(c, c.Param("id")); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	})

	r.GET("/runtime", func(c *gin.Context) {
		runtimeC, errC := lambda.GetRuntimes(c)
		runtimes := []*model.RuntimeM{}

		for {
			select {
			case err := <-errC:
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			case runtime, ok := <-runtimeC:
				if !ok {
					c.JSON(http.StatusOK, runtimes)
					return
				}

				runtimes = append(runtimes, runtime)
			}
		}
	})

	r.POST("/runtime", func(c *gin.Context) {
		runtime := &model.CreateRuntimeM{}
		err := c.ShouldBind(runtime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = model.ValidateCreateRuntimeM(runtime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		runtime, err = lSvc.BootstrapRuntime(c, runtime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, runtime)
	})

	r.Run()
}
