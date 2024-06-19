**TASK 1 (Design the Review Chatbot):**

    - START:
        - Bot starts conversation by thanking user for buying and trusting the store, it addresses the user by name and ask nicely about the user's experience with the product by name (for example Iphone 13) on a manual message.

    - END:
        - We don't want the bot to ask a lot of questions and annoy the user just a simple question and after an answer is given thank him with a manual message and let him know you can assist him with whatever else is on his mind.

    - PROMPT:
        - The bot should have a prompt that it should only answer to messages are that clear and relevant to the store, else it should ask for rephrase - or say something like "I didn't catch that, please be more specific, if I can't help you please contact the human team at human@team.com".



**TASK 2 (Model the Review Chatbot):**
    Relational Database (MySQL) and GOLANG ORM (GORM).
    Create a html page for each table (3 in total).

    - Tables:

        - Customer Information:
            - PRIMARY_KEY : AUTO_GENERATED
            - Name
            - Email
            - Region

        - Chatbot Interations:
            - PRIMARY_KEY : AUTO_GENERATED
            - Name
            - Chat History

        - Customer Reviews:
            - PRIMARY_KEY : AUTO_GENERATED
            - Name
            - Items Bought
            - Star-Review ⭐



**TASK 3 (Implement the Review Chatbot):**
    - 3 Demo Items (Phone, GPU, Monitor)..
    - Check 1 item and click "BUY"..
    - That "BUY" button will trigger the start-message of the bot.



**TASK 4 (Demonstrate the Chatbot):**
    - Each time the button "BUY" is clicked the bot will get the PRIMARY_KEY of the user.

    - After the end message (TASK 1 - END) the bot will check the PRIMARY_KEY of the user and modify the tables:
        - Chatbot Interations:
            - ADDING NEW CHAT HISTORY

        - Customer Reviews:
            - ADDING NEW ITEM-REVIEW



**EXTRA TASK:**
    - On the user's answer about the review, do sentiment analysis instead of Star-Review ⭐