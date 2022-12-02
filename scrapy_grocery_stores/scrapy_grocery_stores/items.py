# Define here the models for your scraped items
#
# See documentation in:
# https://docs.scrapy.org/en/latest/topics/items.html

import scrapy


class GroceryStoreItem(scrapy.Item):
    # TODO(miha): We want to save: item name, price, weight, date of scrape
    id               = scrapy.Field()
    name             = scrapy.Field()
    price            = scrapy.Field()
    discounted_price = scrapy.Field()
    members_only     = scrapy.Field()
    url              = scrapy.Field()
    offer_start      = scrapy.Field()
    offer_end        = scrapy.Field()
    product_image    = scrapy.Field()
    flyer_url        = scrapy.Field()