package bilge

import (
	"log"
	"math/rand"
	"time"
)

type BigleS struct {
	level              chan int
	quit               chan struct{}
	shipSank           chan struct{}
	pumpControl        chan int
	waterlevel         int
	timer              time.Duration
	pumpConsumedEnergy int
}

var bs *BigleS
var _ Ship = (*BigleS)(nil)

func GetShipSubsystems(TimerSeconds int) (Ship, Pump, WaterLevelSensor) {
	bs = &BigleS{
		level:       make(chan int, 1),
		quit:        make(chan struct{}, 1),
		shipSank:    make(chan struct{}, 1),
		pumpControl: make(chan int, 1),
		timer:       time.Second * time.Duration(TimerSeconds),
	}

	go bs.run()
	go bs.pumpWorker()

	return bs, bs, bs
}

func (bs *BigleS) GetWaterLevelInformer() int {
	time.Sleep(10 * time.Nanosecond)
	return <-bs.level
}

func (bs *BigleS) run() {

	var sensorError int
	predWatelewel := bs.waterlevel
	startTime := time.Now()
	var sensor_level int
	ticker := time.NewTicker(time.Microsecond * 10) // Тикер для постоянных обновлений
	defer ticker.Stop()

	for time.Since(startTime) < bs.timer {
		select {
		case <-ticker.C:
			//протечка воды в корпус корабля
			bs.waterlevel++

			//занесение значение датчика уровния с определенной дескритизацией и ошибкой датчика
			sensor_level = bs.waterlevel/1000 + sensorError
			if bs.waterlevel == 2500 && predWatelewel < bs.waterlevel {
				sensorError = rand.Intn(10)
				go log.Println("sensorError:", sensorError)
			}
			if bs.waterlevel == 90000 {
				close(bs.shipSank)
			}

			//обновление значения датчика уровня воды
			select {
			case <-bs.level:
			default:
			}
			bs.level <- sensor_level

			predWatelewel = bs.waterlevel
		default:
		}
	}

	close(bs.quit)
}
func (bs *BigleS) pumpWorker() {
	bs.pumpConsumedEnergy = 0
	var mode int
	var delta int
	lowtime := 0
	ticker := time.NewTicker(time.Microsecond * 2) // Тикер для постоянных обновлений
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			select {
			case val := <-bs.pumpControl:
				if val != mode && val == 1 {
					delta = 300
				}
				mode = val
			default:
			}
			if mode == 1 {
				if delta > 1 {
					delta /= 2
					bs.pumpConsumedEnergy += delta
				}
				bs.waterlevel--
				bs.pumpConsumedEnergy++

				if bs.waterlevel < 20 {
					lowtime++
				}
				if lowtime > 100 {
					log.Println("Pump Brakes")
					return
				}
			}
		default:

		}

	}
}

func (bs *BigleS) StopPump() {
	bs.pumpControl <- 0
}
func (bs *BigleS) RunPump() {
	bs.pumpControl <- 1
}

func (bs *BigleS) ShipBigleQuit() <-chan struct{} {
	return bs.quit
}

func (bs *BigleS) ShipSunk() <-chan struct{} {
	return bs.shipSank
}

func (bs *BigleS) GetPumpConsumedEnergy() int {
	return bs.pumpConsumedEnergy
}
