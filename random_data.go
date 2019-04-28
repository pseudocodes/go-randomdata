// Package randomdata implements a bunch of simple ways to generate (pseudo) random data
package randomdata

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const (
	Male         int = 0
	Female       int = 1
	RandomGender int = 2
)

const (
	Small int = 0
	Large int = 1
)

const (
	FullCountry      = 0
	TwoCharCountry   = 1
	ThreeCharCountry = 2
)

const (
	DateInputLayout  = "2006-01-02"
	DateOutputLayout = "Monday 2 Jan 2006"
)

type jsonContent struct {
	Adjectives          []string `json:"adjectives"`
	Nouns               []string `json:"nouns"`
	FirstNamesFemale    []string `json:"firstNamesFemale"`
	FirstNamesMale      []string `json:"firstNamesMale"`
	LastNames           []string `json:"lastNames"`
	Domains             []string `json:"domains"`
	People              []string `json:"people"`
	StreetTypes         []string `json:"streetTypes"` // Taken from https://github.com/tomharris/random_data/blob/master/lib/random_data/locations.rb
	Paragraphs          []string `json:"paragraphs"`  // Taken from feedbooks.com and www.gutenberg.org
	Countries           []string `json:"countries"`   // Fetched from the world bank at http://siteresources.worldbank.org/DATASTATISTICS/Resources/CLASS.XLS
	CountriesThreeChars []string `json:"countriesThreeChars"`
	CountriesTwoChars   []string `json:"countriesTwoChars"`
	Currencies          []string `json:"currencies"` //https://github.com/OpenBookPrices/country-data
	Cities              []string `json:"cities"`
	States              []string `json:"states"`
	StatesSmall         []string `json:"statesSmall"`
	Days                []string `json:"days"`
	Months              []string `json:"months"`
	FemaleTitles        []string `json:"femaleTitles"`
	MaleTitles          []string `json:"maleTitles"`
	Timezones           []string `json:"timezones"`           // https://en.wikipedia.org/wiki/List_of_tz_database_time_zones
	Locales             []string `json:"locales"`             // https://tools.ietf.org/html/bcp47
	UserAgents          []string `json:"userAgents"`          // http://techpatterns.com/downloads/firefox/useragentswitcher.xml
	CountryCallingCodes []string `json:"countryCallingCodes"` // from https://github.com/datasets/country-codes/blob/master/data/country-codes.csv
	ProvincesGB         []string `json:"provincesGB"`
	StreetNameGB        []string `json:"streetNameGB"`
	StreetTypesGB       []string `json:"streetTypesGB"`
}

type RandData struct {
	*rand.Rand
}

var jsonData = jsonContent{}
var privateRand *RandData

func init() {
	// privateRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	privateRand = NewRandData()
	jsonData = jsonContent{}

	err := json.Unmarshal(data, &jsonData)

	if err != nil {
		log.Fatal(err)
	}
}

func NewRandData() *RandData {
	rd := &RandData{
		rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	return rd
}

func (rd *RandData) CustomRand(randToUse *rand.Rand) {
	rd.Rand = randToUse
}

// Returns a random part of a slice
func (rd RandData) randomFrom(source []string) string {
	return source[rd.Intn(len(source))]
}

// Title returns a random title, gender decides the gender of the name
func (rd RandData) Title(gender int) string {
	var title = ""
	switch gender {
	case Male:
		title = rd.randomFrom(jsonData.MaleTitles)
		break
	case Female:
		title = rd.randomFrom(jsonData.FemaleTitles)
		break
	default:
		title = rd.FirstName(rd.Intn(2))
		break
	}
	return title
}

// FirstName returns a random first name, gender decides the gender of the name
func (rd RandData) FirstName(gender int) string {
	var name = ""
	switch gender {
	case Male:
		name = rd.randomFrom(jsonData.FirstNamesMale)
		break
	case Female:
		name = rd.randomFrom(jsonData.FirstNamesFemale)
		break
	default:
		name = rd.FirstName(rand.Intn(2))
		break
	}
	return name
}

// LastName returns a random last name
func (rd RandData) LastName() string {
	return rd.randomFrom(jsonData.LastNames)
}

// FullName returns a combination of FirstName LastName randomized, gender decides the gender of the name
func (rd RandData) FullName(gender int) string {
	return rd.FirstName(gender) + " " + rd.LastName()
}

// Email returns a random email
func (rd RandData) Email() string {
	return strings.ToLower(rd.FirstName(RandomGender)+rd.LastName()) + rd.StringNumberExt(1, "", 3) + "@" + rd.randomFrom(jsonData.Domains)
}

// Country returns a random country, countryStyle decides what kind of format the returned country will have
func (rd RandData) Country(countryStyle int64) string {
	country := ""
	switch countryStyle {

	default:

	case FullCountry:
		country = rd.randomFrom(jsonData.Countries)
		break

	case TwoCharCountry:
		country = rd.randomFrom(jsonData.CountriesTwoChars)
		break

	case ThreeCharCountry:
		country = rd.randomFrom(jsonData.CountriesThreeChars)
		break
	}
	return country
}

// Currency returns a random currency under ISO 4217 format
func (rd RandData) Currency() string {
	return rd.randomFrom(jsonData.Currencies)
}

// City returns a random city
func (rd RandData) City() string {
	return rd.randomFrom(jsonData.Cities)
}

// ProvinceForCountry returns a randomly selected province (state, county,subdivision ) name for a supplied country.
// If the country is not supported it will return an empty string.
func (rd RandData) ProvinceForCountry(countrycode string) string {
	switch countrycode {
	case "US":
		return rd.randomFrom(jsonData.States)
	case "GB":
		return rd.randomFrom(jsonData.ProvincesGB)
	}
	return ""
}

// State returns a random american state
func (rd RandData) State(typeOfState int) string {
	if typeOfState == Small {
		return rd.randomFrom(jsonData.StatesSmall)
	}
	return rd.randomFrom(jsonData.States)
}

// Street returns a random fake street name
func (rd RandData) Street() string {
	return fmt.Sprintf("%s %s", rd.randomFrom(jsonData.People), rd.randomFrom(jsonData.StreetTypes))
}

// StreetForCountry returns a random fake street name typical to the supplied country.
// If the country is not supported it will return an empty string.
func (rd RandData) StreetForCountry(countrycode string) string {
	switch countrycode {
	case "US":
		return rd.Street()
	case "GB":
		return fmt.Sprintf("%s %s", rd.randomFrom(jsonData.StreetNameGB), rd.randomFrom(jsonData.StreetTypesGB))
	}
	return ""
}

// Address returns an american style address
func (rd RandData) Address() string {
	return fmt.Sprintf("%d %s,\n%s, %s, %s", rd.Number(100), rd.Street(), rd.City(), rd.State(Small), PostalCode("US"))
}

// Paragraph returns a random paragraph
func (rd RandData) Paragraph() string {
	return rd.randomFrom(jsonData.Paragraphs)
}

// Number returns a random number, if only one integer (n1) is supplied it returns a number in [0,n1)
// if a second argument is supplied it returns a number in [n1,n2)
func (rd RandData) Number(numberRange ...int) int {
	nr := 0
	if len(numberRange) > 1 {
		nr = 1
		nr = rd.Intn(numberRange[1]-numberRange[0]) + numberRange[0]
	} else {
		nr = rd.Intn(numberRange[0])
	}
	return nr
}

func (rd RandData) Decimal(numberRange ...int) float64 {
	nr := 0.0
	if len(numberRange) > 1 {
		nr = 1.0
		nr = rd.Float64()*(float64(numberRange[1])-float64(numberRange[0])) + float64(numberRange[0])
	} else {
		nr = rd.Float64() * float64(numberRange[0])
	}

	if len(numberRange) > 2 {
		sf := strconv.FormatFloat(nr, 'f', numberRange[2], 64)
		nr, _ = strconv.ParseFloat(sf, 64)
	}
	return nr
}

func (rd RandData) StringNumberExt(numberPairs int, separator string, numberOfDigits int) string {
	numberString := ""

	for i := 0; i < numberPairs; i++ {
		for d := 0; d < numberOfDigits; d++ {
			numberString += fmt.Sprintf("%d", rd.Number(0, 9))
		}

		if i+1 != numberPairs {
			numberString += separator
		}
	}

	return numberString
}

// StringNumber returns a random number as a string
func (rd RandData) StringNumber(numberPairs int, separator string) string {
	return rd.StringNumberExt(numberPairs, separator, 2)
}

// StringSample returns a random string from a list of strings
func (rd RandData) StringSample(stringList ...string) string {
	str := ""
	if len(stringList) > 0 {
		str = stringList[rd.Number(0, len(stringList))]
	}
	return str
}

func (rd RandData) Boolean() bool {
	nr := rd.Intn(2)
	return nr != 0
}

// Noun returns a random noun
func (rd RandData) Noun() string {
	return rd.randomFrom(jsonData.Nouns)
}

// Adjective returns a random adjective
func (rd RandData) Adjective() string {
	return rd.randomFrom(jsonData.Adjectives)
}

func (rd RandData) uppercaseFirstLetter(word string) string {
	a := []rune(word)
	a[0] = unicode.ToUpper(a[0])
	return string(a)
}

func (rd RandData) lowercaseFirstLetter(word string) string {
	a := []rune(word)
	a[0] = unicode.ToLower(a[0])
	return string(a)
}

// SillyName returns a silly name, useful for randomizing naming of things
func (rd RandData) SillyName() string {
	return rd.uppercaseFirstLetter(rd.Noun()) + rd.Adjective()
}

// IpV4Address returns a valid IPv4 address as string
func (rd RandData) IpV4Address() string {
	blocks := []string{}
	for i := 0; i < 4; i++ {
		number := rd.Intn(255)
		blocks = append(blocks, strconv.Itoa(number))
	}

	return strings.Join(blocks, ".")
}

// IpV6Address returns a valid IPv6 address as net.IP
func (rd RandData) IpV6Address() string {
	var ip net.IP
	for i := 0; i < net.IPv6len; i++ {
		number := uint8(rd.Intn(255))
		ip = append(ip, number)
	}
	return ip.String()
}

// MacAddress returns an mac address string
func (rd RandData) MacAddress() string {
	blocks := []string{}
	for i := 0; i < 6; i++ {
		number := fmt.Sprintf("%02x", rd.Intn(255))
		blocks = append(blocks, number)
	}

	return strings.Join(blocks, ":")
}

// Day returns random day
func (rd RandData) Day() string {
	return rd.randomFrom(jsonData.Days)
}

// Month returns random month
func (rd RandData) Month() string {
	return rd.randomFrom(jsonData.Months)
}

// FullDate returns full date
func (rd RandData) FullDate() string {
	timestamp := time.Now()
	year := timestamp.Year()
	month := rd.Number(1, 13)
	maxDay := time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC).Day()
	day := rd.Number(1, maxDay+1)
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return date.Format(DateOutputLayout)
}

// FullDateInRange returns a date string within a given range, given in the format "2006-01-02".
// If no argument is supplied it will return the result of randomdata.FullDate().
// If only one argument is supplied it is treated as the max date to return.
// If a second argument is supplied it returns a date between (and including) the two dates.
// Returned date is in format "Monday 2 Jan 2006".
func (rd RandData) FullDateInRange(dateRange ...string) string {
	var (
		min        time.Time
		max        time.Time
		duration   int
		dateString string
	)
	if len(dateRange) == 1 {
		max, _ = time.Parse(DateInputLayout, dateRange[0])
	} else if len(dateRange) == 2 {
		min, _ = time.Parse(DateInputLayout, dateRange[0])
		max, _ = time.Parse(DateInputLayout, dateRange[1])
	}
	if !max.IsZero() && max.After(min) {
		duration = rd.Number(int(max.Sub(min))) * -1
		dateString = max.Add(time.Duration(duration)).Format(DateOutputLayout)
	} else if !max.IsZero() && !max.After(min) {
		dateString = max.Format(DateOutputLayout)
	} else {
		dateString = rd.FullDate()
	}
	return dateString
}

func (rd RandData) Timezone() string {
	return rd.randomFrom(jsonData.Timezones)
}

func (rd RandData) Locale() string {
	return rd.randomFrom(jsonData.Locales)
}

func (rd RandData) UserAgentString() string {
	return rd.randomFrom(jsonData.UserAgents)
}

func (rd RandData) PhoneNumber() string {
	str := rd.randomFrom(jsonData.CountryCallingCodes) + " "

	str += Digits(rd.Intn(3) + 1)

	for {
		// max 15 chars
		remaining := 15 - (len(str) - strings.Count(str, " "))
		if remaining < 2 {
			return "+" + str
		}
		str += " " + Digits(rd.Intn(remaining-1)+1)
	}
}
