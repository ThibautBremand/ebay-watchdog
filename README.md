# ebay-watchdog

It automatically and regularly scrapes predefined eBay listing pages, and sends notifications to a configured Telegram 
bot when new listings are found, with basic data.  
Whenever a new listing appears on an eBay search page, you will be notified, on mobile and on desktop.

![watchdoge](https://user-images.githubusercontent.com/9871294/123490445-4d546a80-d614-11eb-9889-520df15e594e.jpg)

This project is based on [this one](https://github.com/samjmckenzie/ebay-monitor), but reworked, simplified and with 
great new features.  

## Quick start
- Make sure Golang is installed on your machine.
- Clone the repository
- Set up the eBay urls in the `config.toml` file (more details below)
- Set up your Telegram credentials in the `.env` file (more details below)
- Run `make build` in order to build the executable.
- Run `make run` to launch the program.

### Config.toml

#### Search URLs
To add a new url to scrape, add the following lines into the `config.toml` file:
```
[[searches]]
url = "copied URL"
```

Add as much as you want:  
```
[[searches]]
url = "https://www.ebay.com/sch/i.html?_from=R40&_nkw=duct+tape&_sacat=0&_sop=10"

[[searches]]
url = "https://www.ebay.com/sch/i.html?_from=R40&_nkw=macbook+pro&_sacat=0&_sop=10"
```

#### (Optional) Domains
You can add multiple domains per search url.
```
[[searches]]
url = "https://www.ebay.com/sch/i.html?_from=R40&_nkw=duct+tape&_sacat=0&_sop=10"
domains = ["com","co.uk"]
```

That way, the scraper will first scrape the ebay.com url, then it will scrape the same url but with the .co.uk domain.
It will then dedupe the listings to keep a list of unique listings.

This is useful when you want to retrieve an exhaustive list of items, when some of them are only available in some 
countries.

If you do not set domains, only the provided url will be scraped as usual.

#### Important notes
- You need to use urls with the **Time: newly listed** enabled, e.g. urls ending with `&_sop=10`.
- You need to use the following urls and domains since these are the only domains that are handled by the program at the 
  moment.
  - ebay.at
  - ebay.ca
  - ebay.ch
  - ebay.co.uk
  - ebay.com
  - ebay.com.au
  - ebay.com.my
  - ebay.com.sg
  - ebay.de
  - ebay.es
  - ebay.fr
  - ebay.ie
  - ebay.it
  - ebay.nl
  - ebay.ph
  - ebay.pl
    

Other parameters:  
- `delay`: period, in seconds, between two scraping loops. Keep it reasonably high.

### Telegram and .env
- First, you need to create a [Telegram account](https://desktop.telegram.org/).
- Then, for the following steps, you need to download and use the desktop version.  

You can create a new Telegram bot [via this link](https://t.me/BotFather). 
- Send the `/newbot` command to BotFather, and
follow the steps to create a new bot. Once the bot is created, you will receive a token.
- Set the `TELEGRAM_TOKEN` variable in the `.env` file with your token.
```
TELEGRAM_TOKEN="1222533313:AAFwNd_HsPtpxBy35vEaZoFzUUB74v5mBpW"
```

Then, you need to find your chat ID.
- Paste the following link in your browser. Replace `<Telegram-token>` with the Telegram token.
```
https://api.telegram.org/bot<Telegram-token>/getUpdates?offset=0
```
- Send a message to your bot in the Telegram application. The message text can be anything. Your chat history must include at least one message to get your chat ID.
- Refresh your browser.
- Identify the numerical chat ID by finding the id inside the chat JSON object. In the example below, the chat ID is 123456789.
```json
{  
   "ok":true,
   "result":[  
      {  
         "update_id":987654321,
         "message":{  
            "message_id":2,
            "from":{  
               "id":123456789,
               "first_name":"Mushroom",
               "last_name":"Kap"
            },
            "chat":{  
               "id":123456789,
               "first_name":"Mushroom",
               "last_name":"Kap",
               "type":"private"
            },
            "date":1487183963,
            "text":"hi"
         }
      }
   ]
}
```
- Set the `TELEGRAM_CHAT_ID` variable in the `.env` file with this value.
```
TELEGRAM_CHAT_ID=123456789
```
