package router

import (
	"bookingBackEnd/dao/mysql"
	"bookingBackEnd/utils"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setup() {
	config_filepath := "../conf/conf.ini"
	utils.InitConfig(config_filepath)
	// 初始化 logger
	utils.NewLoggerHelper()

	// 加载mysql数据库
	err := mysql.Init()
	if err != nil {
		panic(err)
	}
}

func teardown() {

}

func prettyJson(origin_json []byte) (pretty_json []byte, err error) {
	data := map[string]interface{}{}
	err = json.Unmarshal(origin_json, &data)
	if err != nil {
		return
	}
	pretty_json, err = json.MarshalIndent(data, "", "    ")
	return
}

// ParseToStr 将map中的键值对输出成querystring形式
func ParseToStr(mp map[string]string) string {
	values := ""
	for key, val := range mp {
		values += "&" + key + "=" + val
	}
	temp := values[1:]
	values = "?" + temp
	return values
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func TestGetUserAppInfo(t *testing.T) {
	r := NewRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/wechat/login", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	response_json, err := prettyJson(w.Body.Bytes())
	if err != nil {
		t.Fatalf("error is: %v", err)
	}
	t.Log(string(response_json))
}

func TestRegister(t *testing.T) {
	r := NewRouter()

	body := map[string]interface{}{
		"email":    "admin@qq.com",
		"password": "root",
	}
	jsonBytes, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/user/register", bytes.NewReader(jsonBytes))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	response_json, err := prettyJson(w.Body.Bytes())
	if err != nil {
		t.Fatalf("error is: %v", err)
	}
	t.Log(string(response_json))
}

func TestLogin(t *testing.T) {
	r := NewRouter()

	body := map[string]interface{}{
		"email":    "123456789@qq.com",
		"password": "1234567",
	}
	jsonBytes, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/user/login", bytes.NewReader(jsonBytes))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	response_json, err := prettyJson(w.Body.Bytes())
	if err != nil {
		t.Fatalf("error is: %v", err)
	}
	t.Log(string(response_json))
}

func TestGetDetailedClassroomList(t *testing.T) {
	r := NewRouter()

	params := map[string]string{
		"thirdSessionId": "e7a0b2c29fe2421f8b8c113b34fbdff6",
		"pageNum":        "0",
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/classroom/detailedList"+ParseToStr(params), nil)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	response_json, err := prettyJson(w.Body.Bytes())
	if err != nil {
		t.Fatalf("error is: %v", err)
	}
	t.Log(string(response_json))
}

func TestFilterClassroomAndBookingPeriod(t *testing.T) {
	r := NewRouter()

	body := map[string]interface{}{
		"floor":    1,
		"capacity": 1,
		"power":    1,
		"date":     "2023-08-23",
	}
	jsonBytes, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/booking/info?thirdSessionId=e7a0b2c29fe2421f8b8c113b34fbdff6", bytes.NewReader(jsonBytes))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	response_json, err := prettyJson(w.Body.Bytes())
	if err != nil {
		t.Fatalf("error is: %v", err)
	}
	t.Log(string(response_json))
}
