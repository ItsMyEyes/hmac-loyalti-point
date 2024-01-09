package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/urfave/cli"
	"gitlab.iat.id/kiyora/hmac-loyalti/cmd/cli/models"
	"gitlab.iat.id/kiyora/hmac-loyalti/pkg"
)

func GenerateHmacHttp(c *cli.Context) error {
	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	config := c.String("config")
	err := pkg.Load(config)
	if err != nil {
		pkg.ErrorHandling(err)
		return err
	}

	if !pkg.IsValidConfig() {
		err = fmt.Errorf("invalid config, please check your config file")
		pkg.ErrorHandling(err)
		return err
	}

	fmt.Println("http server started")
	_, cancel := context.WithCancel(context.Background())
	go httpServer()

	sig := <-sigC
	fmt.Println("signal received: ", sig)

	defer cancel()
	return nil
}

func httpServer() {
	http.HandleFunc("/hmac", hmac)

	fmt.Println("server started at ", fmt.Sprintf("%s:%s", pkg.GetConfig().Host, pkg.GetConfig().Port))
	http.ListenAndServe(fmt.Sprintf("%s:%s", pkg.GetConfig().Host, pkg.GetConfig().Port), nil)
}

func hmac(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	b, _ := io.ReadAll(req.Body)
	fmt.Println(string(b))

	w.Header().Set("Content-Type", "application/json")
	method := req.Header.Get("X-Method")
	path := req.Header.Get("X-Path")
	timestamp := req.Header.Get("X-Timestamp")

	fmt.Println("method: ", method, "path: ", path, "timestamp: ", timestamp)
	tmp, err := time.Parse("2006-01-02 15:04:05", timestamp)
	var res *models.Result = &models.Result{
		Status:    false,
		Method:    method,
		Path:      path,
		Body:      string(b),
		Message:   "success",
		Timestamp: tmp.Format("20060102150405"),
	}

	if err != nil {
		err = fmt.Errorf("invalid timestamp format, please use YYYY-MM-DD HH:MM:SS")
		w.WriteHeader(http.StatusBadRequest)
		res.Message = err.Error()
		marshal, _ := json.Marshal(res)
		w.Write(marshal)
		return
	}

	hmac := pkg.CalculateHMAC(method, path, b, tmp.Format("20060102150405"))
	if hmac == "" {
		err = fmt.Errorf("hmac cant initialized")
		w.WriteHeader(http.StatusBadRequest)
		res.Message = err.Error()
		marshal, _ := json.Marshal(res)
		w.Write(marshal)
		return
	}

	res.Hmac = hmac
	w.WriteHeader(http.StatusOK)
	marshal, _ := json.Marshal(res)
	w.Write(marshal)
}
