package utils

import (
	"encoding/json"
	"log"
	"strings"
)

const exclude = `[
  {
    "url": "https://advisors.wikimedia.org",
    "dbname": "advisorswiki",
    "code": "advisors",
    "lang": "advisors",
    "sitename": "Advisors",
    "private": true
  },
  {
    "url": "https://advisory.wikimedia.org",
    "dbname": "advisorywiki",
    "code": "advisory",
    "lang": "en",
    "sitename": "Advisory Board",
    "closed": true
  },
  {
    "url": "https://am.wikimedia.org",
    "dbname": "amwikimedia",
    "code": "amwikimedia",
    "lang": "hy",
    "sitename": "Վիքիմեդիա Հայաստան",
    "fishbowl": true
  },
  {
    "url": "https://api.wikimedia.org",
    "dbname": "apiportalwiki",
    "code": "apiportal",
    "lang": "en",
    "sitename": "API Portal"
  },
  {
    "url": "https://ar.wikimedia.org",
    "dbname": "arwikimedia",
    "code": "arwikimedia",
    "lang": "es",
    "sitename": "Wikimedia Argentina"
  },
  {
    "url": "https://arbcom-cs.wikipedia.org",
    "dbname": "arbcom_cswiki",
    "code": "arbcom-cs",
    "lang": "cs",
    "sitename": "Arbitrážní výbor",
    "private": true
  },
  {
    "url": "https://arbcom-de.wikipedia.org",
    "dbname": "arbcom_dewiki",
    "code": "arbcom-de",
    "lang": "de",
    "sitename": "Arbitration Committee",
    "private": true
  },
  {
    "url": "https://arbcom-en.wikipedia.org",
    "dbname": "arbcom_enwiki",
    "code": "arbcom-en",
    "lang": "en",
    "sitename": "Arbitration Committee",
    "private": true
  },
  {
    "url": "https://arbcom-fi.wikipedia.org",
    "dbname": "arbcom_fiwiki",
    "code": "arbcom-fi",
    "lang": "fi",
    "sitename": "Arbitration Committee",
    "private": true
  },
  {
    "url": "https://arbcom-nl.wikipedia.org",
    "dbname": "arbcom_nlwiki",
    "code": "arbcom-nl",
    "lang": "nl",
    "sitename": "Arbitration Committee",
    "private": true
  },
  {
    "url": "https://arbcom-ru.wikipedia.org",
    "dbname": "arbcom_ruwiki",
    "code": "arbcom-ru",
    "lang": "ru",
    "sitename": "Арбитражный комитет русской Википедии",
    "private": true
  },
  {
    "url": "https://auditcom.wikimedia.org",
    "dbname": "auditcomwiki",
    "code": "auditcom",
    "lang": "en",
    "sitename": "Audit Committee",
    "private": true
  },
  {
    "url": "https://bd.wikimedia.org",
    "dbname": "bdwikimedia",
    "code": "bdwikimedia",
    "lang": "bn",
    "sitename": "উইকিমিডিয়া বাংলাদেশ"
  },
  {
    "url": "https://be.wikimedia.org",
    "dbname": "bewikimedia",
    "code": "bewikimedia",
    "lang": "en",
    "sitename": "Wikimedia Belgium"
  },
  {
    "url": "https://beta.wikiversity.org",
    "dbname": "betawikiversity",
    "code": "betawikiversity",
    "lang": "en",
    "sitename": "Wikiversity"
  },
  {
    "url": "https://board.wikimedia.org",
    "dbname": "boardwiki",
    "code": "board",
    "lang": "board",
    "sitename": "Board",
    "private": true
  },
  {
    "url": "https://boardgovcom.wikimedia.org",
    "dbname": "boardgovcomwiki",
    "code": "boardgovcom",
    "lang": "boardgovcom",
    "sitename": "Board Governance Committee",
    "private": true
  },
  {
    "url": "https://br.wikimedia.org",
    "dbname": "brwikimedia",
    "code": "brwikimedia",
    "lang": "pt-BR",
    "sitename": "Wikimedia Brasil"
  },
  {
    "url": "https://ca.wikimedia.org",
    "dbname": "cawikimedia",
    "code": "cawikimedia",
    "lang": "en",
    "sitename": "Wikimedia Canada"
  },
  {
    "url": "https://chair.wikimedia.org",
    "dbname": "chairwiki",
    "code": "chair",
    "lang": "en",
    "sitename": "Wikimedia Board Chair",
    "private": true
  },
  {
    "url": "https://affcom.wikimedia.org",
    "dbname": "chapcomwiki",
    "code": "chapcom",
    "lang": "en",
    "sitename": "Affcom",
    "private": true
  },
  {
    "url": "https://checkuser.wikimedia.org",
    "dbname": "checkuserwiki",
    "code": "checkuser",
    "lang": "en",
    "sitename": "CheckUser Wiki",
    "private": true
  },
  {
    "url": "https://cn.wikimedia.org",
    "dbname": "cnwikimedia",
    "code": "cnwikimedia",
    "lang": "zh",
    "sitename": "中国维基媒体用户组",
    "fishbowl": true
  },
  {
    "url": "https://co.wikimedia.org",
    "dbname": "cowikimedia",
    "code": "cowikimedia",
    "lang": "es",
    "sitename": "Wikimedia Colombia"
  },
  {
    "url": "https://collab.wikimedia.org",
    "dbname": "collabwiki",
    "code": "collab",
    "lang": "en",
    "sitename": "Collab",
    "private": true
  },
  {
    "url": "https://commons.wikimedia.org",
    "dbname": "commonswiki",
    "code": "commons",
    "lang": "commons",
    "sitename": "Wikimedia Commons"
  },
  {
    "url": "https://dk.wikimedia.org",
    "dbname": "dkwikimedia",
    "code": "dkwikimedia",
    "lang": "da",
    "sitename": "Wikimedia Danmark"
  },
  {
    "url": "https://donate.wikimedia.org",
    "dbname": "donatewiki",
    "code": "donate",
    "lang": "en",
    "sitename": "Donate",
    "fishbowl": true
  },
  {
    "url": "https://ec.wikimedia.org",
    "dbname": "ecwikimedia",
    "code": "ecwikimedia",
    "lang": "es",
    "sitename": "Wikimedistas de Ecuador",
    "private": true
  },
  {
    "url": "https://electcom.wikimedia.org",
    "dbname": "electcomwiki",
    "code": "electcom",
    "lang": "electcom",
    "sitename": "Wikipedia",
    "private": true
  },
  {
    "url": "https://ee.wikimedia.org",
    "dbname": "etwikimedia",
    "code": "etwikimedia",
    "lang": "et",
    "sitename": "Wikimedia Eesti"
  },
  {
    "url": "https://exec.wikimedia.org",
    "dbname": "execwiki",
    "code": "exec",
    "lang": "en",
    "sitename": "Wikimedia Executive",
    "private": true
  },
  {
    "url": "https://fdc.wikimedia.org",
    "dbname": "fdcwiki",
    "code": "fdc",
    "lang": "en",
    "sitename": "Wikimedia FDC",
    "private": true
  },
  {
    "url": "https://fi.wikimedia.org",
    "dbname": "fiwikimedia",
    "code": "fiwikimedia",
    "lang": "fi",
    "sitename": "Wikimedia Suomi"
  },
  {
    "url": "https://foundation.wikimedia.org",
    "dbname": "foundationwiki",
    "code": "foundation",
    "lang": "foundation",
    "sitename": "Wikimedia Foundation Governance Wiki",
    "fishbowl": true
  },
  {
    "url": "https://ge.wikimedia.org",
    "dbname": "gewikimedia",
    "code": "gewikimedia",
    "lang": "ka",
    "sitename": "ვიკიმედიის მომხმარებელთა საქართველოს ჯგუფი",
    "fishbowl": true
  },
  {
    "url": "https://gr.wikimedia.org",
    "dbname": "grwikimedia",
    "code": "grwikimedia",
    "lang": "el",
    "sitename": "Wikimedia User Group Greece",
    "fishbowl": true
  },
  {
    "url": "https://grants.wikimedia.org",
    "dbname": "grantswiki",
    "code": "grants",
    "lang": "en",
    "sitename": "Wikimedia Foundation Grants Discussion",
    "private": true
  },
  {
    "url": "https://hi.wikimedia.org",
    "dbname": "hiwikimedia",
    "code": "hiwikimedia",
    "lang": "hi",
    "sitename": "Hindi Wikimedians User Group",
    "fishbowl": true
  },
  {
    "url": "https://id.wikimedia.org",
    "dbname": "idwikimedia",
    "code": "idwikimedia",
    "lang": "id",
    "sitename": "Wikimedia Indonesia",
    "fishbowl": true
  },
  {
    "url": "https://id-internal.wikimedia.org",
    "dbname": "id_internalwikimedia",
    "code": "id-internalwikimedia",
    "lang": "id",
    "sitename": "Wikimedia Indonesia (internal)",
    "private": true
  },
  {
    "url": "https://iegcom.wikimedia.org",
    "dbname": "iegcomwiki",
    "code": "iegcom",
    "lang": "en",
    "sitename": "Individual Engagement Grants Committee Wiki",
    "private": true
  },
  {
    "url": "https://il.wikimedia.org",
    "dbname": "ilwikimedia",
    "code": "ilwikimedia",
    "lang": "he",
    "sitename": "ויקימדיה",
    "private": true
  },
  {
    "url": "https://incubator.wikimedia.org",
    "dbname": "incubatorwiki",
    "code": "incubator",
    "lang": "en",
    "sitename": "Wikimedia Incubator"
  },
  {
    "url": "https://internal.wikimedia.org",
    "dbname": "internalwiki",
    "code": "internal",
    "lang": "internal",
    "sitename": "Internal",
    "private": true,
    "closed": true
  },
  {
    "url": "https://wikitech.wikimedia.org",
    "dbname": "labswiki",
    "code": "labs",
    "lang": "labs",
    "sitename": "Wikipedia",
    "nonglobal": true
  },
  {
    "url": "https://labtestwikitech.wikimedia.org",
    "dbname": "labtestwiki",
    "code": "labtest",
    "lang": "labtest",
    "sitename": "Wikipedia",
    "nonglobal": true
  },
  {
    "url": "https://legalteam.wikimedia.org",
    "dbname": "legalteamwiki",
    "code": "legalteam",
    "lang": "en",
    "sitename": "Legal Team Wiki",
    "private": true
  },
  {
    "url": "https://login.wikimedia.org",
    "dbname": "loginwiki",
    "code": "login",
    "lang": "login",
    "sitename": "Wikimedia Login Wiki"
  },
  {
    "url": "https://mai.wikimedia.org",
    "dbname": "maiwikimedia",
    "code": "maiwikimedia",
    "lang": "mai",
    "sitename": "मैथिली विकिमिडियन्स",
    "fishbowl": true
  },
  {
    "url": "https://www.mediawiki.org",
    "dbname": "mediawikiwiki",
    "code": "mediawiki",
    "lang": "mediawiki",
    "sitename": "MediaWiki"
  },
  {
    "url": "https://meta.wikimedia.org",
    "dbname": "metawiki",
    "code": "meta",
    "lang": "meta",
    "sitename": "Meta"
  },
  {
    "url": "https://mk.wikimedia.org",
    "dbname": "mkwikimedia",
    "code": "mkwikimedia",
    "lang": "mk",
    "sitename": "Викимедија Македонија"
  },
  {
    "url": "https://movementroles.wikimedia.org",
    "dbname": "movementroleswiki",
    "code": "movementroles",
    "lang": "movementroles",
    "sitename": "Movement Roles",
    "private": true
  },
  {
    "url": "https://mx.wikimedia.org",
    "dbname": "mxwikimedia",
    "code": "mxwikimedia",
    "lang": "es",
    "sitename": "Wikimedia México"
  },
  {
    "url": "https://ng.wikimedia.org",
    "dbname": "ngwikimedia",
    "code": "ngwikimedia",
    "lang": "en",
    "sitename": "Wikimedia Nigeria",
    "fishbowl": true
  },
  {
    "url": "https://nl.wikimedia.org",
    "dbname": "nlwikimedia",
    "code": "nlwikimedia",
    "lang": "nl",
    "sitename": "Wikimedia"
  },
  {
    "url": "https://no.wikimedia.org",
    "dbname": "nowikimedia",
    "code": "nowikimedia",
    "lang": "nb",
    "sitename": "Wikimedia Norge"
  },
  {
    "url": "https://noboard-chapters.wikimedia.org",
    "dbname": "noboard_chapterswikimedia",
    "code": "noboard-chapterswikimedia",
    "lang": "nb",
    "sitename": "Wikimedia Norway Internal Board",
    "private": true
  },
  {
    "url": "https://nostalgia.wikipedia.org",
    "dbname": "nostalgiawiki",
    "code": "nostalgia",
    "lang": "nostalgia",
    "sitename": "Wikipedia",
    "fishbowl": true
  },
  {
    "url": "https://nyc.wikimedia.org",
    "dbname": "nycwikimedia",
    "code": "nycwikimedia",
    "lang": "en",
    "sitename": "Wikimedia New York City"
  },
  {
    "url": "https://nz.wikimedia.org",
    "dbname": "nzwikimedia",
    "code": "nzwikimedia",
    "lang": "en",
    "sitename": "Wikimedia",
    "closed": true
  },
  {
    "url": "https://office.wikimedia.org",
    "dbname": "officewiki",
    "code": "office",
    "lang": "en",
    "sitename": "Wikimedia Office",
    "private": true
  },
  {
    "url": "https://ombudsmen.wikimedia.org",
    "dbname": "ombudsmenwiki",
    "code": "ombudsmen",
    "lang": "en",
    "sitename": "Ombuds Wiki",
    "private": true
  },
  {
    "url": "https://otrs-wiki.wikimedia.org",
    "dbname": "otrs_wikiwiki",
    "code": "otrs-wiki",
    "lang": "en",
    "sitename": "OTRS Wiki",
    "private": true
  },
  {
    "url": "https://outreach.wikimedia.org",
    "dbname": "outreachwiki",
    "code": "outreach",
    "lang": "outreach",
    "sitename": "Outreach Wiki"
  },
  {
    "url": "https://pa-us.wikimedia.org",
    "dbname": "pa_uswikimedia",
    "code": "pa-uswikimedia",
    "lang": "en",
    "sitename": "Wikimedia Pennsylvania",
    "closed": true
  },
  {
    "url": "https://pl.wikimedia.org",
    "dbname": "plwikimedia",
    "code": "plwikimedia",
    "lang": "pl",
    "sitename": "Wikimedia"
  },
  {
    "url": "https://projectcom.wikimedia.org",
    "dbname": "projectcomwiki",
    "code": "projectcom",
    "lang": "en",
    "sitename": "Project Grants Committee",
    "private": true
  },
  {
    "url": "https://pt.wikimedia.org",
    "dbname": "ptwikimedia",
    "code": "ptwikimedia",
    "lang": "pt",
    "sitename": "Wikimedia Portugal"
  },
  {
    "url": "https://punjabi.wikimedia.org",
    "dbname": "punjabiwikimedia",
    "code": "punjabiwikimedia",
    "lang": "pa",
    "sitename": "Punjabi Wikimedians",
    "fishbowl": true
  },
  {
    "url": "https://quality.wikimedia.org",
    "dbname": "qualitywiki",
    "code": "quality",
    "lang": "en",
    "sitename": "Wikimedia Quality",
    "closed": true
  },
  {
    "url": "https://romd.wikimedia.org",
    "dbname": "romdwikimedia",
    "code": "romdwikimedia",
    "lang": "ro",
    "sitename": "Wikimedia",
    "fishbowl": true
  },
  {
    "url": "https://rs.wikimedia.org",
    "dbname": "rswikimedia",
    "code": "rswikimedia",
    "lang": "sr",
    "sitename": "Викимедија",
    "fishbowl": true
  },
  {
    "url": "https://ru.wikimedia.org",
    "dbname": "ruwikimedia",
    "code": "ruwikimedia",
    "lang": "ru",
    "sitename": "Викимедиа"
  },
  {
    "url": "https://se.wikimedia.org",
    "dbname": "sewikimedia",
    "code": "sewikimedia",
    "lang": "sv",
    "sitename": "Wikimedia"
  },
  {
    "url": "https://searchcom.wikimedia.org",
    "dbname": "searchcomwiki",
    "code": "searchcom",
    "lang": "en",
    "sitename": "Search Committee",
    "private": true
  },
  {
    "url": "https://wikisource.org",
    "dbname": "sourceswiki",
    "code": "sources",
    "lang": "sources",
    "sitename": "Wikisource"
  },
  {
    "url": "https://spcom.wikimedia.org",
    "dbname": "spcomwiki",
    "code": "spcom",
    "lang": "spcom",
    "sitename": "Spcom",
    "private": true
  },
  {
    "url": "https://species.wikimedia.org",
    "dbname": "specieswiki",
    "code": "species",
    "lang": "species",
    "sitename": "Wikispecies"
  },
  {
    "url": "https://steward.wikimedia.org",
    "dbname": "stewardwiki",
    "code": "steward",
    "lang": "en",
    "sitename": "Steward Wiki",
    "private": true
  },
  {
    "url": "https://strategy.wikimedia.org",
    "dbname": "strategywiki",
    "code": "strategy",
    "lang": "en",
    "sitename": "Strategic Planning",
    "closed": true
  },
  {
    "url": "https://sysop-it.wikipedia.org",
    "dbname": "sysop_itwiki",
    "code": "sysop-it",
    "lang": "it",
    "sitename": "Italian Wikipedia sysops wiki",
    "private": true
  },
  {
    "url": "https://techconduct.wikimedia.org",
    "dbname": "techconductwiki",
    "code": "techconduct",
    "lang": "techconduct",
    "sitename": "CoC committee",
    "private": true
  },
  {
    "url": "https://ten.wikipedia.org",
    "dbname": "tenwiki",
    "code": "ten",
    "lang": "en",
    "sitename": "Wikipedia 10",
    "closed": true
  },
  {
    "url": "https://test.wikipedia.org",
    "dbname": "testwiki",
    "code": "test",
    "lang": "en",
    "sitename": "Wikipedia"
  },
  {
    "url": "https://test2.wikipedia.org",
    "dbname": "test2wiki",
    "code": "test2",
    "lang": "en",
    "sitename": "Wikipedia"
  },
  {
    "url": "https://test-commons.wikimedia.org",
    "dbname": "testcommonswiki",
    "code": "testcommons",
    "lang": "testcommons",
    "sitename": "Test Wikimedia Commons"
  },
  {
    "url": "https://test.wikidata.org",
    "dbname": "testwikidatawiki",
    "code": "testwikidata",
    "lang": "testwikidata",
    "sitename": "Wikipedia"
  },
  {
    "url": "https://thankyou.wikipedia.org",
    "dbname": "thankyouwiki",
    "code": "thankyou",
    "lang": "en",
    "sitename": "Thank You",
    "fishbowl": true
  },
  {
    "url": "https://tr.wikimedia.org",
    "dbname": "trwikimedia",
    "code": "trwikimedia",
    "lang": "tr",
    "sitename": "Wikimedia Türkiye"
  },
  {
    "url": "https://transitionteam.wikimedia.org",
    "dbname": "transitionteamwiki",
    "code": "transitionteam",
    "lang": "en",
    "sitename": "Transition Team Wiki",
    "private": true,
    "closed": true
  },
  {
    "url": "https://ua.wikimedia.org",
    "dbname": "uawikimedia",
    "code": "uawikimedia",
    "lang": "uk",
    "sitename": "Вікімедіа Україна"
  },
  {
    "url": "https://usability.wikimedia.org",
    "dbname": "usabilitywiki",
    "code": "usability",
    "lang": "en",
    "sitename": "Wikimedia Usability Initiative",
    "closed": true
  },
  {
    "url": "https://vote.wikimedia.org",
    "dbname": "votewiki",
    "code": "vote",
    "lang": "en",
    "sitename": "Wikimedia Vote Wiki",
    "fishbowl": true
  },
  {
    "url": "https://wb.wikimedia.org",
    "dbname": "wbwikimedia",
    "code": "wbwikimedia",
    "lang": "bn",
    "sitename": "West Bengal Wikimedians",
    "fishbowl": true
  },
  {
    "url": "https://wg-en.wikipedia.org",
    "dbname": "wg_enwiki",
    "code": "wg-en",
    "lang": "en",
    "sitename": "Wikipedia Working Group",
    "private": true
  },
  {
    "url": "https://www.wikidata.org",
    "dbname": "wikidatawiki",
    "code": "wikidata",
    "lang": "wikidata",
    "sitename": "Wikipedia"
  },
  {
    "url": "https://wikimania.wikimedia.org",
    "dbname": "wikimaniawiki",
    "code": "wikimania",
    "lang": "wikimania",
    "sitename": "Wikipedia"
  },
  {
    "url": "https://wikimania2005.wikimedia.org",
    "dbname": "wikimania2005wiki",
    "code": "wikimania2005",
    "lang": "wikimania2005",
    "sitename": "Wikipedia",
    "closed": true
  },
  {
    "url": "https://wikimania2006.wikimedia.org",
    "dbname": "wikimania2006wiki",
    "code": "wikimania2006",
    "lang": "wikimania2006",
    "sitename": "Wikipedia",
    "closed": true
  },
  {
    "url": "https://wikimania2007.wikimedia.org",
    "dbname": "wikimania2007wiki",
    "code": "wikimania2007",
    "lang": "wikimania2007",
    "sitename": "Wikipedia",
    "closed": true
  },
  {
    "url": "https://wikimania2008.wikimedia.org",
    "dbname": "wikimania2008wiki",
    "code": "wikimania2008",
    "lang": "wikimania2008",
    "sitename": "Wikipedia",
    "closed": true
  },
  {
    "url": "https://wikimania2009.wikimedia.org",
    "dbname": "wikimania2009wiki",
    "code": "wikimania2009",
    "lang": "wikimania2009",
    "sitename": "Wikipedia",
    "closed": true
  },
  {
    "url": "https://wikimania2010.wikimedia.org",
    "dbname": "wikimania2010wiki",
    "code": "wikimania2010",
    "lang": "wikimania2010",
    "sitename": "Wikipedia",
    "closed": true
  },
  {
    "url": "https://wikimania2011.wikimedia.org",
    "dbname": "wikimania2011wiki",
    "code": "wikimania2011",
    "lang": "wikimania2011",
    "sitename": "Wikipedia",
    "closed": true
  },
  {
    "url": "https://wikimania2012.wikimedia.org",
    "dbname": "wikimania2012wiki",
    "code": "wikimania2012",
    "lang": "wikimania2012",
    "sitename": "Wikipedia",
    "closed": true
  },
  {
    "url": "https://wikimania2013.wikimedia.org",
    "dbname": "wikimania2013wiki",
    "code": "wikimania2013",
    "lang": "wikimania2013",
    "sitename": "Wikipedia",
    "closed": true
  },
  {
    "url": "https://wikimania2014.wikimedia.org",
    "dbname": "wikimania2014wiki",
    "code": "wikimania2014",
    "lang": "wikimania2014",
    "sitename": "Wikipedia",
    "closed": true
  },
  {
    "url": "https://wikimania2015.wikimedia.org",
    "dbname": "wikimania2015wiki",
    "code": "wikimania2015",
    "lang": "wikimania2015",
    "sitename": "Wikipedia",
    "closed": true
  },
  {
    "url": "https://wikimania2016.wikimedia.org",
    "dbname": "wikimania2016wiki",
    "code": "wikimania2016",
    "lang": "wikimania2016",
    "sitename": "Wikipedia",
    "closed": true
  },
  {
    "url": "https://wikimania2017.wikimedia.org",
    "dbname": "wikimania2017wiki",
    "code": "wikimania2017",
    "lang": "wikimania2017",
    "sitename": "Wikipedia",
    "closed": true
  },
  {
    "url": "https://wikimania2018.wikimedia.org",
    "dbname": "wikimania2018wiki",
    "code": "wikimania2018",
    "lang": "wikimania2018",
    "sitename": "Wikipedia",
    "closed": true
  },
  {
    "url": "https://wikimaniateam.wikimedia.org",
    "dbname": "wikimaniateamwiki",
    "code": "wikimaniateam",
    "lang": "en",
    "sitename": "WikimaniaTeam",
    "private": true
  }
]`

var excludes = map[string]struct{}{}

func init() {
	specials := []struct {
		DbName string `json:"dbname"`
	}{}

	if err := json.NewDecoder(strings.NewReader(exclude)).Decode(&specials); err != nil {
		log.Panic(err)
	}

	for _, special := range specials {
		excludes[special.DbName] = struct{}{}
	}
}

// Exclude all the projects to be excluded from the streams
func Exclude(dbName string) bool {
	_, ok := excludes[dbName]
	return ok
}
