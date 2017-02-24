package wiki

type RenderContext interface {
	PageNames() []string
}


type Page struct {
	Title string
	Body  []byte
}

