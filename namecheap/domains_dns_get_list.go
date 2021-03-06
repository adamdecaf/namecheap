package namecheap

import (
	"encoding/xml"
	"fmt"
)

type DomainsDNSGetListResponse struct {
	XMLName *xml.Name `xml:"ApiResponse"`
	Errors  *[]struct {
		Message *string `xml:",chardata"`
		Number  *string `xml:"Number,attr"`
	} `xml:"Errors>Error"`
	CommandResponse *DomainsDNSGetListCommandResponse `xml:"CommandResponse"`
}

type DomainsDNSGetListCommandResponse struct {
	DomainDNSGetListResult *DomainDNSGetListResult `xml:"DomainDNSGetListResult"`
}

type DomainDNSGetListResult struct {
	Domain         *string   `xml:"Domain,attr"`
	IsUsingOurDNS  *bool     `xml:"IsUsingOurDNS,attr"`
	IsPremiumDNS   *bool     `xml:"IsPremiumDNS,attr"`
	IsUsingFreeDNS *bool     `xml:"IsUsingFreeDNS,attr"`
	Nameservers    *[]string `xml:"Nameserver"`
}

func (d DomainDNSGetListResult) String() string {
	return fmt.Sprintf("{Domain: %s, IsUsingOurDNS: %t, IsPremiumDNS: %t, IsUsingFreeDNS: %t, Nameservers: %v}",
		*d.Domain, *d.IsUsingOurDNS, *d.IsPremiumDNS, *d.IsUsingFreeDNS, *d.Nameservers,
	)
}

// GetList gets a list of DNS servers associated with the requested domain
//
// Namecheap doc: https://www.namecheap.com/support/api/methods/domains-dns/get-list/
func (dds *DomainsDNSService) GetList(domain string) (*DomainsDNSGetListCommandResponse, error) {
	var response DomainsDNSGetListResponse

	params := map[string]string{
		"Command": "namecheap.domains.dns.getList",
	}

	parsedDomain, err := ParseDomain(domain)
	if err != nil {
		return nil, err
	}

	params["SLD"] = parsedDomain.SLD
	params["TLD"] = parsedDomain.TLD

	req, err := dds.client.NewRequest(params)
	if err != nil {
		return nil, err
	}
	resp, err := dds.client.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	err = decodeBody(resp.Body, &response)
	if err != nil {
		return nil, err
	}
	if response.Errors != nil && len(*response.Errors) > 0 {
		apiErr := (*response.Errors)[0]
		return nil, fmt.Errorf("%s (%s)", *apiErr.Message, *apiErr.Number)
	}

	return response.CommandResponse, nil
}
