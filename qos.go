package cyclonedds

// #cgo pkg-config: CycloneDDS
// #include <stdlib.h>
// #include <dds/dds.h>
import "C"

type Qos struct {
	native *C.dds_qos_t
}

func CreateQos() *Qos {
	cqos := C.dds_create_qos()
	if cqos == nil {
		return nil
	}

	return &Qos{native: cqos}
}

func (qos *Qos) Delete() {
	C.dds_delete_qos(qos.native)
}
