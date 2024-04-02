from dotenv import load_dotenv
import os

load_dotenv('/Users/samuelsvensson/Documents/forecasting_project/forecast_app/.env')

db_name=os.getenv('DB_NAME')
db_user=os.getenv('DB_USER')
db_pass=os.getenv('DB_PASSWORD')
db_host=os.getenv('DB_HOST')
db_port=os.getenv('DB_PORT')