package track

import (
	"bytes"
	"errors"
	"fmt"
	"regexp/syntax"
)

type CaptureGroupHandler func(index int, name string, group *syntax.Regexp, generator Generator, args *GeneratorArgs) string

type generatorFactory func(regexp *syntax.Regexp, args *GeneratorArgs) (*internalGenerator, error)

var generatorFactories map[syntax.Op]generatorFactory

type GeneratorArgs struct {
	Flags               syntax.Flags
	CaptureGroupHandler CaptureGroupHandler
}

func (a *GeneratorArgs) initialize() error {
	if a.CaptureGroupHandler == nil {
		a.CaptureGroupHandler = defaultCaptureGroupHandler
	}
	return nil
}

// Generator generates random strings.
type Generator interface {
	Generate() string
	String() string
}

func NewGenerator(pattern string, inputArgs *GeneratorArgs) (generator Generator, err error) {
	args := GeneratorArgs{}

	// Copy inputArgs so the caller can't change them.
	if inputArgs != nil {
		args = *inputArgs
	}
	if err = args.initialize(); err != nil {
		return nil, err
	}

	var regexp *syntax.Regexp
	regexp, err = syntax.Parse(pattern, args.Flags)
	if err != nil {
		return
	}

	var gen *internalGenerator
	gen, err = newGenerator(regexp, &args)
	if err != nil {
		return
	}

	return gen, nil
}

func init() {
	generatorFactories = map[syntax.Op]generatorFactory{
		syntax.OpLiteral:        opCopyMatch,
		syntax.OpAnyCharNotNL:   opCopyMatch,
		syntax.OpAnyChar:        opCopyMatch,
		syntax.OpQuest:          opCopyMatch,
		syntax.OpStar:           opCopyMatch,
		syntax.OpPlus:           opCopyMatch,
		syntax.OpRepeat:         opCopyMatch,
		syntax.OpCharClass:      opCopyMatch,
		syntax.OpConcat:         opConcat,
		syntax.OpCapture:        opCapture,
		syntax.OpEmptyMatch:     noop,
		syntax.OpAlternate:      noop,
		syntax.OpBeginLine:      noop,
		syntax.OpEndLine:        noop,
		syntax.OpBeginText:      noop,
		syntax.OpEndText:        noop,
		syntax.OpWordBoundary:   noop,
		syntax.OpNoWordBoundary: noop,
	}
}

type internalGenerator struct {
	Name         string
	GenerateFunc func() string
}

func (gen *internalGenerator) Generate() string {
	return gen.GenerateFunc()
}

func (gen *internalGenerator) String() string {
	return gen.Name
}

// Create a new generator for each expression in regexps.
func newGenerators(regexps []*syntax.Regexp, args *GeneratorArgs) ([]*internalGenerator, error) {
	generators := make([]*internalGenerator, len(regexps))
	var err error

	// create a generator for each alternate pattern
	for i, subR := range regexps {
		generators[i], err = newGenerator(subR, args)
		if err != nil {
			return nil, err
		}
	}

	return generators, nil
}

// Create a new generator for r.
func newGenerator(regexp *syntax.Regexp, args *GeneratorArgs) (generator *internalGenerator, err error) {
	simplified := regexp.Simplify()

	factory, ok := generatorFactories[simplified.Op]
	if ok {
		return factory(simplified, args)
	}

	return nil, errors.New("invalid generator pattern")
}

// Generator that does nothing.
func noop(regexp *syntax.Regexp, args *GeneratorArgs) (*internalGenerator, error) {
	return &internalGenerator{regexp.String(), func() string {
		return ""
	}}, nil
}

func opCopyMatch(regexp *syntax.Regexp, args *GeneratorArgs) (*internalGenerator, error) {
	return &internalGenerator{regexp.String(), func() string {
		return runesToString(regexp.Rune...)
	}}, nil
}

func opConcat(regexp *syntax.Regexp, genArgs *GeneratorArgs) (*internalGenerator, error) {
	generators, err := newGenerators(regexp.Sub, genArgs)
	if err != nil {
		return nil, fmt.Errorf("error creating generators for concat pattern /%s/", regexp)
	}

	return &internalGenerator{regexp.String(), func() string {
		var result bytes.Buffer
		for _, generator := range generators {
			result.WriteString(generator.Generate())
		}
		return result.String()
	}}, nil
}

func opCapture(regexp *syntax.Regexp, args *GeneratorArgs) (*internalGenerator, error) {
	if err := enforceSingleSub(regexp); err != nil {
		return nil, err
	}

	groupRegexp := regexp.Sub[0]
	generator, err := newGenerator(groupRegexp, args)
	if err != nil {
		return nil, err
	}

	index := regexp.Cap - 1
	return &internalGenerator{regexp.String(), func() string {
		return args.CaptureGroupHandler(index, regexp.Name, groupRegexp, generator, args)
	}}, nil
}

func defaultCaptureGroupHandler(index int, name string, group *syntax.Regexp, generator Generator, args *GeneratorArgs) string {
	return generator.Generate()
}

// Return an error if r has 0 or more than 1 sub-expression.
func enforceSingleSub(regexp *syntax.Regexp) error {
	if len(regexp.Sub) != 1 {
		return fmt.Errorf("%s expected 1 sub-expression, but got %d: %s", opToString(regexp.Op), len(regexp.Sub), regexp)
	}
	return nil
}

// runesToString converts a slice of runes to the string they represent.
func runesToString(runes ...rune) string {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("RunesToString panicked")
		}
	}()
	var buffer bytes.Buffer
	for _, r := range runes {
		buffer.WriteRune(r)
	}
	return buffer.String()
}

// opToString gets the string name of a regular expression operation.
func opToString(op syntax.Op) string {
	switch op {
	case syntax.OpNoMatch:
		return "OpNoMatch"
	case syntax.OpEmptyMatch:
		return "OpEmptyMatch"
	case syntax.OpLiteral:
		return "OpLiteral"
	case syntax.OpCharClass:
		return "OpCharClass"
	case syntax.OpAnyCharNotNL:
		return "OpAnyCharNotNL"
	case syntax.OpAnyChar:
		return "OpAnyChar"
	case syntax.OpBeginLine:
		return "OpBeginLine"
	case syntax.OpEndLine:
		return "OpEndLine"
	case syntax.OpBeginText:
		return "OpBeginText"
	case syntax.OpEndText:
		return "OpEndText"
	case syntax.OpWordBoundary:
		return "OpWordBoundary"
	case syntax.OpNoWordBoundary:
		return "OpNoWordBoundary"
	case syntax.OpCapture:
		return "OpCapture"
	case syntax.OpStar:
		return "OpStar"
	case syntax.OpPlus:
		return "OpPlus"
	case syntax.OpQuest:
		return "OpQuest"
	case syntax.OpRepeat:
		return "OpRepeat"
	case syntax.OpConcat:
		return "OpConcat"
	case syntax.OpAlternate:
		return "OpAlternate"
	}

	return "Unknown"
}
