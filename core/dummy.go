package core

// Empty ...
type Empty struct {
}

// DummyEmpty ...
func DummyEmpty() interface{} {
	return &Empty{}
}
