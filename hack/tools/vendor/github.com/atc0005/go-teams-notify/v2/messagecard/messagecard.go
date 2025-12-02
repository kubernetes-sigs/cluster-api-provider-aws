// Copyright 2020 Enrico Hoffmann
// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/go-teams-notify
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package messagecard

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	// PotentialActionOpenURIType is the type that must be used for OpenUri
	// potential action.
	PotentialActionOpenURIType = "OpenUri"

	// PotentialActionHTTPPostType is the type that must be used for HttpPOST
	// potential action.
	PotentialActionHTTPPostType = "HttpPOST"

	// PotentialActionActionCardType is the type that must be used for
	// ActionCard potential action.
	PotentialActionActionCardType = "ActionCard"

	// PotentialActionInvokeAddInCommandType is the type that must be used for
	// InvokeAddInCommand potential action.
	PotentialActionInvokeAddInCommandType = "InvokeAddInCommand"

	// PotentialActionActionCardInputTextInputType is the type that must be
	// used for ActionCard TextInput type.
	PotentialActionActionCardInputTextInputType = "TextInput"

	// PotentialActionActionCardInputDateInputType is the type that must be
	// used for ActionCard DateInput type.
	PotentialActionActionCardInputDateInputType = "DateInput"

	// PotentialActionActionCardInputMultichoiceInputType is the type that
	// must be used for ActionCard MultichoiceInput type.
	PotentialActionActionCardInputMultichoiceInputType = "MultichoiceInput"
)

// PotentialActionMaxSupported is the maximum number of actions allowed in a
// PotentialAction collection.
//
// https://docs.microsoft.com/en-us/outlook/actionable-messages/message-card-reference#actions
const PotentialActionMaxSupported = 4

// ErrPotentialActionsLimitReached indicates that the maximum supported number
// of potentialAction collection values has been reached for either a
// MessageCard or a Section.
var ErrPotentialActionsLimitReached = errors.New("potential actions collection limit reached")

// PotentialAction represents potential actions an user can do in a message
// card. See [Legacy actionable message card reference > Actions] for more
// information.
//
// [Legacy actionable message card reference > Actions]: https://docs.microsoft.com/en-us/outlook/actionable-messages/message-card-reference#actions
type PotentialAction struct {
	// Type of the potential action. Can be OpenUri, HttpPOST, ActionCard or
	// InvokeAddInCommand.
	Type string `json:"@type" yaml:"@type"`

	// Name property defines the text that will be displayed on screen for the
	// action.
	Name string `json:"name" yaml:"name"`

	// PotentialActionOpenURI is a set of options for openUri
	// potential action.
	PotentialActionOpenURI

	// PotentialActionHTTPPOST is a set of options for httpPOST
	// potential action.
	PotentialActionHTTPPOST

	// PotentialActionActionCard is a set of options for actionCard
	// potential action.
	PotentialActionActionCard

	// PotentialActionInvokeAddInCommand is a set of options for
	// invokeAddInCommand potential action.
	PotentialActionInvokeAddInCommand
}

// PotentialActionOpenURI represents a OpenUri potential action.
type PotentialActionOpenURI struct {
	// Targets is a collection of name/value pairs that defines one URI per
	// target operating system. Only used for OpenUri action type.
	Targets []PotentialActionOpenURITarget `json:"targets,omitempty" yaml:"targets,omitempty"`
}

// PotentialActionHTTPPOST represents a HttpPOST potential action.
type PotentialActionHTTPPOST struct {
	// Target defines the URL endpoint of the service that implements the
	// action. Only used for HttpPOST action type.
	Target string `json:"target,omitempty" yaml:"target,omitempty"`

	// Headers is a collection of MessageCardPotentialActionHeader objects
	// representing a set of HTTP headers that will be emitted when sending
	// the POST request to the target URL. Only used for HttpPOST action type.
	Headers []PotentialActionHTTPPOSTHeader `json:"headers,omitempty" yaml:"headers,omitempty"`

	// Body is the body of the POST request. Only used for HttpPOST action
	// type.
	Body string `json:"body,omitempty" yaml:"body,omitempty"`

	// BodyContentType is optional and specifies the MIME type of the body in
	// the POST request. Only used for HttpPOST action type.
	BodyContentType string `json:"bodyContentType,omitempty" yaml:"bodyContentType,omitempty"`
}

// PotentialActionActionCard represents an actionCard potential
// action.
type PotentialActionActionCard struct {
	// Inputs is a collection of inputs an user can provide before processing
	// the actions. Only used for ActionCard action type. Three types of
	// inputs are available: TextInput, DateInput and MultichoiceInput
	Inputs []PotentialActionActionCardInput `json:"inputs,omitempty" yaml:"inputs,omitempty"`

	// Actions are the available actions. Only used for ActionCard action
	// type.
	Actions []PotentialActionActionCardAction `json:"actions,omitempty" yaml:"actions,omitempty"`
}

// PotentialActionActionCardAction is used for configuring ActionCard actions
type PotentialActionActionCardAction struct {
	// Type of the action. Can be OpenUri, HttpPOST, ActionCard or
	// InvokeAddInCommand.
	Type string `json:"@type" yaml:"@type"`

	// Name property defines the text that will be displayed on screen for the
	// action.
	Name string `json:"name" yaml:"name"`

	// PotentialActionOpenURI is used to specify a openUri action
	// card's action.
	PotentialActionOpenURI

	// PotentialActionHTTPPOST is used to specify a httpPOST action
	// card's action.
	PotentialActionHTTPPOST
}

// PotentialActionInvokeAddInCommand represents an invokeAddInCommand
// potential action.
type PotentialActionInvokeAddInCommand struct {
	// AddInID specifies the add-in ID of the required add-in. Only used for
	// InvokeAddInCommand action type.
	AddInID string `json:"addInId,omitempty" yaml:"addInId,omitempty"`

	// DesktopCommandID specifies the ID of the add-in command button that
	// opens the required task pane. Only used for InvokeAddInCommand action
	// type.
	DesktopCommandID string `json:"desktopCommandId,omitempty" yaml:"desktopCommandId,omitempty"`

	// InitializationContext is an optional field which provides developers a
	// way to specify any valid JSON object. The value is serialized into a
	// string and made available to the add-in when the action is executed.
	// This allows the action to pass initialization data to the add-in. Only
	// used for InvokeAddInCommand action type.
	InitializationContext interface{} `json:"initializationContext,omitempty" yaml:"initializationContext,omitempty"`
}

// PotentialActionOpenURITarget is used for OpenUri action type.
// It defines one URI per target operating system.
type PotentialActionOpenURITarget struct {
	// OS defines the operating system the target uri refers to. Supported
	// operating system values are default, windows, iOS and android. The
	// default operating system will in most cases simply open the URI in a
	// web browser, regardless of the actual operating system.
	OS string `json:"os,omitempty" yaml:"os,omitempty"`

	// URI defines the URI being called.
	URI string `json:"uri,omitempty" yaml:"uri,omitempty"`
}

// PotentialActionHTTPPOSTHeader defines a HTTP header used for HttpPOST action type.
type PotentialActionHTTPPOSTHeader struct {
	// Name is the header name.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	// Value is the header value.
	Value string `json:"value,omitempty" yaml:"value,omitempty"`
}

// PotentialActionActionCardInput represents an ActionCard input.
type PotentialActionActionCardInput struct {
	// Type of the ActionCard input.
	// Must be either TextInput, DateInput or MultichoiceInput
	Type string `json:"@type" yaml:"@type"`

	// ID uniquely identifies the input so it is possible to reference it in
	// the URL or body of an HttpPOST action.
	ID string `json:"id,omitempty" yaml:"id,omitempty"`

	// Title defines a title for the input.
	Title string `json:"title,omitempty" yaml:"title,omitempty"`

	// Value defines the initial value of the input. For multi-choice inputs,
	// value must be equal to the value property of one of the input's
	// choices.
	Value string `json:"value,omitempty" yaml:"value,omitempty"`

	// MessageCardPotentialActionInputMultichoiceInput must be defined for
	// MultichoiceInput input type.
	PotentialActionActionCardInputMultichoiceInput

	// MessageCardPotentialActionInputTextInput must be defined for InputText
	// input type.
	PotentialActionActionCardInputTextInput

	// MessageCardPotentialActionInputDateInput must be defined for DateInput
	// input type.
	PotentialActionActionCardInputDateInput

	// IsRequired indicates whether users are required to type a value before
	// they are able to take an action that would take the value of the input
	// as a parameter.
	IsRequired bool `json:"isRequired,omitempty" yaml:"isRequired,omitempty"`
}

// PotentialActionActionCardInputTextInput represents a TextInput
// input used for potential action.
type PotentialActionActionCardInputTextInput struct {
	// MaxLength indicates the maximum number of characters that can be
	// entered.
	MaxLength int `json:"maxLength,omitempty" yaml:"maxLength,omitempty"`

	// IsMultiline indicates whether the text input should accept multiple
	// lines of text.
	IsMultiline bool `json:"isMultiline,omitempty" yaml:"isMultiline,omitempty"`
}

// PotentialActionActionCardInputMultichoiceInput represents a
// MultichoiceInput input used for potential action.
type PotentialActionActionCardInputMultichoiceInput struct {
	// Choices defines the values that can be selected for the multichoice
	// input.
	Choices []struct {
		Display string `json:"display,omitempty" yaml:"display,omitempty"`
		Value   string `json:"value,omitempty" yaml:"value,omitempty"`
	} `json:"choices,omitempty" yaml:"choices,omitempty"`

	// Style defines the style of the input. When IsMultiSelect is false,
	// setting the style property to expanded will instruct the host
	// application to try and display all choices on the screen, typically
	// using a set of radio buttons.
	Style string `json:"style,omitempty" yaml:"style,omitempty"`

	// IsMultiSelect indicates whether or not the user can select more than
	// one choice. The specified choices will be displayed as a list of
	// checkboxes. Default value is false.
	IsMultiSelect bool `json:"isMultiSelect,omitempty" yaml:"isMultiSelect,omitempty"`
}

// PotentialActionActionCardInputDateInput represents a DateInput
// input used for potential action.
type PotentialActionActionCardInputDateInput struct {
	// IncludeTime indicates whether the date input should allow for the
	// selection of a time in addition to the date.
	IncludeTime bool `json:"includeTime,omitempty" yaml:"includeTime,omitempty"`
}

// SectionFact represents a section fact entry that is usually displayed in a
// two-column key/value format.
type SectionFact struct {

	// Name is the key for an associated value in a key/value pair
	Name string `json:"name" yaml:"name"`

	// Value is the value for an associated key in a key/value pair
	Value string `json:"value" yaml:"value"`
}

// SectionImage represents an image as used by the heroImage and images
// properties of a section.
type SectionImage struct {

	// Image is the URL to the image.
	Image string `json:"image" yaml:"image"`

	// Title is a short description of the image. Typically, this description
	// is displayed in a tooltip as the user hovers their mouse over the
	// image.
	Title string `json:"title" yaml:"title"`
}

// Section represents a section to include in a message card.
type Section struct {
	// Title is the title property of a section. This property is displayed
	// in a font that stands out, while not as prominent as the card's title.
	// It is meant to introduce the section and summarize its content,
	// similarly to how the card's title property is meant to summarize the
	// whole card.
	Title string `json:"title,omitempty" yaml:"title,omitempty"`

	// Text is the section's text property. This property is very similar to
	// the text property of the card. It can be used for the same purpose.
	Text string `json:"text,omitempty" yaml:"text,omitempty"`

	// ActivityImage is a property used to display a picture associated with
	// the subject of a message card. For example, this might be the portrait
	// of a person who performed an activity that the message card is
	// associated with.
	ActivityImage string `json:"activityImage,omitempty" yaml:"activityImage,omitempty"`

	// ActivityTitle is a property used to summarize the activity associated
	// with a message card.
	ActivityTitle string `json:"activityTitle,omitempty" yaml:"activityTitle,omitempty"`

	// ActivitySubtitle is a property used to show brief, but extended
	// information about an activity associated with a message card. Examples
	// include the date and time the associated activity was taken or the
	// handle of a person associated with the activity.
	ActivitySubtitle string `json:"activitySubtitle,omitempty" yaml:"activitySubtitle,omitempty"`

	// ActivityText is a property used to provide details about the activity.
	// For example, if the message card is used to deliver updates about a
	// topic, then this property would be used to hold the bulk of the content
	// for the update notification.
	ActivityText string `json:"activityText,omitempty" yaml:"activityText,omitempty"`

	// HeroImage is a property that allows for setting an image as the
	// centerpiece of a message card. This property can also be used to add a
	// banner to the message card.
	// Note: heroImage is not currently supported by Microsoft Teams
	// https://stackoverflow.com/a/45389789
	// We use a pointer to this type in order to have the json package
	// properly omit this field if not explicitly set.
	// https://github.com/golang/go/issues/11939
	// https://stackoverflow.com/questions/18088294/how-to-not-marshal-an-empty-struct-into-json-with-go
	// https://stackoverflow.com/questions/33447334/golang-json-marshal-how-to-omit-empty-nested-struct
	HeroImage *SectionImage `json:"heroImage,omitempty" yaml:"heroImage,omitempty"`

	// Facts is a collection of SectionFact values. A section entry
	// usually is displayed in a two-column key/value format.
	Facts []SectionFact `json:"facts,omitempty" yaml:"facts,omitempty"`

	// Images is a property that allows for the inclusion of a photo gallery
	// inside a section.
	// We use a slice of pointers to this type in order to have the json
	// package properly omit this field if not explicitly set.
	// https://github.com/golang/go/issues/11939
	// https://stackoverflow.com/questions/18088294/how-to-not-marshal-an-empty-struct-into-json-with-go
	// https://stackoverflow.com/questions/33447334/golang-json-marshal-how-to-omit-empty-nested-struct
	Images []*SectionImage `json:"images,omitempty" yaml:"images,omitempty"`

	// PotentialActions is a collection of actions for a Section.
	// This is separate from the actions collection for the MessageCard.
	PotentialActions []*PotentialAction `json:"potentialAction,omitempty" yaml:"potentialAction,omitempty"`

	// Markdown represents a toggle to enable or disable Markdown formatting.
	// By default, all text fields in a card and its sections can be formatted
	// using basic Markdown.
	Markdown bool `json:"markdown,omitempty" yaml:"markdown,omitempty"`

	// StartGroup is the section's startGroup property. This property marks
	// the start of a logical group of information. Typically, sections with
	// startGroup set to true will be visually separated from previous card
	// elements.
	StartGroup bool `json:"startGroup,omitempty" yaml:"startGroup,omitempty"`
}

// MessageCard represents a legacy actionable message card used via Office 365
// or Microsoft Teams connectors.
type MessageCard struct {
	// Required; must be set to "MessageCard"
	Type string `json:"@type" yaml:"@type"`

	// Required; must be set to "https://schema.org/extensions"
	Context string `json:"@context" yaml:"@context"`

	// Summary is required if the card does not contain a text property,
	// otherwise optional. The summary property is typically displayed in the
	// list view in Outlook, as a way to quickly determine what the card is
	// all about. Summary appears to only be used when there are sections defined
	Summary string `json:"summary,omitempty" yaml:"summary,omitempty"`

	// Title is the title property of a card. is meant to be rendered in a
	// prominent way, at the very top of the card. Use it to introduce the
	// content of the card in such a way users will immediately know what to
	// expect.
	Title string `json:"title,omitempty" yaml:"title,omitempty"`

	// Text is required if the card does not contain a summary property,
	// otherwise optional. The text property is meant to be displayed in a
	// normal font below the card's title. Use it to display content, such as
	// the description of the entity being referenced, or an abstract of a
	// news article.
	Text string `json:"text,omitempty" yaml:"text,omitempty"`

	// Specifies a custom brand color for the card. The color will be
	// displayed in a non-obtrusive manner.
	ThemeColor string `json:"themeColor,omitempty" yaml:"themeColor,omitempty"`

	// ValidateFunc is an optional user-specified validation function that is
	// responsible for validating a MessageCard. If not specified, default
	// validation is performed.
	ValidateFunc func() error `json:"-" yaml:"-"`

	// Sections is a collection of sections to include in the card.
	Sections []*Section `json:"sections,omitempty" yaml:"sections,omitempty"`

	// PotentialActions is a collection of actions for a MessageCard.
	PotentialActions []*PotentialAction `json:"potentialAction,omitempty" yaml:"potentialAction,omitempty"`

	// payload is a prepared MessageCard in JSON format for submission or
	// pretty printing.
	payload *bytes.Buffer `json:"-" yaml:"-"`
}

// validatePotentialAction inspects the given *PotentialAction
// and returns an error if a value is missing or not known.
func validatePotentialAction(pa *PotentialAction) error {
	if pa == nil {
		return fmt.Errorf("nil PotentialAction received")
	}

	switch pa.Type {
	case PotentialActionOpenURIType,
		PotentialActionHTTPPostType,
		PotentialActionActionCardType,
		PotentialActionInvokeAddInCommandType:

	default:
		return fmt.Errorf("unknown type %s for potential action %s", pa.Type, pa.Name)
	}

	if pa.Name == "" {
		return fmt.Errorf("missing name value for PotentialAction")
	}

	return nil
}

// addPotentialAction adds one or many PotentialAction values to a
// PotentialActions collection.
func addPotentialAction(collection *[]*PotentialAction, actions ...*PotentialAction) error {
	for _, a := range actions {
		if err := validatePotentialAction(a); err != nil {
			return err
		}

		if len(*collection) > PotentialActionMaxSupported {
			return fmt.Errorf("func addPotentialAction: failed to add potential action: %w", ErrPotentialActionsLimitReached)
		}

		*collection = append(*collection, a)
	}

	return nil
}

// AddSection adds one or many additional Section values to a MessageCard.
// Validation is performed to reject invalid values with an error message.
func (mc *MessageCard) AddSection(section ...*Section) error {
	for _, s := range section {
		// bail if a completely nil section provided
		if s == nil {
			return fmt.Errorf("func AddSection: nil Section received")
		}

		// Perform validation of all Section fields in an effort to
		// avoid adding a Section with zero value fields. This is
		// done to avoid generating an empty sections JSON array since the
		// Sections slice for the MessageCard type would technically not be at
		// a zero value state. Due to this non-zero value state, the
		// encoding/json package would end up including the Sections struct
		// field in the output JSON.
		// See also https://github.com/golang/go/issues/11939
		switch {
		// If any of these cases trigger, skip over the `default` case
		// statement and add the section.
		case s.Images != nil:
		case s.Facts != nil:
		case s.HeroImage != nil:
		case s.StartGroup:
		case s.Markdown:
		case s.ActivityText != "":
		case s.ActivitySubtitle != "":
		case s.ActivityTitle != "":
		case s.ActivityImage != "":
		case s.Text != "":
		case s.Title != "":

		default:
			return fmt.Errorf("all fields found to be at zero-value, skipping section")
		}

		mc.Sections = append(mc.Sections, s)
	}

	return nil
}

// AddPotentialAction adds one or many PotentialAction values to a
// PotentialActions collection on a MessageCard.
func (mc *MessageCard) AddPotentialAction(actions ...*PotentialAction) error {
	return addPotentialAction(&mc.PotentialActions, actions...)
}

// Validate performs validation for MessageCard using ValidateFunc if defined,
// otherwise applying default validation.
func (mc *MessageCard) Validate() error {
	if mc.ValidateFunc != nil {
		return mc.ValidateFunc()
	}

	// Falling back to a default implementation
	if (mc.Text == "") && (mc.Summary == "") {
		// This scenario results in:
		// 400 Bad Request
		// Summary or Text is required.
		return fmt.Errorf("invalid message card: summary or text field is required")
	}

	return nil
}

// Prepare handles tasks needed to construct a payload from a MessageCard for
// delivery to an endpoint.
func (mc *MessageCard) Prepare() error {
	jsonMessage, err := json.Marshal(mc)
	if err != nil {
		return fmt.Errorf(
			"error marshalling MessageCard to JSON: %w",
			err,
		)
	}

	switch {
	case mc.payload == nil:
		mc.payload = &bytes.Buffer{}
	default:
		mc.payload.Reset()
	}

	_, err = mc.payload.Write(jsonMessage)
	if err != nil {
		return fmt.Errorf(
			"error updating JSON payload for MessageCard: %w",
			err,
		)
	}

	return nil
}

// Payload returns the prepared MessageCard payload. The caller should call
// Prepare() prior to calling this method, results are undefined otherwise.
func (mc *MessageCard) Payload() io.Reader {
	return mc.payload
}

// PrettyPrint returns a formatted JSON payload of the MessageCard if the
// Prepare() method has been called, or an empty string otherwise.
func (mc *MessageCard) PrettyPrint() string {
	if mc.payload != nil {
		var prettyJSON bytes.Buffer
		_ = json.Indent(&prettyJSON, mc.payload.Bytes(), "", "\t")

		return prettyJSON.String()
	}

	return ""
}

// AddFact adds one or many additional SectionFact values to a
// Section
func (mcs *Section) AddFact(fact ...SectionFact) error {
	for _, f := range fact {
		if f.Name == "" {
			return fmt.Errorf("empty Name field received for new fact: %+v", f)
		}

		if f.Value == "" {
			return fmt.Errorf("empty Value field received for new fact: %+v", f)
		}
	}

	mcs.Facts = append(mcs.Facts, fact...)

	return nil
}

// AddFactFromKeyValue accepts a key and slice of values and converts them to
// SectionFact values
func (mcs *Section) AddFactFromKeyValue(key string, values ...string) error {
	// validate arguments

	if key == "" {
		return errors.New("empty key received for new fact")
	}

	if len(values) < 1 {
		return errors.New("no values received for new fact")
	}

	fact := SectionFact{
		Name:  key,
		Value: strings.Join(values, ", "),
	}

	mcs.Facts = append(mcs.Facts, fact)

	// if we made it this far then all should be well
	return nil
}

// AddPotentialAction adds one or many PotentialAction values to a
// PotentialActions collection on a Section. This is separate from
// the actions collection for the MessageCard.
func (mcs *Section) AddPotentialAction(actions ...*PotentialAction) error {
	return addPotentialAction(&mcs.PotentialActions, actions...)
}

// AddImage adds an image to a MessageCard section. These images are used to
// provide a photo gallery inside a MessageCard section.
func (mcs *Section) AddImage(sectionImage ...SectionImage) error {
	for i := range sectionImage {
		if sectionImage[i].Image == "" {
			return fmt.Errorf("cannot add empty image URL")
		}

		if sectionImage[i].Title == "" {
			return fmt.Errorf("cannot add empty image title")
		}

		mcs.Images = append(mcs.Images, &sectionImage[i])
	}

	return nil
}

// AddHeroImageStr adds a Hero Image to a MessageCard section using string
// arguments. This image is used as the centerpiece or banner of a message
// card.
func (mcs *Section) AddHeroImageStr(imageURL string, imageTitle string) error {
	if imageURL == "" {
		return fmt.Errorf("cannot add empty hero image URL")
	}

	if imageTitle == "" {
		return fmt.Errorf("cannot add empty hero image title")
	}

	heroImage := SectionImage{
		Image: imageURL,
		Title: imageTitle,
	}

	mcs.HeroImage = &heroImage

	// our validation checks didn't find any problems
	return nil
}

// AddHeroImage adds a Hero Image to a MessageCard section using a
// SectionImage argument. This image is used as the centerpiece or
// banner of a message card.
func (mcs *Section) AddHeroImage(heroImage SectionImage) error {
	if heroImage.Image == "" {
		return fmt.Errorf("cannot add empty hero image URL")
	}

	if heroImage.Title == "" {
		return fmt.Errorf("cannot add empty hero image title")
	}

	mcs.HeroImage = &heroImage

	// our validation checks didn't find any problems
	return nil
}

// NewMessageCard creates a new legacy MessageCard with required fields
// predefined.
//
// https://docs.microsoft.com/en-us/outlook/actionable-messages/message-card-reference#card-fields
func NewMessageCard() *MessageCard {
	return &MessageCard{
		Type:    "MessageCard",
		Context: "https://schema.org/extensions",
	}
}

// NewSection creates an empty message card section
func NewSection() *Section {
	msgCardSection := Section{}
	return &msgCardSection
}

// NewSectionFact creates an empty message card section fact
func NewSectionFact() *SectionFact {
	return &SectionFact{}
}

// NewSectionImage creates an empty image for use with message card
// section
func NewSectionImage() *SectionImage {
	return &SectionImage{}
}

// NewPotentialAction creates a new PotentialAction using the provided
// potential action type and name. The name value defines the text that will
// be displayed on screen for the action. An error is returned if invalid
// values are supplied.
func NewPotentialAction(potentialActionType string, name string) (*PotentialAction, error) {
	pa := PotentialAction{
		Type: potentialActionType,
		Name: name,
	}

	if err := validatePotentialAction(&pa); err != nil {
		return nil, err
	}

	return &pa, nil
}
