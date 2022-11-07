# Define your item pipelines here
#
# Don't forget to add your pipeline to the ITEM_PIPELINES setting
# See: https://docs.scrapy.org/en/latest/topics/item-pipeline.html


# useful for handling different item types with a single interface
import pymongo

from itemadapter import ItemAdapter

class ScrapyGroceryStoresPipeline:
    def process_item(self, item, spider):
        return item

    def close_spider(self, spider):
        print("this is from the pipeline")

class AddStatisticsToMongoDB:
    # TODO(miha): Add pymongo to requirements.txt
    def __init__(self):
        self.client = pymongo.MongoClient("localhost", 27017)
        self.db = self.client["stats"]

    def close_spider(self, spider):
        collection = self.db[spider.name]
        stats = spider.crawler.stats.get_stats()
        collection.insert_one(stats)
        self.client.close()