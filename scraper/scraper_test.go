package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestParseItem(t *testing.T) {
	html := `<div class="s-item__wrapper clearfix"><div class="s-item__image-section"><div class="s-item__image"><a tabindex="-1" aria-hidden="true" data-track="{&quot;eventFamily&quot;:&quot;LST&quot;,&quot;eventAction&quot;:&quot;ACTN&quot;,&quot;actionKind&quot;:&quot;NAVSRC&quot;,&quot;actionKinds&quot;:[&quot;NAVSRC&quot;],&quot;operationId&quot;:&quot;2351460&quot;,&quot;flushImmediately&quot;:false,&quot;eventProperty&quot;:{&quot;parentrq&quot;:&quot;44c6e2b217a0ad919717a9bdfffe3b8e&quot;,&quot;pageci&quot;:&quot;ac7ef2a6-d5f0-11eb-9693-e61189cd3eed&quot;,&quot;moduledtl&quot;:&quot;mi:1686|iid:1|li:7400|luid:1|scen:Listings&quot;}}" _sp="p2351460.m1686.l7400" href="https://www.ebay.com/itm/402943017690?hash=item5dd14682da:g:RdAAAOSwudVg0-jc"><div class="s-item__image-wrapper"><div class="s-item__image-helper"></div><img class="s-item__image-img" alt="Puma Powercamp 2.0 Training  Ball Mens Soccer Cleats     - Size 5" src="https://i.ebayimg.com/thumbs/images/g/RdAAAOSwudVg0-jc/s-l225.webp" onload="SITE_SPEED.ATF_TIMER.measure(this); if (performance &amp;&amp; performance.mark) { performance.mark(&quot;first-meaningful-paint&quot;); };if(this.width === 80 &amp;&amp; this.height === 80) {window.SRP.metrics.imageEmptyError.count++;}" onerror="window.SRP.metrics.imageLoadError.count++; " data-atftimer="1624651523805"></div></a></div></div><div class="s-item__info clearfix"><a data-track="{&quot;eventFamily&quot;:&quot;LST&quot;,&quot;eventAction&quot;:&quot;ACTN&quot;,&quot;actionKind&quot;:&quot;NAVSRC&quot;,&quot;actionKinds&quot;:[&quot;NAVSRC&quot;],&quot;operationId&quot;:&quot;2351460&quot;,&quot;flushImmediately&quot;:false,&quot;eventProperty&quot;:{&quot;parentrq&quot;:&quot;44c6e2b217a0ad919717a9bdfffe3b8e&quot;,&quot;pageci&quot;:&quot;ac7ef2a6-d5f0-11eb-9693-e61189cd3eed&quot;,&quot;moduledtl&quot;:&quot;mi:1686|iid:1|li:7400|luid:1|scen:Listings&quot;}}" _sp="p2351460.m1686.l7400" class="s-item__link" href="https://www.ebay.com/itm/402943017690?hash=item5dd14682da:g:RdAAAOSwudVg0-jc"><h3 class="s-item__title">Puma Powercamp 2.0 Training  Ball Mens Soccer Cleats     - Size 5</h3></a><div class="s-item__subtitle"><span class="SECONDARY_INFO">Brand New</span></div><div class="s-item__details clearfix"><div class="s-item__detail s-item__detail--primary"><span class="s-item__price">$19.99</span></div><span class="s-item__detail s-item__detail--secondary"><span class="s-item__dynamic s-item__listingDate"><span class="BOLD">Jun-23 19:07</span></span></span><div class="s-item__detail s-item__detail--primary"><span class="s-item__trending-price">List price: <span class="clipped">Previous Price</span><span class="STRIKETHROUGH">$30.00</span></span>  <span class="s-item__discount s-item__discount"><span class="BOLD">33% off</span></span></div><span class="s-item__detail s-item__detail--secondary"><span class="s-item__gsp-info s-item__gspInfo">Customs services and international tracking provided</span></span><div class="s-item__detail s-item__detail--primary"><span class="s-item__purchase-options-with-icon" aria-label="">Buy It Now</span></div><div class="s-item__detail s-item__detail--primary"><span class="s-item__shipping s-item__logisticsCost">+$33.39 shipping estimate</span></div><div class="s-item__detail s-item__detail--primary"><span class="s-item__location s-item__itemLocation">from United States</span></div><div class="s-item__detail s-item__detail--primary"><span class="s-item__sep"> <span role="text"><span class="s-jre2v01">0</span><span class="s-jre2v01">S</span><span class="s-jre2v01">N</span><span class="s-jre2v01">0</span><span class="s-jre2v01">7</span><span class="s-jre2v01">E</span><span class="s-jre2v01">p</span><span class="s-jre2v01">9</span><span class="s-jre2v01">o</span><span class="s-jre2v01">n</span><span class="s-jre2v01">M</span><span class="s-jre2v01">s</span><span class="s-jre2v01">o</span><span class="s-jre2v01">r</span><span class="s-jre2v01">e</span><span class="s-jre2v01">I</span><span class="s-jre2v01">1</span><span class="s-jre2v01">d</span><span class="s-jre2v01">O</span><span class="s-jre2v01">4</span><span class="s-jre2v01">D</span><span class="s-jre2v01"></span><span class="s-jre2v01"></span><span class="s-jre2v01"></span><span class="s-jre2v01"></span><span class="s-jre2v01"></span><span class="s-jre2v01"></span><span class="s-jre2v01"></span></span></span></div></div><span data-marko-key="@_wbind s0-14-11-6-3-listing1-item-5-1-20-0" class="s-item__watchheart at-corner s-item__watchheart--watch" data-has-widget="false" id="s0-14-11-6-3-listing1-item-5-1-20-0"><a aria-label="watch Puma Powercamp 2.0 Training  Ball Mens Soccer Cleats     - Size 5" _sp="p2351460.m4114.l8480" href="https://www.ebay.com/myb/WatchListAdd?item=402943017690&amp;pt=null&amp;srt=010006000000500af929089e576875d794c7ad4a96f512306854bc293d9413ae055a6d826716c6a2da31c6ef39187493e32238abcf7ffa382450d474ba25570e207b8b1d352cc7a55c185fb436d8b950f69e6a2b1f2068&amp;ru=https%3A%2F%2Fwww.ebay.com%2Fsch%2Fi.html%3F_from%3DR40%26_nkw%3Dsoccer%2Bball%2Bpuma%26_sacat%3D0%26LH_TitleDesc%3D0%26_sop%3D10"><span class="s-item__watchheart-icon"><svg aria-hidden="true" class="svg-icon" width="30px" height="30px"><use xlink:href="#svg-icon--save-circle" class="rest"></use><use xlink:href="#svg-icon--save-circle-hover" class="hover"></use><use xlink:href="#svg-icon--save-circle-active" class="active"></use></svg><span class="clipped"></span></span></a></span></div></div>`

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Errorf("could not create document: %v", err)
	}

	got := make([]Listing, 0)
	doc.Find("div.s-item__info").EachWithBreak(func(i int, sel *goquery.Selection) bool {
		listing, _ := parseItem(sel, nil, "")
		if listing != nil {
			// Keep the test simple by removing dates
			listing.Date = time.Time{}

			got = append(got, *listing)
		}
		return true
	})

	exp := []Listing{
		{
			URL:      "https://www.ebay.com/itm/402943017690?hash=item5dd14682da:g:RdAAAOSwudVg0-jc",
			Title:    "Puma Powercamp 2.0 Training  Ball Mens Soccer Cleats     - Size 5",
			Subtitle: "Brand New",
			Price:    "$19.99",
			Date:     time.Time{},
		},
	}

	if !reflect.DeepEqual(exp, got) {
		t.Errorf("expected %+v but got %+v", exp, got)
	}
}
