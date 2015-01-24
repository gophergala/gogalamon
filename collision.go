package main

//sector int
type sint int8
type sv2 [2]sint

type Overworld struct {
	all     map[Entity]*Colider
	sectors map[sv2]map[Entity]*Colider
}

func (o *Overworld) set(e Entity, x, y, r float32) {
	c, ok := o.all[e]
	if ok {
		for oldsector := range c.sectors() {
			delete(sectors[oldsector], e)
		}
	} else {
		c := new(Colider)
	}
	c.x, c.y, c.r = x, y, r
	for newSector := range c.sectors() {
		s, ok := o.sectors[newSector]
		if !ok {
			s = make(map[Entity]*Colider)
			o.sectors[newSector] = s
		}
		s[e] = c
	}

	s := sv2{sint(x / 100), sint(y / 100)}
	if c, ok := o.all[e]; ok {
		if s != c.s {
			delete(o.sectors[c.s], e)
			sector, ok := o.sectors[s]
			if !ok {
				sector = make(map[Entity]*Colider)
				o.sectors[s] = sector
			}
			sector[e] = c
		}
		c.x, c.y, c.r, c.s = x, y, r, s
	} else {
		c := new(Colider)
		c.x, c.y, c.r, c.s = x, y, r, s
		sector, ok := o.sectors[s]
		if !ok {
			sector = make(map[Entity]*Colider)
			o.sectors[s] = sector
		}
		sector[e] = c
	}
}

func (o *Overworld) remove(e Entity) {
	c, ok := o.all[e]
	if !ok {
		return
	}
	delete(o.all, e)
	delete(o.sectors[c.s], e)
}

func (o *Overworld) query(x, y, r float32) []Entity {
	var entities []Entity
	rs := r * r

	sxmax := sint(x + r/100) ///FIX BUF=G
	for sx := sint(x - r/100); sx <= sxmax; sx++ {

	}

	s := sv2{sint(x / 100), sint(y / 100)}

}

type Colider struct {
	x float32
	y float32
	r float32
}

func (c *Colider) sectors() []sv2 {
	sxmin := sint((c.x - c.r) / 100)
	sxmax := sint((c.x+c.r)/100) + 1
	symin := sint((c.y - c.r) / 100)
	symax := sint((c.y+c.r)/100) + 1

	result := make([]sv2, (sxmax-sxmin)*(symax-symin))
	i := 0
	for j := sxmin; j < sxmax; j++ {
		for k := symin; k < symax; k++ {
			result[i] = sv2{j, k}
			i++
		}
	}
	return result
}
