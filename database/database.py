from pymongo import MongoClient
import config

dbClient = MongoClient(config.DBURI)

mydb = myclient["mydatabase"]

print(myclient.list_database_names())


class Database:

    def __init__(self) -> None:
        pass

    def Connect(self):
        try:
            self.dbClient = MongoClient(config.DBURI)
        except:
            print("error cooured while connecting to DB")
            return False
        
        return True
        
    def Close(self):
        self.dbClient.close()
    
    def Query(self):
        pass

    def Update(self):
        pass

    def Insert(self):
        pass