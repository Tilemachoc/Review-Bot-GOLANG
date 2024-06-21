import openai
import os
from dotenv import load_dotenv
from typing import Text, List
from transformers import AutoTokenizer, AutoModelForSequenceClassification
import torch

load_dotenv()

api_key = os.getenv("OPENAI_API_KEY")
openai.api_key = api_key

client = openai.OpenAI()

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
        history[0].content = instructions
    else:
        history = [{"role": "system", "content": instructions},{"role": "assistant", "content":f"Hi John 👋\nYou recently received the {product} you ordered, can you tell us about your experience?"},{"role": role, "content": message}]
    return history


def response_message(client: openai.OpenAI, history: List) -> str:
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
    return assistant_text


def check_keyword(msg: str) -> str:
    keywords = {"refund": 0, "information": 0, "rating": 0}

    for word in msg.lower().split():
        if word in keywords:
            keywords[word] += 1

    if sum(keywords.values()) == 0:
        return "general"
    elif any(keywords[i] == keywords[j] for i,j in keywords.items() if i != j):
        return "general"
    else:
        return max(keywords, key=keywords.get)


def refundHandler(product: str):
    return f"We're sorry to hear the {product} wasn't a perfect fit!  A full refund will be issued to your original payment method within [number] business days."

def informationHandler():
    return "For info about our products, site, features, or FAQs, head over to our website: http://localhost:8080/buy.  We're happy to help if you have any further questions!"


def ratingHandler(msg: str):
    found = 0
    number = None
    for char in msg:
        if char.isdigit():
            number = int(char)
            if 1 <= number <= 5:
                found += 1
                if found> 1:
                    return get_sentiment_rating(msg)
    return number if found == 1 else get_sentiment_rating(msg)


def get_sentiment_rating(msg: str) -> int:
    tokenizer = AutoTokenizer.from_pretrained("LiYuan/amazon-review-sentiment-analysis")
    model = AutoModelForSequenceClassification.from_pretrained("LiYuan/amazon-review-sentiment-analysis")
    inputs = tokenizer(msg, return_tensors="pt")
    outputs = model(**inputs)
    return torch.argmax(outputs.logits, dim=-1)


def generalHandler(history: List) -> str:
    global client
    return response_message(client,history)


func_map = {
    "refund": refundHandler,
    "information": informationHandler,
    "rating": ratingHandler,
    "general": generalHandler,
}



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
    response_intent = response_message(client, history_intent)
    # have a function that will look for keywords
    keyword = check_keyword(response_intent)
    # for keyword have a map that will map the function for the automatic message or the sentiment analysis
    #we need history to have False intent so the right instructions are placed, so we used a history_intent for the intent response
    history = add_history("user", message, product, False, history)
    if keyword == "refund":
        response_msg = refundHandler(product)
    if keyword == "information":
        response_msg = informationHandler()
    if keyword == "rating":
        response_msg = ratingHandler(message)
    if keyword == "general":
        response_msg = generalHandler(history)
    return {"message": response_msg, "history": history}