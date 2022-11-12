
from train import Traning
import os
from tensorflow import keras

dirname = os.path.dirname(__file__)


class Predict():
    def __init__(self) -> None:
        self.train = Traning()
        pass
    
    def loadModel(self, modelPath):
        return keras.models.load_model(self.modelPath)

    
    def LoadData(self):
        pass
    
    def StartTraning(self):
        pass

    def CheckData(self, modelPath):
        if os.path.exists(modelPath):
            return True
        return False

    def PredictData(self, name):
        stockPath = os.path.join(dirname, 'stock', name)
        modelpath = os.path.join(self.stockPath, name + ".pickel")
        stockCsvPath = os.path.join(self.stockPath, name + ".csv")

        model = None
        if self.CheckData(modelpath):
            model = self.loadModel(modelpath)
        else:
            model = self.train.StartTraning(name, stockPath, modelpath, stockCsvPath)

        ##predict data with the model        
        pass    