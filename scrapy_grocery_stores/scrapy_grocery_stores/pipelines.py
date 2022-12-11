# Define your item pipelines here
#
# Don't forget to add your pipeline to the ITEM_PIPELINES setting
# See: https://docs.scrapy.org/en/latest/topics/item-pipeline.html


# useful for handling different item types with a single interface
import pymongo
import os
import logging
import datetime
from collections import defaultdict

from itemadapter import ItemAdapter

class ScrapyGroceryStoresPipeline:
    def process_item(self, item, spider):
        return item

    #def close_spider(self, spider):
    #    print("this is from the pipeline")

class AddStatisticsToMongoDB:
    # TODO(miha): Get database name from ENV
    def __init__(self):
        self.client = pymongo.MongoClient(os.environ.get('MONGODB_URI'))
        self.db = self.client.get_database("stats")
        self.f_db = self.client.get_database("files")

    def close_spider(self, spider):
        collection = self.db.get_collection("stats1")
        stats = spider.crawler.stats.get_stats()
        stats_obj = defaultdict(int)
        stats_obj.update({"start_time": stats["start_time"],
                     "downloader_request_count": stats["downloader/request_count"],
                     "downloader_response_count": stats["downloader/response_count"],
                     "downloader_response_status_count_404": stats["downloader/response_status_count/404"],
                     "item_scraped_count": stats["item_scraped_count"],
                     "bot_name": spider.name,
                    })
        collection.insert_one(stats_obj)

        # NOTE(miha): Insert file info into collection files.
        file_time = datetime.datetime.strptime(str(stats["start_time"]), "%Y-%m-%d %H:%M:%S.%f").strftime("%Y-%m-%dT%H-%M-%S")
        db_time = datetime.datetime.strptime(str(stats["start_time"]), "%Y-%m-%d %H:%M:%S.%f")
        files_collection = self.f_db.get_collection("files")
        file_name = f"{spider.name}_{file_time}.csv"
        files_collection.insert_one({"bot_name": spider.name, "file_name": file_name, "date": db_time})

        self.client.close()
