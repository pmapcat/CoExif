package main

type Dispatcher struct {
	Instances []ExifProc
}

func (d *Dispatcher) Init() {
	inst := make([]ExifProc, BASE_SIZE)
	for k, _ := range inst {
		e := ExifProc{}
		e.Init()
		e.StartProc()
		inst[k] = e
	}
	d.Instances = inst
}

func (d *Dispatcher) Spawn() ExifProc {
	e := ExifProc{}
	e.Init()
	e.StartProc()
	d.Instances = append(d.Instances, e)
	return e
}

func (d *Dispatcher) Exit() {
	for _, v := range d.Instances {
		v.KillProc()
	}
}

func (d *Dispatcher) GETMeta(path string, probably_filter ...string) IdioticJSON {
	// Will hang until result
	for {
		for _, v := range d.Instances {
			if !v.Busy {
				return v.GETMeta(path, probably_filter...)
			}
		}
		// spawn process
		if AUTO_SPAWN == true {
			method := d.Spawn()
			return method.GETMeta(path, probably_filter...)
		}
	}
}

func (d *Dispatcher) POSTMeta(path string, datum StandartJSON) error {
	// Will hang until result
	for {
		for _, v := range d.Instances {
			if !v.Busy {
				return v.UPDATEMeta(path, datum)
			}
		}
		// spawn process
		if AUTO_SPAWN == true {
			method := d.Spawn()
			return method.UPDATEMeta(path, datum)
		}
	}
}
