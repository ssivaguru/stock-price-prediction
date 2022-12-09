#!/usr/bin/env python


import pandas as pd
import numpy as np
from sklearn.preprocessing import StandardScaler, MinMaxScaler
from keras.models import Sequential
from keras.layers import Dense
from keras.layers import LSTM
from tensorflow import keras
import sys
import yfinance as yf
import os

# To remove the scientific notation from numpy arrays
np.set_printoptions(suppress=True)

TimeSteps=10  # next few day's Price Prediction is based on last how many past day's prices
FutureTimeSteps=5 # How many days in future you want to predict the prices

def remove(stockName):
    os.remove(stockName + ".csv")
    os.remove(stockName + ".h5")

def downloadData(stockName):
    # valid periods: 1d,5d,1mo,3mo,6mo,1y,2y,5y,10y,ytd,max
    # valid intervals: 1m,2m,5m,15m,30m,60m,90m,1h,1d,5d,1wk,1mo,3mo
    # last one minute data is available only for 7 days
    # 5m, 15, and 30m are availavble for 60 days
    data = yf.download(tickers=stockName, period="max", interval="1d")
    data.to_csv(stockName+".csv")


def loadData(stockName):
    StockData=pd.read_csv(stockName + ".csv")
    StockData['TradeDate']=StockData.index
    return StockData[['Close']].values 

def train(stockName):
    global TimeSteps
    global FutureTimeSteps
    # Fetching the data
    FullData = loadData(stockName)
    sc=MinMaxScaler()
    DataScaler = sc.fit(FullData)
    X=DataScaler.transform(FullData)
    X=X.reshape(X.shape[0],)

    # split into samples
    X_samples = list()
    y_samples = list()
    
    NumerOfRows = len(X)
    
    # Iterate thru the values to create combinations
    for i in range(TimeSteps , NumerOfRows-FutureTimeSteps , 1):
        x_sample = X[i-TimeSteps:i]
        y_sample = X[i:i+FutureTimeSteps]
        X_samples.append(x_sample)
        y_samples.append(y_sample)

    # Reshape the Input as a 3D (samples, Time Steps, Features)
    X_data=np.array(X_samples)
    X_data=X_data.reshape(X_data.shape[0],X_data.shape[1], 1)
    
    y_data=np.array(y_samples)

    TestingRecords=10
    
    # Splitting the data into train and test
    X_train=X_data[:-TestingRecords]
    X_test=X_data[-TestingRecords:]
    y_train=y_data[:-TestingRecords]
    y_test=y_data[-TestingRecords:]
    
    TimeSteps=X_train.shape[1]
    TotalFeatures=X_train.shape[2]

    # Initialising the RNN
    regressor = Sequential()
    
    # Adding the First input hidden layer and the LSTM layer
    # return_sequences = True, means the output of every time step to be shared with hidden next layer
    regressor.add(LSTM(units = 10, activation = 'relu', input_shape = (TimeSteps, TotalFeatures), return_sequences=True))
    
    # Adding the Second hidden layer and the LSTM layer
    regressor.add(LSTM(units = 5, activation = 'relu', input_shape = (TimeSteps, TotalFeatures), return_sequences=True))
    
    # Adding the Third hidden layer and the LSTM layer
    regressor.add(LSTM(units = 5, activation = 'relu', return_sequences=False ))
    
    # Adding the output layer
    # Notice the number of neurons in the dense layer is now the number of future time steps 
    # Based on the number of future days we want to predict
    regressor.add(Dense(units = FutureTimeSteps))
    
    # Compiling the RNN
    regressor.compile(optimizer = 'adam', loss = 'mean_squared_error')
    
    # Fitting the RNN to the Training set
    regressor.fit(X_train, y_train, batch_size = 5, epochs = 1)
    regressor.save(stockName + ".h5")
    
    

def predict(stockName):
    global TimeSteps
    global FutureTimeSteps
    regressor = keras.models.load_model(stockName + ".h5")
    FullData = loadData(stockName)
    last_n_days = np.array(FullData[-TimeSteps:])

    # Reshaping the data to (-1,1 )because its a single entry
    last_n_days=last_n_days.reshape(-1, 1)
    sc=MinMaxScaler()
    DataScaler = sc.fit(FullData)
    # Scaling the data on the same level on which model was trained
    X_test=DataScaler.transform(last_n_days)
    
    NumberofSamples=1
    TimeSteps=X_test.shape[0]
    NumberofFeatures=X_test.shape[1]
    # Reshaping the data as 3D input
    X_test=X_test.reshape(NumberofSamples,TimeSteps,NumberofFeatures)
    
    # Generating the predictions for next 5 days
    Next5DaysPrice = regressor.predict(X_test)
    
    # Generating the prices in original scale
    Next5DaysPrice = DataScaler.inverse_transform(Next5DaysPrice)
    return np.array(Next5DaysPrice).reshape(FutureTimeSteps, 1)


##we either train or update
if sys.argv[1] == "predict":
    downloadData(sys.argv[2])
    train(sys.argv[2])
    resp = predict(sys.argv[2])
    print(resp)
else:
    remove(sys.argv[2])
    downloadData(sys.argv[2])
    train(sys.argv[2])
    resp = predict(sys.argv[2])

