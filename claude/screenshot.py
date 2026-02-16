from selenium import webdriver
import time

def take_screenshot(url):
    options = webdriver.ChromeOptions()
    options.add_argument("--headless")
    options.add_argument("--window-size=1920,1080") 
    driver = webdriver.Chrome(options=options)

    driver.get(url)
    time.sleep(5)
    # save the image as screenshot.png
    # send it to the main backend
    # upload to bunny net bucket
    # save the url in project.image
    driver.save_screenshot("/app/screenshot.png")
    driver.quit()
