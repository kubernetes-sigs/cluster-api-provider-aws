package progress

type Manual struct {
	N     int64
	Total int64
	Err   error
}

func (p Manual) Current() int64 {
	return int64(p.N)
}

func (p Manual) Size() int64 {
	return int64(p.Total)
}

func (p Manual) Error() error {
	return p.Err
}

func (p Manual) Progress() Progress {
	return Progress{
		current: p.N,
		size:    p.Total,
		err:     p.Err,
	}
}

func (p *Manual) SetCompleted() {
	p.Err = ErrCompleted
	if p.N > 0 && p.Total <= 0 {
		p.Total = p.N
		return
	}
}
