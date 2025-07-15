package pion

import (
	"fmt"
	"github.com/deleteelf/goframework/webrtc/pion/fmtp"
	"github.com/pion/webrtc/v4"
	"strings"
)

type codecMatchType int

const (
	codecMatchNone    codecMatchType = 0
	codecMatchPartial codecMatchType = 1
	codecMatchExact   codecMatchType = 2
)

func findRTXPayloadType(needle webrtc.PayloadType, haystack []webrtc.RTPCodecParameters) webrtc.PayloadType {
	aptStr := fmt.Sprintf("apt=%d", needle)
	for _, c := range haystack {
		if aptStr == c.SDPFmtpLine {
			return c.PayloadType
		}
	}

	return webrtc.PayloadType(0)
}

func codecParametersFuzzySearch(
	needle webrtc.RTPCodecParameters,
	haystack []webrtc.RTPCodecParameters,
) (webrtc.RTPCodecParameters, codecMatchType) {
	needleFmtp := fmtp.Parse(
		needle.RTPCodecCapability.MimeType,
		needle.RTPCodecCapability.ClockRate,
		needle.RTPCodecCapability.Channels,
		needle.RTPCodecCapability.SDPFmtpLine)

	// First attempt to match on MimeType + ClockRate + Channels + SDPFmtpLine
	for _, c := range haystack {
		cfmtp := fmtp.Parse(
			c.RTPCodecCapability.MimeType,
			c.RTPCodecCapability.ClockRate,
			c.RTPCodecCapability.Channels,
			c.RTPCodecCapability.SDPFmtpLine)

		if needleFmtp.Match(cfmtp) {
			return c, codecMatchExact
		}
	}

	// Fallback to just MimeType + ClockRate + Channels
	for _, c := range haystack {
		if strings.EqualFold(c.RTPCodecCapability.MimeType, needle.RTPCodecCapability.MimeType) &&
			fmtp.ClockRateEqual(c.RTPCodecCapability.MimeType,
				c.RTPCodecCapability.ClockRate,
				needle.RTPCodecCapability.ClockRate) &&
			fmtp.ChannelsEqual(c.RTPCodecCapability.MimeType,
				c.RTPCodecCapability.Channels,
				needle.RTPCodecCapability.Channels) {
			return c, codecMatchPartial
		}
	}

	return webrtc.RTPCodecParameters{}, codecMatchNone
}

// FlattenErrs flattens multiple errors into one.
func FlattenErrs(errs []error) error {
	errs2 := []error{}
	for _, e := range errs {
		if e != nil {
			errs2 = append(errs2, e)
		}
	}
	if len(errs2) == 0 {
		return nil
	}

	return multiError(errs2)
}

type multiError []error //nolint:errname

func (me multiError) Error() string {
	var errstrings []string

	for _, err := range me {
		if err != nil {
			errstrings = append(errstrings, err.Error())
		}
	}

	if len(errstrings) == 0 {
		return "multiError must contain multiple error but is empty"
	}

	return strings.Join(errstrings, "\n")
}
