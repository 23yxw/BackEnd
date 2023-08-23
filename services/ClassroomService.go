package services

import (
	"bookingBackEnd/dao/mysql"
	"bookingBackEnd/model"
	"bookingBackEnd/utils"
	"os"
)

func GetDetailedClassroomList(pageNum int) (results []map[string]interface{}, err error) {
	var detailedClassroomList []model.DetailedClassroomInfo
	err = mysql.ClassroomTableInstance.GetDetailedClassroomList(&detailedClassroomList, pageNum)
	if err != nil {
		utils.ErrorLogger.Errorf("error:%v", err)
		return
	}

	results = make([]map[string]interface{}, len(detailedClassroomList))
	for i, detailedClassroom := range detailedClassroomList {
		imagePath := detailedClassroom.Photo
		// imagePath := "/home/kihensarn/Booking/bookingBackEnd/images/3737366666393835383430643332346265663264656234383032396131366135.jpg"
		// 检查文件是否存在
		_, err = os.Stat(imagePath)
		if os.IsNotExist(err) {
			utils.ErrorLogger.Errorf("error:%v", err)
			return
		}

		// 读取图片文件
		var file *os.File
		file, err = os.Open(imagePath)
		if err != nil {
			utils.ErrorLogger.Errorf("error:%v", err)
			return
		}

		// 获取图片文件信息
		fileInfo, _ := file.Stat()
		fileSize := fileInfo.Size()

		// 创建一个缓冲区来存储图片数据
		buffer := make([]byte, fileSize)

		// 读取图片数据到缓冲区
		_, err = file.Read(buffer)
		if err != nil {
			utils.ErrorLogger.Errorf("error:%v", err)
			return
		}

		detailedClassroomMap := utils.StructToMapWithJson(detailedClassroom.ClassroomInfo)
		results[i] = make(map[string]interface{})
		results[i]["classroomInfo"] = detailedClassroomMap
		results[i]["imageData"] = buffer

		file.Close()
	}
	return
}
