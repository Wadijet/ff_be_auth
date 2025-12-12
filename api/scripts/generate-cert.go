package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

func main() {
	// Đường dẫn file
	certDir := filepath.Join("..", "config", "tls")
	certFile := filepath.Join(certDir, "server.crt")
	keyFile := filepath.Join(certDir, "server.key")

	// Tạo thư mục nếu chưa tồn tại
	if err := os.MkdirAll(certDir, 0755); err != nil {
		fmt.Printf("Lỗi tạo thư mục: %v\n", err)
		os.Exit(1)
	}

	// Kiểm tra xem đã có certificate chưa
	if _, err := os.Stat(certFile); err == nil {
		if _, err := os.Stat(keyFile); err == nil {
			fmt.Printf("Certificate đã tồn tại tại:\n")
			fmt.Printf("  Cert: %s\n", certFile)
			fmt.Printf("  Key:  %s\n", keyFile)
			fmt.Printf("Xóa file cũ và tạo lại? (y/N): ")
			var answer string
			fmt.Scanln(&answer)
			if answer != "y" && answer != "Y" {
				fmt.Println("Bỏ qua tạo certificate mới.")
				os.Exit(0)
			}
			os.Remove(certFile)
			os.Remove(keyFile)
		}
	}

	fmt.Println("Đang tạo self-signed certificate...")

	// Tạo private key
	privKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		fmt.Printf("Lỗi tạo private key: %v\n", err)
		os.Exit(1)
	}

	// Tạo certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Country:            []string{"VN"},
			Province:           []string{"HCM"},
			Locality:           []string{"HoChiMinh"},
			Organization:       []string{"Development"},
			OrganizationalUnit: []string{"Development"},
			CommonName:         "localhost",
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour), // 1 năm
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IPAddresses:           []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		DNSNames:              []string{"localhost", "*.localhost"},
	}

	// Tạo certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privKey.PublicKey, privKey)
	if err != nil {
		fmt.Printf("Lỗi tạo certificate: %v\n", err)
		os.Exit(1)
	}

	// Ghi certificate file
	certOut, err := os.Create(certFile)
	if err != nil {
		fmt.Printf("Lỗi tạo file certificate: %v\n", err)
		os.Exit(1)
	}
	defer certOut.Close()

	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certDER}); err != nil {
		fmt.Printf("Lỗi ghi certificate: %v\n", err)
		os.Exit(1)
	}

	// Ghi private key file
	keyOut, err := os.OpenFile(keyFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Printf("Lỗi tạo file key: %v\n", err)
		os.Exit(1)
	}
	defer keyOut.Close()

	privKeyDER := x509.MarshalPKCS1PrivateKey(privKey)
	if err := pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privKeyDER}); err != nil {
		fmt.Printf("Lỗi ghi private key: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Đã tạo certificate thành công!")
	fmt.Printf("  Cert: %s\n", certFile)
	fmt.Printf("  Key:  %s\n", keyFile)
	fmt.Println("")
	fmt.Println("⚠️  Lưu ý: Đây là self-signed certificate, trình duyệt sẽ cảnh báo.")
	fmt.Println("   Chấp nhận cảnh báo để tiếp tục sử dụng.")
}

