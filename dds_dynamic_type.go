package cyclonedds

// #cgo pkg-config: CycloneDDS
// #include <stdlib.h>
// #include <dds/dds.h>
//
// dds_dynamic_member_descriptor_t dds_dynamic_member_prim(dds_dynamic_type_kind_t type, const char *name) {
//     return DDS_DYNAMIC_MEMBER_PRIM(type, name);
// }
import "C"
import (
	"reflect"
	"unsafe"
)

func kindToCtype(kind reflect.Kind) C.dds_dynamic_type_kind_t {
	switch kind {
	case reflect.Bool:
		return C.DDS_DYNAMIC_BOOLEAN
	// case reflect.
	// 	return DDS_DYNAMIC_BYTE
	case reflect.Int16:
		return C.DDS_DYNAMIC_INT16
	case reflect.Int32:
		return C.DDS_DYNAMIC_INT32
	case reflect.Int64:
		return C.DDS_DYNAMIC_INT64
	case reflect.Uint16:
		return C.DDS_DYNAMIC_UINT16
	case reflect.Uint32:
		return C.DDS_DYNAMIC_UINT32
	case reflect.Uint64:
		return C.DDS_DYNAMIC_UINT64
	case reflect.Float32:
		return C.DDS_DYNAMIC_FLOAT32
	case reflect.Float64:
		return C.DDS_DYNAMIC_FLOAT64
	// case reflect.:
	// 	return DDS_DYNAMIC_FLOAT128
	case reflect.Int8:
		return C.DDS_DYNAMIC_INT8
	case reflect.Uint8:
		return C.DDS_DYNAMIC_UINT8
	default:
		return C.DDS_DYNAMIC_NONE
	}
}

func dynamicTypeCreate(participant Participant, dataType any) C.dds_dynamic_type_t {

	t := reflect.TypeOf(dataType)
	dname := C.CString(t.Name())
	defer C.free(unsafe.Pointer(dname))

	var dtDesc C.dds_dynamic_type_descriptor_t
	dtDesc.kind = C.DDS_DYNAMIC_STRUCTURE
	dtDesc.name = dname

	dStruct := C.dds_dynamic_type_create(participant.getNative(), dtDesc)

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fname := C.CString(f.Name)
		defer C.free(unsafe.Pointer(fname))

		ctype := kindToCtype(f.Type.Kind())
		C.dds_dynamic_type_add_member(&dStruct, C.dds_dynamic_member_prim(ctype, fname))
	}

	return dStruct
}
