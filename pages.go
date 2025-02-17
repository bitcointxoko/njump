//go:generate tmpl bind ./...

package main

import (
	_ "embed"
	"html/template"
	"strings"

	"github.com/nbd-wtf/go-nostr/nip11"
	sdk "github.com/nbd-wtf/nostr-sdk"
	"github.com/tylermmorton/tmpl"
)

type TemplateID int

const (
	Note TemplateID = iota
	Profile
	LongForm
	TelegramInstantView
	FileMetadata
	LiveEvent
	LiveEventMessage
	Other
)

var (
	//go:embed templates/opengraph.html
	tmplOpenGraph     string
	OpenGraphTemplate = tmpl.MustCompile(&OpenGraphPartial{})
)

//tmpl:bind head_common.html
type OpenGraphPartial struct {
	SingleTitle string
	// x (we will always render just the singletitle if we have that)
	Superscript string
	Subscript   string

	BigImage string
	// x (we will always render just the bigimage if we have that)
	Video        string
	VideoType    string
	Image        string
	ProxiedImage string

	// this is the main text we should always have
	Text string
}

func (*OpenGraphPartial) TemplateText() string { return tmplOpenGraph }

var (
	//go:embed templates/head_common.html
	tmplHeadCommon     string
	HeadCommonTemplate = tmpl.MustCompile(&HeadCommonPartial{})
)

//tmpl:bind head_common.html
type HeadCommonPartial struct {
	IsProfile          bool
	TailwindDebugStuff template.HTML
	NaddrNaked         string
	NeventNaked        string
	Oembed             string
}

func (*HeadCommonPartial) TemplateText() string { return tmplHeadCommon }

var (
	//go:embed templates/top.html
	tmplTop     string
	TopTemplate = tmpl.MustCompile(&TopPartial{})
)

//tmpl:bind top.html
type TopPartial struct{}

func (*TopPartial) TemplateText() string { return tmplTop }

var (
	//go:embed templates/details.html
	tmplDetails     string
	DetailsTemplate = tmpl.MustCompile(&DetailsPartial{})
)

//tmpl:bind details.html
type DetailsPartial struct {
	HideDetails     bool
	CreatedAt       string
	EventJSON       template.HTML
	Nevent          string
	Nprofile        string
	Npub            string
	SeenOn          []string
	Kind            int
	KindNIP         string
	KindDescription string

	// kind-specific stuff
	FileMetadata *Kind1063Metadata
	LiveEvent    *Kind30311Metadata
}

func (*DetailsPartial) TemplateText() string { return tmplDetails }

var (
	//go:embed templates/clients.html
	tmplClients     string
	ClientsTemplate = tmpl.MustCompile(&ClientsPartial{})
)

//tmpl:bind clients.html
type ClientsPartial struct {
	Clients []ClientReference
}

func (*ClientsPartial) TemplateText() string { return tmplClients }

var (
	//go:embed templates/footer.html
	tmplFooter     string
	FooterTemplate = tmpl.MustCompile(&FooterPartial{})
)

//tmpl:bind footer.html
type FooterPartial struct{}

func (*FooterPartial) TemplateText() string { return tmplFooter }

var (
	//go:embed templates/telegram_instant_view.html
	tmplTelegramInstantView     string
	TelegramInstantViewTemplate = tmpl.MustCompile(&TelegramInstantViewPage{})
)

type TelegramInstantViewPage struct {
	Video       string
	VideoType   string
	Image       string
	Summary     template.HTML
	Content     template.HTML
	Description string
	Subject     string
	Metadata    sdk.ProfileMetadata
	AuthorLong  string
	CreatedAt   string
	ParentLink  template.HTML
}

func (*TelegramInstantViewPage) TemplateText() string { return tmplTelegramInstantView }

var (
	//go:embed templates/homepage.html
	tmplHomePage     string
	HomePageTemplate = tmpl.MustCompile(&HomePage{})
)

type HomePage struct {
	HeadCommonPartial `tmpl:"head_common"`
	TopPartial        `tmpl:"top"`
	FooterPartial     `tmpl:"footer"`

	Host      string
	Npubs     []string
	LastNotes []string
}

func (*HomePage) TemplateText() string { return tmplHomePage }

var (
	//go:embed templates/archive.html
	tmplArchive     string
	ArchiveTemplate = tmpl.MustCompile(&ArchivePage{})
)

type ArchivePage struct {
	HeadCommonPartial `tmpl:"head_common"`
	TopPartial        `tmpl:"top"`
	FooterPartial     `tmpl:"footer"`

	Title         string
	PathPrefix    string
	Data          []string
	ModifiedAt    string
	PaginationUrl string
	NextPage      int
	PrevPage      int
}

func (*ArchivePage) TemplateText() string { return tmplArchive }

var (
	//go:embed templates/other.html
	tmplOther     string
	OtherTemplate = tmpl.MustCompile(&OtherPage{})
)

type OtherPage struct {
	HeadCommonPartial `tmpl:"head_common"`
	TopPartial        `tmpl:"top"`
	DetailsPartial    `tmpl:"details"`
	FooterPartial     `tmpl:"footer"`

	Kind            int
	KindDescription string
	Alt             string
}

func (*OtherPage) TemplateText() string { return tmplOther }

var (
	//go:embed templates/note.html
	tmplNote     string
	NoteTemplate = tmpl.MustCompile(&NotePage{})
)

type NotePage struct {
	OpenGraphPartial  `tmpl:"opengraph"`
	HeadCommonPartial `tmpl:"head_common"`
	TopPartial        `tmpl:"top"`
	DetailsPartial    `tmpl:"details"`
	ClientsPartial    `tmpl:"clients"`
	FooterPartial     `tmpl:"footer"`

	Content          template.HTML
	CreatedAt        string
	Metadata         sdk.ProfileMetadata
	Npub             string
	NpubShort        string
	ParentLink       template.HTML
	SeenOn           []string
	Subject          string
	TitleizedContent string
}

func (*NotePage) TemplateText() string { return tmplNote }

var (
	//go:embed templates/embedded_note.html
	tmplEmbeddedNote     string
	EmbeddedNoteTemplate = tmpl.MustCompile(&EmbeddedNotePage{})
)

type EmbeddedNotePage struct {
	Content   template.HTML
	CreatedAt string
	Metadata  sdk.ProfileMetadata
	Npub      string
	NpubShort string
	SeenOn    []string
	Subject   string
	Url       string
}

func (*EmbeddedNotePage) TemplateText() string { return tmplEmbeddedNote }

var (
	//go:embed templates/profile.html
	tmplProfile     string
	ProfileTemplate = tmpl.MustCompile(&ProfilePage{})
)

type ProfilePage struct {
	HeadCommonPartial `tmpl:"head_common"`
	TopPartial        `tmpl:"top"`
	DetailsPartial    `tmpl:"details"`
	ClientsPartial    `tmpl:"clients"`
	FooterPartial     `tmpl:"footer"`

	AuthorRelays               []string
	Content                    string
	CreatedAt                  string
	Domain                     string
	LastNotes                  []EnhancedEvent
	Metadata                   sdk.ProfileMetadata
	NormalizedAuthorWebsiteURL string
	RenderedAuthorAboutText    template.HTML
	Nevent                     string
	Npub                       string
	Nprofile                   string
	IsReply                    string
	Proxy                      string
	Title                      string
}

func (*ProfilePage) TemplateText() string { return tmplProfile }

var (
	//go:embed templates/embedded_profile.html
	tmplEmbeddedProfile     string
	EmbeddedProfileTemplate = tmpl.MustCompile(&EmbeddedProfilePage{})
)

type EmbeddedProfilePage struct {
	AuthorRelays               []string
	Content                    string
	CreatedAt                  string
	Domain                     string
	Metadata                   sdk.ProfileMetadata
	NormalizedAuthorWebsiteURL string
	RenderedAuthorAboutText    template.HTML
	Nevent                     string
	Npub                       string
	Nprofile                   string
	Proxy                      string
	Title                      string
}

func (*EmbeddedProfilePage) TemplateText() string { return tmplEmbeddedProfile }

var (
	//go:embed templates/file_metadata.html
	tmplFileMetadata     string
	FileMetadataTemplate = tmpl.MustCompile(&FileMetadataPage{})
)

type FileMetadataPage struct {
	OpenGraphPartial  `tmpl:"opengraph"`
	HeadCommonPartial `tmpl:"head_common"`
	TopPartial        `tmpl:"top"`
	DetailsPartial    `tmpl:"details"`
	ClientsPartial    `tmpl:"clients"`
	FooterPartial     `tmpl:"footer"`

	Content          template.HTML
	CreatedAt        string
	Metadata         sdk.ProfileMetadata
	Npub             string
	NpubShort        string
	ParentLink       template.HTML
	SeenOn           []string
	Style            Style
	Subject          string
	TitleizedContent string
	Alt              string

	FileMetadata Kind1063Metadata
	IsImage      bool
	IsVideo      bool
}

func (*FileMetadataPage) TemplateText() string { return tmplFileMetadata }

var (
	//go:embed templates/live_event.html
	tmplLiveEvent     string
	LiveEventTemplate = tmpl.MustCompile(&LiveEventPage{})
)

type LiveEventPage struct {
	OpenGraphPartial  `tmpl:"opengraph"`
	HeadCommonPartial `tmpl:"head_common"`
	TopPartial        `tmpl:"top"`
	DetailsPartial    `tmpl:"details"`
	ClientsPartial    `tmpl:"clients"`
	FooterPartial     `tmpl:"footer"`

	Content          template.HTML
	CreatedAt        string
	Metadata         sdk.ProfileMetadata
	Npub             string
	NpubShort        string
	ParentLink       template.HTML
	SeenOn           []string
	Style            Style
	Subject          string
	TitleizedContent string
	Alt              string

	LiveEvent Kind30311Metadata
}

func (*LiveEventPage) TemplateText() string { return tmplLiveEvent }

var (
	//go:embed templates/live_event_message.html
	tmplLiveEventMessage     string
	LiveEventMessageTemplate = tmpl.MustCompile(&LiveEventMessagePage{})
)

type LiveEventMessagePage struct {
	OpenGraphPartial  `tmpl:"opengraph"`
	HeadCommonPartial `tmpl:"head_common"`
	TopPartial        `tmpl:"top"`
	DetailsPartial    `tmpl:"details"`
	ClientsPartial    `tmpl:"clients"`
	FooterPartial     `tmpl:"footer"`

	Content          template.HTML
	CreatedAt        string
	Metadata         sdk.ProfileMetadata
	Npub             string
	NpubShort        string
	ParentLink       template.HTML
	SeenOn           []string
	Style            Style
	Subject          string
	TitleizedContent string
	Alt              string

	LiveEventMessage Kind1311Metadata
}

func (*LiveEventMessagePage) TemplateText() string { return tmplLiveEventMessage }

var (
	//go:embed templates/relay.html
	tmplRelay     string
	RelayTemplate = tmpl.MustCompile(&RelayPage{})
)

type RelayPage struct {
	HeadCommonPartial `tmpl:"head_common"`
	TopPartial        `tmpl:"top"`
	ClientsPartial    `tmpl:"clients"`
	FooterPartial     `tmpl:"footer"`

	Info       *nip11.RelayInformationDocument
	Hostname   string
	Proxy      string
	LastNotes  []EnhancedEvent
	ModifiedAt string
}

func (*RelayPage) TemplateText() string { return tmplRelay }

var (
	//go:embed templates/sitemap.xml
	tmplSitemap     string
	SitemapTemplate = tmpl.MustCompile(&SitemapPage{})
)

type SitemapPage struct {
	Host       string
	ModifiedAt string

	// for the profile sitemap
	Npub string

	// for the relay sitemap
	RelayHostname string
	Info          *nip11.RelayInformationDocument

	// for the profile and relay sitemaps
	LastNotes []EnhancedEvent

	// for the archive sitemap
	PathPrefix string
	Data       []string
}

func (*SitemapPage) TemplateText() string { return tmplSitemap }

var (
	//go:embed templates/rss.xml
	tmplRSS     string
	RSSTemplate = tmpl.MustCompile(&RSSPage{})
)

type RSSPage struct {
	Host       string
	ModifiedAt string
	Title      string

	// for the profile RSS
	Npub     string
	Metadata sdk.ProfileMetadata

	// for the relay RSS
	RelayHostname string
	Info          *nip11.RelayInformationDocument

	// for the profile and relay RSSs
	LastNotes []EnhancedEvent

	// for the archive RSS
	PathPrefix string
	Data       []string
}

func (*RSSPage) TemplateText() string { return tmplRSS }

var (
	//go:embed templates/error.html
	tmplError     string
	ErrorTemplate = tmpl.MustCompile(&ErrorPage{})
)

type ErrorPage struct {
	HeadCommonPartial `tmpl:"head_common"`
	TopPartial        `tmpl:"top"`
	FooterPartial     `tmpl:"footer"`
	Message           template.HTML
	Errors            string
}

func (e *ErrorPage) TemplateText() string {
	if e.Message != "" {
		return tmplError
	}
	switch {
	case strings.Contains(e.Errors, "invalid checksum"):
		e.Message = "It looks like you entered an invalid event code.<br> Check if you copied it fully, a good idea is compare the first and the last characters."
	case strings.Contains(e.Errors, "couldn't find this"):
		e.Message = "Can't find the event in the relays. Try getting an `nevent1` code with relay hints."
	case strings.Contains(e.Errors, "invalid bech32 string length"), strings.Contains(e.Errors, "invalid separator"):
		e.Message = "You have typed a wrong event code, we need a URL path that starts with /npub1, /nprofile1, /nevent1, /naddr1, or something like /name@domain.com (or maybe just /domain.com) or an event id as hex (like /aef8b32af...)"
	default:
		e.Message = "I can't give any suggestions to solve the problem.<br> Please tag <a href='/dtonon.com'>daniele</a> and <a href='/fiatjaf.com'>fiatjaf</a> and complain!"
	}
	return tmplError
}
