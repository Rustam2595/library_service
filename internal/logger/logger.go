package logger

import (
	"github.com/rs/zerolog"
	"os"
	"strconv"
	"sync"
)

var once sync.Once //такая переменная, что отработает 1 раз

var log zerolog.Logger

func Get(flags ...bool) zerolog.Logger {
	once.Do(func() {
		zerolog.TimestampFieldName = "Time"
		zerolog.LevelFieldName = "Level"
		zerolog.CallerMarshalFunc = func(_ uintptr, file string, line int) string {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
			return file + ":" + strconv.Itoa(line)
		} //добавили номер строки и название файла
		if flags[0] {
			log = zerolog.New(os.Stdout).
				Level(zerolog.DebugLevel).
				With().
				Timestamp().
				Caller().
				Logger().
				Output(zerolog.ConsoleWriter{Out: os.Stderr})
		} else {
			log = zerolog.New(os.Stdout).
				Level(zerolog.InfoLevel).
				With().
				Timestamp().
				Caller().
				Logger()
		} //если flag поля Config.Debug=true,все логи уровня дебаг и выше выводим на экран
		//если flag поля Config.Debug=false,все логи уровня info и выше выводим на экран(без debug)
	})
	return log
}
