package util

import (
	"bytes"
	"carrot-backyard/config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type qqSendResponse struct {
	Status string `json:"status"`
}

func QQSendAndFindWhetherSuccess(userId int64, message string) bool {
	sendMsgUrl := fmt.Sprintf("%s%s", config.C.QQBot.Host, "/send_private_msg")
	req, err := http.NewRequest("POST", sendMsgUrl, bytes.NewBuffer(nil))
	if err != nil {
		ErrorPrint(err, userId, "message send")
		return false
	}

	q := req.URL.Query()
	q.Add("user_id", strconv.FormatInt(userId, 10))
	q.Add("message", packageMessage(message))
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		ErrorPrint(err, userId, "message send")
		return false
	}
	body, _ := ioutil.ReadAll(resp.Body)
	qqSendResp := qqSendResponse{}
	if err = json.Unmarshal(body, &qqSendResp); err != nil || qqSendResp.Status == "failed" {
		return false
	}
	return true
}
