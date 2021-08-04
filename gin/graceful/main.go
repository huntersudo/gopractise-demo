package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		c.String(http.StatusOK, "Welcome Gin Server")
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

// Go 1.8 版本或者更新的版本，http.Server 内置的 Shutdown 方法，已经实现了优雅关闭。
// 上面的示例中，需要把 srv.ListenAndServe 放在 goroutine 中执行，这样才不会阻塞到 srv.Shutdown 函数。
// 因为我们把 srv.ListenAndServe 放在了 goroutine 中，所以需要一种可以让整个进程常驻的机制。
//
// 这里，我们借助了无缓冲 channel，并且调用 signal.Notify 函数将该 channel 绑定到 SIGINT、SIGTERM 信号上。
// 这样，收到 SIGINT、SIGTERM 信号后，quilt 通道会被写入值，从而结束阻塞状态，程序继续运行，执行 srv.Shutdown(ctx)，优雅关停 HTTP 服务。