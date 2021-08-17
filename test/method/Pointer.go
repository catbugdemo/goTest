package method

type Point struct {
	X float64
	Y float64
}

func (p *Point) ScaleBy(factor float64)  {
	p.X *= factor
}