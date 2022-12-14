package linkapi

// Ex: "https://en.wikipedia.org/w/api.php?action=query&format=json&formatversion=2&gapcontinue=%22Boss%22_Tweed&gapfilterredir=nonredirects&gaplimit=max&gapnamespace=0&generator=allpages"
const AllLinksTemplate = `https://%s.wikipedia.org/w/api.php?action=query&generator=allpages&gapnamespace=0&gaplimit=max&format=json&formatversion=2&gapfilterredir=nonredirects`

type AllLinksStruct struct {
	Batchcomplete bool             `json:"batchcomplete"`
	Continue      AllLinksContinue `json:"continue"`
	Limits        AllLinksLimits   `json:"limits"`
	Query         AllLinksQuery    `json:"query"`
}
type AllLinksContinue struct {
	Gapcontinue string `json:"gapcontinue"`
	Continue    string `json:"continue"`
}
type AllLinksLimits struct {
	Allpages int `json:"allpages"`
}
type AllLinksPages struct {
	Pageid int    `json:"pageid"`
	Ns     int    `json:"ns"`
	Title  string `json:"title"`
}
type AllLinksQuery struct {
	Pages []AllLinksPages `json:"pages"`
}

// Ex: https://en.wikipedia.org/w/api.php?action=query&format=json&formatversion=2&plcontinue=37751%7C0%7CPainted_turtle&pllimit=max&plnamespace=0&prop=links&titles=Turtle
const LinksToTemplate = `https://%s.wikipedia.org/w/api.php?action=query&format=json&formatversion=2&pllimit=max&plnamespace=0&prop=links&titles=%s`

type LinksToStruct struct {
	Continue LinksToContinue `json:"continue"`
	Query    LinksToQuery    `json:"query"`
	Limits   LinksToLimits   `json:"limits"`
}
type LinksToContinue struct {
	Plcontinue string `json:"plcontinue"`
	Continue   string `json:"continue"`
}
type LinksToLinks struct {
	Ns    int    `json:"ns"`
	Title string `json:"title"`
}
type LinksToPages struct {
	Pageid int            `json:"pageid"`
	Ns     int            `json:"ns"`
	Title  string         `json:"title"`
	Links  []LinksToLinks `json:"links"`
}
type LinksToQuery struct {
	Pages []LinksToPages `json:"pages"`
}
type LinksToLimits struct {
	Links int `json:"links"`
}
