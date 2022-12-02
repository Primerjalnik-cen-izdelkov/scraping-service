import scrapy
from scrapy import Request
from scrapy.loader import ItemLoader
from scrapy_grocery_stores.items import GroceryStoreItem

headers = {
    'user-agent': 'Lidl/4.33.4(#197) Android',
    'x-device': 'samsung SM-P610',
    'x-os': 'Android 12',
    'content-type': 'application/json; charset=UTF-8',
    'accept': 'application/json',
    'authorization': 'Basic cHJvZGFwaXVzZXI6bGlkbGxvaG50c2ljaDIwMTQ=',
}

api = "https://endpoints.leaflets.schwarz/v3/{catalogue_id}/flyer.json?regionCode=0"
api_catalogue = "https://www.lidl.si/spletni-katalog"

api_products = "https://mobile.lidl.de/Mobile-Server/service/78/containerService/SI/campaign/sl/{section_id}/0/50"  # ponudba
api_product_categories = "https://mobile.lidl.de/Mobile-Server/service/78/containerService/SI/root/sl"
api_single_product = "https://mobile.lidl.de/Mobile-Server/service/78/productService/SI/product/sl/{product_id}"

class LidlSpider(scrapy.Spider):
    name = 'lidl'
    allowed_domains = ['endpoints.leaflets.schwarz', 'lidl.si', 'mobile.lidl.de']

    ids = set()
    flyer_ids = set()
    fields_to_export = ["store_name", "id", "name", "price", "discounted_price", "product_url",
                        "product_image", "offer_start", "offer_end", "flyer_url"]

    def start_requests(self):
        yield Request(api_catalogue, callback=self.get_catalogues)
        yield Request(api_product_categories, headers=headers, callback=self.parse_categories)

    def get_catalogues(self, response):
        catalogues = response.xpath('//*[@class="leafletcontainer__content"]/a/@href').extract()

        for catalogue in catalogues:
            slug = catalogue.split("katalog/")[-1].split("/")[0]
            if slug not in self.flyer_ids:
                self.flyer_ids.add(slug)
                yield Request(api.format(catalogue_id=slug), callback=self.parse)

    def parse_categories(self, response):
        json_data = response.json()

        for container in json_data["RootContainer"]["ContainerItems"]:
            if container["GroupedContainer"]["groupedContainerLanguageSet"]["GroupedContainerLanguageSet"]["title"] == "Ponudba":
                for section in container["GroupedContainer"]["ContainerItems"]:  # Bio izdelki, mesnica, novi izdelki itd.
                    yield Request(api_products.format(section_id=section["campaignId"]), headers=headers, callback=self.parse_category)

    def parse_category(self, response):
        json_data = response.json()["Campaign"]

        for item in json_data["ContainerItems"]:
            product_id = item["Product"]["productId"]
            yield Request(api_single_product.format(product_id=product_id), headers=headers, callback=self.parse_product)

        if json_data.get("nextPagingStepDataPath"):
            yield Request(response.urljoin(json_data["nextPagingStepDataPath"]), headers=headers, callback=self.parse_category)

    def parse_product(self, response):
        item = response.json()["Product"]
        self.ids.add(item["productId"])

        item_loader = ItemLoader(item=GroceryStoreItem())
        item_loader.add_value("store_name", self.name)
        item_loader.add_value("id", item["productId"])
        item_loader.add_value("name", item["productLanguageSet"]["ProductLanguageSet"]["title"])
        item_loader.add_value("discounted_price", item["price"])
        item_loader.add_value("price", item.get("strokePrice", item["price"]))  # izdelki v ponudbi mogoče nimajo discounted price
        item_loader.add_value("product_image", item["mainImage"]["Image"]["smallUrl"])
        item_loader.add_value("url", item["shareUrl"])
        # item_loader.add_value("flyer_url", json_data["flyerUrlAbsolute"]) # ni letakov tu
        # item_loader.add_value("category_url", item["primaryCampaign"]["Campaign"]["url"])
        # item_loader.add_value("offer_start", json_data["offerStartDate"])  # stalna ponudba
        # item_loader.add_value("offer_end", json_data["offerEndDate"])  # stalna ponudba
        yield item_loader.load_item()

    def parse(self, response):
        json_data = response.json()['flyer']

        for item in json_data.get("products"):
            if item.get("productId") and item["productId"] not in self.ids:
                self.ids.add(item["productId"])

                item_loader = ItemLoader(item=GroceryStoreItem())
                item_loader.add_value("store_name", self.name)
                item_loader.add_value("id", item["productId"])
                item_loader.add_value("name", item["title"])
                item_loader.add_value("discounted_price", item["price"])
                item_loader.add_value("price", item["strokePrice"])
                item_loader.add_value("product_image", item["image"])
                item_loader.add_value("url", item["url"])
                item_loader.add_value("flyer_url", json_data["flyerUrlAbsolute"])
                item_loader.add_value("offer_start", json_data["offerStartDate"])
                item_loader.add_value("offer_end", json_data["offerEndDate"])
                yield item_loader.load_item()

        for flyer in json_data.get("relatedFlyers"):
            if flyer["slug"] not in self.flyer_ids:
                self.flyer_ids.add(flyer["slug"])
                yield Request(api.format(catalogue_id=flyer["slug"]), callback=self.parse)