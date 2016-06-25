package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"

	"golang.org/x/net/context"

	"github.com/romainmenke/universal-notifier/pkg/wercker"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

func main() {

	fmt.Print("Starting Hearts")

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return
	}

	fmt.Printf("Listening on port : %s" port)

	s := grpc.NewServer()
	wercker.RegisterNotificationServiceServer(s, &server{})
	s.Serve(lis)

	fmt.Print("Ready to serve clients")

}

// server is used to implement helloworld.GreeterServer.
type server struct{}

func (s *server) Notify(ctx context.Context, in *wercker.WerckerMessage) (*wercker.WerckerResponse, error) {
	fmt.Print(in)
	return &wercker.WerckerResponse{Success: true}, nil
}

func setupUser() error {
	user := "heartsbot"

	user = fmt.Sprintf("'%s'", user)

	cmd := exec.Command("git", "config", "user.name", user)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}
	fmt.Print(out.String())
	return nil

}

func setupEmail() error {
	email := "romainmenke+heartsbot@gmail.com"

	email = fmt.Sprintf("'%s'", email)

	cmd := exec.Command("git", "config", "user.email", email)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}
	fmt.Print(out.String())
	return nil
}

func add() error {

	cmd := exec.Command("git", "add", ".")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}
	fmt.Print(out.String())
	return nil
}

func commit() error {

	cmd := exec.Command("git", "commit", "-m", "'entry'")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}
	fmt.Print(out.String())
	return nil
}

func push() error {

	user := os.Getenv("USER")
	password := os.Getenv("PASS")
	url := fmt.Sprintf("https://%s:%s@github.com/romainmenke/githeartsdb.git", user, password)

	cmd := exec.Command("git", "push", url)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}
	fmt.Print(out.String())
	return nil
}

func write(path string) error {

	dirPath := fmt.Sprintf("db/%s", path)
	filePath := fmt.Sprintf("db/%sshield.svg", path)

	err := createDir(dirPath)
	if err != nil {
		return err
	}

	err = removeFile(filePath)
	if err != nil {
		return err
	}

	err = createFile(filePath)
	if err != nil {
		return err
	}

	d1 := []byte(svg(2))
	err = ioutil.WriteFile(filePath, d1, 0644)
	if err != nil {
		return err
	}
	return nil
}

func createDir(path string) error {
	cmd := exec.Command("mkdir", "-p", path)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func createFile(path string) error {
	cmd := exec.Command("touch", path)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func removeFile(path string) error {
	cmd := exec.Command("rm", "-f", path)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func svg(hearts int) string {
	switch hearts {
	case 1:
		return svgOne()
	case 2:
		return svgTwo()
	case 3:
		return svgThree()
	default:
		return svgThree()
	}
}

func svgOne() string {
	return "<svg width=\"30\" height=\"20\" xmlns=\"http://www.w3.org/2000/svg\"><!-- Created with Method Draw - http://github.com/duopixel/Method-Draw/ --><g><title>background</title><rect fill=\"#fff\" id=\"canvas_background\" height=\"22\" width=\"32\" y=\"-1\" x=\"-1\"/><g display=\"none\" overflow=\"visible\" y=\"0\" x=\"0\" height=\"100%\" width=\"100%\" id=\"canvasGrid\"><rect fill=\"url(#gridpattern)\" stroke-width=\"0\" y=\"0\" x=\"0\" height=\"100%\" width=\"100%\"/></g></g><g><title>Layer 1</title><rect id=\"svg_1\" height=\"147\" width=\"432\" y=\"96\" x=\"77\" stroke-width=\"1.5\" stroke=\"#000\" fill=\"#fff\"/><text xml:space=\"preserve\" text-anchor=\"start\" font-family=\"Helvetica, Arial, sans-serif\" font-size=\"24\" id=\"svg_2\" y=\"188\" x=\"266.5\" stroke-width=\"0\" stroke=\"#000\" fill=\"#000000\">♥♥♥</text><text transform=\"matrix(1 0 0 1 0 0)\" stroke=\"#000\" xml:space=\"preserve\" text-anchor=\"start\" font-family=\"Helvetica, Arial, sans-serif\" font-size=\"15\" id=\"svg_3\" y=\"15.16415\" x=\"10.56461\" stroke-width=\"0\" fill=\"#ff002e\">♥</text></g></svg>"
}

func svgTwo() string {
	return "<svg width=\"30\" height=\"20\" xmlns=\"http://www.w3.org/2000/svg\"><!-- Created with Method Draw - http://github.com/duopixel/Method-Draw/ --><g><title>background</title><rect fill=\"#fff\" id=\"canvas_background\" height=\"22\" width=\"32\" y=\"-1\" x=\"-1\"/><g display=\"none\" overflow=\"visible\" y=\"0\" x=\"0\" height=\"100%\" width=\"100%\" id=\"canvasGrid\"><rect fill=\"url(#gridpattern)\" stroke-width=\"0\" y=\"0\" x=\"0\" height=\"100%\" width=\"100%\"/></g></g><g><title>Layer 1</title><rect id=\"svg_1\" height=\"147\" width=\"432\" y=\"96\" x=\"77\" stroke-width=\"1.5\" stroke=\"#000\" fill=\"#fff\"/><text xml:space=\"preserve\" text-anchor=\"start\" font-family=\"Helvetica, Arial, sans-serif\" font-size=\"24\" id=\"svg_2\" y=\"188\" x=\"266.5\" stroke-width=\"0\" stroke=\"#000\" fill=\"#000000\">♥♥♥</text><text transform=\"matrix(1 0 0 1 0 0)\" stroke=\"#000\" xml:space=\"preserve\" text-anchor=\"start\" font-family=\"Helvetica, Arial, sans-serif\" font-size=\"15\" id=\"svg_3\" y=\"15.16415\" x=\"6.12711\" stroke-width=\"0\" fill=\"#ff002e\">♥♥</text></g></svg>"
}

func svgThree() string {
	return "<svg width=\"30\" height=\"20\" xmlns=\"http://www.w3.org/2000/svg\"><!-- Created with Method Draw - http://github.com/duopixel/Method-Draw/ --><g><title>background</title><rect fill=\"#fff\" id=\"canvas_background\" height=\"22\" width=\"32\" y=\"-1\" x=\"-1\"/><g display=\"none\" overflow=\"visible\" y=\"0\" x=\"0\" height=\"100%\" width=\"100%\" id=\"canvasGrid\"><rect fill=\"url(#gridpattern)\" stroke-width=\"0\" y=\"0\" x=\"0\" height=\"100%\" width=\"100%\"/></g></g><g><title>Layer 1</title><rect id=\"svg_1\" height=\"147\" width=\"432\" y=\"96\" x=\"77\" stroke-width=\"1.5\" stroke=\"#000\" fill=\"#fff\"/><text xml:space=\"preserve\" text-anchor=\"start\" font-family=\"Helvetica, Arial, sans-serif\" font-size=\"24\" id=\"svg_2\" y=\"188\" x=\"266.5\" stroke-width=\"0\" stroke=\"#000\" fill=\"#000000\">♥♥♥</text><text stroke=\"#000\" xml:space=\"preserve\" text-anchor=\"start\" font-family=\"Helvetica, Arial, sans-serif\" font-size=\"15\" id=\"svg_3\" y=\"15.16579\" x=\"1.69757\" stroke-width=\"0\" fill=\"#ff002e\">♥♥♥</text></g></svg>"
}
