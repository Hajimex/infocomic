package main

import (	
	"encoding/json"
	"log"
	// "fmt"
	"os"
	// "path/filepath"

	"github.com/fogleman/gg"
	"github.com/urfave/cli"

	"time"
	"io"
	"io/ioutil"
	
	"net/http"

	"bytes"
	"github.com/aws/aws-sdk-go/aws" 
	// "github.com/aws/aws-sdk-go/aws/awsutil" 
	// "github.com/aws/aws-sdk-go/aws/credentials" 
	"github.com/aws/aws-sdk-go/service/s3" 
	"github.com/aws/aws-sdk-go/aws/session" 
	"github.com/aws/aws-lambda-go/lambda"
)


// ConfigHolder defines the input for configuration
type ConfigHolder struct {
	FirstHeader string `json:"first_header,omitempty"`

	SecondHeader string `json:"second_header,omitempty"`

	LeftFirstContainer struct {
		First  string `json:"first,omitempty"`
		Second string `json:"second,omitempty"`
		Third  string `json:"third,omitempty"`
	} `json:"left_first_container,omitempty"`

	LeftSecondContainer struct {
		First  string `json:"first,omitempty"`
		Second string `json:"second,omitempty"`
		Third  string `json:"third,omitempty"`
	} `json:"left_second_container,omitempty"`

	LeftThirdContainer struct {
		First  string `json:"first,omitempty"`
		Second string `json:"second,omitempty"`
		Third  string `json:"third,omitempty"`
		Forth  string `json:"forth,omitempty"`
	} `json:"left_third_container,omitempty"`

	LeftForthContainer struct {
		First  string `json:"first,omitempty"`
		Second string `json:"second,omitempty"`
		Third  string `json:"third,omitempty"`
		Forth  string `json:"forth,omitempty"`
	} `json:"left_forth_container,omitempty"`

	RightFirstContainer struct {
		First  string `json:"first,omitempty"`
		Second string `json:"second,omitempty"`
		Third  string `json:"third,omitempty"`
	} `json:"right_first_container,omitempty"`

	RightSecondContainer struct {
		First  string `json:"first,omitempty"`
		Second string `json:"second,omitempty"`
		Third  string `json:"third,omitempty"`
		Forth  string `json:"forth,omitempty"`
	} `json:"right_second_container,omitempty"`

	RightThirdContainer struct {
		First  string `json:"first,omitempty"`
		Second string `json:"second,omitempty"`
		Third  string `json:"third,omitempty"`
	} `json:"right_third_container,omitempty"`

	RightForthContainer struct {
		First  string `json:"first,omitempty"`
		Second string `json:"second,omitempty"`
		Third  string `json:"third,omitempty"`
	} `json:"right_forth_container,omitempty"`
}

type MyEvent struct {
  Input string `json:"Input"`
  Config string `json:"Config"`
}

type coordinate struct {
	x     float64
	y     float64
	width float64
}

// Config holds the value of the configuration loaded from the configuration file
var Config *ConfigHolder = &ConfigHolder{}

// Coordinates defines the coordinates of main points/headers
var Coordinates = map[string]coordinate{
	"FirstHeader":  {98, 22, 593},
	"SecondHeader": {724, 24, 1200},

	"LeftFirstCFirst":  {124, 131, 572},
	"LeftFirstCSecond": {51, 248, 192},
	"LeftFirstCThird":  {434, 265, 582},

	"LeftSecondCFirst":  {147, 547, 581},
	"LeftSecondCSecond": {63, 664, 199},
	"LeftSecondCThird":  {414, 669, 546},

	"LeftThirdCFirst":  {118, 964, 361},
	"LeftThirdCSecond": {389, 978, 560},
	"LeftThirdCThird":  {105, 1059, 214},
	"LeftThirdCForth":  {399, 1195, 559},

	"LeftForthCFirst":  {140, 1380, 457},
	"LeftForthCSecond": {475, 1435, 572},
	"LeftForthCThird":  {71, 1474, 239},
	"LeftForthCForth":  {456, 1588, 571},

	"RightFirstCFirst":  {716, 150, 1178},
	"RightFirstCSecond": {655, 285, 762},
	"RightFirstCThird":  {1039, 306, 1178},

	"RightSecondCFirst":  {729, 561, 1177},
	"RightSecondCSecond": {788, 667, 1009},
	"RightSecondCThird":  {657, 780, 773},
	"RightSecondCForth":  {1032, 754, 1168},

	"RightThirdCFirst":  {733, 975, 1180},
	"RightThirdCSecond": {885, 1044, 1098},
	"RightThirdCThird":  {818, 1167, 992},

	"RightForthCFirst":  {766, 1369, 948},
	"RightForthCSecond": {811, 1541, 904},
	"RightForthCThird":  {1021, 1582, 1167},
}

const (
    S3_REGION = "ap-northeast-1"
    S3_BUCKET = "kuwayama"
)

// Args
var appArgs = struct {
	ImageInput  string
	ImageOutput string
	Font        string
	Config      string
}{}

func AddFileToS3(s *session.Session, fileDir string) error {

    // Open the file for use
    file, err := os.Open(fileDir)
    if err != nil {
        return err
    }
    defer file.Close()

    // Get file size and read the file content into a buffer
    fileInfo, _ := file.Stat()
    var size int64 = fileInfo.Size()
    buffer := make([]byte, size)
    file.Read(buffer)

    // Config settings: this is where you choose the bucket, filename, content-type etc.
    // of the file you're uploading.
    _, err = s3.New(s).PutObject(&s3.PutObjectInput{
        Bucket:               aws.String(S3_BUCKET),
        Key:                  aws.String(fileDir),
        ACL:                  aws.String("public-read"),
        Body:                 bytes.NewReader(buffer),
        ContentLength:        aws.Int64(size),
        ContentType:          aws.String(http.DetectContentType(buffer)),
        ContentDisposition:   aws.String("attachment"),
        ServerSideEncryption: aws.String("AES256"),
    })
    return err
}

// main 
func infocomic(event MyEvent) {
	app := cli.NewApp()
	// APP NAME
	app.Name = "Solver - a dedicated script to automate drawing japanese scripts on images"
	log.Println("event.Input1: %s", event.Input)
	// APP FLAGS
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "input",
			Value:       event.Input,
			Destination: &appArgs.ImageInput,
			Usage:       "Input image: The source",
		},
		cli.StringFlag{
			Name:        "output",
			Value:       "output.png",
			Destination: &appArgs.ImageOutput,
			Usage:       "Output image: The destination after putting words",
		},
		cli.StringFlag{
			Name:        "config",
			Value:       event.Config,
			Destination: &appArgs.Config,
			Usage:       "The directory of configuration file",
		},
		cli.StringFlag{
			Name:        "font",
			Value:       "wqy-zenhei.ttf",
			Destination: &appArgs.Font,
			Usage:       "The directory of font source",
		},
	}

	app.Action = func(c *cli.Context) error {
		log.Println("Starting the solver script")
		Load(appArgs.Config)
		PutText()
		return nil
	}
	
	app.Run(os.Args)

	s, err := session.NewSession(&aws.Config{Region: aws.String(S3_REGION)})
    if err != nil {
        log.Fatal(err)
    }	
    str := "/tmp/" + appArgs.ImageOutput
	err = AddFileToS3(s, str)
	if err != nil {
        log.Fatal(err)
    }

}

// Load loads configuration from a file
func Load(filePath string) error {

	// confPath, _ := filepath.Abs(filePath)
	url := filePath
	spaceClient := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}
	log.Printf("Loading configuration file from %s \n", filePath)


	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "spacecount-tutorial")
	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	// f, err := os.Open(confPath)
	// if err != nil {
	// 	log.Println("Failed to open the configuration file")
	// 	return err
	// }

	// defer f.Close()

	err = json.Unmarshal([]byte(body), &Config)
	// decoder := json.NewDecoder(f)

	// err = decoder.Decode(&Config)
	if err != nil {
		log.Println("Failed to decode configuration file", err)
		return err
	}

	return nil
}

// PutText puts texts on a given image
func PutText() {
    url := appArgs.ImageInput
    // don't worry about errors
    response, e := http.Get(url)
    if e != nil {
        log.Fatal(e)
    }
    defer response.Body.Close()

    //open a file for writing
    file, err := os.Create("/tmp/asdf.jpg")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    // Use io.Copy to just dump the response body to the file. This supports huge files
    _, err = io.Copy(file, response.Body)
    if err != nil {
        log.Fatal(err)
    }
    
	log.Println("Loading input image file from", appArgs.ImageInput)

	im, err := gg.LoadImage("/tmp/asdf.jpg")
	if err != nil {
		log.Fatal(err)
	}
	dc := gg.NewContext(1240, 1796)	
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)

	dc.LoadFontFace(appArgs.Font, 40)

	dc.DrawImage(im, 0, 0)
	// TOP BOX 1
	AddLine(dc, Config.FirstHeader, Coordinates["FirstHeader"])
	log.Println("Config.FirstHeader: %s", Config.FirstHeader)

	// TOP BOX 2
	AddLine(dc, Config.SecondHeader, Coordinates["SecondHeader"])

	dc.LoadFontFace(appArgs.Font, 28)

	// (FOR LEFT COLUMN)
	// (CONTAINER 1)
	AddLine(dc, Config.LeftFirstContainer.First, Coordinates["LeftFirstCFirst"])
	AddLine(dc, Config.LeftFirstContainer.Second, Coordinates["LeftFirstCSecond"])
	AddLine(dc, Config.LeftFirstContainer.Third, Coordinates["LeftFirstCThird"])

	// (CONTAINER 2)
	AddLine(dc, Config.LeftSecondContainer.First, Coordinates["LeftSecondCFirst"])
	AddLine(dc, Config.LeftSecondContainer.Second, Coordinates["LeftSecondCSecond"])
	AddLine(dc, Config.LeftSecondContainer.Third, Coordinates["LeftSecondCThird"])

	// (CONTAINER 3)
	AddLine(dc, Config.LeftThirdContainer.First, Coordinates["LeftThirdCFirst"])
	AddLine(dc, Config.LeftThirdContainer.Second, Coordinates["LeftThirdCSecond"])
	AddLine(dc, Config.LeftThirdContainer.Third, Coordinates["LeftThirdCThird"])
	AddLine(dc, Config.LeftThirdContainer.Forth, Coordinates["LeftThirdCForth"])

	// (CONTAINER 4)
	AddLine(dc, Config.LeftForthContainer.First, Coordinates["LeftForthCFirst"])
	AddLine(dc, Config.LeftForthContainer.Second, Coordinates["LeftForthCSecond"])
	AddLine(dc, Config.LeftForthContainer.Third, Coordinates["LeftForthCThird"])
	AddLine(dc, Config.LeftForthContainer.Forth, Coordinates["LeftForthCForth"])

	// (FOR RIGHT COLUMN)
	// (CONTAINER 1)
	AddLine(dc, Config.RightFirstContainer.First, Coordinates["RightFirstCFirst"])
	AddLine(dc, Config.RightFirstContainer.Second, Coordinates["RightFirstCSecond"])
	AddLine(dc, Config.RightFirstContainer.Third, Coordinates["RightFirstCThird"])

	// (CONTAINER 2)
	AddLine(dc, Config.RightSecondContainer.First, Coordinates["RightSecondCFirst"])
	AddLine(dc, Config.RightSecondContainer.Second, Coordinates["RightSecondCSecond"])
	AddLine(dc, Config.RightSecondContainer.Third, Coordinates["RightSecondCThird"])
	AddLine(dc, Config.RightSecondContainer.Forth, Coordinates["RightSecondCForth"])

	// (CONTAINER 3)
	AddLine(dc, Config.RightThirdContainer.First, Coordinates["RightThirdCFirst"])
	AddLine(dc, Config.RightThirdContainer.Second, Coordinates["RightThirdCSecond"])
	AddLine(dc, Config.RightThirdContainer.Third, Coordinates["RightThirdCThird"])

	// (CONTAINER 4)
	AddLine(dc, Config.RightForthContainer.First, Coordinates["RightForthCFirst"])
	AddLine(dc, Config.RightForthContainer.Second, Coordinates["RightForthCSecond"])
	AddLine(dc, Config.RightForthContainer.Third, Coordinates["RightForthCThird"])

	dc.Clip()
	str := "/tmp/" + appArgs.ImageOutput
	dc.SavePNG(str)
	log.Println("Saving final image to ", str)
}

// AddLine adds a new line to gg context with certain coordinates
// and wrap the words
func AddLine(dc *gg.Context, input string, cord coordinate) {
	_, h := dc.MeasureString(input)
	stringsofar := ""
	heightsofar := cord.y + h
	for _, ri := range input {
		widthnow, _ := dc.MeasureString(stringsofar + string(ri))
		if widthnow+cord.x > cord.width {
			dc.DrawString(stringsofar, cord.x, heightsofar)
			heightsofar, stringsofar = heightsofar+30.0, ""
		}
		stringsofar += string(ri)
	}
	dc.DrawString(stringsofar, cord.x, heightsofar)
}
 
func main() {
    lambda.Start(infocomic)
}
