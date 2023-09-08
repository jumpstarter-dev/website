package jumpstarter_board

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
)

const BASE_DISKSBYID = "/dev/disk/by-id/"

const BLOCK_SIZE = 32 * 1024 * 1024

const WAIT_TIME_USB_STORAGE = 6 * time.Second
const WAIT_TIME_USB_STORAGE_OFF = 2 * time.Second

type StorageTarget int

const (
	HOST StorageTarget = iota
	DUT
	OFF
)

func (d *JumpstarterDevice) SetDiskImage(path string) error {

	fmt.Print("🔍 Detecting USB storage device and connecting to host: ")
	diskPath, err := d.detectStorageDevice()
	if err != nil {
		return fmt.Errorf("SetDiskImage: %w", err)
	}
	fmt.Println("done")

	fmt.Printf("📋 %s -> %s: \n", path, diskPath)

	if err := writeImageToDisk(path, diskPath); err != nil {
		return fmt.Errorf("SetDiskImage: %w", err)
	}

	if err := d.connectStorageTo(OFF); err != nil {
		return fmt.Errorf("SetDiskImage: %w", err)
	}

	return nil
}

func (d *JumpstarterDevice) AttachStorage(connected bool) error {
	var err error
	switch connected {
	case true:
		err = d.connectStorageTo(DUT)
	case false:
		err = d.connectStorageTo(OFF)
	}
	if err != nil {
		return fmt.Errorf("ConnectDiskImage(%v): %w", connected, err)
	}
	return nil
}

func (d *JumpstarterDevice) connectStorageTo(target StorageTarget) error {
	if err := d.ensureSerial(); err != nil {
		return fmt.Errorf("connectStorageTo: %w", err)
	}

	if err := d.exitConsole(); err != nil {
		return fmt.Errorf("connectStorageTo: %w", err)
	}

	var cmd string
	var response string
	switch target {
	case HOST:
		cmd = "host"
		response = "connected to host"
	case DUT:
		cmd = "dut"
		response = "connected to device"
	case OFF:
		cmd = "off"
		response = "storage disconnected"
	default:
		return fmt.Errorf("connectStorageTo: invalid target %v", target)
	}

	if err := d.sendAndExpect("storage "+cmd, response); err != nil {
		return fmt.Errorf("connectStorageTo(%q): %w", cmd, err)
	}
	return nil
}

func scanForStorageDevices(prefix string) (*mapset.Set[string], error) {

	diskSet := mapset.NewSet[string]()

	err := filepath.Walk(BASE_DISKSBYID, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && info.Name() == "devices" {
			return nil
		}

		if info.Mode()&os.ModeSymlink != 0 {
			baseName := filepath.Base(path)
			re := regexp.MustCompile(`part\d+$`)
			if strings.HasPrefix(baseName, prefix) && !re.MatchString(baseName) {
				diskSet.Add(path)
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("scanForStorageDevices: %w", err)
	}

	return &diskSet, nil
}

func (d *JumpstarterDevice) detectStorageDevice() (string, error) {
	if err := d.connectStorageTo(OFF); err != nil {
		return "", fmt.Errorf("detectStorageDevice: %w", err)
	}
	time.Sleep(WAIT_TIME_USB_STORAGE_OFF)

	diskSetOff, err := scanForStorageDevices("usb-")
	if err != nil {
		return "", fmt.Errorf("detectStorageDevice: %w", err)
	}

	if err := d.connectStorageTo(HOST); err != nil {
		return "", fmt.Errorf("detectStorageDevice: %w", err)
	}

	// get current timestamp so we can measure how long it takes to detect the new disk
	start := time.Now()

	var diskSetOn *mapset.Set[string]

	for {
		time.Sleep(500 * time.Millisecond)
		diskSetOn, err = scanForStorageDevices("usb-")
		if err != nil {
			return "", fmt.Errorf("detectStorageDevice: %w", err)
		}
		newDiskSet := (*diskSetOn).Difference(*diskSetOff)
		if newDiskSet.Cardinality() == 1 {
			diskPath, _ := newDiskSet.Pop()
			return diskPath, nil
		}
		if time.Since(start) > WAIT_TIME_USB_STORAGE {
			if newDiskSet.Cardinality() > 1 {
				return "", fmt.Errorf("detectStorageDevice: more than one new disk detected")
			}
			return "", fmt.Errorf("detectStorageDevice: no new disk detected after 30 seconds")
		}
	}

}

func writeImageToDisk(imagePath string, diskPath string) error {
	inputFile, err := os.OpenFile(imagePath, os.O_RDONLY, 0666)
	if err != nil {
		return fmt.Errorf("writeImageToDisk: %w", err)
	}
	defer inputFile.Close()

	outputFile, err := os.OpenFile(diskPath, os.O_WRONLY|os.O_SYNC, 0666)
	if err != nil {
		return fmt.Errorf("writeImageToDisk: %w", err)
	}

	buffer := make([]byte, BLOCK_SIZE)
	bytesCopied := 0
	start := time.Now()
	for {
		n, err := inputFile.Read(buffer)
		if err != nil && err != io.EOF {
			outputFile.Close()
			return fmt.Errorf("writeImageToDisk: %w", err)
		}
		if n == 0 {
			break
		}
		if _, err := outputFile.Write(buffer[:n]); err != nil {
			outputFile.Close()
			return fmt.Errorf("writeImageToDisk: %w", err)
		}
		elapsed := time.Since(start)

		bytesCopied += n
		MBCopied := bytesCopied / 1024 / 1024
		MBPerSec := float64(MBCopied) / elapsed.Seconds()
		fmt.Printf("\r💾 %d MB copied %.2f MB/s         ", MBCopied, MBPerSec)
	}
	outputFile.Close()
	fmt.Println()

	if err := exec.Command("sync").Run(); err != nil {
		return fmt.Errorf("writeImageToDisk: sync %w", err)
	}

	time.Sleep(WAIT_TIME_USB_STORAGE * 2)
	cmd := exec.Command("udisksctl", "power-off", "-b", diskPath)
	var errb bytes.Buffer
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("writeImageToDisk: %w %s", err, errb.String())
	}
	time.Sleep(WAIT_TIME_USB_STORAGE)
	return nil
}
