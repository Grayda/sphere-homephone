package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/ninjasphere/go-ninja/api"
	"github.com/ninjasphere/go-ninja/model"
)

// PhoneDevice holds info about our phone device.
type PhoneDevice struct {
	driver    ninja.Driver
	info      *model.Device
	sendEvent func(event string, payload interface{}) error // For pasing info back to the API. Use this to send configs and such
}

// NewPhoneDevice may not even be used in this driver?
func NewPhoneDevice(driver ninja.Driver) *PhoneDevice {
	name := "Phone"

	device := &PhoneDevice{
		driver: driver,

		info: &model.Device{
			NaturalID:     fmt.Sprintf("phone"),
			NaturalIDType: "phone",
			Name:          &name,
			Signatures: &map[string]string{ // This stuff appears in our JSON when browsing the "things" part of REST, plus also in MQTT
				"ninja:manufacturer": "Phone",
				"ninja:productName":  "FixedPhone",
				"ninja:productType":  "Phone", // I think thingType (and maybe productType) is stored in a redis database
				"ninja:thingType":    "phone",
			},
		},
	}

	return device
}

// GetDeviceInfo just returns some info back to the API (I think?) about what type of thing it is etc., to show the proper icon in the Sphere app
func (d *PhoneDevice) GetDeviceInfo() *model.Device {
	return d.info
}

// GetDriver does a similar thing as GetDeviceInfo. Because it's got a capital letter at the start, it's exported by Go, so other packages can access it.
func (d *PhoneDevice) GetDriver() ninja.Driver {
	return d.driver
}

// SetEventHandler is something I have no idea of. Looks like it's just handing a sendEvent off to the PhoneDevice. Necessary?
func (d *PhoneDevice) SetEventHandler(sendEvent func(event string, payload interface{}) error) {
	d.sendEvent = sendEvent
}

// Regex that finds a-z (lowercase) and 0-9. Used when creating a safe name for our socket
var reg, _ = regexp.Compile("[^a-z0-9]")

// SetName might not be used by us, but the Sphere might do this for us (though I doubt it?). If we rename a device in the Sphere app, make it into a name that won't cause crashes
func (d *PhoneDevice) SetName(name *string) (*string, error) {

	log.Printf("Setting device name to %s", *name)

	safe := reg.ReplaceAllString(strings.ToLower(*name), "")
	if len(safe) > 16 {
		safe = safe[0:16]
	}

	log.Printf("We can only set 5 lowercase alphanum. Name now: %s", safe)
	d.info.Name = &safe
	d.sendEvent("renamed", safe)

	return &safe, nil
}
