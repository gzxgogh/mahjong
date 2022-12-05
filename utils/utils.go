package utils

import (
	"encoding/json"
	"mahjong/model"
	"math/rand"
	"strings"
	"time"
)

func Error(code int, msg string) model.Result {
	return model.Result{
		Status: code,
		Msg:    msg,
		Data:   nil,
	}
}

func Success(data interface{}) model.Result {
	return model.Result{
		Status: 200,
		Msg:    "success",
		Data:   data,
	}
}

func GetRandomWithAll(min, max int) int64 {
	rand.Seed(time.Now().UnixNano())
	return int64(rand.Intn(max-min+1) + min)
}

func FromJSON(j string, o interface{}) *interface{} {
	err := json.Unmarshal([]byte(j), &o)
	if err != nil {
		return nil
	} else {
		return &o
	}
}

func ToJSON(o interface{}) string {
	j, err := json.Marshal(o)
	if err != nil {
		return "{}"
	} else {
		js := string(j)
		js = strings.Replace(js, "\\u003c", "<", -1)
		js = strings.Replace(js, "\\u003e", ">", -1)
		js = strings.Replace(js, "\\u0026", "&", -1)
		return js
	}
}
