
from sklearn.preprocessing import MinMaxScaler
from keras.models import Sequential
from keras.layers import Dense
import keras.backend as K
from keras.callbacks import EarlyStopping
from keras.optimizers import Adam
from keras.models import load_model
from keras.layers import LSTM
import numpy as np
import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
import matplotlib
from sklearn.preprocessing import MinMaxScaler
from keras.layers import LSTM,Dense,Dropout
from sklearn.model_selection import TimeSeriesSplit
from sklearn.metrics import mean_squared_error, r2_score
import matplotlib.dates as mdates
from sklearn import linear_model
import yfinance as yf
import os
import shutil
from tensorflow import keras

dirname = os.path.dirname(__file__)

class Traning:
    def __init__(self) -> None:
        pass

    def preProcessData(self, ):
        pass

    def GenerateModel(self, y):
        K.clear_session()
        model_lstm = Sequential()
        model_lstm.add(LSTM(16, input_shape=(1, y), activation='relu', return_sequences=False))
        model_lstm.add(Dense(1))
        model_lstm.compile(loss='mean_squared_error', optimizer='adam')
        return model_lstm

    def SaveData(self, model):
        model.save(self.modelpath)
    
    def FetchData(self, name):
        data = yf.download(tickers=name, period="max", interval="1d")
        self.stockCsvPath = os.path.join(self.stockPath, name + ".csv")
        if os.path.exists(self.stockCsvPath):
            os.remove(self.stockCsvPath)
        data.to_csv(self.stockCsvPath)
    
    ##reuse data if its alread trained
    def CheckDataExists(self, name):
        self.stockPath = os.path.join(dirname, 'stock', name)
        self.modelpath = os.path.join(self.stockPath, name + ".pickel")
        if os.path.exists(self.modelpath):
            return True
        os.mkdir(self.stockPath)
        return False
    
    def loadModel(self):
        return keras.models.load_model(self.modelpath)

    def Train(self, name):
        
        ##traning part is over
        if self.CheckDataExists(name):
            return self.loadModel()

        self.FetchData(name)
        ##load data
        df_final = pd.read_csv(self.stockCsvPath,na_values=['null'],index_col='Date',parse_dates=True,infer_datetime_format=True)
        test = df_final
        target_adj_close = pd.DataFrame(test['Adj Close'])
        # selecting Feature Columns
        feature_columns = ['Open', 'High', 'Low', 'Volume']
        scaler = MinMaxScaler()
        feature_minmax_transform_data = scaler.fit_transform(test[feature_columns])
        feature_minmax_transform = pd.DataFrame(columns=feature_columns, data=feature_minmax_transform_data, index=test.index)
        feature_minmax_transform.head()
        # Shift target array because we want to predict the n + 1 day value
        target_adj_close = target_adj_close.shift(-1)
        validation_y = target_adj_close[-90:-1]
        target_adj_close = target_adj_close[:-90]

        # Taking last 90 rows of data to be validation set
        validation_X = feature_minmax_transform[-90:-1]
        feature_minmax_transform = feature_minmax_transform[:-90]
        ts_split= TimeSeriesSplit(n_splits=10)
        for train_index, test_index in ts_split.split(feature_minmax_transform):
                X_train, X_test = feature_minmax_transform[:len(train_index)], feature_minmax_transform[len(train_index): (len(train_index)+len(test_index))]
                y_train, y_test = target_adj_close[:len(train_index)].values.ravel(), target_adj_close[len(train_index): (len(train_index)+len(test_index))].values.ravel()
        X_train =np.array(X_train)
        X_test =np.array(X_test)

        X_tr_t = X_train.reshape(X_train.shape[0], 1, X_train.shape[1])
        X_tst_t = X_test.reshape(X_test.shape[0], 1, X_test.shape[1])

        early_stop = EarlyStopping(monitor='loss', patience=5, verbose=1)

        model_lstm = self.GenerateModel(X_train.shape[1])

        model_lstm.fit(X_tr_t, y_train, epochs=200, batch_size=8, verbose=1, shuffle=False, callbacks=[early_stop])

        score_lstm= model_lstm.evaluate(X_tst_t, y_test, batch_size=1)

        print('LSTM: %f'%score_lstm)
        self.SaveData(model_lstm)
        return model_lstm

train = Traning()

train.Train("MARUTI.NS")