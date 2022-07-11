package main

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	api "github.com/hedlx/doless/client"
	"github.com/hedlx/doless/manager/lambda"
	"github.com/hedlx/doless/manager/logger"
	"github.com/hedlx/doless/manager/model"
	"github.com/hedlx/doless/manager/task"
	"github.com/hedlx/doless/manager/util"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	tSvc := task.CreateTaskService()
	lSvc, err := lambda.CreateLambdaService()

	if err != nil {
		panic(err)
	}

	controlSrv, err := StartControlServer(ctx, lSvc, tSvc)
	if err != nil {
		panic(err)
	}

	<-ctx.Done()
	stop()

	logger.L.Info("Shutting down gracefully, press Ctrl+C again to force")

	srvCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := controlSrv.Shutdown(srvCtx); err != nil {
		logger.L.Fatal("Server forced to shutdown", zap.Error(err))
	}

	svcCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	lSvc.Stop(svcCtx)
}

func StartControlServer(ctx context.Context, lSvc lambda.LambdaService, tSvc task.TaskService) (*http.Server, error) {
	r := gin.Default()

	r.POST("/upload", func(c *gin.Context) {
		fileHeader, err := c.FormFile("file")
		if err != nil {
			logger.L.Error("internal server error", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			logger.L.Error("internal server error", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		id, err := lambda.UploadTmp(c, file)
		if err != nil {
			logger.L.Error("internal server error", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"id": id})
	})

	r.GET("/lambda", func(c *gin.Context) {
		lambdas, err := lambda.GetLambdas(c)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, lambdas)
	})

	r.GET("/lambda/:id", func(c *gin.Context) {
		lambda, err := lambda.GetLambda(c, c.Param("id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, lambda)
	})

	r.POST("/lambda", func(c *gin.Context) {
		cLambda := &api.CreateLambda{}
		err := c.ShouldBind(cLambda)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = model.ValidateCreateLambda(cLambda)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		lambda, err := lSvc.BootstrapLambda(c, cLambda)
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
		runtimes, err := lambda.GetRuntimes(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, runtimes)
	})

	r.GET("/runtime/:id", func(c *gin.Context) {
		runtime, err := lambda.GetRuntime(c, c.Param("id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, runtime)
	})

	r.POST("/runtime", func(c *gin.Context) {
		cRuntime := &api.CreateRuntime{}
		err := c.ShouldBind(cRuntime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = model.ValidateCreateRuntime(cRuntime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		runtime, err := lSvc.BootstrapRuntime(c, cRuntime)
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
		Addr:    fmt.Sprintf(":%d", util.GetIntVar("PORT")),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.L.Error("Failed to start server", zap.Error(err))
		}
	}()

	return srv, nil
}
