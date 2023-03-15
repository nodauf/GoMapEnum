package linkedin

import (
	"GoMapEnum/src/logger"
	"GoMapEnum/src/utils"
)

var log *logger.Logger

// Options for linkedin module
type Options struct {
	Format     string
	Email      bool
	ExactMatch bool
	Cookie     string
	CompanyID  int32
	Company    string
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
		RecipeType   string `json:"$recipeType"`
		FeatureUnion struct {
			SimpleText struct {
				RecipeType    string        `json:"$recipeType"`
				AttributesV2  []interface{} `json:"attributesV2"`
				Text          string        `json:"text"`
				TextDirection string        `json:"textDirection"`
			} `json:"simpleText"`
		} `json:"featureUnion"`
		Results []struct {
			RecipeType               string `json:"$recipeType"`
			AddEntityToSearchHistory bool   `json:"addEntityToSearchHistory"`
			BadgeText                struct {
				RecipeType        string        `json:"$recipeType"`
				AccessibilityText string        `json:"accessibilityText"`
				AttributesV2      []interface{} `json:"attributesV2"`
				Text              string        `json:"text"`
				TextDirection     string        `json:"textDirection"`
			} `json:"badgeText"`
			EntityCustomTrackingInfo struct {
				RecipeType     string `json:"$recipeType"`
				MemberDistance string `json:"memberDistance"`
				NameMatch      bool   `json:"nameMatch"`
			} `json:"entityCustomTrackingInfo"`
			EntityUrn string `json:"entityUrn"`
			Image     struct {
				RecipeType        string `json:"$recipeType"`
				AccessibilityText string `json:"accessibilityText"`
				Attributes        []struct {
					RecipeType string `json:"$recipeType"`
					DetailData struct {
						ProfilePicture struct {
							AntiAbuseAnnotations []struct {
								AttributeID int64 `json:"attributeId"`
								EntityID    int64 `json:"entityId"`
							} `json:"$anti_abuse_annotations"`
							RecipeType     string `json:"$recipeType"`
							EntityUrn      string `json:"entityUrn"`
							ProfilePicture struct {
								RecipeType            string `json:"$recipeType"`
								DisplayImageReference struct {
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
								} `json:"displayImageReference"`
							} `json:"profilePicture"`
						} `json:"profilePicture"`
					} `json:"detailData"`
					DetailDataUnion struct {
						ProfilePicture string `json:"profilePicture"`
					} `json:"detailDataUnion"`
				} `json:"attributes"`
			} `json:"image"`
			Insights []struct {
				SimpleInsight struct {
					RecipeType string `json:"$recipeType"`
					Image      struct {
						RecipeType string `json:"$recipeType"`
						Attributes []struct {
							RecipeType string `json:"$recipeType"`
							DetailData struct {
								ProfilePicture struct {
									AntiAbuseAnnotations []struct {
										AttributeID int64 `json:"attributeId"`
										EntityID    int64 `json:"entityId"`
									} `json:"$anti_abuse_annotations"`
									RecipeType     string `json:"$recipeType"`
									EntityUrn      string `json:"entityUrn"`
									ProfilePicture struct {
										RecipeType            string `json:"$recipeType"`
										DisplayImageReference struct {
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
										} `json:"displayImageReference"`
									} `json:"profilePicture"`
								} `json:"profilePicture"`
							} `json:"detailData"`
							DetailDataUnion struct {
								ProfilePicture string `json:"profilePicture"`
							} `json:"detailDataUnion"`
						} `json:"attributes"`
					} `json:"image"`
					NavigationURL    string `json:"navigationUrl"`
					SearchActionType string `json:"searchActionType"`
					Title            struct {
						RecipeType    string        `json:"$recipeType"`
						AttributesV2  []interface{} `json:"attributesV2"`
						Text          string        `json:"text"`
						TextDirection string        `json:"textDirection"`
					} `json:"title"`
				} `json:"simpleInsight"`
			} `json:"insights"`
			LazyLoadedActions struct {
				RecipeType string `json:"$recipeType"`
				EntityUrn  string `json:"entityUrn"`
			} `json:"lazyLoadedActions"`
			LazyLoadedActionsUrn string `json:"lazyLoadedActionsUrn"`
			NavigationContext    struct {
				RecipeType string `json:"$recipeType"`
				URL        string `json:"url"`
			} `json:"navigationContext"`
			NavigationURL   string `json:"navigationUrl"`
			PrimarySubtitle struct {
				RecipeType    string        `json:"$recipeType"`
				AttributesV2  []interface{} `json:"attributesV2"`
				Text          string        `json:"text"`
				TextDirection string        `json:"textDirection"`
			} `json:"primarySubtitle"`
			SecondarySubtitle struct {
				RecipeType    string        `json:"$recipeType"`
				AttributesV2  []interface{} `json:"attributesV2"`
				Text          string        `json:"text"`
				TextDirection string        `json:"textDirection"`
			} `json:"secondarySubtitle"`
			Title struct {
				RecipeType        string        `json:"$recipeType"`
				AccessibilityText string        `json:"accessibilityText"`
				AttributesV2      []interface{} `json:"attributesV2"`
				Text              string        `json:"text"`
				TextDirection     string        `json:"textDirection"`
			} `json:"title"`
			TrackingID  string `json:"trackingId"`
			TrackingUrn string `json:"trackingUrn"`
		} `json:"results"`
	} `json:"elements"`
	Metadata struct {
		RecipeType         string `json:"$recipeType"`
		BlockedQuery       bool   `json:"blockedQuery"`
		FilterAppliedCount int64  `json:"filterAppliedCount"`
		PrimaryResultType  string `json:"primaryResultType"`
		SearchID           string `json:"searchId"`
		TotalResultCount   int64  `json:"totalResultCount"`
	} `json:"metadata"`
	Paging struct {
		RecipeType string        `json:"$recipeType"`
		Count      int64         `json:"count"`
		Links      []interface{} `json:"links"`
		Start      int64         `json:"start"`
		Total      int64         `json:"total"`
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
