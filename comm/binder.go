package comm

import (
	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/irz/data"
)

func moveArrasIsGood() bool {
	nArrays := data.DataValue.Arrays
	pos := 0
	for pos < len(areaPriv) {
		nomer := areaPriv[pos]
		pos++
		size := areaPriv[pos]
		area := make([]int, size)
		pos++
		for i := 0; i < int(size); i++ {
			area[i] = int(areaPriv[pos])
			pos++
		}
		logger.Info.Printf("%d:%v", nomer, area)
	}
	err := nArrays.IsCorrect()
	if err != nil {
		logger.Info.Printf("Ошибка %s", err.Error())
		return false
	}
	data.Arrays <- nArrays
	return true
}
