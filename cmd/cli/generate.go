package cli

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli"
	"gitlab.iat.id/kiyora/hmac-loyalti/cmd/cli/models"
	"gitlab.iat.id/kiyora/hmac-loyalti/pkg"
)

func GenerateHmacCli(ctx *cli.Context) error {
	getPath := ctx.String("config")
	path := ctx.String("path")
	method := ctx.String("method")
	body := ctx.String("body")
	timestamp := ctx.String("timestamp")
	tmp, err := time.Parse("2006-01-02 15:04:05", timestamp)
	if err != nil {
		err = fmt.Errorf("invalid timestamp format, please use YYYY-MM-DD HH:MM:SS")
		pkg.ErrorHandling(err)
		return err
	}

	err = pkg.Load(getPath)
	if err != nil {
		pkg.ErrorHandling(err)
		return err
	}

	if !pkg.IsValidConfig() {
		err = fmt.Errorf("invalid config, please check your config file")
		pkg.ErrorHandling(err)
		return err
	}

	if body == "" {
		pkg.ErrorHandling(fmt.Errorf("body cannot be empty"))
	}

	bodySplit := strings.Split(body, ",")
	bodyPure := ""
	for _, v := range bodySplit {
		bodySplitDot := strings.Split(v, ":")
		if strings.Contains(bodySplitDot[0], "{") {
			bodySplitDot[0] = strings.Replace(bodySplitDot[0], "{", "", -1)
		}
		if strings.Contains(bodySplitDot[1], "}") {
			bodySplitDot[1] = strings.Replace(bodySplitDot[1], "}", "", -1)
		}
		checkIfInt, err := strconv.Atoi(bodySplitDot[1])
		var result string
		if err == nil && pkg.StringInSlice(bodySplitDot[0], []string{"amount", "point"}) {
			result = fmt.Sprintf(`"%s":%d`, bodySplitDot[0], checkIfInt)
		} else {
			result = fmt.Sprintf(`"%s":"%s"`, bodySplitDot[0], bodySplitDot[1])
		}

		bodyPure += result + ","
	}
	bodyPure = strings.TrimSuffix(bodyPure, ",")
	body = fmt.Sprintf("{%s}", bodyPure)

	hmac := pkg.CalculateHMAC(method, path, []byte(body), tmp.Format("20060102150405"))
	if hmac == "" {
		err = fmt.Errorf("hmac cant initialized")
		pkg.ErrorHandling(err)
		return err
	}

	var res *models.Result = &models.Result{
		Method:    method,
		Path:      path,
		Body:      body,
		Timestamp: tmp.Format("20060102150405"),
		Hmac:      hmac,
	}

	marsh, err := json.Marshal(res)
	if err != nil {
		err = fmt.Errorf("error marshalling result, %s", err)
		pkg.ErrorHandling(err)
		return err
	}
	if pkg.GetConfig().ResultCLi == "" || pkg.GetConfig().ResultCLi == "text" {
		fmt.Println(hmac)
	} else if pkg.GetConfig().ResultCLi == "json" {
		fmt.Println(string(marsh))
	}
	return nil
}
