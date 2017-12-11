package media

type Media interface {
	GetName() string
	GetContent() ([]byte, error)
}
