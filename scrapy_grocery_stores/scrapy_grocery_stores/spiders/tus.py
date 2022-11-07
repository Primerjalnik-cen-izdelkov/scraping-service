import scrapy
from scrapy import Request
from scrapy.loader import ItemLoader
from scrapy_grocery_stores.items import GroceryStoreItem
import json
import string

class TusSpider(scrapy.Spider):
    name = 'tus'
    allowed_domains = ['tus.si', 'hitrinakup.com']

    # TODO(miha): There is also tusdrogerija.si website.

    api = "https://www.tus.si/aktualno/akcijska-ponudba/aktualno-iz-kataloga/page/{page}"
    SI_api = "https://www.tus.si/page/{page}/?s={char}&post_type=product"
    graphql_api = "https://hitrinakup.com/graphql"

    ids = set()

    category_headers = {
        'authority': 'hitrinakup.com',
        'accept': '*/*',
        'apiversion': '3.2',
        'content-type': 'application/json',
        'cookie': 'TSSO=Im1uamExZTA5NTMzNmY3MzQ2MmY5Nzk1YzYxNmM5MTY1NjY5Ig.Ovh0GqFb8jGWIA6jqf1vOielIP4; ph_phc_1FsDfc0x2bt2bhPxYAf0E2IurHeXa9crA07U4I1RIkY_posthog=%7B%22distinct_id%22%3A%226c54a5e2-6e95-4070-8247-d3d8d76ce2a5%22%2C%22%24device_id%22%3A%226c54a5e2-6e95-4070-8247-d3d8d76ce2a5%22%2C%22%24referrer%22%3A%22https%3A%2F%2Fhitrinakup.com%2Flanding%22%2C%22%24referring_domain%22%3A%22hitrinakup.com%22%2C%22%24sesid%22%3A%5B1666159941876%2C%22183eedc02e5656-08b19dbbeb6ed5-26021c51-384000-183eedc02e6e36%22%2C1666159739619%5D%7D',
        'origin': 'https://hitrinakup.com',
        'referer': 'https://hitrinakup.com/subcategories_with_items/Sadje%20in%20zelenjava',
        'sessionid': '1a324292-42e2-4968-8fb7-984b5a867248',
        'user-agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36',
    }

    category_json_data = {
        'operationName': 'getCategories',
        'variables': {
            'userId': '',
            'storeId': '5994',
            'limit': 0,
            'cypherQuery': '2d078df3-117a-4b05-be35-cfb99105fb77',
        },
        'query': 'query getCategories($userId: String, $limit: Int, $cypherQuery: String) {\n  getCategories(userId: $userId, limit: $limit, cypherQuery: $cypherQuery) {\n    name\n    key\n    image\n    children {\n      name\n      key\n      children {\n        name\n        key\n        __typename\n      }\n      __typename\n    }\n    __typename\n  }\n}\n',
    }

    sub_category_json_data = {
        'operationName': 'getSubcategoriesWithItems',
        'variables': {
            'storeId': '5994',
            'categoriesLimit': 10,
            'cypherQuery': '3a4d5cb9-52a9-4446-8d2b-6a8f93282934',
            'categoryName': '{category_name}',
            'filterProperties': None,
            'skip': 0,
            'limit': 3,
            'date': 'Wed Oct 19 2022',
        },
        'query': 'query getSubcategoriesWithItems($storeId: String, $categoriesLimit: Int, $cypherQuery: String, $categoryName: String, $filterProperties: FilterUpdateInput, $limit: Int, $skip: Int, $date: String) {\n  getSubcategoriesWithItems(\n    storeId: $storeId\n    categoriesLimit: $categoriesLimit\n    cypherQuery: $cypherQuery\n    categoryName: $categoryName\n    filterProperties: $filterProperties\n    skip: $skip\n    limit: $limit\n    date: $date\n  ) {\n    name\n    filters {\n      general {\n        name\n        itemCount\n        __typename\n      }\n      brands {\n        name\n        itemCount\n        __typename\n      }\n      categories {\n        name\n        itemCount\n        __typename\n      }\n      allergens {\n        name\n        itemCount\n        __typename\n      }\n      __typename\n    }\n    items {\n      weighing\n      weight\n      itemId\n      EAN\n      name\n      orderLimit\n      price\n      discountedPrice\n      lowQuantity\n      promotionDisplayPrice\n      promotionType\n      promotionDisplayProcentage\n      promotionProcentage\n      alcohol\n      group\n      quantity\n      itemWeighingChangeQuantityStep\n      inStock\n      priceEm\n      em\n      badges\n      type\n      comment\n      discountEan\n      id\n      displayName\n      itemBackground\n      bannerColor\n      bannerTextBold\n      bannerTextNormal\n      bannerTextColor\n      img\n      brand\n      category\n      __typename\n    }\n    __typename\n  }\n}\n',
    }

    item_json_data = {
        'operationName': 'getItemsForSelectedSubCategory',
        'variables': {
            'storeId': '5994',
            'categoriesLimit': 10,
            'categoriesSkip': 0,
            'categoryName': '{category_name}',
            'subcategoryName': '{sub_category_name}',
            'cypherQuery': '123a17d4-b159-4976-ac50-8590059b48c0',
            'date': 'Wed Oct 19 2022',
            'filterProperties': None,
        },
        'query': 'query getItemsForSelectedSubCategory($storeId: String, $categoriesLimit: Int, $categoriesSkip: Int, $subcategoryName: String, $cypherQuery: String, $filterProperties: FilterUpdateInput, $categoryName: String, $date: String) {\n  getItemsForSelectedSubCategory(\n    storeId: $storeId\n    categoriesLimit: $categoriesLimit\n    categoriesSkip: $categoriesSkip\n    subcategoryName: $subcategoryName\n    filterProperties: $filterProperties\n    cypherQuery: $cypherQuery\n    categoryName: $categoryName\n    date: $date\n  ) {\n    name\n    filters {\n      general {\n        name\n        itemCount\n        __typename\n      }\n      brands {\n        name\n        itemCount\n        __typename\n      }\n      categories {\n        name\n        itemCount\n        __typename\n      }\n      allergens {\n        name\n        itemCount\n        __typename\n      }\n      __typename\n    }\n    items {\n      weighing\n      weight\n      itemId\n      EAN\n      name\n      orderLimit\n      price\n      discountedPrice\n      lowQuantity\n      promotionDisplayPrice\n      promotionType\n      promotionDisplayProcentage\n      promotionProcentage\n      alcohol\n      group\n      quantity\n      itemWeighingChangeQuantityStep\n      inStock\n      priceEm\n      em\n      badges\n      type\n      comment\n      discountEan\n      id\n      displayName\n      itemBackground\n      bannerColor\n      bannerTextBold\n      bannerTextNormal\n      bannerTextColor\n      img\n      brand\n      category\n      __typename\n    }\n    __typename\n  }\n}\n',
    }

    def parse_graphql(self, response):
        json_data = response.json()
        print("MY DATA")
        for category in json_data["data"]["getCategories"]:
            print(category["name"])
            myjson = self.sub_category_json_data.copy()
            myjson["variables"]["categoryName"] = str(category["name"])
            yield Request(self.graphql_api, callback= self.parse_sub_category, method="POST", body=json.dumps(myjson), headers=self.category_headers, cb_kwargs={'category':category["name"]})

    def parse_sub_category(self, response, category):
        json_data = response.json()
        print("SUB CATEGORY DATA")
        for sub_category in json_data["data"]["getSubcategoriesWithItems"]:
            print(sub_category["name"])
            myjson = self.item_json_data.copy()
            myjson["variables"]["categoryName"] = str(category)
            myjson["variables"]["subcategoryName"] = str(sub_category["name"])
            yield Request(self.graphql_api, callback= self.parse_item, method="POST", body=json.dumps(myjson), headers=self.category_headers)

    def parse_item(self, response):
        json_data = response.json()
        print("ITEM DATA")
        for item in json_data["data"]["getItemsForSelectedSubCategory"]["items"]:
            print(item["name"])
            if item["id"] not in self.ids:
                self.ids.add(item["id"])
                item_loader = ItemLoader(item=GroceryStoreItem())
                item_loader.add_value("id", item["id"])
                item_loader.add_value("name", item["name"])
                item_loader.add_value("price", item["price"])
                yield item_loader.load_item()

    def start_requests(self):
        for c in string.ascii_lowercase:
            yield Request(self.SI_api.format(page=1, char=c), callback=self.SI_parse)
        yield Request(self.api.format(page=1))
        yield Request(self.graphql_api, method="POST", body=json.dumps(self.category_json_data), headers=self.category_headers, callback=self.parse_graphql)

    """
    sample div card:
    <div class="card card-product" data-product-ean="9002975301558"> <div
    class="hover"> <figure> <a
    href="https://www.tus.si/izdelki/bonboni-haribo-zlati-medo-100-g/"> <img
    alt="Bonboni Haribo, zlati medo, 100 g"
    data-src="https://www.tus.si/app/images/9002975301558.jpg" class="thumbnail
    lazyload"
    src="data:image/gif;base64,R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==">
    <noscript><img src="https://www.tus.si/app/images/9002975301558.jpg"
    alt="Bonboni Haribo, zlati medo, 100 g" class="thumbnail"></noscript>
    </a> <figcaption> <a href="#" class="btn green add-to-list">Dodaj na
    nakupovalni listek</a> <a
    href="https://www.tus.si/izdelki/bonboni-haribo-zlati-medo-100-g/"
    class="padding-left-right-15px480">Več o izdelku</a> </figcaption>
    </figure> <h3> <a
    href="https://www.tus.si/izdelki/bonboni-haribo-zlati-medo-100-g/">Bonboni
    Haribo, zlati medo, 100 g</a> </h3> </div> <span class="price">0,79 €
    </span> <p class="price-regular price-regular-desktop"
    style="margin:0;"><label>Redna cena: 0,97 € </label></p> <p
    class="price-regular-mobile" style="margin:0;"><label><s>0,97 €
    </s></label></p> <p class="price-regular price-regular-desktop
    price-saving" style="margin-top: -10px;"><label>Prihranek: 0,18
    €</label></p> <p class="price-regular-mobile price-saving"
    style="margin-top: -10px;"><label>Prihranek: 0,18 €</label></p> <span
    class="discount" style="background-image:
        url('https://www.tus.si/app/themes/tus-theme/dist/images/bg-discount-yellow.svg')">
        <span>18</span> </span> </div>
    """
    def parse(self, response, page=1):
        for el in response.xpath("//div[contains(@class, 'card')]"):
            id = el.xpath(".//attribute::data-product-ean").get()
            name = el.xpath(".//attribute::alt").get()
            # TODO(miha): We need to save discounted price too!
            price = el.xpath(".//s/text()").get()

            if id not in self.ids:
                self.ids.add(id)
                item_loader = ItemLoader(item=GroceryStoreItem())
                item_loader.add_value("id", id)
                item_loader.add_value("name", name)
                item_loader.add_value("price", price)
                yield item_loader.load_item()

        if page < 18:
            yield Request(self.api.format(page=page+1), cb_kwargs={"page": page+1})

    def SI_parse(self, response, page=1, char='a'):
        for el in response.xpath("//div[contains(@class, 'card')]"):
            id = el.xpath(".//attribute::data-product-ean").get()
            name = el.xpath(".//attribute::alt").get()
            # TODO(miha): We need to save discounted price too!
            price = el.xpath(".//s/text()").get()

            if id not in self.ids:
                self.ids.add(id)
                item_loader = ItemLoader(item=GroceryStoreItem())
                item_loader.add_value("id", id)
                item_loader.add_value("name", name)
                item_loader.add_value("price", price)
                yield item_loader.load_item()

        if page < 18:
            yield Request(self.SI_api.format(page=page+1, char=char), cb_kwargs={"page": page+1, "char":char}, callback=self.SI_parse)
