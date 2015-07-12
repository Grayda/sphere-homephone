package main

import (
	"encoding/json"
	"fmt"

	"github.com/Grayda/go-phone"
	"github.com/ninjasphere/go-ninja/model"
	"github.com/ninjasphere/go-ninja/suit"
)

// This file contains most of the code for the UI (i.e. what appears in the Labs)

type configService struct {
	driver *PhoneDriver
}

// This function is common across all UIs, and is called by the Sphere. Shows our menu option on the main Labs screen
// The "c" bit at the start means this func is an extension of the configService struct (like prototyping, I think?)
func (c *configService) GetActions(request *model.ConfigurationRequest) (*[]suit.ReplyAction, error) {
	// What we're going to show
	var screen []suit.ReplyAction
	screen = append(screen, suit.ReplyAction{
		Name:        "",
		Label:       "Landline Phone",
		DisplayIcon: "phone",
	},
	)

	// Return our screen to the sphere-ui for rendering
	return &screen, nil
}

// When you click on a ReplyAction button (e.g. the "Configure AllOne" button defined above), Configure is called. requests.Action == the "Name" portion of the ReplyAction
func (c *configService) Configure(request *model.ConfigurationRequest) (*suit.ConfigurationScreen, error) {
	fmt.Sprintf("Incoming configuration request. Action:%s Data:%s", request.Action, string(request.Data))

	switch request.Action {

	case "setcomport": // Blasting IR codes
		// Make a map of strings
		var vals map[string]string
		// Take our json response from sphere-ui and place it into our vals map
		json.Unmarshal(request.Data, &vals)
		// Because ActionListOption can't return more than one lot of data and we need to let our code know
		// WHAT code to blast, and what AllOne to shoot it from, we use a pipe to mash data together
		driver.config.COMPort = vals["comport"]
		phone.Stop()
		phone.Start(driver.config.COMPort)
		// c.list creates a list of AllOne IR codes and sends them back to sphere-ui / suits for displaying
		return c.list()
	case "": // Coming in from the main menu
		return c.list()

	default: // Everything else

		// return c.list()
		return c.error(fmt.Sprintf("Unknown action: %s", request.Action), true)
	}

	// If this code runs, then we done fucked up, because default: didn't catch. When this code runs, the universe melts into a gigantic heap. But
	// removing this violates Apple guidelines and ensures the downfall of humanity (probably) so I don't want to risk it.
	// Then again, I could be making all this up. Do you want to remove it and try? ( ͡° ͜ʖ ͡°)
	return nil, nil
}

// So this function (which is an extension of the configService struct that suit (or Sphere-UI) requires) creates a box with a single "Okay" button and puts in a title and text
func (c *configService) confirm(title string, description string) (*suit.ConfigurationScreen, error) {
	// We create a new suit.ConfigurationScreen which is a whole page within the UI
	screen := suit.ConfigurationScreen{
		Title: title,
		Sections: []suit.Section{ // The UI lets us create sections for separating options. This line creates an array of sections
			suit.Section{ // And within that array of sections, a single section
				Contents: []suit.Typed{ // The contents of that section. I don't know what suit.Typed is. It's an interface, but asides from that, I don't know much else just yet
					suit.StaticText{ // Create some static text
						Title: "About this screen",
						Value: description,
					},
				},
			},
		},
		Actions: []suit.Typed{ // This configuration screen can show actionable buttons at the bottom. ReplyAction, as shown above, calls Configure. There is also CloseAction for cancel buttons
			suit.ReplyAction{
				Label:        "Okay",
				Name:         "list",
				DisplayClass: "success", // These are bootstrap classes (or rather, font-awesome classes). They are basically btn-*, where * is DisplayClass (e.g. btn-success)
				DisplayIcon:  "ok",      // Same as above. If you want to show fa-open-folder, you'd set DisplayIcon to "open-folder"
			},
		},
	}

	return &screen, nil
}

// Error! Same as above. It's a function that is added on to configService and displays an error message
func (c *configService) error(message string, cancel bool) (*suit.ConfigurationScreen, error) {

	var action []suit.Typed
	if cancel == true {

		action = []suit.Typed{
			suit.ReplyAction{ // Shows a button we can click on. Takes us back to c.Configuration (reply.Action will be "list")
				Label:        "OK",
				Name:         "list",
				DisplayClass: "success",
				DisplayIcon:  "ok",
			},
		}
	} else {
		action = []suit.Typed{
			suit.CloseAction{ // Shows a button we can click on. Takes us back to c.Configuration (reply.Action will be "list")
				Label: "OK",
			},
		}

	}

	return &suit.ConfigurationScreen{
		Sections: []suit.Section{
			suit.Section{
				Contents: []suit.Typed{
					suit.Alert{
						Title:        "Error",
						Subtitle:     message,
						DisplayClass: "danger",
					},
				},
			},
		},

		Actions: action,
	}, nil
}

// The meat of our UI. Shows a list of IR codes to be blasted. This could show anything you like, really.
func (c *configService) list() (*suit.ConfigurationScreen, error) {

	if phone.Connected == false {
		return c.error("Modem isn't connected! It should be connected to "+driver.config.COMPort+". Please green reset the Sphere and try again", false)
	}

	// Sections, for logical grouping
	var sections []suit.Section

	for _, num := range driver.config.Numbers {
		sections = append(sections, suit.Section{ // Append a new suit.Section into our sections variable
			Contents: []suit.Typed{ // Again, dunno what this means
				suit.StaticText{
					Title: "Phone call at " + num.Time,
					Value: "Number: " + num.Number,
				},
			},
		})

	}

	// Now that we've looped and got our sections, it's time to build the actual screen
	screen := suit.ConfigurationScreen{
		Title:    "Landline Phone",
		Sections: sections, // Our sections. Contains all the buttons and everything!
		Actions: []suit.Typed{ // Actiosn you can take on this page
			suit.CloseAction{ // Here we go! This takes a label and probably a DisplayIcon and DisplayClass and just takes you back to the main screen. Not YOUR main screen though, so use a ReplyAction with a "" name to go back to YOUR menu
				Label: "Done",
			},
		},
	}

	return &screen, nil
}

// Aye-aye, captain.
// Not actually needed (?)
func i(i int) *int {
	return &i
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
