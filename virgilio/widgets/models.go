package widgets

type Layout struct {
	Widgets []Widget `json:"layout"`
}

type Serie struct {
	Name string `json:"name"`
	Data []int  `json:"data"`
}

type Widget struct {
	X    int    `json:"x"`
	Y    int    `json:"y"`
	I    int    `json:"i"`
	Id   string `json:"id"`
	Type string `json:"type"`
}

type Label struct {
	Type     string   `json:"type"`
	Title    string   `json:"title"`
	MinerId  string   `json:"minerId"`
	Value    string   `json:"value"`
	Position int      `json:"position"`
	Tags     []string `json:"tags"`
	Snapshot bool     `json:"snapshot"`
}

type Counter struct {
	Type     string   `json:"type"`
	Title    string   `json:"title"`
	MinerId  string   `json:"minerId"`
	Value    float64  `json:"value"`
	Position int      `json:"position"`
	Tags     []string `json:"tags"`
	Snapshot bool     `json:"snapshot"`
}

type Chart struct {
	Type       string   `json:"type"`
	ChartType  string   `json:"chartType"`
	Title      string   `json:"title"`
	MinerId    string   `json:"minerId"`
	Position   int      `json:"position"`
	Tags       []string `json:"tags"`
	Categories []string `json:"categories"`
	Series     []Series `json:"series"`
	Snapshot   bool     `json:"snapshot"`
}

type Table struct {
	Type         string      `json:"type"`
	Title        string      `json:"title"`
	MinerId      string      `json:"minerId"`
	Position     int         `json:"position"`
	Tags         []string    `json:"tags"`
	Unit         string      `json:"bytes"`
	ColumnHeader []string    `json:"columnHeader"`
	RowHeader    []string    `json:"rowHeader"`
	Rows         [][]float64 `json:"rows"`
	Snapshot     bool        `json:"snapshot"`
}

type Series struct {
	Name string    `json:"name"`
	I18n string    `json:"i18n"`
	Data []float64 `json:"data"`
}
