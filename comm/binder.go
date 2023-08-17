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
		area := make([]int, size+4)
		area[0] = int(nomer)
		area[2] = int(nomer)
		area[3] = int(size)
		pos++
		for i := 0; i < int(size); i++ {
			area[i+4] = int(areaPriv[pos])
			pos++
		}
		switch nomer {
		case 14:
			nArrays.StatDefine.FromBuffer(area)
			logger.Info.Printf("statdefine %v", nArrays.StatDefine)
		case 15:
			nArrays.PointSet.FromBuffer(area)
			logger.Info.Printf("pointSet %v", nArrays.PointSet)
		case 16:
			nArrays.UseInput.FromBuffer(area)
			logger.Info.Printf("useInput %v", nArrays.UseInput)
		case 21:
			nArrays.TimeDivice.FromBuffer(area)
			logger.Info.Printf("timeDevice %v", nArrays.TimeDivice)
		case 7:
			area[0] = 40
			err := nArrays.SetupDK.FromBuffer(area)
			if err != nil {
				logger.Error.Print(err.Error())
			}
			logger.Info.Printf("setupDK %v", nArrays.SetupDK)
		case 20:
			area[0] = 157
			err := nArrays.SetTimeUse.FromBuffer(area)
			if err != nil {
				logger.Error.Print(err.Error())
			}
			logger.Info.Printf("setTimeUse %v", nArrays.SetTimeUse)
		case 23:
			area[0] = 148
			err := nArrays.SetTimeUse.FromBuffer(area)
			if err != nil {
				logger.Error.Print(err.Error())
			}
			logger.Info.Printf("setTimeUse %v", nArrays.SetTimeUse)
		case 24:
			area[0] = 149
			err := nArrays.SetCtrl.FromBuffer(area)
			if err != nil {
				logger.Error.Print(err.Error())
			}
			logger.Info.Printf("setCtrl %v", nArrays.SetCtrl)
		case 8:
			//Недельные карты
			area[0] = 45 + (area[4] - 1)
			err := nArrays.WeekSets.FromBuffer(area)
			if err != nil {
				logger.Error.Print(err.Error())
			}
			logger.Info.Printf("setWeek %v", nArrays.WeekSets.WeekSets[area[0]-45])
		case 9:
			//Суточные карты
			area[0] = 65 + (area[4] - 1)
			area[2] = 137
			err := nArrays.DaySets.FromBuffer(area)
			if err != nil {
				logger.Error.Print(err.Error())
			}
			logger.Info.Printf("daySet %v", nArrays.DaySets.DaySets[area[0]-65])
		case 22:
			//Годовая карты
			area[0] = 85 + (area[4] - 1)
			err := nArrays.MonthSets.FromBuffer(area)
			if err != nil {
				logger.Error.Print(err.Error())
			}
			logger.Info.Printf("monthSet %v", nArrays.MonthSets.MonthSets[area[0]-85])
		case 5:
			//Планы координации
			area[0] = 100 + (area[4] - 1)
			area[2] = 133
			err := nArrays.SetDK.FromBuffer(area)
			if err != nil {
				logger.Error.Print(err.Error())
			}
			logger.Info.Printf("setDK %v", nArrays.SetDK.DK[area[0]-100])

		}
		logger.Info.Printf("%d:%v", nomer, area)

	}
	err := nArrays.IsCorrect()
	if err != nil {
		logger.Info.Printf("Ошибка %s", err.Error())
		return false
	}
	if !data.DataValue.Connect {
		data.DataValue.SetArrays(nArrays)
	} else {
		data.Arrays <- nArrays
	}
	return true
}
