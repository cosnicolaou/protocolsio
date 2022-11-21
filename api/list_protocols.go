// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    protocolsV3, err := UnmarshalProtocolsV3(bytes)
//    bytes, err = protocolsV3.Marshal()

package api

import (
	"bytes"
	"encoding/json"
	"errors"
)

func UnmarshalProtocolsV3(data []byte) (ProtocolsV3, error) {
	var r ProtocolsV3
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ProtocolsV3) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type ProtocolsV3 struct {
	Extras       Extras     `json:"extras"`
	Items        []Item     `json:"items"`
	Pagination   Pagination `json:"pagination"`
	StatusCode   int64      `json:"status_code"`
	Total        int64      `json:"total"`
	TotalPages   int64      `json:"total_pages"`
	TotalResults int64      `json:"total_results"`
}

type Extras struct {
	The2 []The2 `json:"2"`
}

type The2 struct {
	ID          int64            `json:"id"`
	URI         string           `json:"uri"`
	Title       string           `json:"title"`
	Image       PlaceholderClass `json:"image"`
	TechSupport TechSupport      `json:"tech_support"`
	IsMember    int64            `json:"is_member"`
	Description *string          `json:"description,omitempty"`
	Request     *Request         `json:"request,omitempty"`
}

type PlaceholderClass struct {
	Source      string `json:"source"`
	Placeholder string `json:"placeholder"`
}

type Request struct {
	ID                  int64            `json:"id"`
	URI                 string           `json:"uri"`
	Title               string           `json:"title"`
	Image               PlaceholderClass `json:"image"`
	TechSupport         TechSupport      `json:"tech_support"`
	IsMember            int64            `json:"is_member"`
	Description         interface{}      `json:"description"`
	ResearchInterests   interface{}      `json:"research_interests"`
	Website             interface{}      `json:"website"`
	Location            interface{}      `json:"location"`
	Affiliation         interface{}      `json:"affiliation"`
	Status              Status           `json:"status"`
	Stats               RequestStats     `json:"stats"`
	UserStatus          UserStatus       `json:"user_status"`
	JoinLink            interface{}      `json:"join_link"`
	Token               interface{}      `json:"token"`
	Owner               Owner            `json:"owner"`
	IsProtocolRequested int64            `json:"is_protocol_requested"`
	IsGroupRequested    int64            `json:"is_group_requested"`
	IsMy                bool             `json:"is_my"`
	IsRequest           bool             `json:"is_request"`
	IsConfirmed         int64            `json:"is_confirmed"`
	IsDeclined          int64            `json:"is_declined"`
	Requester           Requester        `json:"requester"`
	Protocol            Protocol         `json:"protocol"`
	CreatedOn           int64            `json:"created_on"`
	ResolveOn           int64            `json:"resolve_on"`
	ResolvedUser        ResolvedUser     `json:"resolved_user"`
	Shared              bool             `json:"shared"`
}

type Owner struct {
	Name              string           `json:"name"`
	Affiliation       interface{}      `json:"affiliation"`
	Affiliations      []interface{}    `json:"affiliations"`
	Username          interface{}      `json:"username"`
	Note              interface{}      `json:"note"`
	Link              interface{}      `json:"link"`
	Image             PlaceholderClass `json:"image"`
	Badges            []interface{}    `json:"badges"`
	Verified          int64            `json:"verified"`
	IsVerifiedUser    bool             `json:"is_verified_user"`
	ResearchInterests interface{}      `json:"research_interests"`
	BlockedByYou      bool             `json:"blocked_by_you"`
	BlockedYou        bool             `json:"blocked_you"`
	HideFollowing     bool             `json:"hide_following"`
}

type Protocol struct {
	ID                int64            `json:"id"`
	Title             string           `json:"title"`
	TitleHTML         interface{}      `json:"title_html"`
	Image             PlaceholderClass `json:"image"`
	Doi               interface{}      `json:"doi"`
	DoiStatus         int64            `json:"doi_status"`
	URI               string           `json:"uri"`
	TypeID            int64            `json:"type_id"`
	TemplateID        int64            `json:"template_id"`
	PublishedOn       interface{}      `json:"published_on"`
	Stats             ProtocolStats    `json:"stats"`
	ParentProtocols   []interface{}    `json:"parent_protocols"`
	ParentCollections []interface{}    `json:"parent_collections"`
	CitedProtocols    []interface{}    `json:"cited_protocols"`
}

type ProtocolStats struct {
	NumberOfViews     int64 `json:"number_of_views"`
	NumberOfSteps     int64 `json:"number_of_steps"`
	NumberOfBookmarks int64 `json:"number_of_bookmarks"`
	NumberOfComments  int64 `json:"number_of_comments"`
	NumberOfExports   int64 `json:"number_of_exports"`
	NumberOfRuns      int64 `json:"number_of_runs"`
	NumberOfVotes     int64 `json:"number_of_votes"`
	IsVoted           int64 `json:"is_voted"`
}

type Requester struct {
	Name           string      `json:"name"`
	Affiliation    interface{} `json:"affiliation"`
	AffiliationURL interface{} `json:"affiliation_url"`
	Username       interface{} `json:"username"`
	Link           interface{} `json:"link"`
}

type ResolvedUser struct {
	Name              string            `json:"name"`
	Affiliation       interface{}       `json:"affiliation"`
	Affiliations      []interface{}     `json:"affiliations"`
	Username          interface{}       `json:"username"`
	Note              interface{}       `json:"note"`
	Link              interface{}       `json:"link"`
	Image             ResolvedUserImage `json:"image"`
	Badges            []interface{}     `json:"badges"`
	Verified          int64             `json:"verified"`
	IsVerifiedUser    bool              `json:"is_verified_user"`
	ResearchInterests interface{}       `json:"research_interests"`
	BlockedByYou      bool              `json:"blocked_by_you"`
	BlockedYou        bool              `json:"blocked_you"`
	HideFollowing     bool              `json:"hide_following"`
}

type ResolvedUserImage struct {
	Source      *Placeholder `json:"source"`
	Placeholder *Placeholder `json:"placeholder"`
}

type RequestStats struct {
	Files             []interface{} `json:"files"`
	TotalMembers      int64         `json:"total_members"`
	TotalFollowers    int64         `json:"total_followers"`
	TotalChildGroups  int64         `json:"total_child_groups"`
	TotalParentGroups int64         `json:"total_parent_groups"`
	HasCollaborations int64         `json:"has_collaborations"`
}

type Status struct {
	IsVisible   bool  `json:"is_visible"`
	AccessLevel int64 `json:"access_level"`
}

type TechSupport struct {
	Email       interface{} `json:"email"`
	Phone       interface{} `json:"phone"`
	HideContact int64       `json:"hide_contact"`
	UseEmail    int64       `json:"use_email"`
	URL         interface{} `json:"url"`
}

type UserStatus struct {
	IsMember    bool `json:"is_member"`
	IsConfirmed bool `json:"is_confirmed"`
	IsInvited   bool `json:"is_invited"`
	IsOwner     bool `json:"is_owner"`
	IsAdmin     bool `json:"is_admin"`
	IsFollowing bool `json:"is_following"`
}

type Item struct {
	ID                  int64            `json:"id"`
	Title               string           `json:"title"`
	TitleHTML           string           `json:"title_html"`
	Image               PlaceholderClass `json:"image"`
	Doi                 *string          `json:"doi"`
	DoiStatus           int64            `json:"doi_status"`
	URI                 string           `json:"uri"`
	TypeID              int64            `json:"type_id"`
	TemplateID          int64            `json:"template_id"`
	PublishedOn         int64            `json:"published_on"`
	Stats               ProtocolStats    `json:"stats"`
	ParentProtocols     []interface{}    `json:"parent_protocols"`
	ParentCollections   []interface{}    `json:"parent_collections"`
	CitedProtocols      []interface{}    `json:"cited_protocols"`
	VersionID           int64            `json:"version_id"`
	VersionData         VersionData      `json:"version_data"`
	CreatedOn           int64            `json:"created_on"`
	ModifiedOn          interface{}      `json:"modified_on"`
	Categories          interface{}      `json:"categories"`
	Public              int64            `json:"public"`
	IsUnlisted          int64            `json:"is_unlisted"`
	Creator             Creator          `json:"creator"`
	Journal             interface{}      `json:"journal"`
	JournalName         interface{}      `json:"journal_name"`
	JournalLink         interface{}      `json:"journal_link"`
	ArticleCitation     interface{}      `json:"article_citation"`
	HasVersions         int64            `json:"has_versions"`
	Link                string           `json:"link"`
	TotalCollections    int64            `json:"total_collections"`
	NumberOfSteps       int64            `json:"number_of_steps"`
	Authors             []Creator        `json:"authors"`
	Versions            []Version        `json:"versions"`
	Groups              []The2           `json:"groups"`
	IsOwner             int64            `json:"is_owner"`
	HasSubprotocols     int64            `json:"has_subprotocols"`
	IsSubprotocol       int64            `json:"is_subprotocol"`
	IsBookmarked        int64            `json:"is_bookmarked"`
	CanClaimAuthorship  int64            `json:"can_claim_authorship"`
	CanAcceptAuthorship int64            `json:"can_accept_authorship"`
	CanBeCopied         int64            `json:"can_be_copied"`
	CanRemoveFork       int64            `json:"can_remove_fork"`
	ForkID              interface{}      `json:"fork_id"`
	URL                 string           `json:"url"`
	ForksCount          ForksCount       `json:"forks_count"`
	Access              map[string]int64 `json:"access"`
}

type Creator struct {
	Name              string        `json:"name"`
	Affiliation       *string       `json:"affiliation"`
	Affiliations      []Affiliation `json:"affiliations"`
	Username          string        `json:"username"`
	Note              *string       `json:"note"`
	Link              *string       `json:"link"`
	Image             CreatorImage  `json:"image"`
	Badges            []Badge       `json:"badges"`
	Verified          int64         `json:"verified"`
	IsVerifiedUser    bool          `json:"is_verified_user"`
	ResearchInterests interface{}   `json:"research_interests"`
	BlockedByYou      bool          `json:"blocked_by_you"`
	BlockedYou        bool          `json:"blocked_you"`
	HideFollowing     bool          `json:"hide_following"`
}

type Affiliation struct {
	Affiliation *string `json:"affiliation"`
	URL         *string `json:"url"`
	IsDefault   int64   `json:"is_default"`
}

type Badge struct {
	ID    int64             `json:"id"`
	Image ResolvedUserImage `json:"image"`
	Name  string            `json:"name"`
}

type CreatorImage struct {
	Source      string `json:"source"`
	Placeholder string `json:"placeholder"`
	WebpSource  string `json:"webp_source"`
}

type ForksCount struct {
	Private int64 `json:"private"`
	Public  int64 `json:"public"`
}

type VersionData struct {
	ID                      int64   `json:"id"`
	Code                    string  `json:"code"`
	ParentID                int64   `json:"parent_id"`
	ParentURI               *string `json:"parent_uri"`
	HasPendingMergeRequest  int64   `json:"has_pending_merge_request"`
	HasApprovedMergeRequest int64   `json:"has_approved_merge_request"`
}

type Version struct {
	ID                int64             `json:"id"`
	Title             string            `json:"title"`
	TitleHTML         string            `json:"title_html"`
	Image             ResolvedUserImage `json:"image"`
	Doi               *string           `json:"doi"`
	DoiStatus         int64             `json:"doi_status"`
	URI               string            `json:"uri"`
	TypeID            int64             `json:"type_id"`
	TemplateID        int64             `json:"template_id"`
	PublishedOn       int64             `json:"published_on"`
	Stats             ProtocolStats     `json:"stats"`
	ParentProtocols   []interface{}     `json:"parent_protocols"`
	ParentCollections []interface{}     `json:"parent_collections"`
	CitedProtocols    []interface{}     `json:"cited_protocols"`
	VersionID         int64             `json:"version_id"`
	VersionData       VersionData       `json:"version_data"`
	CreatedOn         int64             `json:"created_on"`
	ModifiedOn        int64             `json:"modified_on"`
	Categories        interface{}       `json:"categories"`
	Public            int64             `json:"public"`
	IsUnlisted        int64             `json:"is_unlisted"`
	Creator           Creator           `json:"creator"`
	Journal           interface{}       `json:"journal"`
	JournalName       interface{}       `json:"journal_name"`
	JournalLink       interface{}       `json:"journal_link"`
	ArticleCitation   interface{}       `json:"article_citation"`
}

type Pagination struct {
	CurrentPage  int64       `json:"current_page"`
	TotalPages   int64       `json:"total_pages"`
	TotalResults int64       `json:"total_results"`
	NextPage     string      `json:"next_page"`
	PrevPage     interface{} `json:"prev_page"`
	PageSize     int64       `json:"page_size"`
	First        int64       `json:"first"`
	Last         int64       `json:"last"`
	ChangedOn    interface{} `json:"changed_on"`
}

type Placeholder struct {
	PlaceholderClass *PlaceholderClass
	String           *string
}

func (x *Placeholder) UnmarshalJSON(data []byte) error {
	x.PlaceholderClass = nil
	var c PlaceholderClass
	object, err := unmarshalUnion(data, nil, nil, nil, &x.String, false, nil, true, &c, false, nil, false, nil, true)
	if err != nil {
		return err
	}
	if object {
		x.PlaceholderClass = &c
	}
	return nil
}

func (x *Placeholder) MarshalJSON() ([]byte, error) {
	return marshalUnion(nil, nil, nil, x.String, false, nil, x.PlaceholderClass != nil, x.PlaceholderClass, false, nil, false, nil, true)
}

func unmarshalUnion(data []byte, pi **int64, pf **float64, pb **bool, ps **string, haveArray bool, pa interface{}, haveObject bool, pc interface{}, haveMap bool, pm interface{}, haveEnum bool, pe interface{}, nullable bool) (bool, error) {
	if pi != nil {
		*pi = nil
	}
	if pf != nil {
		*pf = nil
	}
	if pb != nil {
		*pb = nil
	}
	if ps != nil {
		*ps = nil
	}

	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	tok, err := dec.Token()
	if err != nil {
		return false, err
	}

	switch v := tok.(type) {
	case json.Number:
		if pi != nil {
			i, err := v.Int64()
			if err == nil {
				*pi = &i
				return false, nil
			}
		}
		if pf != nil {
			f, err := v.Float64()
			if err == nil {
				*pf = &f
				return false, nil
			}
			return false, errors.New("Unparsable number")
		}
		return false, errors.New("Union does not contain number")
	case float64:
		return false, errors.New("Decoder should not return float64")
	case bool:
		if pb != nil {
			*pb = &v
			return false, nil
		}
		return false, errors.New("Union does not contain bool")
	case string:
		if haveEnum {
			return false, json.Unmarshal(data, pe)
		}
		if ps != nil {
			*ps = &v
			return false, nil
		}
		return false, errors.New("Union does not contain string")
	case nil:
		if nullable {
			return false, nil
		}
		return false, errors.New("Union does not contain null")
	case json.Delim:
		if v == '{' {
			if haveObject {
				return true, json.Unmarshal(data, pc)
			}
			if haveMap {
				return false, json.Unmarshal(data, pm)
			}
			return false, errors.New("Union does not contain object")
		}
		if v == '[' {
			if haveArray {
				return false, json.Unmarshal(data, pa)
			}
			return false, errors.New("Union does not contain array")
		}
		return false, errors.New("Cannot handle delimiter")
	}
	return false, errors.New("Cannot unmarshal union")

}

func marshalUnion(pi *int64, pf *float64, pb *bool, ps *string, haveArray bool, pa interface{}, haveObject bool, pc interface{}, haveMap bool, pm interface{}, haveEnum bool, pe interface{}, nullable bool) ([]byte, error) {
	if pi != nil {
		return json.Marshal(*pi)
	}
	if pf != nil {
		return json.Marshal(*pf)
	}
	if pb != nil {
		return json.Marshal(*pb)
	}
	if ps != nil {
		return json.Marshal(*ps)
	}
	if haveArray {
		return json.Marshal(pa)
	}
	if haveObject {
		return json.Marshal(pc)
	}
	if haveMap {
		return json.Marshal(pm)
	}
	if haveEnum {
		return json.Marshal(pe)
	}
	if nullable {
		return json.Marshal(nil)
	}
	return nil, errors.New("Union must not be null")
}
