package types

type CodeModelGCCInvocation struct {
	GCC       string
	Arguments []string
}

type CodeModelLibrary struct {
	Name            string
	SourceDirectory string
	Invocations     []CodeModelGCCInvocation
}

type CodeModelBuilder struct {
	MergedSketchFile string
	Core             CodeModelLibrary
	Libraries        []CodeModelLibrary
}
