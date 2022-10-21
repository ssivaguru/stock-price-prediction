

'''
install the followign 
pip install yfinance
pip install pendulum
pip install matplotlib
'''



##Sample code to download minute level data
# Import package
from pathlib import Path
import yfinance as yf
import os.path



filepath = '../spark/train.csv'
# valid periods: 1d,5d,1mo,3mo,6mo,1y,2y,5y,10y,ytd,max
# valid intervals: 1m,2m,5m,15m,30m,60m,90m,1h,1d,5d,1wk,1mo,3mo
# last one minute data is available only for 7 days
# 5m, 15, and 30m are availavble for 60 days

# Get the data
data = yf.download(tickers="HDFCBANK.NS", period="3mo", interval="60m")

#print data
print(data.head())
myfile = Path(filepath)
if myfile.is_file():
    os.remove(filepath)
# save data

data.to_csv(filepath)