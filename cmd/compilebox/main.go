package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/frenata/compilebox"
)

type CodeSubmission struct {
	Language string   `json:"language"`
	Code     string   `json:"code"`
	Stdins   []string `json:"stdins"`
}

func (s CodeSubmission) String() string {
	return fmt.Sprintf("( <CodeSubmission> {Language: %s, Code: Hidden, Stdins: %s} )", s.Language, s.Stdins)
}

type ExecutionResult struct {
	Stdouts []string           `json:"stdouts"`
	Message compilebox.Message `json:"message"`
}

const (
	CompilersFile = "data/compilers.json"
)

// type LanguagesResponse struct {
// 	Languages map[string]compilebox.Language `json:"languages"`
// }

var box compilebox.Interface

func main() {
	var (
		err error
	)

	// on spinup, run a smoke test against the compilebox Docker container
	err = runSmokeTest()
	if err != nil {
		panic(err)
	}

	port := getEnv("COMPILEBOX_API_SERVER_PORT", "31337")

	box = compilebox.New(CompilersFile)

	http.HandleFunc("/languages/", getLangs)
	http.HandleFunc("/eval/", evalCode)

	log.Println("testbox listening on " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	log.Printf("Environment variable %s not found, setting to %s", key, fallback)
	os.Setenv(key, fallback)
	return fallback
}

func evalCode(w http.ResponseWriter, r *http.Request) {
	log.Println("Received code subimssion...")
	decoder := json.NewDecoder(r.Body)
	var submission CodeSubmission
	err := decoder.Decode(&submission)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	// fmt.Printf("...along with %d stdin inputs\n", len(submission.Stdins))
	fmt.Println(submission)
	stdouts, msg := box.EvalWithStdins(submission.Language, submission.Code, submission.Stdins)
	log.Println(stdouts, msg)

	if len(stdouts) == 0 {
		log.Println("Code produced no output")
		stdouts = append(stdouts, "ZERO OUTPUTS")
	}

	buf, _ := json.MarshalIndent(ExecutionResult{
		Stdouts: stdouts,
		Message: msg,
	}, "", "   ")

	w.Header().Set("Content-Type", "application/json")
	w.Write(buf)
}

func getLangs(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Received languages request...")
	workingLangs := make(map[string]compilebox.Language)

	// make a list of currently supported languages
	for k, v := range box.LanguageMap {
		if v.Disabled != "true" {
			workingLangs[k] = v
		}
	}

	fmt.Printf("currently supporting %d of %d known languages\n", len(workingLangs), len(box.LanguageMap))

	// add boilerplate and comment info
	// log.Println(workingLangs)

	// encode language list
	buf, _ := json.MarshalIndent(workingLangs, "", "   ")

	// write working language list back to client
	w.Header().Set("Content-Type", "application/json")
	w.Write(buf)
}

// executes a hello world program for all supported languages inside the compilerbox
// Docker container.
func runSmokeTest() error {
	testBox := compilebox.New(CompilersFile)

	// currently passing:
	compilerTests := make(map[string]string)
	compilerTests["C++"] = "#include <iostream>\nusing namespace std;\n\nint main() {\n\tcout<<\"Hello\";\n\treturn 0;\n}"
	compilerTests["Java"] = "\n\nimport java.io.*;\n\nclass myCode\n{\n\tpublic static void main (String[] args) throws java.lang.Exception\n\t{\n\t\t\n\t\tSystem.out.println(\"Hello\");\n\t}\n}"
	compilerTests["C#"] = "using System;\n\npublic class Challenge\n{\n\tpublic static void Main()\n\t{\n\t\t\tConsole.WriteLine(\"Hello\");\n\t}\n}"
	compilerTests["Clojure"] = "(println \"Hello\")"
	compilerTests["Perl"] = "use strict;\nuse warnings\n;use v5.14; say 'Hello';"
	compilerTests["Golang"] = "package main\nimport \"fmt\"\n\nfunc main(){\n  \n\tfmt.Printf(\"Hello\")\n}"
	compilerTests["JavaScript"] = "console.log(\"Hello\");"
	compilerTests["Python"] = "print(\"Hello\")"
	compilerTests["Ruby"] = "puts \"Hello\""
	compilerTests["Bash"] = "echo 'Hello' "
	compilerTests["PHP"] = "<?php\n$ho = fopen('php://stdout', \"w\");\n\nfwrite($ho, \"Hello\");\n\n\nfclose($ho);\n"

	// currently broken:
	// Haskell ghc missing, maybe need to rebuild docker file
	// compilerTests["Haskell"] = "module Main where\nmain = putStrLn \"Hello\""
	//
	// Scala: don't understand the error this generates
	// compilerTests["Scala"] = "object HelloWorld {def main(args: Array[String]) = println(\"Hello\")}"
	//
	// Rust seems to be missing and there's a problem setting environment variables
	// compilerTests["Rust"] = "fn main() {\n\tprintln!(\"Hello\");\n}"

	stdin := ""
	expected := "Hello"
	langResults := make(map[string]string)

	// run tests for each language
	for lang, code := range compilerTests {
		stdouts, msg := testBox.EvalWithStdins(lang, code, []string{stdin})

		log.Println(stdouts[0], msg)

		if stdouts[0] == expected {
			log.Printf("%s passed 'Hello' test.", lang)
			langResults[lang] = "Pass"
		} else {
			log.Println(stdouts)
			log.Printf("%s failed 'Hello' test.", lang)
			langResults[lang] = "Fail"
		}

		fmt.Println("-----------------------------------------------------")
	}

	for lang, result := range langResults {
		fmt.Printf("%s -> %s\n", lang, result)
	}

	return nil
}
