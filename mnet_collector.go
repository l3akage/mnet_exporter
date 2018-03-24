package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

const prefix = "mnet_aircon_"

var (
	upDesc            *prometheus.Desc
	acUpDesc          *prometheus.Desc
	acCurrentTempDesc *prometheus.Desc
	acTargetTempDesc  *prometheus.Desc
)

func init() {
	upDesc = prometheus.NewDesc(prefix+"up", "Scrape was successful", nil, nil)
	acUpDesc = prometheus.NewDesc(prefix+"ac_up", "Ac turned on", []string{"id"}, nil)
	acCurrentTempDesc = prometheus.NewDesc(prefix+"ac_current_temp", "Current temp measured by the ac", []string{"id", "mode", "speed", "direction"}, nil)
	acTargetTempDesc = prometheus.NewDesc(prefix+"ac_target_temp", "Target temp set", []string{"id", "mode", "speed", "direction"}, nil)
}

type mnetCollector struct {
}

func (c mnetCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- upDesc
	ch <- acUpDesc
	ch <- acCurrentTempDesc
	ch <- acTargetTempDesc
}

func loadData() (Packet, error) {
	var x = Packet{}
	body := "<?xml version=\"1.0\" encoding=\"UTF-8\"?><Packet><Command>getRequest</Command><DatabaseManager>"
	for _, id := range strings.Split(*acIds, ",") {
		body += "<Mnet Group=\"" + id + "\" Bulk=\"*\" EnergyControl=\"*\" SetbackControl=\"*\" ScheduleAvail=\"*\" />"
	}
	body += "</DatabaseManager></Packet>"

	url := "http://" + *address + "/servlet/MIMEReceiveServlet"
	resp, err := http.Post(url, "text/xml", strings.NewReader(body))
	if err != nil {
		return x, err
	}
	defer resp.Body.Close()

	rawbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return x, err
	}

	err = xml.Unmarshal([]byte(rawbody), &x)
	if err != nil {
		return x, err
	}

	return x, nil
}

func convertAc(packet Packet) ([]ac, error) {
	var acs []ac
	for _, mnet := range packet.DatabaseManager.Mnet {
		var values []int
		var value string
		for _, r := range mnet.Bulk {
			value += string(r)
			if len(value) == 2 {
				x, err := strconv.ParseInt(value, 16, 64)
				if err != nil {
					return nil, err
				}
				values = append(values, int(x)&0xFF)
				value = ""
			}
		}

		a := ac{
			Group:       mnet.Group,
			State:       convertDrive(values[1]) == "ON",
			Mode:        convertMode(values[2]),
			Speed:       convertFanSpeed(values[8]),
			Direction:   convertAirDirection(values[7]),
			TargetTemp:  convertTemp(values[3], values[4]),
			TargetTemp1: convertTempX(values[83], values[84]),
			TargetTemp2: convertTempX(values[85], values[86]),
			CurrentTemp: convertTempX(values[5], values[6]),
			UseTemp:     convertUseTemp(values[136]),
		}
		acs = append(acs, a)
	}

	return acs, nil
}

func (c mnetCollector) Collect(ch chan<- prometheus.Metric) {
	r, err := loadData()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error loading data", err)
		ch <- prometheus.MustNewConstMetric(upDesc, prometheus.GaugeValue, 0)
	} else {
		acs, err := convertAc(r)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error converting data", err)
			ch <- prometheus.MustNewConstMetric(upDesc, prometheus.GaugeValue, 0)
		} else {
			ch <- prometheus.MustNewConstMetric(upDesc, prometheus.GaugeValue, 1)
			for _, ac := range acs {
				l := []string{ac.Group}
				state := 0
				if ac.State {
					state = 1
				}
				ch <- prometheus.MustNewConstMetric(acUpDesc, prometheus.GaugeValue, float64(state), l...)

				l = append(l, ac.Mode, ac.Speed, ac.Direction)
				ch <- prometheus.MustNewConstMetric(acCurrentTempDesc, prometheus.GaugeValue, ac.CurrentTemp, l...)

				var temp float64
				if ac.UseTemp {
					temp = ac.TargetTemp
				} else if ac.Mode == "HEAT" {
					temp = ac.TargetTemp2
				} else {
					temp = ac.TargetTemp1
				}

				ch <- prometheus.MustNewConstMetric(acTargetTempDesc, prometheus.GaugeValue, temp, l...)
			}
		}
	}
}
