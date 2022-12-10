# stock-price-prediction
An application that uses Spark AI to predict the price of a stock

tf 
    -contains the Traning and Prediction Code

db
    -will store the trained data and model

client

    - will handle the predicted data

server
    -will handle the communication with the client
    -and organize the communication between db and tf

data 
    -will contain the methods to fetch data


Default Configuration 
    DB-Server :- port 8080
    Prediction-API-Server := port 8010