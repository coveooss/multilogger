package multilogger

import (
	"testing"
	"time"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestFormatConcurrency(t *testing.T) {
	t.Parallel()

	logger := New("my_module")
	formatter := logger.Formatter()

	numberOfThreads := 1000
	resultChannel := make(chan string, numberOfThreads)
	for i := 0; i < numberOfThreads; i++ {
		go func() {
			entry := &logrus.Entry{
				Message: color.BlueString("test"),
				Level:   logrus.InfoLevel,
				Time:    time.Date(2019, 12, 1, 10, 10, 11, 0, time.UTC),
				Logger:  logger.Logger,
				Data:    map[string]interface{}{moduleFieldName: "my_module"},
			}
			result, err := formatter.Format(entry)
			assert.Nil(t, err)
			resultChannel <- string(result)

		}()
	}

	for i := 0; i < numberOfThreads; i++ {
		assert.Equal(t, "[my_module] 2019/12/01 10:10:11.000 INFO     test\n", <-resultChannel)
	}

}
