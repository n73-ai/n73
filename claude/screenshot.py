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

take_screenshot("https://project-d88c4024-e85e-4928-850d-a52fc4b2d01a.pages.dev/")
