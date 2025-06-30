package linkedin

import (
	"GoMapEnum/src/logger"
	"GoMapEnum/src/utils"
)

var log *logger.Logger

// Options for linkedin module
type Options struct {
	Format        string
	Email         bool
	ExactMatch    bool
	CookieSession string
	CookieCSRF    string
	CompanyID     int32
	Company       string
	utils.BaseOptions
}

type linkedinListCompany struct {
	Elements []struct {
		RecipeType       string `json:"$recipeType"`
		AutoFill         bool   `json:"autoFill"`
		EntityLockupView struct {
			RecipeType string `json:"$recipeType"`
			Image      struct {
				RecipeType string `json:"$recipeType"`
				Attributes []struct {
					RecipeType      string `json:"$recipeType"`
					DetailDataUnion struct {
						Icon        string `json:"icon"`
						VectorImage struct {
							RecipeType string `json:"$recipeType"`
							Artifacts  []struct {
								RecipeType                    string `json:"$recipeType"`
								ExpiresAt                     int64  `json:"expiresAt"`
								FileIdentifyingURLPathSegment string `json:"fileIdentifyingUrlPathSegment"`
								Height                        int64  `json:"height"`
								Width                         int64  `json:"width"`
							} `json:"artifacts"`
							RootURL string `json:"rootUrl"`
						} `json:"vectorImage"`
					} `json:"detailDataUnion"`
				} `json:"attributes"`
			} `json:"image"`
			NavigationURL string `json:"navigationUrl"`
			Subtitle      struct {
				RecipeType    string `json:"$recipeType"`
				Text          string `json:"text"`
				TextDirection string `json:"textDirection"`
			} `json:"subtitle"`
			Title struct {
				RecipeType   string `json:"$recipeType"`
				AttributesV2 []struct {
					RecipeType      string `json:"$recipeType"`
					DetailDataUnion struct {
						Style string `json:"style"`
					} `json:"detailDataUnion"`
					Length int64 `json:"length"`
					Start  int64 `json:"start"`
				} `json:"attributesV2"`
				Text          string `json:"text"`
				TextDirection string `json:"textDirection"`
			} `json:"title"`
			TrackingID  string `json:"trackingId"`
			TrackingUrn string `json:"trackingUrn"`
		} `json:"entityLockupView"`
		Icon string `json:"icon"`
	} `json:"elements"`
	Metadata struct {
		RecipeType string `json:"$recipeType"`
		SearchID   string `json:"searchId"`
	} `json:"metadata"`
	Paging struct {
		RecipeType string        `json:"$recipeType"`
		Count      int64         `json:"count"`
		Links      []interface{} `json:"links"`
		Start      int64         `json:"start"`
	} `json:"paging"`
}

type linkedinListPeople struct {
	Elements []struct {
		_RecipeType string `json:"$recipeType"`
		Items       []struct {
			_RecipeType string `json:"$recipeType"`
			ItemUnion   struct {
				EntityResult struct {
					_RecipeType              string `json:"$recipeType"`
					AddEntityToSearchHistory bool   `json:"addEntityToSearchHistory"`
					BadgeText                struct {
						_RecipeType       string        `json:"$recipeType"`
						AccessibilityText string        `json:"accessibilityText"`
						AttributesV2      []interface{} `json:"attributesV2"`
						Text              string        `json:"text"`
						TextDirection     string        `json:"textDirection"`
					} `json:"badgeText"`
					EntityCustomTrackingInfo struct {
						_RecipeType    string `json:"$recipeType"`
						MemberDistance string `json:"memberDistance"`
						NameMatch      bool   `json:"nameMatch"`
					} `json:"entityCustomTrackingInfo"`
					EntityUrn string `json:"entityUrn"`
					Image     struct {
						_RecipeType       string `json:"$recipeType"`
						AccessibilityText string `json:"accessibilityText"`
						Attributes        []struct {
							_RecipeType     string `json:"$recipeType"`
							DetailDataUnion struct {
								NonEntityProfilePicture struct {
									_RecipeType string `json:"$recipeType"`
									Profile     struct {
										_RecipeType string `json:"$recipeType"`
										EntityUrn   string `json:"entityUrn"`
									} `json:"profile"`
									ProfileUrn string `json:"profileUrn"`
								} `json:"nonEntityProfilePicture"`
							} `json:"detailDataUnion"`
						} `json:"attributes"`
					} `json:"image"`
					Insights []struct {
						SimpleInsight struct {
							_RecipeType string `json:"$recipeType"`
							Image       struct {
								_RecipeType string `json:"$recipeType"`
								Attributes  []struct {
									_RecipeType     string `json:"$recipeType"`
									DetailDataUnion struct {
										NonEntityProfilePicture struct {
											_RecipeType string `json:"$recipeType"`
											Profile     struct {
												_RecipeType string `json:"$recipeType"`
												EntityUrn   string `json:"entityUrn"`
											} `json:"profile"`
											ProfileUrn  string `json:"profileUrn"`
											VectorImage struct {
												_RecipeType string `json:"$recipeType"`
												Artifacts   []struct {
													_RecipeType                   string `json:"$recipeType"`
													ExpiresAt                     int    `json:"expiresAt"`
													FileIdentifyingUrlPathSegment string `json:"fileIdentifyingUrlPathSegment"`
													Height                        int    `json:"height"`
													Width                         int    `json:"width"`
												} `json:"artifacts"`
												RootUrl string `json:"rootUrl"`
											} `json:"vectorImage"`
										} `json:"nonEntityProfilePicture"`
									} `json:"detailDataUnion"`
								} `json:"attributes"`
							} `json:"image"`
							NavigationUrl    string `json:"navigationUrl"`
							SearchActionType string `json:"searchActionType"`
							Title            struct {
								_RecipeType   string        `json:"$recipeType"`
								AttributesV2  []interface{} `json:"attributesV2"`
								Text          string        `json:"text"`
								TextDirection string        `json:"textDirection"`
							} `json:"title"`
						} `json:"simpleInsight"`
					} `json:"insights"`
					LazyLoadedActions struct {
						_RecipeType string `json:"$recipeType"`
						EntityUrn   string `json:"entityUrn"`
					} `json:"lazyLoadedActions"`
					LazyLoadedActionsUrn string `json:"lazyLoadedActionsUrn"`
					NavigationContext    struct {
						_RecipeType string `json:"$recipeType"`
						Url         string `json:"url"`
					} `json:"navigationContext"`
					NavigationUrl   string `json:"navigationUrl"`
					PrimarySubtitle struct {
						_RecipeType   string        `json:"$recipeType"`
						AttributesV2  []interface{} `json:"attributesV2"`
						Text          string        `json:"text"`
						TextDirection string        `json:"textDirection"`
					} `json:"primarySubtitle"`
					SecondarySubtitle struct {
						_RecipeType   string        `json:"$recipeType"`
						AttributesV2  []interface{} `json:"attributesV2"`
						Text          string        `json:"text"`
						TextDirection string        `json:"textDirection"`
					} `json:"secondarySubtitle"`
					Title struct {
						_RecipeType       string        `json:"$recipeType"`
						AccessibilityText string        `json:"accessibilityText"`
						AttributesV2      []interface{} `json:"attributesV2"`
						Text              string        `json:"text"`
						TextDirection     string        `json:"textDirection"`
					} `json:"title"`
					TrackingId  string `json:"trackingId"`
					TrackingUrn string `json:"trackingUrn"`
				} `json:"entityResult"`
			} `json:"itemUnion"`
			Position int `json:"position"`
		} `json:"items"`
		Position   int    `json:"position"`
		TrackingId string `json:"trackingId"`
	} `json:"elements"`
	Metadata struct {
		_RecipeType        string `json:"$recipeType"`
		BlockedQuery       bool   `json:"blockedQuery"`
		FilterAppliedCount int    `json:"filterAppliedCount"`
		PrimaryResultType  string `json:"primaryResultType"`
		SearchId           string `json:"searchId"`
		TotalResultCount   int    `json:"totalResultCount"`
	} `json:"metadata"`
	Paging struct {
		_RecipeType string        `json:"$recipeType"`
		Count       int           `json:"count"`
		Links       []interface{} `json:"links"`
		Start       int           `json:"start"`
		Total       int           `json:"total"`
	} `json:"paging"`
}

type linkedinGetCompany struct {
	BasicCompanyInfo struct {
		FollowingInfo struct {
			DashFollowingStateUrn string `json:"dashFollowingStateUrn"`
			EntityUrn             string `json:"entityUrn"`
			Following             bool   `json:"following"`
			FollowingType         string `json:"followingType"`
			TrackingUrn           string `json:"trackingUrn"`
		} `json:"followingInfo"`
		Headquarters string `json:"headquarters"`
		MiniCompany  struct {
			Active         bool   `json:"active"`
			DashCompanyUrn string `json:"dashCompanyUrn"`
			EntityUrn      string `json:"entityUrn"`
			Logo           struct {
				Com_linkedin_common_VectorImage struct {
					Artifacts []struct {
						ExpiresAt                     int64  `json:"expiresAt"`
						FileIdentifyingURLPathSegment string `json:"fileIdentifyingUrlPathSegment"`
						Height                        int64  `json:"height"`
						Width                         int64  `json:"width"`
					} `json:"artifacts"`
					RootURL string `json:"rootUrl"`
				} `json:"com.linkedin.common.VectorImage"`
			} `json:"logo"`
			Name          string `json:"name"`
			ObjectUrn     string `json:"objectUrn"`
			Showcase      bool   `json:"showcase"`
			TrackingID    string `json:"trackingId"`
			UniversalName string `json:"universalName"`
		} `json:"miniCompany"`
	} `json:"basicCompanyInfo"`
	CompanyType        string `json:"companyType"`
	Description        string `json:"description"`
	EmployeeCountRange string `json:"employeeCountRange"`
	EntityInfo         struct {
		ObjectUrn  string `json:"objectUrn"`
		TrackingID string `json:"trackingId"`
	} `json:"entityInfo"`
	EntityUrn   string `json:"entityUrn"`
	FoundedDate struct {
		Year int64 `json:"year"`
	} `json:"foundedDate"`
	Industries  []string `json:"industries"`
	Specialties []string `json:"specialties"`
	WebsiteURL  string   `json:"websiteUrl"`
}
