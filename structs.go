package main

type Transport struct{}

type Config struct {
	Query		   string
	Platform   	   string
	Language 	   string
	FolderTemplate string
	Goodies		   bool
	OutPath        string
	PlatformIDs	   string
}

type Args struct {
	Query    	   string `arg:"positional"`
	Platform 	   string `arg:"-p, --platform" help:"Item platform. windows/win, linux, mac/osx."`
	Language 	   string `arg:"-l, --language" help:"Item language.\n\t\t\t en, cz, de, es, fr, it, hu, nl, pl, pt, br, sv, tr, uk, ru, ar, ko, cn, jp, all."`
	FolderTemplate string `arg:"-t, --template" help:"Game folder naming template. title, titlePeriods.\n\t\t\t Ex: {{.title}} [GOG], {{.titlePeriods}}.GOG"`
	Goodies 	   bool	  `arg:"-g, --goodies" help:"Include goodies."`
	OutPath  	   string `arg:"-o, --out-path" help:"Where to download to. Path will be made if it doesn't already exist."`
}

type Cookie struct {
	Domain         string  `json:"domain"`
	ExpirationDate float64 `json:"expirationDate,omitempty"`
	HostOnly       bool    `json:"hostOnly"`
	HTTPOnly       bool    `json:"httpOnly"`
	Name           string  `json:"name"`
	Path           string  `json:"path"`
	SameSite       string  `json:"sameSite"`
	Secure         bool    `json:"secure"`
	Session        bool    `json:"session"`
	StoreID        string  `json:"storeId"`
	Value          string  `json:"value"`
	ID             int     `json:"id"`
}

type UserData struct {
	Country    string `json:"country"`
	Currencies []struct {
		Code   string `json:"code"`
		Symbol string `json:"symbol"`
	} `json:"currencies"`
	SelectedCurrency struct {
		Code   string `json:"code"`
		Symbol string `json:"symbol"`
	} `json:"selectedCurrency"`
	PreferredLanguage struct {
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"preferredLanguage"`
	RatingBrand string `json:"ratingBrand"`
	IsLoggedIn  bool   `json:"isLoggedIn"`
	Checksum    struct {
		Cart         interface{} `json:"cart"`
		Games        string      `json:"games"`
		Wishlist     interface{} `json:"wishlist"`
		ReviewsVotes interface{} `json:"reviews_votes"`
		GamesRating  interface{} `json:"games_rating"`
	} `json:"checksum"`
	Updates struct {
		Messages              int `json:"messages"`
		PendingFriendRequests int `json:"pendingFriendRequests"`
		UnreadChatMessages    int `json:"unreadChatMessages"`
		Products              int `json:"products"`
		Total                 int `json:"total"`
	} `json:"updates"`
	UserID        string `json:"userId"`
	Username      string `json:"username"`
	GalaxyUserID  string `json:"galaxyUserId"`
	Email         string `json:"email"`
	Avatar        string `json:"avatar"`
	WalletBalance struct {
		Currency string `json:"currency"`
		Amount   int    `json:"amount"`
	} `json:"walletBalance"`
	PurchasedItems struct {
		Games  int `json:"games"`
		Movies int `json:"movies"`
	} `json:"purchasedItems"`
	WishlistedItems           int           `json:"wishlistedItems"`
	Friends                   []interface{} `json:"friends"`
	PersonalizedProductPrices []interface{} `json:"personalizedProductPrices"`
	PersonalizedSeriesPrices  []interface{} `json:"personalizedSeriesPrices"`
}

type Product struct {
	IsGalaxyCompatible bool          `json:"isGalaxyCompatible"`
	Tags               []interface{} `json:"tags"`
	ID                 int           `json:"id"`
	Availability       struct {
		IsAvailable          bool `json:"isAvailable"`
		IsAvailableInAccount bool `json:"isAvailableInAccount"`
	} `json:"availability"`
	Title   string `json:"title"`
	Image   string `json:"image"`
	URL     string `json:"url"`
	WorksOn struct {
		Windows bool `json:"Windows"`
		Mac     bool `json:"Mac"`
		Linux   bool `json:"Linux"`
	} `json:"worksOn"`
	Category     string `json:"category"`
	Rating       int    `json:"rating"`
	IsComingSoon bool   `json:"isComingSoon"`
	IsMovie      bool   `json:"isMovie"`
	IsGame       bool   `json:"isGame"`
	Slug         string `json:"slug"`
	Updates      int    `json:"updates"`
	IsNew        bool   `json:"isNew"`
	DlcCount     int    `json:"dlcCount"`
	ReleaseDate  struct {
		Date         string `json:"date"`
		TimezoneType int    `json:"timezone_type"`
		Timezone     string `json:"timezone"`
	} `json:"releaseDate"`
	IsBaseProductMissing bool          `json:"isBaseProductMissing"`
	IsHidingDisabled     bool          `json:"isHidingDisabled"`
	IsInDevelopment      bool          `json:"isInDevelopment"`
	ExtraInfo            []interface{} `json:"extraInfo"`
	IsHidden             bool          `json:"isHidden"`
}

type Search struct {
	SortBy                     string      `json:"sortBy"`
	Page                       int         `json:"page"`
	TotalProducts              int         `json:"totalProducts"`
	TotalPages                 int         `json:"totalPages"`
	ProductsPerPage            int         `json:"productsPerPage"`
	ContentSystemCompatibility interface{} `json:"contentSystemCompatibility"`
	MoviesCount                int         `json:"moviesCount"`
	Tags                       []struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		ProductCount string `json:"productCount"`
	} `json:"tags"`
	Products []Product
	UpdatedProductsCount       int `json:"updatedProductsCount"`
	HiddenUpdatedProductsCount int `json:"hiddenUpdatedProductsCount"`
	AppliedFilters             struct {
		Tags interface{} `json:"tags"`
	} `json:"appliedFilters"`
	HasHiddenProducts bool `json:"hasHiddenProducts"`
}

type GameMeta struct {
	Title                  string          `json:"title"`
	BackgroundImage        string          `json:"backgroundImage"`
	CdKey                  string          `json:"cdKey"`
	TextInformation        string          `json:"textInformation"`
	Downloads              [][]interface{} `json:"downloads"`
	GalaxyDownloads        []interface{}   `json:"galaxyDownloads"`
	Extras          	   []*Download      `json:"extras"`
	Dlcs                   []interface{}   `json:"dlcs"`
	Tags                   []interface{}   `json:"tags"`
	IsPreOrder             bool            `json:"isPreOrder"`
	ReleaseTimestamp       int             `json:"releaseTimestamp"`
	Messages               []interface{}   `json:"messages"`
	Changelog              string          `json:"changelog"`
	ForumLink              string          `json:"forumLink"`
	IsBaseProductMissing   bool            `json:"isBaseProductMissing"`
	MissingBaseProduct     interface{}     `json:"missingBaseProduct"`
	Features               []interface{}   `json:"features"`
	SimpleGalaxyInstallers []struct {
		Path string `json:"path"`
		Os   string `json:"os"`
	} `json:"simpleGalaxyInstallers"`
}

type Download struct {
	ManualURL string `json:"manualUrl"`
	Name      string `json:"name"`
	Version   string `json:"version"`
	Date      string `json:"date"`
	Size      string `json:"size"`
	Type      string `json:"type"`
}

type WriteCounter struct {
	Total      int64
	TotalStr   string
	Downloaded int64
	Percentage int
	StartTime  int64
}