import scrapy
from scrapy import Request
from scrapy.loader import ItemLoader
from scrapy_grocery_stores.items import GroceryStoreItem

class MercatorSpider(scrapy.Spider):
    name = 'mercator'
    allowed_domains = ['mercator.si']
    start_urls = ['https://trgovina.mercator.si/market/brskaj']

    api = "https://trgovina.mercator.si/market/products/browseProducts/getProducts?limit=1000&offset={page}"
    ids = set()

    def start_requests(self):
        yield Request(self.api.format(page=0))

    def parse(self, response, page=0):
        json_data = response.json()

        for item in json_data:
            if item.get("data") and item["data"]["cinv"] not in self.ids:
                self.ids.add(item["data"]["cinv"])

                item_loader = ItemLoader(item=GroceryStoreItem())
                item_loader.add_value("id", item["data"]["cinv"])
                item_loader.add_value("name", item["data"]["name"])
                item_loader.add_value("price", item["data"]["current_price"])
                yield item_loader.load_item()

        if page < 14:
            yield Request(self.api.format(page=page+1), cb_kwargs={"page":page+1})

    def closed(self, reason):
        if reason == "finished":
            stats = self.crawler.stats.get_stats()
            print("from spider_closed:")
            print(stats)