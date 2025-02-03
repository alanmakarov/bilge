package bilge

// interface
type Ship interface {
	ShipBigleQuit() <-chan struct{}
	ShipSunk() <-chan struct{}
}
type Pump interface {
	RunPump()
	StopPump()
	GetPumpConsumedEnergy() int
}
type WaterLevelSensor interface {
	GetWaterLevelInformer() int
}
