from dotenv import load_dotenv
import os
import psycopg2
from datetime import datetime

load_dotenv('/Users/samuelsvensson/Documents/forecasting_project/posts/')

db_name=os.getenv('DB_NAME')
db_user=os.getenv('DB_USER')
db_pass=os.getenv('DB_PASSWORD')
db_host=os.getenv('DB_HOST')
db_port=os.getenv('DB_PORT')

# idea for how to do this: load file names into an array and do a for loop that takes each file and uploads it into database
def upload_blog_post(title, post, created_date, summary, slug, related_forecasts):
    with psycopg2.connect(dbname=db_name, user=db_user, password=db_pass, host=db_host) as conn:
        cursor=conn.cursor()
        cursor.execute('''INSERT INTO blogposts 
                       (title, post, created_date, summary, slug, related_forecasts) 
                       VALUES (%s, %s, %s, %s, %s)''', 
                       (title, post, created_date,summary,slug,related_forecasts))

# load files into array
creation_date = datetime.now().date()
while True:
    file_path = input('Enter the file path: ')
    post_title = input('Enter the post title: ')
    post_summary=input('Enter the post summary: ')
    post_slug=input('Enter the slug: ')
    post_related_forecasts=input('Enter the forecast ids of the related forecasts: ')
    
    # use file_path to load file
    with open(file_path, 'r', encoding='utf-8') as file:
        markdown_content = file.read()
    
    upload_blog_post(post_title, markdown_content, creation_date,post_summary, post_slug,post_related_forecasts)

