package utils

import (
	"fmt"
	"go01_4/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var mlog = logger.GetLogger()

func Succ(c *gin.Context, message string, data interface{}) {
	mlog.Infoln(message)
	c.JSON(http.StatusOK, gin.H{
		"message": message,
		"data":    data,
	})
}

func Error(c *gin.Context, status int, message string) {
	mlog.Errorln(message)
	c.AbortWithStatusJSON(status, gin.H{
		"message": message,
	})
	c.Abort()
}

func StrIdToUint(id string) uint {
	id_uint64, err := strconv.ParseUint(id, 10, 0)
	if err != nil {
		panic(fmt.Sprintf("id(%s)转换出错。 %v", id, err))
	}
	return uint(id_uint64)
}

func StrIdToInt(id string) int {
	id_uint64, err := strconv.ParseUint(id, 10, 0)
	if err != nil {
		panic(fmt.Sprintf("id(%s)转换出错。 %v", id, err))
	}
	return int(id_uint64)
}
