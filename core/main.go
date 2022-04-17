package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/hedlx/doless/core/lambda"
	"github.com/hedlx/doless/core/logger"
	"github.com/hedlx/doless/core/model"
	"github.com/hedlx/doless/core/task"
	"github.com/hedlx/doless/core/util"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	tSvc := task.CreateTaskService()
	lSvc, err := lambda.CreateLambdaService()

	if err != nil {
		panic(err)
	}

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

	r.GET("/lambda", func(c *gin.Context) {
		lambdasC, errC := lambda.GetLambdas(c)
		lambdas := []*model.LambdaM{}

		for {
			select {
			case err := <-errC:
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			case lambda, ok := <-lambdasC:
				if !ok {
					c.JSON(http.StatusOK, lambdas)
					return
				}

				lambdas = append(lambdas, lambda)
			}
		}
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
		id := util.UUID()
		tSvc.Add(id)

		go func() {
			ctx := context.TODO()

			if err := lSvc.Start(ctx, c.Param("id")); err != nil {
				tSvc.Failed(id, struct {
					Error string `json:"error"`
				}{Error: err.Error()})
				return
			}

			tSvc.Succeeded(id, nil)
		}()

		c.JSON(http.StatusAccepted, gin.H{"task": id})
	})

	r.POST("/lambda/:id/destroy", func(c *gin.Context) {
		id := util.UUID()
		tSvc.Add(id)

		go func() {
			ctx := context.TODO()

			if err := lSvc.Destroy(ctx, c.Param("id")); err != nil {
				tSvc.Failed(id, struct {
					Error string `json:"error"`
				}{Error: err.Error()})
				return
			}

			tSvc.Succeeded(id, nil)
		}()

		c.JSON(http.StatusAccepted, gin.H{"task": id})
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

	r.GET("/task/:id", func(c *gin.Context) {
		status := tSvc.Get(c.Param("id"))

		if status == nil {
			c.Status(http.StatusNotFound)
			return
		}

		c.JSON(http.StatusOK, task.PrepareStatus(status))
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.L.Error("Failed to start server", zap.Error(err))
		}
	}()

	<-ctx.Done()
	stop()

	logger.L.Info("Shutting down gracefully, press Ctrl+C again to force")

	srvCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(srvCtx); err != nil {
		logger.L.Fatal("Server forced to shutdown", zap.Error(err))
	}

	svcCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	lSvc.Stop(svcCtx)
}
