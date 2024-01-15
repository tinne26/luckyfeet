package utils

type Blinker struct {
	minValue float64
	maxValue float64
	speed float64
	osc float64
	oscIncreasing bool
}

func NewBlinker(minValue, maxValue, speed float64) *Blinker {
	if speed > 1.0 { panic("speed can't be >1.0") }
	if speed <= 0 { panic("speed can't be <=0") }
	if maxValue <= minValue { panic("maxValue must be strictly greater than minValue") }
	return &Blinker{
		minValue: minValue,
		maxValue: maxValue,
		speed: speed,
		oscIncreasing: true,
	}
}

func (self *Blinker) Update() {
	if self.oscIncreasing {
		self.osc += self.speed
		if self.osc > 1.0 {
			self.oscIncreasing = false
			self.osc = 1.0 - (self.osc - 1.0)
		}
	} else {
		self.osc -= self.speed
		if self.osc < 0.0 {
			self.oscIncreasing = true
			self.osc = -self.osc
		}
	}
}

func (self *Blinker) Value() float64 {
	return self.minValue + (self.maxValue - self.minValue)*self.osc
}
