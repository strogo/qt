package templater

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/therecipe/qt/internal/binding/converter"
	"github.com/therecipe/qt/internal/binding/parser"
)

func goFunction(function *parser.Function) string {
	var output = fmt.Sprintf("%v{\n%v\n}", goFunctionHeader(function), goFunctionBody(function))
	if functionIsSupported(parser.CurrentState.ClassMap[function.ClassName()], function) {
		if UseStub() {
			if function.SignalMode != parser.CALLBACK {
				return output
			}
		} else {
			return output
		}
	}
	return ""
}

func goFunctionHeader(function *parser.Function) string {
	return fmt.Sprintf("func %v %v(%v)%v",
		func() string {
			if function.Static || function.Meta == parser.CONSTRUCTOR || function.SignalMode == parser.CALLBACK {
				return ""
			}
			return fmt.Sprintf("(ptr *%v)", function.ClassName())
		}(),

		converter.GoHeaderName(function),

		converter.GoHeaderInput(function),

		converter.GoHeaderOutput(function),
	)
}

func goFunctionBody(function *parser.Function) string {
	var bb = new(bytes.Buffer)
	defer bb.Reset()

	if strings.ToLower(os.Getenv("QT_DEBUG")) == "true" {
		fmt.Fprintf(bb, "defer qt.Recover(\"\t%v%v%v(%v) %v\")\n",
			function.ClassName(),

			func() string {
				return strings.Repeat(" ", 45-len(function.ClassName()))
			}(),

			converter.GoHeaderName(function), converter.GoHeaderInput(function), converter.GoHeaderOutput(function))

		fmt.Fprintf(bb, "qt.Debug(\"\t%v%v%v(%v) %v\")\n",

			function.ClassName(),

			func() string {
				return strings.Repeat(" ", 45-len(function.ClassName()))
			}(),

			converter.GoHeaderName(function), converter.GoHeaderInput(function), converter.GoHeaderOutput(function))
	}

	if parser.CurrentState.ClassMap[function.ClassName()].Stub {
		if converter.GoHeaderOutput(function) != "" {
			return fmt.Sprintf("\nreturn %v", converter.GoOutputParametersFromCFailed(function))
		}
		return ""
	}

	if !(function.Static || function.Meta == parser.CONSTRUCTOR || function.SignalMode == parser.CALLBACK) {
		fmt.Fprintf(bb, "if ptr.Pointer() != nil {\n")
	}

	for _, parameter := range function.Parameters {
		if parameter.Value == "..." || (parameter.Value == "T" && parser.CurrentState.ClassMap[function.ClassName()].Module == "QtAndroidExtras" && function.TemplateModeJNI == "") {
			for i := 0; i < 10; i++ {
				if parameter.Value == "T" {
					fmt.Fprintf(bb, "var p%v, d%v = assertion(%v, %v)\n", i, i, i, parameter.Name)
				} else {
					fmt.Fprintf(bb, "var p%v, d%v = assertion(%v, v...)\n", i, i, i)
				}
				fmt.Fprintf(bb, "if d%v != nil {\ndefer d%v()\n}\n", i, i)

				if parameter.Value == "T" {
					break
				}
			}
		}
	}

	if ((function.Meta == parser.PLAIN && function.SignalMode == "") ||
		(function.Meta == parser.SLOT && function.SignalMode == "") ||
		(function.Meta == parser.CONSTRUCTOR || function.Meta == parser.DESTRUCTOR) && function.SignalMode == "") ||
		(function.Meta == parser.SIGNAL && (function.SignalMode == "" || function.SignalMode == parser.CONNECT || function.SignalMode == parser.DISCONNECT)) ||
		(function.Meta == parser.GETTER || function.Meta == parser.SETTER) {

		//TODO:
		if functionIsSupported(parser.CurrentState.ClassMap[function.ClassName()], function) {
			cppFunction(function)
			if functionIsSupported(parser.CurrentState.ClassMap[function.ClassName()], function) {

				for _, alloc := range converter.GoInputParametersForCAlloc(function) {
					fmt.Fprint(bb, alloc)
				}

				var body = converter.GoOutputParametersFromC(function, fmt.Sprintf("C.%v(%v)", converter.CppHeaderName(function), converter.GoInputParametersForC(function)))
				fmt.Fprint(bb, func() string {
					if converter.GoHeaderOutput(function) != "" {
						switch {
						case function.NeedsFinalizer && classIsSupported(parser.CurrentState.ClassMap[converter.CleanValue(function.Output)]) || function.Meta == parser.CONSTRUCTOR && !(classNeedsCallbackFunctions(parser.CurrentState.ClassMap[function.Name]) || parser.CurrentState.ClassMap[function.Name].IsSubClassOfQObject()):
							{
								return fmt.Sprintf("var tmpValue = %v\nruntime.SetFinalizer(tmpValue, (%v).Destroy%v)\nreturn tmpValue%v",

									body,

									func() string {
										if function.TemplateModeJNI != "" {
											return fmt.Sprintf("*%v", converter.CleanValue(function.Output))
										}
										return converter.GoHeaderOutput(function)
									}(),

									func() string {
										if function.Meta == parser.CONSTRUCTOR {
											return function.Name
										}
										return converter.CleanValue(function.Output)
									}(),

									func() string {
										if function.TemplateModeJNI == "String" {
											return ".ToString()"
										}
										return ""
									}())
							}

						case parser.CurrentState.ClassMap[converter.CleanValue(function.Output)].IsSubClassOfQObject() && converter.GoHeaderOutput(function) != "unsafe.Pointer" || function.Meta == parser.CONSTRUCTOR && parser.CurrentState.ClassMap[converter.CleanValue(function.Name)].IsSubClassOfQObject():
							{
								return fmt.Sprintf("var tmpValue = %v\nif !qt.ExistsSignal(fmt.Sprint(tmpValue.Pointer()), \"QObject::destroyed\") {\ntmpValue.ConnectDestroyed(func(%v){ tmpValue.SetPointer(nil) })\n}\nreturn tmpValue",

									body,

									func() string {
										if parser.CurrentState.ClassMap[function.ClassName()].Module == "QtCore" {
											return "*QObject"
										}
										return "*core.QObject"
									}())
							}

						default:
							{
								return fmt.Sprintf("return %v", body)
							}
						}
					}
					return body
				}())

			}
			function.Access = "public"
		}

	}

	switch function.SignalMode {
	case parser.CALLBACK:
		{
			fmt.Fprintf(bb, "%vif signal := qt.GetSignal(fmt.Sprint(ptr), \"%v::%v%v\"); signal != nil {\n",
				func() string {
					if function.Meta != parser.SLOT {
						return "\n"
					}
					return ""
				}(), function.ClassName(), function.Name, function.OverloadNumber)

			if converter.GoHeaderOutput(function) == "" {
				fmt.Fprintf(bb, "signal.(%v)(%v)", converter.GoHeaderInputSignalFunction(function), converter.GoInputParametersForCallback(function))
			} else {
				fmt.Fprintf(bb, "return %v", converter.GoInput(fmt.Sprintf("signal.(%v)(%v)", converter.GoHeaderInputSignalFunction(function), converter.GoInputParametersForCallback(function)), function.Output, function))
			}

			fmt.Fprintf(bb, "\n}%v\n",
				func() string {
					if converter.GoHeaderOutput(function) == "" {
						if function.Virtual == parser.IMPURE {
							return "else{"
						}
					}
					return ""
				}(),
			)

			if converter.GoHeaderOutput(function) == "" {
				var class, _ = function.Class()
				if class.Module == parser.MOC && function.PureBaseFunction {

				} else {
					if function.Virtual == parser.IMPURE && functionIsSupportedDefault(function) {
						fmt.Fprintf(bb, "New%vFromPointer(ptr).%v%vDefault(%v)", strings.Title(function.ClassName()), strings.Replace(strings.Title(function.Name), parser.TILDE, "Destroy", -1), function.OverloadNumber, converter.GoInputParametersForCallback(function))
					}
				}
			} else {
				var class, _ = function.Class()
				if class.Module == parser.MOC && function.PureBaseFunction {
					fmt.Fprintf(bb, "\nreturn %v", converter.GoInput(converter.GoOutputParametersFromCFailed(function), function.Output, function))
				} else {
					if function.Virtual == parser.IMPURE && functionIsSupportedDefault(function) {
						fmt.Fprintf(bb, "\nreturn %v", converter.GoInput(fmt.Sprintf("New%vFromPointer(ptr).%v%vDefault(%v)", strings.Title(function.ClassName()), strings.Replace(strings.Title(function.Name), parser.TILDE, "Destroy", -1), function.OverloadNumber, converter.GoInputParametersForCallback(function)), function.Output, function))
					} else {
						fmt.Fprintf(bb, "\nreturn %v", converter.GoInput(converter.GoOutputParametersFromCFailed(function), function.Output, function))
					}
				}
			}

			fmt.Fprintf(bb, "%v",
				func() string {
					if converter.GoHeaderOutput(function) == "" {
						if function.Virtual == parser.IMPURE {
							return "\n}"
						}
					}
					return ""
				}(),
			)

		}

	case parser.CONNECT, parser.DISCONNECT:
		{
			fmt.Fprintf(bb, "\nqt.%vSignal(fmt.Sprint(ptr.Pointer()), \"%v::%v%v\"%v)",
				function.SignalMode,

				function.ClassName(),

				function.Name,

				function.OverloadNumber,

				func() string {
					if function.SignalMode == parser.CONNECT {
						return ", f"
					}
					return ""
				}(),
			)
		}
	}

	if (function.Meta == parser.DESTRUCTOR || strings.Contains(function.Name, "deleteLater") || strings.HasPrefix(function.Name, parser.TILDE)) && function.SignalMode == "" {
		if classNeedsCallbackFunctions(parser.CurrentState.ClassMap[function.ClassName()]) || parser.CurrentState.ClassMap[function.ClassName()].IsSubClassOfQObject() {
			fmt.Fprint(bb, "\nqt.DisconnectAllSignals(fmt.Sprint(ptr.Pointer()))")
		}
		fmt.Fprint(bb, "\nptr.SetPointer(nil)")
	}

	if !(function.Static || function.Meta == parser.CONSTRUCTOR || function.SignalMode == parser.CALLBACK) {
		fmt.Fprint(bb, "\n}")
		if converter.GoHeaderOutput(function) != "" {
			fmt.Fprintf(bb, "\nreturn %v", converter.GoOutputParametersFromCFailed(function))
		}
	}

	return bb.String()
}
