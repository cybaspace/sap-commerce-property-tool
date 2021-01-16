package main

import (
	"bufio"
	_ "bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const localPropertiesFile = "local.properties"
const mylocalPropertiesFile = "mylocal.properties"
const filesListFile = "property-files"

type property struct {
	value            string
	firstSeen        string
	firstSeenLine    int
	overriddenBy     string
	overriddenByLine int
}

var properties map[string]property
var orderedProperties []string

var configPath string
var propertyFiles []string
var fileListPtr *string
var system string
var verbose bool
var output string

func main() {

	flag.StringVar(&configPath, "path", "./", "Path where property files are stored")
	fileListPtr = flag.String("files", "", "Comma separated list of property and csv files")
	flag.StringVar(&system, "system", "", "Name of system for which properties should be considered")
	flag.BoolVar(&verbose, "v", false, "Print information to the console")
	flag.StringVar(&output, "output", localPropertiesFile, "Send output to a file or with '<console>' to console (default='local.properties')")
	flag.Parse()

	if !strings.HasSuffix(configPath, "/") {
		configPath = configPath + "/"
	}

	if len(flag.Args()) < 1 {
		showHelp()
		return
	}

	switch strings.ToLower(flag.Args()[0]) {
	case "generate":
		generate(flag.Args())
	case "get":
		getProperty(flag.Args())
	case "diff":
		diffFiles(flag.Args())
	case "list":
		listSystems()
	default:
		showHelp()
	}
}

// Help

func showHelp() {
	fmt.Println("SAP Commerce property tool")
	fmt.Println("  usage: yprops <command> [<args>]")
	fmt.Println()
	fmt.Println("  commands")
	fmt.Println("    generate")
	fmt.Println("      generate local.properties file out of many property files")
	fmt.Println("    get")
	fmt.Println("      get one property value out of many property files")
	fmt.Println("    diff")
	fmt.Println("      show difference between two property files")
	fmt.Println("    help")
	fmt.Println("      show this help page")
	fmt.Println()
	fmt.Println("  use: props <command> help")
	fmt.Println("  to get help about the specific command")
}

func showGenerateHelp(systems []string, invalidSystem string) {
	fmt.Println("SAP Commerce property tool")
	fmt.Println("  usage: yprops generate <system>")
	if invalidSystem != "" {
		fmt.Println()
		fmt.Println("    requested system '" + invalidSystem + "' not found")
	}
	fmt.Println()
	fmt.Println("    defined systems")
	for _, system := range systems {
		fmt.Println("      " + system)
	}
}

func showGetHelp(systems []string, invalidSystem string) {
	fmt.Println("SAP Commerce property tool")
	fmt.Println("  usage: yprops get <system> <property>")
	if invalidSystem != "" {
		fmt.Println()
		fmt.Println("    requested system '" + invalidSystem + "' not found")
	}
	fmt.Println()
	fmt.Println("    defined systems:")
	for _, system := range systems {
		fmt.Println("    - " + system)
	}
}

func showDiffHelp(file1 string, file2 string) {
	fmt.Println("SAP Commerce property tool")
	fmt.Println("  usage: yprops diff <property-file-1> <property-file-2>")
	if file1 != "<ok>" {
		fmt.Println()
		if file1 == "" {
			fmt.Println("    Specify a property file 1")
		} else {
			fmt.Println("    property-file-1 '" + file1 + "' not found")
		}
	}
	if file2 != "<ok>" {
		fmt.Println()
		if file2 == "" {
			fmt.Println("    Specify a property file 2")
		} else {
			fmt.Println("    property-file-2 '" + file2 + "' not found")
		}
	}
}

// generate
func generate(args []string) {

	if system == "" {
		systems := readSystemNames()
		showGenerateHelp(systems, "")
		return
	}

	if system != "<all>" {
		evaluateAllPropertyFiles()
		createOutput()
	} else {
		expandFileList()
		systems := readSystemNames()
		outputOriginal := output
		for _, system = range systems {
			evaluateAllPropertyFiles()
			if outputOriginal != "<console>" {
				if outputOriginal == localPropertiesFile {
					output = "local-" + system + ".properties"
				} else {
					output = outputOriginal + "-" + system
				}
			}
			createOutput()
		}
	}
}

// get
func getProperty(args []string) {

	var requestedProperty = args[1]

	if system == "" {
		systems := readSystemNames()
		showGenerateHelp(systems, "")
		return
	}

	evaluateAllPropertyFiles()

	if !strings.Contains(requestedProperty, ",") {
		if value, ok := properties[requestedProperty]; ok {
			fmt.Print(value.value)
		} else {
			os.Exit(1)
		}
	} else {
		allRequestedProperties := strings.Split(requestedProperty, ",")
		for _, property := range allRequestedProperties {
			if value, ok := properties[property]; ok {
				fmt.Println(property + "=" + value.value)
			} else {
				fmt.Println("+++ Missing value for property key: " + property)
				os.Exit(1)
			}
		}
	}
}

// Diff
func diffFiles(args []string) {
	if len(strings.Split(*fileListPtr, ",")) != 2 && (system == "" || !strings.Contains(system, ",")) {
		showDiffHelp("", "")
		return
	}

	var diffArg1 string
	var diffArg2 string

	var props1 map[string]string
	var props2 map[string]string

	if system != "" && strings.Contains(system, ",") {
		systems := strings.Split(system, ",")

		props1 = evaluatePropertiesForSystem(systems[0])
		props2 = evaluatePropertiesForSystem(systems[1])

		diffArg1 = "system " + systems[0]
		diffArg2 = "system " + systems[1]
	} else {
		expandFileList()

		diffArg1 = propertyFiles[0]
		diffArg2 = propertyFiles[1]

		file1, err1 := os.Open(diffArg1)
		if err1 == nil {
			defer file1.Close()
		}

		file2, err2 := os.Open(diffArg2)
		if err2 == nil {
			defer file2.Close()
		}

		if err1 != nil || err2 != nil {
			showDiffHelp(filenameOk(diffArg1, err1), filenameOk(diffArg2, err2))
			return
		}

		props1 = readPropertiesFromFile(file1)
		props2 = readPropertiesFromFile(file2)
	}

	var diffs []string

	for key, value := range props1 {
		if value2, ok := props2[key]; ok {
			if value != value2 {
				diffs = append(diffs, fmt.Sprintf("≈ %v: [%v] ≈ [%v]", key, value, value2))
			}
		} else {
			diffs = append(diffs, fmt.Sprintf("- %v: [%v]", key, value))
		}
	}
	for key, value := range props2 {
		if _, ok := props1[key]; !ok {
			diffs = append(diffs, fmt.Sprintf("+ %v: [%v]", key, value))
		}
	}
	if len(diffs) < 1 {
		fmt.Printf("Found no difference in properties between %v and %v\n", diffArg1, diffArg2)
	} else {
		fmt.Printf("There are differences in properties between %v and %v\n", diffArg1, diffArg2)
		fmt.Printf("\n≈ means the value differ between %v and %v\n", diffArg1, diffArg2)
		fmt.Println("- This property has been found only in " + diffArg1)
		fmt.Println("+ This property has been found only in " + diffArg2 + "\n")
		for _, diff := range diffs {
			fmt.Println(diff)
		}
	}
}

// List Systems
func listSystems() {
	expandFileList()

	for _, filepath := range propertyFiles {
		if strings.HasSuffix(filepath, ".csv") {
			var systemnames = readSystemNames()
			for _, systemname := range systemnames {
				fmt.Println(systemname)
			}
			return
		}
	}
	logError("No system found")
}

// Common functions

func expandFileList() {
	var filelist []string
	if *fileListPtr != "" {
		filelist = strings.Split(*fileListPtr, ",")
	} else {
		filelist = *readFileList()
	}
	for _, filename := range filelist {
		var filenameLower = strings.ToLower(filename)
		if !(strings.HasSuffix(filenameLower, ".properties") || strings.HasSuffix(filenameLower, ".csv")) {
			logError("Only .properties and .csv files are supported - file is: " + filename)
		}
		var filePath string
		if strings.Contains(filename, "/") {
			filePath = filename
		} else {
			filePath = configPath + filename
		}
		if _, err := os.Stat(filePath); err != nil {
			logError("File " + filePath + " not found:\n  " + err.Error())
		}
		propertyFiles = append(propertyFiles, filePath)
	}

	var mylocalFilePath string = configPath + mylocalPropertiesFile
	if _, err := os.Stat(mylocalFilePath); err == nil {
		propertyFiles = append(propertyFiles, mylocalFilePath)
	}
}

func filenameOk(filename string, err error) string {
	if err == nil {
		return "<ok>"
	}
	return filename
}

func evaluatePropertiesForSystem(systemName string) map[string]string {
	system = systemName
	evaluateAllPropertyFiles()
	var props = make(map[string]string)
	for key, prop := range properties {
		props[key] = prop.value
	}
	return props
}

func readPropertiesFromFile(file *os.File) map[string]string {

	var props = make(map[string]string)

	scanner := bufio.NewScanner(file)
	for lineNr := 1; scanner.Scan(); lineNr++ {
		line := scanner.Text()
		if strings.TrimSpace(line) != "" && !strings.HasPrefix(line, "#") && strings.Contains(line, "=") {
			parts := strings.Split(line, "=")
			props[parts[0]] = strings.Join(parts[1:], "=")
		}
	}
	return props
}

func evaluateAllPropertyFiles() {
	expandFileList()

	var systems = readSystemNames()

	var columnNr = detectSystemColumn(systems, system)
	columnNr++

	properties = make(map[string]property)

	addHeader(system)

	for _, propfile := range propertyFiles {
		logInfo(fmt.Sprintf("evaluating property file %v \n", propfile))
		evaluatePropertiesFromFile(propfile)
	}
}

func isCsvFile(filepath string) bool {
	return strings.HasSuffix(strings.ToLower(filepath), ".csv")
}

func readSystemNames() []string {

	for _, filepath := range propertyFiles {
		if isCsvFile(filepath) {
			file, err := os.Open(filepath)
			if err != nil {
				logError(err.Error())
			}
			defer file.Close()

			csv := csv.NewReader(file)
			csv.Comma = ';'
			csv.FieldsPerRecord = -1

			systems, err := csv.Read()

			if err != nil {
				logError(err.Error())
			}
			return systems[1:]
		}
	}
	return make([]string, 0)
}

func detectSystemColumn(systems []string, system string) int {
	if len(systems) > 0 && system == "" {
		// generelle Hilfe
		showGenerateHelp(systems, "")
		os.Exit(1)
	}
	for i, s := range systems {
		if strings.EqualFold(s, system) {
			return i
		}
	}
	// generelle Hilfe
	showGenerateHelp(systems, system)
	os.Exit(1)
	return -1
}

func printResultLine(outfile *os.File, line string) {
	if outfile == nil {
		fmt.Println(line)
	} else {
		outfile.WriteString(line + "\n")
	}
}

func createOutput() {
	var outfile *os.File = nil
	if output != "<console>" {
		var err error
		outfile, err = os.Create(configPath + output)
		check(err)
		defer outfile.Close()
	}

	for _, propertyKey := range orderedProperties {
		if strings.TrimSpace(propertyKey) == "" || strings.HasPrefix(propertyKey, "#") {
			printResultLine(outfile, propertyKey)
			continue
		}
		property := properties[propertyKey]
		if strings.EqualFold(">remove<", property.value) {
			printResultLine(outfile, propRemoveString(propertyKey, property))
		} else {
			if property.overriddenBy != "" {
				printResultLine(outfile, fmt.Sprintf(
					"# + defined in %v:%v - overwritten with %v:%v",
					property.firstSeen, property.firstSeenLine, property.overriddenBy, property.overriddenByLine))
			}
			printResultLine(outfile, propertyKey+"="+property.value)
		}
	}

	if outfile != nil {
		outfile.Sync()
	}
}

func propRemoveString(propertyKey string, property property) string {
	template := "# - removed [%v] in %v:%v\n"
	if property.overriddenBy == "" {
		return fmt.Sprintf(template, propertyKey, property.firstSeen, strconv.Itoa(property.firstSeenLine))
	}
	return fmt.Sprintf(template, propertyKey, property.overriddenBy, strconv.Itoa(property.overriddenByLine))
}

func addHeader(systemName string) {
	addToList("#")
	if systemName == "" {
		addToList("# generated local.properties at " + time.Now().Format("2006-01-02 15:04:05"))
	} else {
		addToList("# generated local.properties for system '" + systemName + "' at " + time.Now().Format("2006-01-02 15:04:05"))
	}
	addToList("#")
}

func addFileHeader(filename string) {
	addToList("")
	addToList("#")
	addToList("# --- Properties of file " + filename)
}

func addFileFooter(filename string) {
	addToList("# --- End of properties of file " + filename)
}

func evaluatePropertiesFromFile(filepath string) {
	addFileHeader(filepath)

	if isCsvFile(filepath) {
		evaluatePropsForSystem(filepath)
	} else {
		evaluatePropertiesFile(filepath)
	}

	addFileFooter(filepath)
}

func evaluatePropertiesFile(filepath string) {

	file, err := os.Open(filepath)
	if err != nil {
		logError(err.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for lineNr := 1; scanner.Scan(); lineNr++ {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
			addToList(line)
		} else {
			parts := strings.Split(line, "=")
			if len(parts) == 1 {
				logWarn(fmt.Sprintf("line %v:%v converting to comment: [%v] \n", filepath, lineNr, parts[0]))
				addToList("# +++ illegal line " + filepath + ":" + strconv.Itoa(lineNr) + " automatic converted to comment:")
				addToList("# +++   " + parts[0])
			} else {
				addProperty(filepath, lineNr, parts[0], strings.Join(parts[1:], "="))
			}
		}

	}
}

func evaluatePropsForSystem(filepath string) {

	file, err := os.Open(filepath)
	if err != nil {
		logError(err.Error())
	}
	defer file.Close()

	csv := csv.NewReader(file)
	csv.Comma = ';'
	csv.FieldsPerRecord = -1

	column := detectColumnNumberForSystem(csv, system)

	fields, err := csv.ReadAll()

	if err != nil {
		logError(err.Error())
	}

	for i := 0; i < len(fields); i++ {
		key := fields[i][0]
		if strings.TrimSpace(key) == "" || strings.HasPrefix(key, "#") {
			addToList(key)
		} else {
			value := findSystemValue(fields[i], column)
			addProperty(filepath, i+1, key, value)
		}
	}
}

func findSystemValue(values []string, column int) string {
	if column > len(values)-1 {
		return lookForParentValue(values, column-1)
	}
	value := strings.TrimSpace(values[column])
	if value == "*=" {
		return strings.TrimSpace(values[column-1])
	}
	if value == "*>" {
		return strings.TrimSpace(values[column-1])
	}
	if value == "" {
		return lookForParentValue(values, column-1)
	}
	return value
}

func lookForParentValue(values []string, column int) string {
	if column == 1 {
		return ""
	}
	if column > len(values)-1 {
		return lookForParentValue(values, column-1)
	}
	value := strings.TrimSpace(values[column])
	if value == "*>" {
		return strings.TrimSpace(values[column-1])
	}
	if value == "" {
		return lookForParentValue(values, column-1)
	}
	return ""
}

func detectColumnNumberForSystem(csv *csv.Reader, system string) int {

	systems, err := csv.Read()

	if err != nil {
		logError(err.Error())
	}

	for i, s := range systems {
		if s == system {
			return i
		}
	}
	logError(fmt.Sprintf("System '%v' not defined in properties csv file", system))
	return -1
}

func readFileList() *[]string {
	file, err := os.Open(configPath + filesListFile)
	if err != nil {
		logError(err.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var filelist []string
	for scanner.Scan() {
		filelist = append(filelist, scanner.Text())

		if err := scanner.Err(); err != nil {
			logError(err.Error())
		}
	}

	return &filelist
}

func addToList(value string) {
	orderedProperties = append(orderedProperties, value)
}

func addProperty(propFilename string, lineNr int, key string, value string) {
	_, existing := properties[key]
	if !existing {
		addToList(key)
		properties[key] = property{value, propFilename, lineNr, "", 0}
	} else {
		var prop = properties[key]
		prop.value = value
		prop.overriddenBy = propFilename
		prop.overriddenByLine = lineNr
		properties[key] = prop
	}
}

func check(e error) {
	if e != nil {
		logError(e.Error())
	}
}

func logError(msg string) {
	fmt.Println("ERROR - " + msg)
	os.Exit(1)
}

func logWarn(msg string) {
	if verbose {
		fmt.Println("WARN - " + msg)
	}
}

func logInfo(msg string) {
	if verbose {
		fmt.Println("Info - " + msg)
	}
}
