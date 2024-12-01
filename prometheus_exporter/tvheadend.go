package prometheus_exporter

import (
	"encoding/json"
	"fmt"
	dac "github.com/Snawoot/go-http-digest-auth-client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

var tvLabels = []string{"decoder"}
var decoderSubs = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "tvheadend_decoder_subs"}, tvLabels)
var decoderWeight = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "tvheadend_decoder_weight"}, tvLabels)
var decoderSignal = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "tvheadend_decoder_signal"}, tvLabels)
var decoderSignalScale = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "tvheadend_decoder_signal_scale"}, tvLabels)
var decoderBer = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "tvheadend_decoder_ber"}, tvLabels)
var decoderSnr = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "tvheadend_decoder_snr"}, tvLabels)
var decoderSnrScale = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "tvheadend_decoder_snr_scale"}, tvLabels)
var decoderUnc = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "tvheadend_decoder_unc"}, tvLabels)
var decoderBps = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "tvheadend_decoder_bps"}, tvLabels)
var decoderTe = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "tvheadend_decoder_te"}, tvLabels)
var decoderCc = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "tvheadend_decoder_cc"}, tvLabels)
var decoderECBit = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "tvheadend_decoder_ec_bit"}, tvLabels)
var decoderTCBit = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "tvheadend_decoder_tc_bit"}, tvLabels)
var decoderECBlock = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "tvheadend_decoder_ec_block"}, tvLabels)
var decoderTCBlock = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "tvheadend_decoder_tc_block"}, tvLabels)
var decoderLastUpdate = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "tvheadend_decoder_last_update"}, tvLabels)

type TvheadendConfig struct {
	Host     string
	Username string
	Password string
}

type TVInputsStatus struct {
	Entries []struct {
		Uuid        string    `json:"uuid"`
		Input       string    `json:"input"`
		Stream      string    `json:"stream"`
		Subs        float64   `json:"subs"`
		Weight      float64   `json:"weight"`
		Pids        []float64 `json:"pids"`
		Signal      float64   `json:"signal"`
		SignalScale float64   `json:"signal_scale"`
		Ber         float64   `json:"ber"`
		Snr         float64   `json:"snr"`
		SnrScale    float64   `json:"snr_scale"`
		Unc         float64   `json:"unc"`
		Bps         float64   `json:"bps"`
		Te          float64   `json:"te"`
		Cc          float64   `json:"cc"`
		EcBit       float64   `json:"ec_bit"`
		TcBit       float64   `json:"tc_bit"`
		EcBlock     float64   `json:"ec_block"`
		TcBlock     float64   `json:"tc_block"`
	} `json:"entries"`
	TotalCount int `json:"totalCount"`
}

func getTvheadendInputsStatus(config TvheadendConfig) (*TVInputsStatus, error) {
	url := fmt.Sprintf("%s/api/status/inputs", config.Host)

	client := &http.Client{
		Transport: dac.NewDigestTransport(config.Username, config.Password, http.DefaultTransport),
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var status TVInputsStatus
	err = json.Unmarshal(body, &status)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

func collectMetrics(logger *slog.Logger, config TvheadendConfig) {
	logger.Info("collecting tvheadend metrics")

	status, err := getTvheadendInputsStatus(config)
	if err != nil {
		logger.Error("couldn't collect status", "error", err)
		return
	}

	for _, entry := range status.Entries {
		labels := prometheus.Labels{"decoder": entry.Input}

		decoderSubs.With(labels).Set(entry.Subs)
		decoderWeight.With(labels).Set(entry.Weight)
		decoderSignal.With(labels).Set(entry.Signal)
		decoderSignalScale.With(labels).Set(entry.SignalScale)
		decoderBer.With(labels).Set(entry.Ber)
		decoderSnr.With(labels).Set(entry.Snr)
		decoderSnrScale.With(labels).Set(entry.SnrScale)
		decoderUnc.With(labels).Set(entry.Unc)
		decoderBps.With(labels).Set(entry.Bps)
		decoderTe.With(labels).Set(entry.Te)
		decoderCc.With(labels).Set(entry.Cc)
		decoderECBit.With(labels).Set(entry.EcBit)
		decoderTCBit.With(labels).Set(entry.TcBit)
		decoderECBlock.With(labels).Set(entry.EcBlock)
		decoderTCBlock.With(labels).Set(entry.TcBlock)
		decoderLastUpdate.With(labels).SetToCurrentTime()
	}
}

func CollectTvheadendMetrics(logger *slog.Logger, config TvheadendConfig, wg *sync.WaitGroup, quitChannel chan bool) {
	defer wg.Done()

	ticker := time.NewTicker(15 * time.Second)
	for {
		select {
		case <-quitChannel:
			logger.Info("closing collecting metrics gracefully")
			return
		case <-ticker.C:
			collectMetrics(logger, config)
		}
	}
}
