package category

// Assets is the list of basic assets to seed into a new database.
var Assets = []Category{
	{
		Name:        "bank_accounts",
		Description: "Liquid assets held in various bank accounts like checking and savings accounts.",
		Type:        Asset,
	},
	{
		Name:        "cash",
		Description: "Physical cash that you have on hand or easily accessible for daily expenses.",
		Type:        Asset,
	},
	{
		Name:        "investments",
		Description: "Assets with potential returns over time, including stocks, bonds, and mutual funds.",
		Type:        Asset,
	},
	{
		Name:        "real_state",
		Description: "Properties you own, including your primary residence and rental properties.",
		Type:        Asset,
	},
	{
		Name:        "vehicles",
		Description: "Various types of vehicles you own, such as cars, motorcycles, or boats.",
		Type:        Asset,
	},
	{
		Name:        "personal_property",
		Description: "Valuable possessions like jewelry, artwork, and collectibles.",
		Type:        Asset,
	},
	{
		Name:        "intangible_assets",
		Description: "Assets with no physical form, such as patents, copyrights, and trademarks.",
		Type:        Asset,
	},
	{
		Name:        "business_assets",
		Description: "Assets related to a business you own, including equipment and inventory.",
		Type:        Asset,
	},
	{
		Name:        "retirement",
		Description: "Accounts specifically for retirement savings, like 401(k) and IRAs.",
		Type:        Asset,
	},
	{
		Name:        "digital_assets",
		Description: "Digital assets stored electronically, such as cryptocurrencies and domain names.",
		Type:        Asset,
	},
	{
		Name:        "others",
		Description: "A catch-all category for miscellaneous or uncommon assets.",
		Type:        Asset,
	},
}
