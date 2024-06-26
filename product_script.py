import openai
import os
from dotenv import load_dotenv
from typing import Text, List
from transformers import AutoTokenizer, AutoModelForSequenceClassification
import torch
import requests
import json
import logging


# logging helping functions


def logging_variable(name, variable):
    logging.basicConfig(filename="variable_logs.log", level=logging.ERROR)
    logging.info("%s= %s" % (name,variable))


def logging_function(filename):
    def decorator(func):
        def wrapper(*args, **kwargs):
            result = func(*args, **kwargs)
            with open(filename, 'a') as f:
                f.write(f"Function '{func.__name__}' returned: {result}\n")
            
            return result
        return wrapper
    return decorator


load_dotenv()

api_key = os.getenv("OPENAI_API_KEY")
openai.api_key = api_key

client = openai.OpenAI()

#@logging_function("function_logs.log")
def add_history(role: str, message: Text, product: str, intent: bool=True, history: List = None):
    if intent:
        instructions = """
        You are an AI bot that extracts the main intent from customer messages. Based on the provided message, identify if the customer is:

        Providing a rating (1 to 5 or any other rating) - respond with: RATING
        Asking for a refund - respond with: REFUND
        Requesting information - respond with: INFORMATION
        Any other type of message - respond with: GENERAL
        Respond with only one word, and ensure it accurately reflects the customer's intent.

        Example Messages:
        "I would like to return my purchase and get my money back." -> REFUND


        "Can you tell me the operating hours of your store?" -> INFORMATION


        "I had a great experience! 5 stars!" -> RATING


        "I need help with my order." -> GENERAL


        """
    else:
        instructions = f"You are iStore AI assistant bot assisting collection information about from customers about their experience with the {product} and answering any questions they might have. Please provide relevant information or assistance on these topics. Ensure that the responses are clear, concise, and helpful to the customer.If the customer's question or statement is about irrelevant topics, politely acknowledge it with a brief apology.In order to assist the user further ask them to provide a star review rating for the {product}"
    if history:
        history.append({"role": role, "content": message})
        history[0]["content"] = instructions
    else:
        history = [{"role": "system", "content": instructions},{"role": "assistant", "content":f"Hi John 👋\nYou recently received the {product} you ordered, can you tell us about your experience?"},{"role": role, "content": message}]
    return history


@logging_function("function_logs.log")
def response_message(client: openai.OpenAI, history: List, role: int) -> str:
    try:
        response = client.chat.completions.create(
            model="gpt-3.5-turbo",
            messages = history,
            temperature=0,
            max_tokens=50,
            top_p=0.2,
            frequency_penalty=0,
            presence_penalty=0,
            n=1
        )

        assistant_text = response.choices[0].message.content.strip()
    except Exception as e:
        print("Error at response message:",e)
        if role == "intent-assistant":
            # TEST intents: REFUND, INFORMATION, RATING, GENERAL
            return "rating is  lool"
        elif role == "assistant":
            return "Hello! Thanks for messaging me!"
    return assistant_text


@logging_function("function_logs.log")
def check_keyword(msg: str) -> str:
    try:
        keywords = {"refund": 0, "information": 0, "rating": 0, "general": 0}

        for word in msg.lower().split():
            if word in keywords:
                keywords[word] += 1
        if sum(keywords.values()) == 0:
            return "general"
        
        sorted_values = sorted(keywords.values(), reverse=True)
        if sorted_values[0] == sorted_values[1]:
            return "general"
    
        else:
            #logging_variable("max(keywords, key=keywords.get)", max(keywords, key=keywords.get))
            return max(keywords, key=keywords.get)
    except Exception as e:
        logging.error("Error in check_keyword function:", e)
        return "err"


@logging_function("function_logs.log")
def refundHandler(product: str):
    return f"We're sorry to hear the {product} wasn't a perfect fit! \n\nA full refund will be issued to your original payment method within 14 business days."


@logging_function("function_logs.log")
def informationHandler():
    return "For info about our products, site, features, or FAQs, head over to our website: http://localhost:8080/buy.  We're happy to help if you have any further questions!"


@logging_function("function_logs.log")
def ratingHandler(msg: str):
    found = 0
    number = None
    for char in msg:
        if char.isdigit():
            number = int(char)
            if 1 <= number <= 5:
                found += 1
                if found> 1:
                    return "Thank you for your rating!", get_sentiment_rating(msg)
    return "Thank you for your rating!", number if found == 1 else get_sentiment_rating(msg)


@logging_function("function_logs.log")
def get_sentiment_rating(msg: str) -> int:
    try:
        tokenizer = AutoTokenizer.from_pretrained("LiYuan/amazon-review-sentiment-analysis")
        model = AutoModelForSequenceClassification.from_pretrained("LiYuan/amazon-review-sentiment-analysis")
        inputs = tokenizer(msg, return_tensors="pt")
        outputs = model(**inputs)
    except Exception as e:
        print("Error at get_sentiment_rating:",e)
        return 3
    return torch.argmax(outputs.logits, dim=-1).item()+1 #+1 because it return index


@logging_function("function_logs.log")
def generalHandler(history: List) -> str:
    global client
    return response_message(client,history,"assistant")


func_map = {
    "refund": refundHandler,
    "information": informationHandler,
    "rating": ratingHandler,
    "general": generalHandler,
}


#@logging_function("function_logs.log")
def get_response(message, product, history=None) -> dict:
    global client
    global func_map
    #---Pseudocoding---
    #We get the message and add it to the history
    #We ask the model if the user's message shows: rating-refund-information
    #Ask it to answer with function-calling format but we will extract the key_work with python work
    #Rating -> exract number if there is in the msg or another bot to rate from 1 to 5
    #Refund -> Automatic message to refund {product}
    #Information -> Automatic message for recommanding List - product: http://localhost:8080/buy
    #General -> automatic message, bot with a base prompt
    #---Psuedocoding---
    #add history
    history_intent = add_history("user", message, product, True, history)
    #get ai answer
    response_intent = response_message(client, history_intent, "intent-assistant")
    # have a function that will look for keywords
    keyword = check_keyword(response_intent)
    # for keyword have a map that will map the function for the automatic message or the sentiment analysis
    #we need history to have False intent so the right instructions are placed, so we used a history_intent for the intent response
    history = add_history("user", message, product, False, history)
    if keyword == "refund":
        response_msg = refundHandler(product)
    elif keyword == "information":
        response_msg = informationHandler()
    elif keyword == "rating":
        response_msg, rating = ratingHandler(message)
        respond_code = send_rating(rating=rating, user_message=message)
        if respond_code == 201:
            send_orderitem(product)
    elif keyword == "general":
        response_msg = generalHandler(history)
    else:
        response_msg = "I'm sorry, I didn't understand your request."
    history = add_history("assistant", response_msg, product, True, history)
    return {"message": response_msg, "history": history, "product": product}


#It would be better practice to instead get all the information from golang, maybe cookies but it's much more difficult
@logging_function("function_logs.log")
def send_rating(rating: int, user_message, Id = 2, orderitemid = 2, user_id = 1):
    url = "http://localhost:8080/api/reviews"
    rating_data = {
        "review_id": Id,
        "order_item_id": orderitemid,
        "user_id": user_id,
        "rating": rating,
        "review_text": user_message
    }
    logging_variable("rating_data", rating_data)
    headers = {
        "Content-Type": "application/json"
    }

    response = requests.post(url, data=json.dumps(rating_data), headers=headers)
    logging_variable("send_rating_response", response)
    return response.status_code




@logging_function("function_logs.log")
def send_orderitem(product: str, orderitem_id:int = 3, order_id:int = 1):
    ItemPriceMap = {
        "iphone-13": 649.00,
        "gpu": 2247.49,
        "monitor": 899.20,
    }

    url = "http://localhost:8080/api/orderitems"
    orderitem_data = {
        "orderitem_id":orderitem_id,
        "order_id": order_id,
        "item_name": product,
        "item_price": ItemPriceMap.get(product, 0)
    }

    headers = {
        "Content-Type": "application/json"
    }
    
    headers = {
        "Content-Type": "application/json"
    }

    response = requests.post(url, data=json.dumps(orderitem_data), headers=headers)
    logging_variable("send_orderitem_response", response)
    return response.status_code