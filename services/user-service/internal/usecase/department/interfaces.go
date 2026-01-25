package department

type EventPublisher interface {
	Publish(subject string, data interface{}) error
}
