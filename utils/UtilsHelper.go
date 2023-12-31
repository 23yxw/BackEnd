package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// StructToMapWithJson 使用json的marshal和unmarshal将struct转化为map
func StructToMapWithJson(obj interface{}) (ret map[string]interface{}) {
	objByte, _ := json.Marshal(obj)
	_ = json.Unmarshal(objByte, &ret)
	return
}

func StructAddToMap(obj interface{}, dst map[string]interface{}) (err error) {
	// 获取结构体类型
	v := reflect.ValueOf(obj)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		firstChar := strings.ToLower(string(fieldType.Name[0]))
		fieldName := firstChar + fieldType.Name[1:]

		zero := reflect.Zero(fieldType.Type)
		if !reflect.DeepEqual(field.Interface(), zero.Interface()) {
			dst[fieldName] = field.Interface()
		}
	}
	return
}

// StructArrayToMapWithJson 使用json的marshal和unmarshal将struct array转化为map
func StructArrayToMapWithJson(obj []interface{}) (ret []map[string]interface{}) {
	for k, objV := range obj {
		ret[k] = StructToMapWithJson(objV)
	}
	return
}

func GetMd5() (ret string) {
	data := []byte(time.Now().String() + "_" + strconv.Itoa(GenerateRangeNum(1, 100)))
	h := md5.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func GenerateRangeNum(min, max int) int {
	rand.Seed(time.Now().Unix())
	randNum := rand.Intn(max-min) + min
	return randNum
}
