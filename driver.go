package main

import (
	"fmt"
	"log"  // For outputting stuff to the screen
	"time" // Used to set "ringing" to false if no ring detected within 2 seconds since the last ring

	"github.com/Grayda/go-phone"          // Our phone library
	"github.com/ninjasphere/go-ninja/api" // Ninja Sphere API
	"github.com/ninjasphere/go-ninja/model"
	"github.com/ninjasphere/go-ninja/support"
)

// package.json is required, otherwise the app just exits and doesn't show any output
var info = ninja.LoadModuleInfo("./package.json")
var driver *PhoneDriver // So we can access this in our configuration.go file

// Are we ready to rock? This is sphere-Phone only code by the way. You don't need to do this in your own driver?
var ready = false
var started = false // Stops us from running theloop twice

// PhoneDriver holds info about our driver, including our configuration
type PhoneDriver struct {
	support.DriverSupport
	config *PhoneDriverConfig // This is how we save and load call logs and such. Call this by using driver.config
	conn   *ninja.Connection  // For connecting to MQTT?
}

// PhoneDriverConfig holds config info.
type PhoneDriverConfig struct {
	Initialised bool // Has our driver run once before?
	COMPort     string
	Numbers     []PhoneCallRecord
}

// PhoneCallRecord holds info about our phonecall, including the time, and the number, if caller ID is enabled
type PhoneCallRecord struct {
	Time   string
	Number string
}

// No config provided? Set up some defaults
func defaultConfig() *PhoneDriverConfig {
	return &PhoneDriverConfig{
		Initialised: false,
		COMPort:     "/dev/ttyACM0",
	}
}

// NewDriver does what it says on the tin: makes a new driver for us to run. This is called through main.go
func NewDriver() (*PhoneDriver, error) {

	// Make a new PhoneDriver. Ampersand means to make a new copy, not reference the parent one (so A = new B instead of A = new B, C = A)
	driver = &PhoneDriver{}
	// Initialize our driver. Throw back an error if necessary. Remember, := is basically a short way of saying "var blah string = 'abcd'"
	err := driver.Init(info)

	if err != nil {
		log.Fatalf("Failed to initialize Phone driver: %s", err)
	}

	// Now we export the driver so the Sphere can find it
	err = driver.Export(driver)

	if err != nil {
		log.Fatalf("Failed to export Phone driver: %s", err)
	}

	// NewDriver returns two things, PhoneDriver, and an error if present
	return driver, nil
}

// Start is where the fun and magic happens! The driver is fired up and connects to our serial port
func (d *PhoneDriver) Start(config *PhoneDriverConfig) error {

	log.Printf("Driver Starting with config %v", config)

	d.config = config // Load our config

	if len(d.config.COMPort) == 0 { // No config loaded? Make one
		d.config = defaultConfig()
	}

	// This tells the API that we're going to expose a UI, and to run GetActions() in configuration.go
	d.Conn.MustExportService(&configService{d}, "$driver/"+info.ID+"/configure", &model.ServiceAnnouncement{
		Schema: "/protocol/configuration",
	})

	// Not quite working yet, but shows an icon on the LED matrix
	//showPane()

	// Some sample data. Remove these for production
	driver.config.Numbers = append(driver.config.Numbers, PhoneCallRecord{
		Time:   "Yesterday",
		Number: "1234567890",
	})

	driver.config.Numbers = append(driver.config.Numbers, PhoneCallRecord{
		Time:   "Today",
		Number: "0987654321",
	})

	driver.config.Numbers = append(driver.config.Numbers, PhoneCallRecord{
		Time:   "Tomorrow",
		Number: "0123498765",
	})

	var err error
	// If we've not started the driver
	if started == false {
		// Start a loop that handles everything this driver does (waiting for data, events etc.)
		// We put it in its own loop to keep the code neat
		err = theloop(d, d.config)
	}
	if err != nil {
		fmt.Println("Cannot start the serial port. Reason is:", err)
	}
	return d.SendEvent("config", config)
}

func theloop(d *PhoneDriver, config *PhoneDriverConfig) error {
	// Connect to our serial port. Return an error if there is one
	err := phone.Start(config.COMPort)

	if err != nil {
		// TO-DO: Show an icon when the modem isn't detected?
		return err
	}

	// Run this concurrently to ensure the rest of the driver isn't held up on an infinite loop
	go func() {
		// If started = true, theloop isn't called twice
		started = true
		fmt.Println("Calling theloop. COM port is:", config.COMPort)

		for { // Loop forever
			select { // This lets us do non-blocking channel reads. If we have a message, process it. If not, loop
			case msg := <-phone.Events: // If there is an event waiting
				switch msg.Name { // What event is it?
				case "READY": // Connected to COM port, ready to start sending / receiving data!
					fmt.Println("Connected to serial. Waiting for data!")
					phone.Read()
				case "RING": // The phone is ringing.
					fmt.Println("Phone is ringing!")
					// TO-DO: Show an icon on the Sphere
					phone.Read()
				case "OTHER": // Another event. Usually "OK" or "FAIL" messages from our modem. We don't care about these
					phone.Read()
				case "NMBR": // Phone has detected a number
					if msg.Message == "P" { // P == Private Number
						t := time.Now()
						config.Numbers = append(config.Numbers, PhoneCallRecord{
							Time:   t.Format("Monday _2 Jan 2006 15:04:05 2006"),
							Number: "Private Number",
						})
					} else {
						t := time.Now()
						config.Numbers = append(config.Numbers, PhoneCallRecord{
							Time:   t.Format("Monday _2 Jan 2006 15:04:05 2006"),
							Number: msg.Message,
						})

					}
					fmt.Println("Number detected:", msg.Message)
					phone.Read()
				default: // If there are no messages to parse, look for more bytes from our port, then try again
					phone.Read()
				}
			}
		}
	}()
	return nil
}

// Stop does what it says on the tin. Stops our driver. Stops the remote pane, and close the serial port (necessary?)
func (d *PhoneDriver) Stop() error {
	phone.Stop()
	return nil
}
