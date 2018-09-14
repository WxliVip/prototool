// Copyright (c) 2018 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// Package exec brings all the functionality of Prototool together in a format
// easily consumable by CLI libraries. It is effectively the glue between
// internal/cmd and all other packages.
package exec

import (
	"io"

	"go.uber.org/zap"
)

// ExitError is an error that signals to exit with a certain code.
type ExitError struct {
	Code    int
	Message string
}

// Error implements error.
func (e *ExitError) Error() string {
	return e.Message
}

// Runner runs commands.
//
// The args given are the args from the command line.
// Each additional parameter generally refers to a command-specific flag.
type Runner interface {
	Init(args []string, uncomment bool) error
	Create(args []string, pkg string) error
	Version() error
	Download() error
	Clean() error
	Files(args []string) error
	Compile(args []string, dryRun bool) error
	Gen(args []string, dryRun bool) error
	DescriptorProto(args []string) error
	FieldDescriptorProto(args []string) error
	ServiceDescriptorProto(args []string) error
	Lint(args []string, listAllLinters bool, listLinters bool) error
	ListLintGroup(group string) error
	ListAllLintGroups() error
	Format(args []string, overwrite, diffMode, lintMode, fix bool) error
	BinaryToJSON(args []string) error
	JSONToBinary(args []string) error
	All(args []string, disableFormat, disableLint, fix bool) error
	GRPC(args, headers []string, address, method, data, callTimeout, connectTimeout, keepaliveTime string, stdin bool) error
}

// RunnerOption is an option for a new Runner.
type RunnerOption func(*runner)

// RunnerWithLogger returns a RunnerOption that uses the given logger.
//
// The default is to use zap.NewNop().
func RunnerWithLogger(logger *zap.Logger) RunnerOption {
	return func(runner *runner) {
		runner.logger = logger
	}
}

// RunnerWithCachePath returns a RunnerOption that uses the given cache path.
func RunnerWithCachePath(cachePath string) RunnerOption {
	return func(runner *runner) {
		runner.cachePath = cachePath
	}
}

// RunnerWithJSON returns a RunnerOption that will print failures as JSON.
func RunnerWithJSON() RunnerOption {
	return func(runner *runner) {
		runner.json = true
	}
}

// RunnerWithNoCache returns a RunnerOption that will disable the protoc cache.
func RunnerWithNoCache() RunnerOption {
	return func(runner *runner) {
		runner.noCache = true
	}
}

// RunnerWithPrintFields returns a RunnerOption that uses the given colon-separated
// print fields. The default is filename:line:column:message.
func RunnerWithPrintFields(printFields string) RunnerOption {
	return func(runner *runner) {
		runner.printFields = printFields
	}
}

// RunnerWithProtocURL returns a RunnerOption that uses the given protoc zip file URL.
func RunnerWithProtocURL(protocURL string) RunnerOption {
	return func(runner *runner) {
		runner.protocURL = protocURL
	}
}

// NewRunner returns a new Runner.
//
// workDirPath should generally be the current directory.
// input and output generally refer to stdin and stdout.
func NewRunner(workDirPath string, input io.Reader, output io.Writer, options ...RunnerOption) Runner {
	return newRunner(workDirPath, input, output, options...)
}
