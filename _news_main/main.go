package main

import (
	"fmt"
	"os"

	"github.com/pavlo67/punctum/basis/filelib"
	"github.com/pavlo67/punctum/starter"
	"github.com/pavlo67/punctum/starter/config"
	"github.com/pavlo67/punctum/starter/logger"

	"github.com/pavlo67/punctum/processor/importer/importer_rss"
	"github.com/pavlo67/punctum/processor/news"

	"log"

	"github.com/pavlo67/punctum/_news_main/news_router"
	"github.com/pavlo67/punctum/_news_main/news_starters"
)

//var setup = flag.Bool("setup", false, "recreate structures for the selected (or all if no) component")

func main() {
	err := logger.Init(logger.Config{LogLevel: logger.DebugLevel})
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("can't logger.Init(logger.Config{LogLevel: logger.DebugLevel}): %s", err))
		os.Exit(1)
	}
	l := logger.Get()

	cfgPath := filelib.CurrentPath() + "../cfg.json5"
	conf, err := config.Get(cfgPath, l)
	if err != nil {
		l.Fatalf("can't config.Get(%s): %s", cfgPath, err)
	}

	// flag.Parse()

	starters, label := news_starters.Starters()
	joiner, err := starter.Run(conf, starters, label, nil)
	if err != nil {
		l.Fatal(err)

	}
	defer joiner.CloseAll()

	newsOp, ok := joiner.Interface(news.InterfaceKey).(news.Operator)
	if !ok {
		log.Fatalf("no news.Operator with key %s", news.InterfaceKey)

	}

	rssURLs := []string{
		"https://rss.unian.net/site/news_ukr.rss",
		"https://ua.censor.net.ua/includes/news_uk.xml",
		"https://ua.censor.net.ua/includes/resonance_uk.xml",
		"https://ua.censor.net.ua/includes/events_uk.xml",
		"https://ua.interfax.com.ua/news/last.rss",
		"https://www.pravda.com.ua/rss/",
		"https://gazeta.ua/rss",
		"http://tyzhden.ua/RSS/All/",
		"https://day.kyiv.ua/uk/news-rss.xml",
		"https://krytyka.com/ua/rss/all",
		"https://krytyka.com/ua/rss/journals",
	}

	rssOp := &importer_rss.RSS{}

	for _, rssURL := range rssURLs {
		l.Info(rssURL)

		num, numNew, errs := news_router.Load(rssOp, rssURL, newsOp)

		l.Infof("total: %d, with errors: %d, added new: %d", num, len(errs), numNew)

		if len(errs) > 0 {
			l.Errorf("%#v", errs)
		}
	}

}