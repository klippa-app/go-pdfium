package responses

type JavaScriptAction struct {
	Name   string
	Script string
}

type GetJavaScriptActions struct {
	JavaScriptActions []JavaScriptAction
}
