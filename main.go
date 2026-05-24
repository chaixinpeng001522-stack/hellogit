package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := newRouter()

	addr := addrFromEnv()
	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		log.Printf("http server starting addr=%s", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		log.Printf("shutdown signal received")
	case err := <-errCh:
		if err != nil {
			log.Printf("http server error: %v", err)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
		_ = srv.Close()
	}
	log.Printf("http server stopped")
}

func newRouter() *gin.Engine {
	r := gin.Default()
	if err := r.SetTrustedProxies(nil); err != nil {
		log.Printf("set trusted proxies failed: %v", err)
	}
	if err := r.SetTrustedProxies(nil); err != nil {
		log.Printf("set trusted proxies failed: %v", err)
	}
	if err := r.SetTrustedProxies(nil); err != nil {
		log.Printf("set trusted proxies failed: %v", err)
	}

	r.Use(requestIDMiddleware())

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	
	r.GET("/hello", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Hello Gin</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        }
        .card {
            background: white;
            border-radius: 20px;
            padding: 48px 64px;
            box-shadow: 0 20px 60px rgba(0,0,0,0.3);
            text-align: center;
            max-width: 480px;
        }
        .emoji { font-size: 72px; margin-bottom: 24px; }
        h1 {
            font-size: 42px;
            color: #1a1a2e;
            margin-bottom: 12px;
        }
        .subtitle {
            color: #666;
            font-size: 18px;
            margin-bottom: 32px;
        }
        .info {
            background: #f8f9fa;
            border-radius: 12px;
            padding: 20px;
            text-align: left;
        }
        .info-item {
            display: flex;
            justify-content: space-between;
            padding: 8px 0;
            border-bottom: 1px solid #eee;
        }
        .info-item:last-child { border-bottom: none; }
        .info-label { color: #888; font-size: 14px; }
        .info-value { color: #333; font-family: monospace; font-size: 14px; }
        .status {
            display: inline-block;
            background: #10b981;
            color: white;
            padding: 4px 12px;
            border-radius: 20px;
            font-size: 12px;
            font-weight: 600;
        }
    </style>
</head>
<body>
    <div class="card">
        <div class="emoji">👋</div>
        <h1>Hello, Gin!</h1>
        <p class="subtitle">欢迎来到 Gin Web 服务</p>
        <div class="info">
            <div class="info-item">
                <span class="info-label">状态</span>
                <span class="status">运行中</span>
            </div>
            <div class="info-item">
                <span class="info-label">请求 ID</span>
                <span class="info-value">`+c.GetString(requestIDKey)+`</span>
            </div>
            <div class="info-item">
                <span class="info-label">时间</span>
                <span class="info-value">`+time.Now().Format("2006-01-02 15:04:05")+`</span>
            </div>
        </div>
    </div>
</body>
</html>`)
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
	})

	return r
}

const requestIDHeader = "X-Request-Id"
const requestIDKey = "request_id"

func requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader(requestIDHeader)
		if rid == "" {
			rid = newRequestID()
		}
		c.Set(requestIDKey, rid)
		c.Header(requestIDHeader, rid)
		c.Next()
	}
}

func newRequestID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 10)
	}
	return hex.EncodeToString(b)
}

func addrFromEnv() string {
	host := os.Getenv("HOST")
	portStr := os.Getenv("PORT")

	port := 8080
	if portStr != "" {
		p, err := strconv.Atoi(portStr)
		if err == nil && p > 0 && p <= 65535 {
			port = p
		} else {
			log.Printf("invalid PORT=%q, fallback to %d", portStr, port)
		}
	}

	if host == "" {
		return ":" + strconv.Itoa(port)
	}
	return host + ":" + strconv.Itoa(port)
}