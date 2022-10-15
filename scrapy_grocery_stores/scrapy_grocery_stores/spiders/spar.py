import scrapy
from scrapy import Request
from scrapy.loader import ItemLoader
from scrapy_grocery_stores.items import GroceryStoreItem

class SparSpider(scrapy.Spider):
    name = 'spar'
    allowed_domains = ['spar.si', 'spar-ics.com']

    api = 'https://search-spar.spar-ics.com/fact-finder/rest/v4/search/products_lmos_si?query=*&q=*&page={page}&hitsPerPage=1000'
    ids = set()

    def start_requests(self):
        # NOTE(miha): In spar's website page argument starts with 1!
        yield Request(self.api.format(page=1))

    def parse(self, response, page=1):
        json_data = response.json()

        for item in json_data["hits"]:
            if item["id"] not in self.ids:
                self.ids.add(item["id"])

                item_loader = ItemLoader(item=GroceryStoreItem())
                item_loader.add_value("id", item["id"])
                item_loader.add_value("name", item["masterValues"]["name"])
                item_loader.add_value("price", item["masterValues"]["price"])
                yield item_loader.load_item()

        if page < 18:
            yield Request(self.api.format(page=page+1), cb_kwargs={"page":page+1})
