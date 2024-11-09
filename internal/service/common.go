package service

type LocalCache interface {
	Set(string, interface{})
	Get(string) (interface{}, bool)
	Del(string)
	GetAll() map[string]interface{}
	Reset()
}

type KafkaProducer interface {
	Send(string, interface{}) error
}
