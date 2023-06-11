package cyclonedds

// #cgo pkg-config: CycloneDDS
// #include <stdlib.h>
// #include <dds/dds.h>
import "C"
import (
	"bytes"
	"encoding/gob"
	"errors"
	"unsafe"
)

type DomainID uint32

type Participant struct {
	native C.int32_t
}

func (e Participant) getNative() C.int32_t {
	return e.native
}

type Topic struct {
	native C.int32_t
}

func (e Topic) getNative() C.int32_t {
	return e.native
}

type Publisher struct {
	native C.int32_t
}

func (e Publisher) getNative() C.int32_t {
	return e.native
}

type Subscriber struct {
	native C.int32_t
}

func (e Subscriber) getNative() C.int32_t {
	return e.native
}

type DataWriter struct {
	native C.int32_t
}

func (e DataWriter) getNative() C.int32_t {
	return e.native
}

type DataReader struct {
	native C.int32_t
}

func (e DataReader) getNative() C.int32_t {
	return e.native
}

type ParticipantOrPublisher interface {
	Participant | Publisher
	getNative() C.int32_t
}

type ParticipantOrSubsciber interface {
	Participant | Subscriber
	getNative() C.int32_t
}

type Entity interface {
	Participant | Topic | Publisher | Subscriber | DataWriter | DataReader
	getNative() C.int32_t
}

func checkEntity[E Entity](e E) (E, error) {
	if e.getNative() > 0 {
		return e, nil
	}
	return e, errors.New(C.GoString(C.dds_strretcode(-e.getNative())))
}

func CreateParticipant(domain DomainID, qos *Qos) (Participant, error) {
	e := Participant{native: C.dds_create_participant(C.uint32_t(domain), nil, nil)}
	return checkEntity(e)
}

func (p Participant) Delete() error {
	rc := C.dds_delete(p.getNative())
	if rc != C.DDS_RETCODE_OK {
		return errors.New(C.GoString(C.dds_strretcode(-rc)))
	}
	return nil
}

func CreateTopic(participant Participant, name string, dataType any, qos *Qos) (Topic, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	dStruct := dynamicTypeCreate(participant, dataType)
	defer C.dds_dynamic_type_unref(&dStruct)

	var typeInfo *C.dds_typeinfo_t
	rc := C.dds_dynamic_type_register(&dStruct, &typeInfo)
	if rc != C.DDS_RETCODE_OK {
		return Topic{native: 0}, errors.New(C.GoString(C.dds_strretcode(-rc)))
	}
	defer C.dds_free_typeinfo(typeInfo)

	var topicDesc *C.dds_topic_descriptor_t
	rc = C.dds_create_topic_descriptor(C.DDS_FIND_SCOPE_LOCAL_DOMAIN, participant.getNative(), typeInfo, 0, &topicDesc)
	if rc != C.DDS_RETCODE_OK {
		return Topic{native: 0}, errors.New(C.GoString(C.dds_strretcode(-rc)))
	}
	defer C.dds_delete_topic_descriptor(topicDesc)

	e := Topic{native: C.dds_create_topic(participant.getNative(), topicDesc, cname, nil, nil)}
	return checkEntity(e)
}

func CreatePublisher(participant Participant, qos *Qos) (Publisher, error) {
	e := Publisher{native: C.dds_create_publisher(C.int32_t(participant.native), nil, nil)}
	return checkEntity(e)
}

func CreateSubscriber(participant Participant, qos *Qos) (Subscriber, error) {
	e := Subscriber{native: C.dds_create_subscriber(C.int32_t(participant.native), nil, nil)}
	return checkEntity(e)
}

func CreateWtiter[P ParticipantOrPublisher](participantOrPublisher P, topic Topic) (DataWriter, error) {
	e := DataWriter{native: C.dds_create_writer(participantOrPublisher.getNative(), topic.native, nil, nil)}
	return checkEntity(e)
}

func (w DataWriter) Write(data any) error {

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(data); err != nil {
		return err
	}
	msg := buf.Bytes()

	rc := C.dds_write(w.getNative(), unsafe.Pointer(&msg[0]))
	if rc != C.DDS_RETCODE_OK {
		return errors.New(C.GoString(C.dds_strretcode(-rc)))
	}
	return nil
}

func CreateReader[P ParticipantOrSubsciber](participantOrSubscriber P, topic Topic) (DataReader, error) {
	e := DataReader{native: C.dds_create_reader(participantOrSubscriber.getNative(), topic.native, nil, nil)}
	return checkEntity(e)
}
