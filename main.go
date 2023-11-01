package main

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
	"syscall"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type Server struct {
	PreferredDNS string
	AlternateDNS string
}

var (
	Shecan       = Server{PreferredDNS: "178.22.122.100", AlternateDNS: "185.51.200.2"}
	Online403    = Server{PreferredDNS: "10.202.10.202", AlternateDNS: "10.202.10.102"}
	Electro      = Server{PreferredDNS: "78.157.42.100", AlternateDNS: "78.157.42.101"}
	CurrentPref  string
	CurrentAlter string
)

const (
	AdapterName = "Wi-Fi"
)

func main() {
	a := app.New()
	w := a.NewWindow("SpiritDns")
	pW := widget.NewLabel("...")
	aW := widget.NewLabel("...")
	SetDnsText(pW, aW)
	shecanBtn := widget.NewButton("Shecan", func() {
		setDNS(AdapterName, Shecan.PreferredDNS, Shecan.AlternateDNS)
		SetDnsText(pW, aW)
	})
	online403Btn := widget.NewButton("403", func() {
		setDNS(AdapterName, Online403.PreferredDNS, Online403.AlternateDNS)
		SetDnsText(pW, aW)
	})
	electro := widget.NewButton("Electro", func() {
		setDNS(AdapterName, Electro.PreferredDNS, Electro.AlternateDNS)
		SetDnsText(pW, aW)
	})
	clear := widget.NewButton("Clear Dns", func() {
		setDNS(AdapterName, "", "")
		SetDnsText(pW, aW)
	})

	dns := container.New(layout.NewGridLayout(2), pW, aW)
	continer := container.New(layout.NewGridLayout(1), shecanBtn, online403Btn, electro, clear, dns)
	w.CenterOnScreen()
	w.Resize(fyne.Size{250, 200})
	w.SetContent(continer)
	w.ShowAndRun()
}

func setDNS(adapterName, preferredDNS, alternateDNS string) error {
	cmd := exec.Command("netsh", "interface", "ipv4", "set", "dns", `name="`+adapterName+`"`, "source=static", "address="+preferredDNS)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error setting preferred DNS: %v", err)
	}

	cmd = exec.Command("netsh", "interface", "ipv4", "add", "dns", `name="`+adapterName+`"`, "address="+alternateDNS, "index=2")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error setting alternate DNS: %v", err)
	}

	return nil
}

func getDNSAddressesForAdapter(adapterName string) (preferredDNS, alternateDNS string, err error) {
	cmd := exec.Command("netsh", "interface", "ipv4", "show", "dnsservers", `name="`+adapterName+`"`)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", "", err
	}

	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")

	for i, line := range lines {
		if i == 2 {
			s := strings.Replace(line, "Statically Configured DNS Servers:", " ", 1)
			preferredDNS = strings.TrimSpace(s)
		}
		if i == 3 {
			alternateDNS = strings.TrimSpace(line)
		}
	}

	return preferredDNS, alternateDNS, nil
}

func SetDnsText(pW, aW *widget.Label) {
	pI, aI, err := getDNSAddressesForAdapter(AdapterName)
	if err != nil {
		fmt.Errorf("error getting DNS: %v", err)
	}
	ipP := net.ParseIP(pI)
	ipA := net.ParseIP(aI)
	if ipP != nil && ipA != nil {
		pW.SetText(pI)
		aW.SetText(aI)
	} else {
		pW.SetText("empty")
		aW.SetText("empty")
	}

}
