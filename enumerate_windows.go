package main

import (
	"bytes"
	"encoding/csv"
	"errors"
	"log"
	"os/exec"
)

//TODO: Call DLL directly
func enumerateDevices() ([]device, error) {
	if err := checkExe("powershell"); err != nil {
		return nil, err
	}
	cmd := exec.Command("powershell", "-command", `gwmi Win32_USBControllerDevice|%{[wmi]($_.Dependent)}|Select-Object Description,DeviceID|ConvertTo-CSV -notypeinformation|select -skip 1`)
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	csvReader := csv.NewReader(bytes.NewReader(out))
	deviceList := []device{}
	parsed, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	for _, d := range parsed {
		if d[0] != "" && d[1] != "" {
			deviceList = append(deviceList, device{d[0], d[1]})
		} else {
			return nil, errors.New("incorrect device output")
		}
	}
	return deviceList, nil
}