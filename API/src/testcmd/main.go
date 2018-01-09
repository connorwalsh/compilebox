package main

import (
	"log"
	"testbox"
)

func main() {
	box := testbox.New("data/compilers.json")
	compilerTests := make(map[string]string)

	// need test for Haskell
	// FIXME no output
	compilerTests["C++"] = "#include <iostream>\nusing namespace std;\n\nint main() {\n\tcout<<\"Hello\";\n\treturn 0;\n}"
	// FIXME
	// compilerTests["Java"] = "\n\nimport java.io.*;\n\nclass myCode\n{\n\tpublic static void main (String[] args) throws java.lang.Exception\n\t{\n\t\t\n\t\tSystem.out.println(\"Hello\");\n\t}\n}"
	// FIXME
	// compilerTests["C#"] = "using System;\n\npublic class Challenge\n{\n\tpublic static void Main()\n\t{\n\t\t\tConsole.WriteLine(\"Hello\");\n\t}\n}"
	// FIXME
	//compilerTests["Scala"] = "object HelloWorld {def main(args: Array[String]) = println(\"Hello\")}"
	// FIXME compiler
	//compilerTests["Rust"] = "fn main() {\n\tprintln!(\"Hello\");\n}"
	// FIXME
	//compilerTests["PHP"] = "<?php\n$ho = fopen('php://stdout', \"w\");\n\nfwrite($ho, \"Hello\");\n\n\nfclose($ho);\n"
	// compilerTests["Clojure"] = "(println \"Hello\")"
	// compilerTests["Perl"] = "use strict;\nuse warnings\n;use v5.14; say 'Hello';"
	// compilerTests["Golang"] = "package main\nimport \"fmt\"\n\nfunc main(){\n  \n\tfmt.Printf(\"Hello\")\n}"
	// compilerTests["JavaScript"] = "console.log(\"Hello\");"
	// compilerTests["Python"] = "print(\"Hello\")"
	// compilerTests["Ruby"] = "puts \"Hello\""
	// compilerTests["Bash"] = "echo 'Hello' "

	/*
		"MySQL":"create table myTable(name varchar(10));\ninsert into myTable values(\"Hello\");\nselect * from myTable;",
		"Objective-C": "#include <Foundation/Foundation.h>\n\n@interface Challenge\n+ (const char *) classStringValue;\n@end\n\n@implementation Challenge\n+ (const char *) classStringValue;\n{\n    return \"Hey!\";\n}\n@end\n\nint main(void)\n{\n    printf(\"%s\\n\", [Challenge classStringValue]);\n    return 0;\n}",
		"VB.NET": "Imports System\n\nPublic Class Challenge\n\tPublic Shared Sub Main() \n    \tSystem.Console.WriteLine(\"Hello\")\n\tEnd Sub\nEnd Class",
	*/

	stdin := "" + testbox.Seperator
	expected := "Hello" + testbox.Seperator
	langResults := make(map[string]string)
	for lang, code := range compilerTests {
		out, msg := box.CompileAndChallenge(lang, code, stdin, expected)
		// oOut, oMsg := box.CompileAndPrint(lang, code, "test")
		log.Println(out, msg)
		// log.Println(oOut, oMsg)
		langResults[lang] = out[""]
		if out[""] == "Pass" {
			log.Printf("%s passed 'Hello' test.", lang)
		} else {
			log.Println(out)
			log.Printf("%s failed 'Hello' test.", lang)
		}
	}

	for lang, result := range langResults {
		log.Printf("%s -> %s", lang, result)
	}
}
