package kdm

type Kdm struct {
	Area0x0000 Area0x0000
	Area0x001B Area0x001B
	Area0x00E0 Area0x00E0
	Area0x00F0 Area0x00F0
	Area0x0100 Area0x0100
	Area0x0200 Area0x0200
	Area0x0300 Area0x0300
	Area0x0400 Area0x0400
	Area0x0440 Area0x0440
	Area0x0500 Area0x0500
	Area0x0700 Area0x0700
	Area0x0800 Area0x0800
	Area0x080A Area0x080A
	Area0x0900 Area0x0900
	Area0x0A00 Area0x0A00
	Area0x0B00 Area0x0B00
	Area0x0C00 Area0x0C00
	Area0x0D00 Area0x0D00
	Area0x0F00 Area0x0F00
	Area0x1000 Area0x1000
}
type Area0x0000 struct {
	OutKeys     int          //Выходные ключи
	InStatus    int          //Состояние входов и управление реле
	Takt        int          //Номер исполняемого такта
	Time        int          //и время до его завершения
	Status      int          // Статус
	BadKeys     int          // Аварии ключей
	Sync        int          // Синхронизация
	SyncRTime   int          // Синхронизация по часам реального времени
	Programm    int          // Принудительный выбор программы
	BadKeyExt   int          // Аварии ключей платы расширения
	SwitchOff   int          // Выключить контроллер
	PhaseRU     int          // Фаза РУ
	GoToRU      int          // Перейти в РУ
	YellowBlink int          // Режим ЖМ
	BadKZ       int          // Аварии КЗ ключей 1-8
	TKyes       [8]int       // Мгновенный ток нагрузки на ключ
	buffer      [0x16]uint16 // Буфер состояния
}
