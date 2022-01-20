package collector

type file struct {
	date       string
	sourceId   string
	sourceType string
	paths      []*filePath
}
type filePath struct {
	account string
	pwd     string
	host    string
	path    string
	kind  	string
	saveTo  string
}

type api struct {
	date       string
	sourceId   string
	sourceType string
	urls       []apiUrl
}
type apiUrl struct {
	url      string
	params   map[string]string
	fields   []string
	nextStep *apiUrl
}

type mongoDb struct {
	date       string
	sourceId   string
	sourceType string
	queries    []dbQuery
}

type mysql struct {
	date       string
	sourceId   string
	sourceType string
	queries    []dbQuery
}

type dbQuery struct {
	query  string
	params map[string]string
}
