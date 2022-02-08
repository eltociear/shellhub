package api

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Header struct {
	Key   string
	Value string
}

type Case struct {
	Path     string
	Query    string
	Header   []Header
	Body     string
	Auth     Auth
	Expected int
}

type Test struct {
	Description string
	Method      string
	Url         string
	Kind        string
	Cases       []Case
}
