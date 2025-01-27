package utils

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/prachin77/pkr/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

// ANSI escape codes for background colors
const (
	ColorReset      = "\033[0m"
	BackgroundGreen = "\033[42m" // Success
	BackgroundRed   = "\033[41m" // Error
	BackgroundBlue  = "\033[44m" // Method
)

func StructuredLoggerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Log start time and method
		start := time.Now()
		method := info.FullMethod

		// Get client IP address
		p, ok := peer.FromContext(ctx)
		var clientIP string
		if ok {
			clientIP, _, _ = net.SplitHostPort(p.Addr.String())
		} else {
			clientIP = "unknown"
		}

		// Handle the request
		resp, err := handler(ctx, req)

		// Log response status and duration with background color coding
		duration := time.Since(start).Milliseconds()
		status := "Success"            // Default status
		statusColor := BackgroundGreen // Default color
		if err != nil {
			status = "Failed "
			statusColor = BackgroundRed
		}

		// Format log output
		methodField := fmt.Sprintf("%-40s", method)     // Align method column
		statusField := fmt.Sprintf("%-7s", status)      // Align status column
		durationField := fmt.Sprintf("%4dms", duration) // Align duration column

		log.Printf("%s %s %s | %s | %s |%s %s %s",
			statusColor,
			statusField,
			ColorReset,
			durationField,
			clientIP,
			BackgroundBlue,
			methodField,
			ColorReset,
		)

		return resp, err
	}
}

func ClearScreen() {
	var cmd *exec.Cmd

	// Check the operating system to determine the appropriate clear command
	switch runtime.GOOS {
	case "linux", "darwin":
		cmd = exec.Command("clear") // for Unix-like systems
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls") // for Windows
	default:
		fmt.Println("Unsupported platform.")
		return
	}

	// Execute the clear command
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func GetHostPublicKey() (string, string, error) {
	public_key_data, err := os.ReadFile(models.PUBLIC_KEY_FILE)
	if err != nil {
		return "", "", err
	}

	path, err := filepath.Abs(models.PUBLIC_KEY_FILE)
	if err != nil {
		fmt.Println("error retrieving host public key file path : ", err)
		return "", "", err
	}
	return string(public_key_data), path, nil
}

func GetHostPrivateKeys() (string, string, error) {
	private_key_data, err := os.ReadFile(models.PRIVATE_KEY_FILE)
	if err != nil {
		return "", "", err
	}

	path, err := filepath.Abs(models.PRIVATE_KEY_FILE)
	if err != nil {
		fmt.Println("error retrieving host public key file path : ", err)
		return "", "", err
	}

	return string(private_key_data), path, nil
}

// ParseBytesToPublicKey converts a PEM-encoded public key (as []byte) into an *rsa.PublicKey
func ParseBytesToPublicKey(publicKeyBytes []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(publicKeyBytes)
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return nil, errors.New("failed to decode PEM block containing public key")
	}

	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSA public key: %v", err)
	}

	return publicKey, nil
}

