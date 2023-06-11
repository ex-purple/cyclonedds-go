package cyclonedds

import (
	"testing"
)

func TestCreateParticipant(t *testing.T) {
	p, err := CreateParticipant(0, nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = p.Delete()
	if err != nil {
		t.Fatal(err.Error())
	}
}

type MyData struct {
	X int32
}

func TestCreateTopic(t *testing.T) {

	var data MyData

	p, _ := CreateParticipant(0, nil)
	defer p.Delete()

	_, err := CreateTopic(p, "test", data, nil)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestCreatePublisher(t *testing.T) {
	p, _ := CreateParticipant(0, nil)
	defer p.Delete()

	_, err := CreatePublisher(p, nil)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestCreateWriter(t *testing.T) {
	var data MyData

	p, _ := CreateParticipant(0, nil)
	defer p.Delete()

	topic, _ := CreateTopic(p, "test", data, nil)
	w, err := CreateWtiter(p, topic)
	if err != nil {
		t.Fatal(err.Error())
	}

	err = w.Write(data)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestCreateSubscriber(t *testing.T) {
	p, _ := CreateParticipant(0, nil)
	defer p.Delete()

	_, err := CreateSubscriber(p, nil)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestCreateReader(t *testing.T) {
	var data MyData

	p, _ := CreateParticipant(0, nil)
	defer p.Delete()

	topic, _ := CreateTopic(p, "test", data, nil)
	_, err := CreateReader(p, topic)
	if err != nil {
		t.Fatal(err.Error())
	}
}
