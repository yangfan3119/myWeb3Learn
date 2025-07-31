package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// Log 是对外暴露的日志实例
var mlog *logrus.Logger

// 自定义写入器，用于处理日志文件的每日轮转
type dailyRotatingWriter struct {
	currentDate string     // 当前日期信息
	logFile     *os.File   //file句柄完成信息写入
	mutex       sync.Mutex //互斥锁解决写入安全
	logDir      string     // 日志保存目录，一般为log，此处也可以改
	filePrefix  string     //日志文件前缀，前缀_日期.log
}

// NewDailyRotatingWriter 创建一个新的每日轮转写入器
func newDailyRotatingWriter(logDir, filePrefix string) (*dailyRotatingWriter, error) {
	// 创建日志目录（如果不存在）
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("无法创建日志目录: %v", err)
	}

	writer := &dailyRotatingWriter{
		logDir:     logDir,
		filePrefix: filePrefix,
	}

	// 初始化日志文件
	if err := writer.rotateIfNeeded(); err != nil {
		return nil, err
	}

	return writer, nil
}

// 检查是否需要轮转日志文件
func (w *dailyRotatingWriter) rotateIfNeeded() error {
	today := time.Now().Format("2006-01-02")

	w.mutex.Lock()
	defer w.mutex.Unlock()

	// 如果日期未变且文件已打开，则不需要轮转
	if today == w.currentDate && w.logFile != nil {
		return nil
	}

	// 关闭当前日志文件（如果存在）
	if w.logFile != nil {
		if err := w.logFile.Close(); err != nil {
			return fmt.Errorf("关闭日志文件失败: %v", err)
		}
		w.logFile = nil
	}

	// 生成新的日志文件名
	logFileName := filepath.Join(w.logDir, fmt.Sprintf("%s_%s.log", w.filePrefix, today))

	// 打开新的日志文件（创建+写入+追加模式）
	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("打开日志文件失败: %v", err)
	}

	w.currentDate = today
	w.logFile = file
	return nil
}

// Write 实现io.Writer接口
func (w *dailyRotatingWriter) Write(p []byte) (n int, err error) {
	// 每次写入前检查是否需要轮转
	if err := w.rotateIfNeeded(); err != nil {
		return 0, err
	}

	w.mutex.Lock()
	defer w.mutex.Unlock()

	// 写入日志内容
	return w.logFile.Write(p)
}

// 初始化日志配置
func Init() {
	if mlog != nil {
		return // 如果日志已经初始化，则直接返回
	}

	fileWriter, err := newDailyRotatingWriter("log", "blog")
	if err != nil {
		panic(fmt.Sprintf("初始化日志写入器失败: %v", err))
	}

	// 创建多写入器：同时写入文件和标准输出
	multiWriter := io.MultiWriter(fileWriter, os.Stdout)

	// 初始化logrus实例
	mlog = logrus.New()
	mlog.SetOutput(multiWriter)
	mlog.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000",
	})
	mlog.SetLevel(logrus.InfoLevel) // 设置日志级别

	// 打印不同级别的日志
	// mlog.Trace("这是一条trace级别的日志（最详细）")
	// mlog.Debug("这是一条debug级别的日志（调试信息）")
	// mlog.Info("这是一条info级别的日志（普通信息）")
	// mlog.Warn("这是一条warn级别的日志（警告信息）")
	// mlog.Error("这是一条error级别的日志（错误信息）")
}

func GetLogger() *logrus.Logger {
	if mlog == nil {
		Init()
		if mlog == nil {
			panic("日志实例未初始化，请先调用 Init() 方法")
		}
	}
	return mlog
}
