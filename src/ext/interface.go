package ext

type ParallelMap interface {
	Set(k, v interface{}) bool
	Get(k interface{}) interface{}
	Delete(k interface{}) bool
	Len() int

	//DirtySet	after DirtySet Get maybe not get value
}
