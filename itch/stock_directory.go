package itch

import (
	"encoding/binary"
	"fmt"
	"strings"
	"time"
)

var (
	StockMap  = make(map[string]uint16)
	Directory = make(map[uint16]StockDirectory)
)

type MarketCategory uint8
type FinancialStatusIndicator uint8
type IssueClassification uint8
type IssueSubType string
type Authenticity uint8

const (
	MKTCTG_NASDAQ_GLOBAL_SELECT MarketCategory = 'Q'
	MKTCTG_NASDAQ_GLOBAL        MarketCategory = 'G'
	MKTCTG_NASDAQ_CAPITAL       MarketCategory = 'S'
	MKTCTG_NYSE                 MarketCategory = 'N'
	MKTCTG_NYSE_AMERICAN        MarketCategory = 'A'
	MKTCTG_NYSE_ARCA            MarketCategory = 'P'
	MKTCTG_BATS_Z               MarketCategory = 'Z'
	MKTCTG_INVESTORS_EXCHANGE   MarketCategory = 'V'
	MKTCTG_NOT_AVAILABLE        MarketCategory = ' '

	FSI_DEFICIENT                       FinancialStatusIndicator = 'D'
	FSI_DELINQUENT                      FinancialStatusIndicator = 'E'
	FSI_BANKRUPT                        FinancialStatusIndicator = 'Q'
	FSI_SUSPENDED                       FinancialStatusIndicator = 'S'
	FSI_DEFICIENT_AND_BANKRUPT          FinancialStatusIndicator = 'G'
	FSI_DEFICIENT_AND_DELINQUENT        FinancialStatusIndicator = 'H'
	FSI_DELINQUENT_AND_BANKRUPT         FinancialStatusIndicator = 'J'
	FSI_DEFICIENT_DELINQUENT_BANKRUPT   FinancialStatusIndicator = 'K'
	FSI_CREATIONS_REDEMPTIONS_SUSPENDED FinancialStatusIndicator = 'C'
	FSI_NORMAL                          FinancialStatusIndicator = 'N'
	FSI_NOT_AVAILABLE                   FinancialStatusIndicator = ' '

	IC_AMERICAN_DEPOSITORY_SHARE  IssueClassification = 'A'
	IC_BOND                       IssueClassification = 'B'
	IC_COMMON_STOCK               IssueClassification = 'C'
	IC_DEPOSITORY_RECEIPT         IssueClassification = 'F'
	IC_144A                       IssueClassification = 'I'
	IC_LIMITED_PARTNERSHIP        IssueClassification = 'L'
	IC_NOTES                      IssueClassification = 'N'
	IC_ORDINARY_SHARE             IssueClassification = 'O'
	IC_PREFERRED_STOCK            IssueClassification = 'P'
	IC_OTHER_SECURITIES           IssueClassification = 'Q'
	IC_RIGHT                      IssueClassification = 'R'
	IC_SHARES_BENEFICIAL_INTEREST IssueClassification = 'S'
	IC_CONVERTIBLE_DEBENTURE      IssueClassification = 'T'
	IC_UNIT                       IssueClassification = 'U'
	IC_UNITS_BENIF                IssueClassification = 'V'
	IC_WARRANT                    IssueClassification = 'W'

	ICS_PREFERRED_TRUST_SECURITIES                       IssueSubType = "A"
	ICS_ALPHA_INDEX_ETN                                  IssueSubType = "AI"
	ICS_INDEX_BASED_DERIVATIVE                           IssueSubType = "B"
	ICS_COMMON_SHARES                                    IssueSubType = "C"
	ICS_COMMODITY_BASED_TRUST_SHARES                     IssueSubType = "CB"
	ICS_COMMODITY_FUTURES_TRUST_SHARES                   IssueSubType = "CF"
	ICS_COMMODITY_LINKED_SECURITIES                      IssueSubType = "CL"
	ICS_COMMODITY_INDEX_TRUST_SHARES                     IssueSubType = "CM"
	ICS_COLLATERALIZED_MORTGAGE_OBLIGATION               IssueSubType = "CO"
	ICS_CURRENCY_TRUST_SHARES                            IssueSubType = "CT"
	ICS_COMMODITY_CURRENCY_LINKED_SECURITIES             IssueSubType = "CU"
	ICS_CURRENCY_WARRANTS                                IssueSubType = "CW"
	ICS_GLOBAL_DEPOSITORY_SHARES                         IssueSubType = "D"
	ICS_ETF_PORTFOLIO_DEPOSITARY_RECEIPT                 IssueSubType = "E"
	ICS_EQUITY_GOLD_SHARES                               IssueSubType = "EG"
	ICS_ETN_EQUITY_INDEX_LINKED_SECURITIES               IssueSubType = "EI"
	ICS_NEXTSHARES_EXCHANGE_TRADED_MANAGED_FUND          IssueSubType = "EM"
	ICS_EXCHANGE_TRADED_NOTES                            IssueSubType = "EN"
	ICS_EQUITY_UNITS                                     IssueSubType = "EU"
	ICS_HOLDRS                                           IssueSubType = "F"
	ICS_ETN_FIXED_INCOME_LINKED_SECURITIES               IssueSubType = "FI"
	ICS_ETN_FUTURES_LINKED_SECURITIES                    IssueSubType = "FL"
	ICS_GLOBAL_SHARES                                    IssueSubType = "G"
	ICS_ETF_INDEX_FUND_SHARES                            IssueSubType = "I"
	ICS_INTEREST_RATE                                    IssueSubType = "IR"
	ICS_INDEX_WARRANT                                    IssueSubType = "IW"
	ICS_INDEX_LINKED_EXCHANGEABLE_NOTES                  IssueSubType = "IX"
	ICS_CORPORATE_BACKED_TRUST_SECURITY                  IssueSubType = "J"
	ICS_CONTINGENT_LITIGATION_RIGHT                      IssueSubType = "L"
	ICS_LLC                                              IssueSubType = "LL"
	ICS_EQUITY_BASED_DERIVATIVE                          IssueSubType = "M"
	ICS_MANAGED_FUND_SHARES                              IssueSubType = "MF"
	ICS_ETN_MULTI_FACTOR_INDEX_LINKED_SECURITIES         IssueSubType = "ML"
	ICS_MANAGED_TRUST_SECURITIES                         IssueSubType = "MT"
	ICS_NY_REGISTRY_SHARES                               IssueSubType = "N"
	ICS_OPEN_ENDED_MUTUAL_FUND                           IssueSubType = "O"
	ICS_PRIVATELY_HELD_SECURITY                          IssueSubType = "P"
	ICS_POISON_PILL                                      IssueSubType = "PP"
	ICS_PARTNERSHIP_UNITS                                IssueSubType = "PU"
	ICS_CLOSED_END_FUNDS                                 IssueSubType = "Q"
	ICS_REG_S                                            IssueSubType = "R"
	ICS_COMMODITY_REDEEMABLE_COMMODITY_LINKED_SECURITIES IssueSubType = "RC"
	ICS_ETN_REDEEMABLE_FUTURES_LINKED_SECURITIES         IssueSubType = "RF"
	ICS_REIT                                             IssueSubType = "RT"
	ICS_COMMODITY_REDEEMABLE_CURRENCY_LINKED_SECURITIES  IssueSubType = "RU"
	ICS_SEED                                             IssueSubType = "S"
	ICS_SPOT_RATE_CLOSING                                IssueSubType = "SC"
	ICS_SPOT_RATE_INTRADAY                               IssueSubType = "SI"
	ICS_TRACKING_STOCK                                   IssueSubType = "T"
	ICS_TRUST_CERTIFICATES                               IssueSubType = "TC"
	ICS_TRUST_UNITS                                      IssueSubType = "TU"
	ICS_PORTAL                                           IssueSubType = "U"
	ICS_CONTINGENT_VALUE_RIGHT                           IssueSubType = "V"
	ICS_TRUST_ISSUED_RECEIPTS                            IssueSubType = "W"
	ICS_WORLD_CURRENCY_OPTION                            IssueSubType = "WC"
	ICS_TRUST                                            IssueSubType = "X"
	ICS_OTHER                                            IssueSubType = "Y"
	ICS_NOT_APPLICABLE                                   IssueSubType = "Z"

	AUTHENTICITY_LIVE Authenticity = 'P'
	AUTHENTICITY_TEST Authenticity = 'T'
)

type StockDirectory struct {
	StockLocate                 uint16
	TrackingNumber              uint16
	Timestamp                   time.Duration
	Stock                       string
	MarketCategory              MarketCategory
	FinancialStatusIndicator    FinancialStatusIndicator
	RoundLotSize                uint32
	RoundLotsOnly               bool
	IssueClassification         IssueClassification
	IssueSubType                IssueSubType
	Authenticity                Authenticity
	ShortSaleThresholdIndicator string
	IpoFlag                     string
	LuldReferencePriceTier      string
	EtpFlag                     string
	EtpLeverageFactor           uint32
	InverseIndicator            bool
}

func MakeStockDirectory(data []byte) Message {
	locate := binary.BigEndian.Uint16(data[1:3])
	tracking := binary.BigEndian.Uint16(data[3:5])
	data[3] = 0
	data[4] = 0
	t := binary.BigEndian.Uint64(data[3:11])

	roundLotsOnly := false
	if data[25] == 'Y' {
		roundLotsOnly = true
	}

	inverseIndicator := false
	if data[38] == 'Y' {
		inverseIndicator = true
	}

	stock := strings.TrimSpace(string(data[11:19]))

	StockMap[stock] = locate

	Directory[locate] = StockDirectory{
		StockLocate:                 locate,
		TrackingNumber:              tracking,
		Timestamp:                   time.Duration(t),
		Stock:                       stock,
		MarketCategory:              MarketCategory(data[19]),
		FinancialStatusIndicator:    FinancialStatusIndicator(data[20]),
		RoundLotSize:                binary.BigEndian.Uint32(data[21:25]),
		RoundLotsOnly:               roundLotsOnly,
		IssueClassification:         IssueClassification(data[26]),
		IssueSubType:                IssueSubType(strings.TrimSpace(string(data[27:29]))),
		Authenticity:                Authenticity(data[29]),
		ShortSaleThresholdIndicator: string(data[30]),
		IpoFlag:                     string(data[31]),
		LuldReferencePriceTier:      string(data[32]),
		EtpFlag:                     string(data[33]),
		EtpLeverageFactor:           binary.BigEndian.Uint32(data[34:38]),
		InverseIndicator:            inverseIndicator,
	}

	return Directory[locate]
}

func (e StockDirectory) String() string {
	return fmt.Sprintf("[Stock Directory]\n"+
		"Stock Locate: %v\n"+
		"Tracking Number: %v\n"+
		"Timestamp: %v\n"+
		"Stock: %v\n"+
		"Market Category: %v\n"+
		"Financial Status Indicator: %v\n"+
		"Round Lot Size: %v\n"+
		"Round Lots Only: %v\n"+
		"Issue Classification: %v\n"+
		"Issue Sub-Type: %v\n"+
		"Authenticity: %v\n"+
		"Short Sale Threshold Indicator: %v\n"+
		"IPO Flag: %v\n"+
		"LULD Reference Price Tier: %v\n"+
		"ETP Flag: %v\n"+
		"ETP Leverage Factor: %v\n"+
		"Inverse Indicator: %v\n",
		e.StockLocate, e.TrackingNumber, e.Timestamp, e.Stock, e.MarketCategory, e.FinancialStatusIndicator, e.RoundLotSize,
		e.RoundLotsOnly, e.IssueClassification, e.IssueSubType, e.Authenticity, e.ShortSaleThresholdIndicator,
		e.IpoFlag, e.LuldReferencePriceTier, e.EtpFlag, e.EtpLeverageFactor, e.InverseIndicator,
	)
}

func (c MarketCategory) String() string {
	switch c {
	case MKTCTG_NASDAQ_GLOBAL_SELECT:
		return "Nasdaq Global Select Market"
	case MKTCTG_NASDAQ_GLOBAL:
		return "Nasdaq Global Market"
	case MKTCTG_NASDAQ_CAPITAL:
		return "Nasdaq Capital Market"
	case MKTCTG_NYSE:
		return "New York Stock Exchange (NYSE)"
	case MKTCTG_NYSE_AMERICAN:
		return "NYSE American"
	case MKTCTG_NYSE_ARCA:
		return "NYSE Arca"
	case MKTCTG_BATS_Z:
		return "BATS Z Exchange"
	case MKTCTG_INVESTORS_EXCHANGE:
		return "Investors' Exchange, LLC"
	case MKTCTG_NOT_AVAILABLE:
		return "Not Available"
	}

	return "Unknown Market Category"
}

func (i FinancialStatusIndicator) String() string {
	switch i {
	case FSI_DEFICIENT:
		return "Deficient"
	case FSI_DELINQUENT:
		return "Delinquent"
	case FSI_BANKRUPT:
		return "Bankrupt"
	case FSI_SUSPENDED:
		return "Suspended"
	case FSI_DEFICIENT_AND_BANKRUPT:
		return "Deficient and Bankrupt"
	case FSI_DEFICIENT_AND_DELINQUENT:
		return "Deficient and Delinquent"
	case FSI_DELINQUENT_AND_BANKRUPT:
		return "Delinquent and Bankrupt"
	case FSI_DEFICIENT_DELINQUENT_BANKRUPT:
		return "Deficient, Delinquent and Bankrupt"
	case FSI_CREATIONS_REDEMPTIONS_SUSPENDED:
		return "Creation and/or Redemptions Suspended"
	case FSI_NORMAL:
		return "Normal"
	case FSI_NOT_AVAILABLE:
		return "Not Available"
	}

	return "Unknown Financial Status Indicator"
}

func (c IssueClassification) String() string {
	switch c {
	case IC_AMERICAN_DEPOSITORY_SHARE:
		return "American Depository Share"
	case IC_BOND:
		return "Bond"
	case IC_COMMON_STOCK:
		return "Common Stock"
	case IC_DEPOSITORY_RECEIPT:
		return "Depository Receipt"
	case IC_144A:
		return "144A"
	case IC_LIMITED_PARTNERSHIP:
		return "Limited Partnership"
	case IC_NOTES:
		return "Notes"
	case IC_ORDINARY_SHARE:
		return "Ordinary Share"
	case IC_PREFERRED_STOCK:
		return "Preferred Stock"
	case IC_OTHER_SECURITIES:
		return "Other Securities"
	case IC_RIGHT:
		return "Right"
	case IC_SHARES_BENEFICIAL_INTEREST:
		return "Shares of Beneficial Interest"
	case IC_CONVERTIBLE_DEBENTURE:
		return "Convertible Debenture"
	case IC_UNIT:
		return "Unit"
	case IC_UNITS_BENIF:
		return "Units/Benif Int"
	case IC_WARRANT:
		return "Warrant"
	}

	return "Unknown Issue Classification"
}

func (i IssueSubType) String() string {
	switch i {
	case ICS_PREFERRED_TRUST_SECURITIES:
		return "Preferred Trust Securities"
	case ICS_ALPHA_INDEX_ETN:
		return "Alpha Index ETNs"
	case ICS_INDEX_BASED_DERIVATIVE:
		return "Index Based Derivative"
	case ICS_COMMON_SHARES:
		return "Common Shares"
	case ICS_COMMODITY_BASED_TRUST_SHARES:
		return "Commodity Based Trust Shares"
	case ICS_COMMODITY_FUTURES_TRUST_SHARES:
		return "Commodity Futures Trust Shares"
	case ICS_COMMODITY_LINKED_SECURITIES:
		return "Commodity Linked Securities"
	case ICS_COMMODITY_INDEX_TRUST_SHARES:
		return "Commodity Index Trust Shares"
	case ICS_COLLATERALIZED_MORTGAGE_OBLIGATION:
		return "Collateralized Mortage Obligation"
	case ICS_CURRENCY_TRUST_SHARES:
		return "Currency Trust Shares"
	case ICS_COMMODITY_CURRENCY_LINKED_SECURITIES:
		return "Commodity Currency Linked Securities"
	case ICS_CURRENCY_WARRANTS:
		return "Currency Warrants"
	case ICS_GLOBAL_DEPOSITORY_SHARES:
		return "Global Depository Shares"
	case ICS_ETF_PORTFOLIO_DEPOSITARY_RECEIPT:
		return "ETF Portfolio Depositary Receipt"
	case ICS_EQUITY_GOLD_SHARES:
		return "Equity Gold Shares"
	case ICS_ETN_EQUITY_INDEX_LINKED_SECURITIES:
		return "ETN Equity Index Linked Securities"
	case ICS_NEXTSHARES_EXCHANGE_TRADED_MANAGED_FUND:
		return "NextShares Exchanged Traded Managed Fund"
	case ICS_EXCHANGE_TRADED_NOTES:
		return "Exchange Traded Notes"
	case ICS_EQUITY_UNITS:
		return "Equity Units"
	case ICS_HOLDRS:
		return "HOLDRS"
	case ICS_ETN_FIXED_INCOME_LINKED_SECURITIES:
		return "ETN Fixed Income Linked Securities"
	case ICS_ETN_FUTURES_LINKED_SECURITIES:
		return "ETN Futures Linked Securities"
	case ICS_GLOBAL_SHARES:
		return "Global Shares"
	case ICS_ETF_INDEX_FUND_SHARES:
		return "ETF Index Fund Shares"
	case ICS_INTEREST_RATE:
		return "Interest Rate"
	case ICS_INDEX_WARRANT:
		return "Index Warrant"
	case ICS_INDEX_LINKED_EXCHANGEABLE_NOTES:
		return "Index Linked Exchangeable Notes"
	case ICS_CORPORATE_BACKED_TRUST_SECURITY:
		return "Corporate Backed Trust Security"
	case ICS_CONTINGENT_LITIGATION_RIGHT:
		return "Contingent Litigation Right"
	case ICS_LLC:
		return "Securities of a company set up as a Limited Liability Company (LLC)"
	case ICS_EQUITY_BASED_DERIVATIVE:
		return "Equity Based Derivative"
	case ICS_MANAGED_FUND_SHARES:
		return "Managed Fund Shares"
	case ICS_ETN_MULTI_FACTOR_INDEX_LINKED_SECURITIES:
		return "ETN Multi Factor Index Linked Securities"
	case ICS_MANAGED_TRUST_SECURITIES:
		return "Managed Trust Securities"
	case ICS_NY_REGISTRY_SHARES:
		return "NY Registry Shares"
	case ICS_OPEN_ENDED_MUTUAL_FUND:
		return "Open Ended Mutual Fund"
	case ICS_PRIVATELY_HELD_SECURITY:
		return "Privately Held Security"
	case ICS_POISON_PILL:
		return "Poison Pill"
	case ICS_PARTNERSHIP_UNITS:
		return "Partnership Units"
	case ICS_CLOSED_END_FUNDS:
		return "Closed End Funds"
	case ICS_REG_S:
		return "Reg S"
	case ICS_COMMODITY_REDEEMABLE_COMMODITY_LINKED_SECURITIES:
		return "Commodity Redeemable Commodity Linked Securities"
	case ICS_ETN_REDEEMABLE_FUTURES_LINKED_SECURITIES:
		return "ETN Redeemable Futures Linked Securities"
	case ICS_REIT:
		return "REIT"
	case ICS_COMMODITY_REDEEMABLE_CURRENCY_LINKED_SECURITIES:
		return "Commodity Redeemable Currency Linked Securities"
	case ICS_SEED:
		return "SEED"
	case ICS_SPOT_RATE_CLOSING:
		return "Spot Rate Closing"
	case ICS_SPOT_RATE_INTRADAY:
		return "Spot Rate Intraday"
	case ICS_TRACKING_STOCK:
		return "Tracking Stock"
	case ICS_TRUST_CERTIFICATES:
		return "Trust Certificates"
	case ICS_TRUST_UNITS:
		return "Trust Units"
	case ICS_PORTAL:
		return "Portal"
	case ICS_CONTINGENT_VALUE_RIGHT:
		return "Contingent Value Right"
	case ICS_TRUST_ISSUED_RECEIPTS:
		return "Trust Issued Receipts"
	case ICS_WORLD_CURRENCY_OPTION:
		return "World Currency Option"
	case ICS_TRUST:
		return "Trust"
	case ICS_OTHER:
		return "Other"
	case ICS_NOT_APPLICABLE:
		return "Not Applicable"
	}
	return "Unkown Issue-Sub Type"
}

func (a Authenticity) String() string {
	switch a {
	case AUTHENTICITY_LIVE:
		return "Live/Production"
	case AUTHENTICITY_TEST:
		return "Test"
	}

	return "Unkown Authenticity"
}
