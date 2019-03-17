# load csv from Text API
from collections import defaultdict
import requests
import datetime
import csv

PROD_URL = "http://35.198.123.101:5000/"
DATA_PATH = "./data/"


def main():
    r = requests.get(PROD_URL)
    json_response = r.json()
    now = datetime.datetime.now()
    file_name = now.strftime("%d_%B_%Y__%H_%M") + ".csv"
    save_csv(json_response, file_name)


def save_csv(json_response, file_name):
    idk_dict = defaultdict(int)
    hs_dict = defaultdict(int)
    not_hs_dict = defaultdict(int)
    all_keys = []
    for text in json_response:
        content = text["Content"].lstrip()
        all_keys.append(content)
        is_hs = text["IsHS"]
        is_not_hs = text["IsNotHS"]
        idk = text["Idk"]
        if not is_hs and not is_not_hs and not idk:
            continue
        else:
            if idk:
                idk_dict[content] += 1
            elif is_hs:
                hs_dict[content] += 1
            elif is_not_hs:
                not_hs_dict[content] += 1
    with open(DATA_PATH+file_name, 'w') as data_file:
        for key in all_keys:
            row = [key, str(hs_dict[key]), str(not_hs_dict[key]), str(idk_dict[key])]
            wr = csv.writer(data_file, quoting=csv.QUOTE_ALL)
            wr.writerow(row)


if __name__ == "__main__":
    main()
